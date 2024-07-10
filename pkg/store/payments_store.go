package store

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentsStore struct {
	db *pgxpool.Pool
}

func NewPaymentsStore(db *pgxpool.Pool) *PaymentsStore {
	return &PaymentsStore{
		db: db,
	}
}

func (s *PaymentsStore) CreateAccount(acc *types.Account) error {
	row := s.db.QueryRow(context.Background(), "INSERT INTO accounts (user_id, currency_id, balance) VALUES ($1, $2, $3) RETURNING *;", acc.UserId, acc.CurrencyId, acc.Balance)
	if err := row.Scan(&acc.Id, &acc.UserId, &acc.CurrencyId, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
		return fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return nil
}

func (s *PaymentsStore) GetAccount(userId int32, currencyTicker string) (*types.Account, error) {
	var acc = new(types.Account)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM accounts INNER JOIN currencies ON accounts.currency_id = currencies.id WHERE user_id = $1 AND ticker = $2;", userId, currencyTicker)
	if err := row.Scan(&acc.Id, &acc.UserId, &acc.CurrencyId, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
		return &types.Account{}, fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return acc, nil
}

func (s *PaymentsStore) UpdateAccountBalance(accountId int32, balance int64) error {
	_, err := s.db.Exec(context.Background(), "UPDATE accounts SET balance = $1 WHERE id = $2;", balance, accountId)
	if err != nil {
		return fmt.Errorf("failed executing update query: %w", err)
	}

	return nil
}

func (s *PaymentsStore) CreateDeposit(dep *types.Deposit) error {
	row := s.db.QueryRow(context.Background(), "INSERT INTO deposits (account_id, currency_id, amount) VALUES ($1, $2, $3) RETURNING *;", dep.AccountId, dep.CurrencyId, dep.Amount)
	if err := row.Scan(&dep.Id, &dep.AccountId, &dep.CurrencyId, &dep.Amount, &dep.CreatedAt); err != nil {
		return fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return nil
}

func (s *PaymentsStore) CreateWithdrawal(withdrawal *types.Withdrawal) error {
	row := s.db.QueryRow(context.Background(), "INSERT INTO withdrawals (account_id, currency_id, amount) VALUES ($1, $2, $3) RETURNING *;", withdrawal.AccountId, withdrawal.CurrencyId, withdrawal.Amount)
	if err := row.Scan(&withdrawal.Id, &withdrawal.AccountId, &withdrawal.CurrencyId, &withdrawal.Amount, &withdrawal.CreatedAt); err != nil {
		return fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return nil
}

func (s *PaymentsStore) CreateTransaction(tx *types.Transaction) error {

	return nil
}

func (s *PaymentsStore) GetCurrency(ticker string) (*types.Currency, error) {
	curr := new(types.Currency)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM currencies WHERE ticker = $1", ticker)
	if err := row.Scan(&curr.Id, &curr.Name, &curr.Ticker, &curr.Decimals); err != nil {
		return &types.Currency{}, fmt.Errorf("failed scanning row into currency struct: %w", err)
	}

	return curr, nil
}
