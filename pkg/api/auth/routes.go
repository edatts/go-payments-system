package auth

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"github.com/edatts/go-payment-system/pkg/config"
	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/edatts/go-payment-system/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
}

func (h *Handler) handleRegister(rw http.ResponseWriter, req *http.Request) {
	var registerUserReq types.RegisterUserRequest

	// Receives JSON payload
	if err := utils.ReadRequestJSON(req, &registerUserReq); err != nil {
		log.Printf("failed reading request json: %s", err)
		utils.WriteError(rw, http.StatusBadRequest)
		return
	}

	// Validate
	if err := utils.Validate.Struct(registerUserReq); err != nil {
		errs := err.(validator.ValidationErrors)
		log.Printf("validation errors: %v", errs)
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrFailedValidation(errs))
		return
	}

	// Checks if the user already exists
	//	- Creates account if user not exists
	// 	- Returns error to caller is user exists
	user, err := h.store.GetUserByEmail(registerUserReq.Email)
	if err == nil && user.Username == registerUserReq.Username {
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrUsernameExists(user.Username))
		return
	}

	if err == nil && user.Email == registerUserReq.Email {
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrEmailExists(user.Email))
		return
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("failed getting user by email: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	hashedPassword, err := HashPassword(registerUserReq.Password)
	if err != nil {
		log.Printf("failed hashing password: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	user = &types.User{
		FirstName: registerUserReq.FirstName,
		LastName:  registerUserReq.LastName,
		Username:  registerUserReq.Username,
		Email:     registerUserReq.Email,
		Password:  hashedPassword,
	}

	if err := h.store.CreateUser(user); err != nil {
		log.Printf("failed creating user: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
	}

	if err := utils.WriteJSON(rw, http.StatusCreated, nil); err != nil {
		log.Printf("failed writing create user response: %s", err)
	}

}

func (h *Handler) handleLogin(rw http.ResponseWriter, req *http.Request) {
	var loginReq types.LoginRequest

	// Receives JSON payload
	if err := utils.ReadRequestJSON(req, &loginReq); err != nil {
		log.Printf("failed reading request json: %s", err)
		utils.WriteError(rw, http.StatusBadRequest)
		return
	}

	// Validate
	if err := utils.Validate.Struct(loginReq); err != nil {
		errs := err.(validator.ValidationErrors)
		log.Printf("validation errors: %v", errs)
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrFailedValidation(errs))
		return
	}

	var (
		user = new(types.User)
		err  error
	)
	switch {
	case loginReq.Username != "":
		user, err = h.store.GetUser(loginReq.Username)
	case loginReq.Email != "":
		user, err = h.store.GetUserByEmail(loginReq.Email)
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows) && loginReq.Username != "":
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrUsernameNotExists(loginReq.Username))
		return
	case errors.Is(err, pgx.ErrNoRows) && loginReq.Email != "":
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrEmailNotExists(loginReq.Email))
		return
	case err != nil:
		log.Printf("failed getting user from db: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	err = VerifyPassword(loginReq.Password, user.Password)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrWrongPassword)
		return
	}

	if err != nil {
		log.Printf("failed verifying password: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	secret, err := base64.StdEncoding.DecodeString(config.Envs.JWT_SECRET)
	if err != nil {
		log.Printf("failed decoding jwt secret from base64: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	token, err := CreateJWT(secret, user.Id)
	if err != nil {
		log.Printf("failed creating json web token: %s", err)
		utils.WriteError(rw, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(rw, http.StatusOK, map[string]string{"JWT": token})
}
