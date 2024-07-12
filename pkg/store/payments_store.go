package store

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentsQueries struct {
	db Executor
}

func NewPaymentsQueries(db Executor) *PaymentsQueries {
	return &PaymentsQueries{
		db: db,
	}
}

type PaymentsStore struct {
	connPool *pgxpool.Pool
	*PaymentsQueries
}

func NewPaymentsStore(db *pgxpool.Pool) *PaymentsStore {
	return &PaymentsStore{
		connPool:        db,
		PaymentsQueries: NewPaymentsQueries(db),
	}
}

func (p *PaymentsStore) BeginTx() (pgx.Tx, error) {
	tx, err := p.connPool.Begin(context.Background())
	if err != nil {
		return tx, fmt.Errorf("failed starting transaction: %s", err)
	}

	return tx, nil
}

func (p *PaymentsStore) WithTx(tx pgx.Tx) PaymentsQueryHandler {
	return &PaymentsQueries{
		db: tx,
	}
}

func (s *PaymentsQueries) GetUser(username string) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE username = $1;", username)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}

func (p *PaymentsQueries) GetUserById(userId int32) (*types.User, error) {
	var user = new(types.User)

	row := p.db.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1;", userId)

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}

func (p *PaymentsQueries) CreateAccount(acc *types.Account) error {
	row := p.db.QueryRow(context.Background(), "INSERT INTO accounts (user_id, currency_id, balance) VALUES ($1, $2, $3) RETURNING *;", acc.UserId, acc.CurrencyId, acc.Balance)
	if err := row.Scan(&acc.Id, &acc.UserId, &acc.CurrencyId, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
		return fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return nil
}

func (p *PaymentsQueries) GetAccount(userId int32, currencyTicker string) (*types.Account, error) {
	var acc = new(types.Account)

	row := p.db.QueryRow(context.Background(), "SELECT * FROM accounts INNER JOIN currencies ON accounts.currency_id = currencies.id WHERE user_id = $1 AND ticker = $2;", userId, currencyTicker)
	if err := row.Scan(&acc.Id, &acc.UserId, &acc.CurrencyId, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
		return &types.Account{}, fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return acc, nil
}

func (p *PaymentsQueries) UpdateAccountBalance(accountId int32, balance int64) error {
	_, err := p.db.Exec(context.Background(), "UPDATE accounts SET balance = $1 WHERE id = $2;", balance, accountId)
	if err != nil {
		return fmt.Errorf("failed executing update query: %w", err)
	}

	return nil
}

func (p *PaymentsQueries) CreateDeposit(dep *types.Deposit) error {
	row := p.db.QueryRow(context.Background(), "INSERT INTO deposits (account_id, currency_id, amount) VALUES ($1, $2, $3) RETURNING *;", dep.AccountId, dep.CurrencyId, dep.Amount)
	if err := row.Scan(&dep.Id, &dep.AccountId, &dep.CurrencyId, &dep.Amount, &dep.CreatedAt); err != nil {
		return fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return nil
}

func (p *PaymentsQueries) CreateWithdrawal(withdrawal *types.Withdrawal) error {
	row := p.db.QueryRow(context.Background(), "INSERT INTO withdrawals (account_id, currency_id, amount) VALUES ($1, $2, $3) RETURNING *;", withdrawal.AccountId, withdrawal.CurrencyId, withdrawal.Amount)
	if err := row.Scan(&withdrawal.Id, &withdrawal.AccountId, &withdrawal.CurrencyId, &withdrawal.Amount, &withdrawal.CreatedAt); err != nil {
		return fmt.Errorf("failed scanning row into account struct: %w", err)
	}

	return nil
}

func (p *PaymentsQueries) CreateTransfer(transfer *types.Transfer) error {
	row := p.db.QueryRow(context.Background(), "INSERT INTO transfers (sender_id, recipient_id, currency_id, amount) VALUES ($1, $2, $3, $4) RETURNING *;", transfer.SenderId, transfer.RecipientId, transfer.CurrencyId, transfer.Amount)
	if err := row.Scan(&transfer.Id, &transfer.SenderId, &transfer.RecipientId, &transfer.CurrencyId, &transfer.Amount, &transfer.CreatedAt); err != nil {
		return fmt.Errorf("failed scanning row into transfer struct: %w", err)
	}

	return nil
}

func (p *PaymentsQueries) GetCurrency(ticker string) (*types.Currency, error) {
	curr := new(types.Currency)

	row := p.db.QueryRow(context.Background(), "SELECT * FROM currencies WHERE ticker = $1", ticker)
	if err := row.Scan(&curr.Id, &curr.Name, &curr.Ticker, &curr.Decimals); err != nil {
		return &types.Currency{}, fmt.Errorf("failed scanning row into currency struct: %w", err)
	}

	return curr, nil
}
