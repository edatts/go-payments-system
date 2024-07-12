package store

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	connPool *pgxpool.Pool
	*UserQueries
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{
		connPool:    db,
		UserQueries: NewUserQueries(db),
	}
}

func (u *UserStore) BeginTx() (pgx.Tx, error) {
	tx, err := u.connPool.Begin(context.Background())
	if err != nil {
		return tx, fmt.Errorf("failed starting transaction: %s", err)
	}

	return tx, nil
}

func (u *UserStore) WithTx(tx pgx.Tx) *UserQueries {
	return NewUserQueries(tx)
}

type UserQueries struct {
	db Executor
}

func NewUserQueries(db Executor) *UserQueries {
	return &UserQueries{
		db: db,
	}
}

func (s *UserQueries) CreateUser(user *types.User) error {
	_, err := s.db.Exec(
		context.Background(),
		`INSERT INTO users (
			first_name,
			last_name,
			username,
			email,
			password
		) VALUES (
			$1,$2,$3,$4,$5
		);`,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Password,
	)
	if err != nil {
		return fmt.Errorf("failed inserting user into db: %w", err)
	}

	return nil
}

func (s *UserQueries) GetUser(username string) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE username = $1;", username)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}

func (s *UserQueries) GetUserByEmail(email string) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1;", email)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}

func (s *UserQueries) GetUserById(id int32) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1;", id)

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}
