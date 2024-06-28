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

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow(context.Background(), "SELECT * FROM users WHERE email = ?", email)

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return &types.User{}, fmt.Errorf("failed scanning row: %w", err)
	}

	return user, nil
}

func (s *Store) GetUserById(id int32) (*types.User, error) {

	return &types.User{}, nil
}
