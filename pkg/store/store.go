package store

import (
	"context"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Executor interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type PaymentsQueryHandler interface {
	GetUser(username string) (*types.User, error)
	GetUserById(userId int32) (*types.User, error)
	GetCurrency(ticker string) (*types.Currency, error)
	GetAccount(userId int32, currencyTicker string) (*types.Account, error)
	CreateAccount(*types.Account) error
	UpdateAccountBalance(accountId int32, balance int64) error
	CreateDeposit(*types.Deposit) error
	CreateWithdrawal(*types.Withdrawal) error
	CreateTransfer(*types.Transfer) error
}

type PaymentsStorer interface {
	BeginTx() (pgx.Tx, error)
	WithTx(tx pgx.Tx) PaymentsQueryHandler
	PaymentsQueryHandler
}

var _ PaymentsQueryHandler = (*PaymentsQueries)(nil)
var _ PaymentsStorer = (*PaymentsStore)(nil)

type QueryHandler interface {
	CreateUser(*types.User) error
	GetUser(username string) (*types.User, error)
	GetUserByEmail(email string) (*types.User, error)
	GetUserById(id int32) (*types.User, error)
}

type UserStorer interface {
	BeginTx() (pgx.Tx, error)
	QueryHandler
}

var _ QueryHandler = (*UserQueries)(nil)
var _ UserStorer = (*UserStore)(nil)
