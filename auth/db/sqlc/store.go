package db

import "database/sql"

type TxStore interface {
	Querier
}

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) TxStore {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}
