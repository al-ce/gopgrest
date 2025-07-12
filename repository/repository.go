package repository

import (
	"database/sql"
)

// Repository handles database transactions
type Repository struct {
	db *sql.DB
}

// NewRepository returns a new Repository
func NewRepository(db *sql.DB) Repository {
	return Repository{
		db: db,
	}
}
