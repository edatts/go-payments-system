package payments

import (
	"errors"
	"log"
	"net/http"

	"github.com/edatts/go-payment-system/pkg/api/auth"
	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/edatts/go-payment-system/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	store types.PaymentsStore
}

func NewHandler(store types.PaymentsStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/deposit", h.HandleDeposit).Methods("POST")
	router.HandleFunc("/withdraw", h.HandleWithdraw).Methods("POST")
	router.HandleFunc("/pay", h.HandlePay).Methods("POST")
}

func (h *Handler) HandleDeposit(rw http.ResponseWriter, req *http.Request) {
	// Determine the UserId from the request context, this is set
	// during authentication.
	userId, ok := auth.GetUserIdFromContext(req.Context())
	if !ok {
		log.Printf("failed getting user id from context in deposit handler")
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	var depositReq types.DepositRequest
	if err := utils.ReadRequestJSON(req, &depositReq); err != nil {
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

	// Get account
	acc, err := h.store.GetAccount(userId, depositReq.CurrencyTicker)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("failed getting account from database: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// Create account

	}

	// Calculate new balance
	var newBalance int64

	// Make deposit
	if err := h.store.UpdateAccountBalance(acc.Id, newBalance); err != nil {
		log.Printf("failed updating balance for account (%v): %s", acc.Id, err)
		utils.WriteError(rw, http.StatusInternalServerError)
	}
}

func (h *Handler) HandleWithdraw(rw http.ResponseWriter, req *http.Request) {

}

func (h *Handler) HandlePay(rw http.ResponseWriter, req *http.Request) {

}
