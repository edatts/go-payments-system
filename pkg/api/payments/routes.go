package payments

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/edatts/go-payment-system/pkg/jwt"
	"github.com/edatts/go-payment-system/pkg/store"
	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/edatts/go-payment-system/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

// type CtxKey string
// const JWTPubCtxKey = CtxKey("JWTPublicKey")

type Handler struct {
	// store types.PaymentsStore
	store       store.PaymentsStorer
	jwtPubkey   []byte
	jwtPubkeyMu *sync.RWMutex
}

// func NewHandler(store types.PaymentsStore) *Handler {
// 	return &Handler{
// 		store:       store,
// 		jwtPubkeyMu: &sync.RWMutex{},
// 	}
// }

func NewHandler(store *store.PaymentsStore) *Handler {
	// func NewHandler(store types.PaymentsStore) *Handler {
	return &Handler{
		store:       store,
		jwtPubkeyMu: &sync.RWMutex{},
	}
}

func (h *Handler) Init() *Handler {
	var retries = 5
	for i := 0; i < retries; i++ {
		time.Sleep(1 * time.Second)
		if err := h.updateJWTPublicKey(); err != nil {
			log.Printf("failed updating jwt public key: %s", err)
			continue
		}

		return h
	}

	log.Fatalf("failed getting jwt public key after %d retries", retries)
	return h
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/deposit", h.withJWTPubkey(jwt.WithJWTAuth(h.HandleDeposit, h.store))).Methods("POST")
	router.HandleFunc("/withdraw", h.HandleWithdraw).Methods("POST")
	router.HandleFunc("/transfer", h.HandleTransfer).Methods("POST")
}

