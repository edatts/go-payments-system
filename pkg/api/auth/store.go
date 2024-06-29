package auth

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(user *types.User) error {
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

func (s *Store) GetUser(username string) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE uesrname = $1;", username)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1;", email)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}

func (s *Store) GetUserById(id int32) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1;", id)

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}
