package types

import (
	"time"
)

type RegisterUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=120"`
}

type User struct {
	Id        int32
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type LoginRequest struct {
	Username string `json:"username" validate:"required_without=Email,excluded_with=Email"`
	Email    string `json:"email" validate:"required_without=Username,excluded_with=Username"`
	Password string `json:"password" validate:"required"`
}

type Account struct {
	Id         int32
	UserId     int32
	CurrencyId int32
	Balance    int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Currency struct {
	Id       int32
	Name     string
	Ticker   string
	Decimals string
}

type TransferRequest struct {
	RecipientUsername string `json:"recipientUsername" validate:"required"`
	CurrencyTicker    string `json:"currencyTicker" validate:"required"`
	Amount            int64  `json:"amount" validate:"required"`
}

type Transfer struct {
	Id          int64
	SenderId    int32
	RecipientId int32
	CurrencyId  int32
	Amount      int64
	CreatedAt   time.Time
}

type DepositRequest struct {
	CurrencyTicker string `json:"currencyTicker" validate:"required"`
	Amount         int64  `json:"amount" validate:"required"`
}

type Deposit struct {
	Id         int64
	AccountId  int32
	CurrencyId int32
	Amount     int64
	CreatedAt  time.Time
}

type WithdrawalRequest struct {
	CurrencyTicker string `json:"currencyTicker" validate:"required"`
	Amount         int64  `json:"amount" validate:"required"`
}

type Withdrawal struct {
	Id         int64
	AccountId  int32
	CurrencyId int32
	Amount     int64
	CreatedAt  time.Time
}

type GetJWTPublicKeyResponse struct {
	PublicKey string
}