func (h *Handler) HandleDeposit(rw http.ResponseWriter, req *http.Request) {
	// Determine the UserId from the request context, this is set
	// during authentication.
	userId, ok := jwt.GetUserIdFromContext(req.Context())
	if !ok {
		log.Printf("failed getting user id from context in deposit handler")
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	var depositReq = &types.DepositRequest{}
	if err := utils.ReadRequestJSON(req, depositReq); err != nil {
		log.Printf("failed reading deposit request json: %s", err)
		utils.WriteError(rw, http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(depositReq); err != nil {
		errs := err.(validator.ValidationErrors)
		log.Printf("validation errors: %v", errs)
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrFailedValidation(errs))
		return
	}

	tx, err := h.store.BeginTx()
	if err != nil {
		log.Printf("failed starting database transaction: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := tx.Rollback(context.Background()); err != nil {
			log.Printf("failed rolling back database transaction: %s", err)
		}
	}()

	txQueries := h.store.WithTx(tx)

	// Get currency
	currency, err := txQueries.GetCurrency(depositReq.CurrencyTicker)
	if err != nil {
		log.Printf("failed getting currency from database: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	// Get account
	acc, err := txQueries.GetAccount(userId, depositReq.CurrencyTicker)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("failed getting account from database: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// Create account
		acc = &types.Account{
			UserId:     userId,
			CurrencyId: currency.Id,
			Balance:    0,
		}

		if err := txQueries.CreateAccount(acc); err != nil {
			log.Printf("failed creating new account in the database: %s", err)
			utils.WriteError(rw, http.StatusInternalServerError)
			return
		}
	}

	// Calculate new balance
	var newBalance int64 = acc.Balance + depositReq.Amount

	// Update deposit
	if err := txQueries.UpdateAccountBalance(acc.Id, newBalance); err != nil {
		log.Printf("failed updating balance for account (%v): %s", acc.Id, err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	var deposit = &types.Deposit{
		AccountId:  acc.Id,
		CurrencyId: currency.Id,
		Amount:     depositReq.Amount,
	}

	// Insert deposit
	if err := txQueries.CreateDeposit(deposit); err != nil {
		log.Printf("failed inserting deposit into database: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Printf("failed commiting database transaction: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if err := utils.WriteJSON(rw, http.StatusOK, nil); err != nil {
		log.Printf("failed writing response json: %s", err)
	}
}

func (h *Handler) HandleWithdraw(rw http.ResponseWriter, req *http.Request) {
	// Get UserId from the request context
	userId, ok := jwt.GetUserIdFromContext(req.Context())
	if !ok {
		log.Printf("failed getting user id from context in deposit handler")
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	var withdrawalReq = &types.WithdrawalRequest{}
	if err := utils.ReadRequestJSON(req, withdrawalReq); err != nil {
		log.Printf("failed reading withdrawal request json: %s", err)
		utils.WriteError(rw, http.StatusBadRequest)
		return
	}

	// Validate request
	if err := utils.Validate.Struct(withdrawalReq); err != nil {
		errs := err.(validator.ValidationErrors)
		log.Printf("failed validating withdrawal request: %s", err)
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrFailedValidation(errs))
		return
	}

	tx, err := h.store.BeginTx()
	if err != nil {
		log.Printf("failed starting database transaction: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := tx.Rollback(context.Background()); err != nil {
			log.Printf("failed rolling back database transaction: %s", err)
		}
	}()

	txQueries := h.store.WithTx(tx)

	currency, err := txQueries.GetCurrency(withdrawalReq.CurrencyTicker)
	if err != nil {
		log.Printf("failed getting currency from database: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	// Get account
	acc, err := txQueries.GetAccount(userId, withdrawalReq.CurrencyTicker)
	if errors.Is(err, pgx.ErrNoRows) {
		// Account not exists
		log.Printf("account not found in database: %s", err)
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrAccountNotExists)
		return
	}

	if err != nil {
		log.Printf("failed getting account from database: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if withdrawalReq.Amount > acc.Balance {
		log.Printf("not enough funds for withdrawal")
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrNotEnoughFunds)
		return
	}

	var newBalance int64 = acc.Balance - withdrawalReq.Amount

	// Update balance
	if err := txQueries.UpdateAccountBalance(acc.Id, newBalance); err != nil {
		log.Printf("failed updating account balance: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	var withdrawal = &types.Withdrawal{
		AccountId:  acc.Id,
		CurrencyId: currency.Id,
		Amount:     withdrawalReq.Amount,
	}

	// Insert withdrawal
	if err := txQueries.CreateWithdrawal(withdrawal); err != nil {
		log.Printf("failed inserting withdrawal into db: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Printf("failed commiting database transaction: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if err := utils.WriteJSON(rw, http.StatusOK, nil); err != nil {
		log.Printf("failed writing json response: %s", err)
	}
}

func (h *Handler) HandleTransfer(rw http.ResponseWriter, req *http.Request) {

	senderUserId, ok := jwt.GetUserIdFromContext(req.Context())
	if !ok {
		log.Printf("failed getting user id from context in deposit handler")
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	var transferReq = &types.TransferRequest{}
	if err := utils.ReadRequestJSON(req, transferReq); err != nil {
		log.Printf("failed reading request json: %s", err)
		utils.WriteError(rw, http.StatusBadRequest)
		return
	}

	if err := utils.Validate.Struct(transferReq); err != nil {
		log.Printf("failed validation for transfer request: %s", err)
		utils.WriteError(rw, http.StatusBadRequest)
		return
	}

	tx, err := h.store.BeginTx()
	if err != nil {
		log.Printf("failed starting database transaction: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := tx.Rollback(context.Background()); err != nil {
			log.Printf("failed rolling back database transaction: %s", err)
		}
	}()

	txQueries := h.store.WithTx(tx)

	// Get accounts of both parties
	senderAcc, err := txQueries.GetAccount(senderUserId, transferReq.CurrencyTicker)
	if err != nil {
		log.Printf("failed getting sender account: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
	}

	recipientUser, err := txQueries.GetUser(transferReq.RecipientUsername)
	if err != nil {
		log.Printf("failed getting receiving user from database: %s", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("user not found: %s", err)
			utils.WriteCustomError(rw, http.StatusBadRequest, ErrRecipientNotExsits)
			return
		}
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	recipientAcc, err := txQueries.GetAccount(recipientUser.Id, transferReq.CurrencyTicker)
	if err != nil {
		log.Printf("failed getting recipient account form database: %s", err)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("account not found: %s", err)
			utils.WriteCustomError(rw, http.StatusBadRequest, ErrRecipientAccountNotExists)
			return
		}
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	// Check balance of sender
	if senderAcc.Balance < transferReq.Amount {
		log.Printf("insufficient balance to complete transaction")
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrNotEnoughFunds)
		return
	}

	// Insert transfer into db
	var transfer = &types.Transfer{
		SenderId:    senderUserId,
		RecipientId: recipientUser.Id,
		CurrencyId:  senderAcc.CurrencyId,
		Amount:      transferReq.Amount,
	}

	if err := txQueries.CreateTransfer(transfer); err != nil {
		log.Printf("failed inserting transfer into database: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	// Update balances of parties
	newSenderBalance := senderAcc.Balance - transfer.Amount
	newRecipientBalance := recipientAcc.Balance + transfer.Amount

	if err := txQueries.UpdateAccountBalance(senderAcc.Id, newSenderBalance); err != nil {
		log.Printf("failed updating sender balance: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if err := txQueries.UpdateAccountBalance(recipientAcc.Id, newRecipientBalance); err != nil {
		log.Printf("failed updating recipient balance: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Printf("failed committing database transaction: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if err := utils.WriteJSON(rw, http.StatusOK, nil); err != nil {
		log.Printf("failed writing json response: %s", err)
	}
}

func (h *Handler) updateJWTPublicKey() error {

	res, err := http.Get("auth:4444/jwt-public-key")
	if err != nil {
		return fmt.Errorf("failed getting jwt public key from auth service: %w", err)
	}

	jsonBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed reading response body: %w", err)
	}

	var getPubRes = new(types.GetJWTPublicKeyResponse)
	if err := json.Unmarshal(jsonBytes, getPubRes); err != nil {
		return fmt.Errorf("failed unmarshalling json bytes: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(getPubRes.PublicKey)
	if err != nil {
		return fmt.Errorf("failed decoding base64 string: %w", err)
	}

	h.jwtPubkey = decoded

	return nil
}

func (h *Handler) getJWTPublicKey() []byte {
	h.jwtPubkeyMu.Lock()
	defer h.jwtPubkeyMu.Unlock()
	return h.jwtPubkey
}

func (h *Handler) withJWTPubkey(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// Add public key to request context.
		ctx := context.WithValue(req.Context(), jwt.JWTCtxKey, h.getJWTPublicKey())
		req = req.WithContext(ctx)
		handlerFunc(rw, req)
	}
}
