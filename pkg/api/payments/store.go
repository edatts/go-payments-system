package payments

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	db *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{
		db: dbConn,
	}
}

func (s *Store) CreateAccount(acc *types.Account) error {

	return nil
}

func (s *Store) GetAccount(userId int32, currencyTicker string) (*types.Account, error) {
	var acc = new(types.Account)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM accounts INNER JOIN currency ON accounts.currency_id = currency.id WHERE user_id = $1 AND ticker = $2;", userId, currencyTicker)
	if err := row.Scan(acc); err != nil {
		return &types.Account{}, fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return acc, nil
}

func (s *Store) GetAccountBalance(accountId int32) (int64, error) {

	return 0, nil
}

func (s *Store) UpdateAccountBalance(accountId int32, balance int64) error {

	return nil
}

func (s *Store) CreateDeposit(*types.Deposit) error {

	return nil
}

func (s *Store) CreateWithdrawal(*types.Withdrawal) error {

	return nil
}

func (s *Store) CreateTransaction(tx *types.Transaction) error {

	return nil
}
