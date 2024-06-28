package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/edatts/go-payment-system/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
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
	if err == nil {
		utils.WriteCustomError(rw, http.StatusBadRequest, ErrUserExists(user.Email))
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

}
