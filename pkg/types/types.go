package types

import "time"

type RegisterUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=120"`
}

type User struct {
	Id        int32     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStore interface {
	CreateUser(*User) error
	GetUser(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int32) (*User, error)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required_without=Email"`
	Email    string `json:"email" validate:"required_without=Username,email"`
	Password string `json:"password" validate:"required"`
}

type Account struct {
	Id         int32     `json:"id"`
	UserId     int32     `json:"userId"`
	CurrencyId int32     `json:"currencyId"`
	Balance    int64     `json:"balance"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Currency struct {
	Id       int32
	Name     string
	Ticker   string
	Decimals string
}

type Transaction struct {
	Id         int64 `json:"id"`
	SenderId   int32 `json:"senderId"`
	ReceiverId int32 `json:"receiverId"`
	// SenderName   string   `json:"senderName"`
	// ReceiverName string   `json:"receiverName"`
	CurrencyId int32     `json:"currencyId"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"createdAt"`
}

type DepositRequest struct {
	// AccountId      uint64 `json:"accountId" validate:"required"`
	CurrencyTicker string `json:"currencyTicker" validate:"required"`
	Amount         int64  `json:"amount" validate:"required"`
}

type Deposit struct {
	Id         int64
	AccountId  int32
	CurrencyId int32
	Amount     uint64
	CreatedAt  time.Time
}

type WithdrawalRequest struct {
	// AccountId      uint64 `json:"accountId" validate:"required"`
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

type PaymentsStore interface {
	CreateAccount(*Account) error
	GetAccount(userId int32, currencyTicker string) (*Account, error)
	GetAccountBalance(accountId int32) (int64, error)
	UpdateAccountBalance(accountId int32, balance int64) error
	CreateTransaction(*Transaction) error
	CreateDeposit(*Deposit) error
	CreateWithdrawal(*Withdrawal) error
}
