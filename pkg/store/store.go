package store

import (
	"context"
	"fmt"

	"github.com/edatts/go-payment-system/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
	*UserStore
	*PaymentsStore
}

// type Store struct {
// 	db *pgxpool.Pool
// }

type Tx struct {
	tx pgx.Tx
}

func (t Tx) Commit() error {
	return t.tx.Commit(context.Background())
}

func (t Tx) Rollback() error {
	return t.tx.Rollback(context.Background())
}

func NewStore(db *pgxpool.Pool) types.Store {
	return &Store{
		db:            db,
		UserStore:     NewUserStore(db),
		PaymentsStore: NewPaymentsStore(db),
	}
}

func (s *Store) BeginTx() (types.Tx, error) {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return &Tx{tx: tx}, fmt.Errorf("failed starting transaction: %s", err)
	}

	return &Tx{tx: tx}, nil
}
