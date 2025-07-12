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

// ListSets retrieves the list of all sets from the exercise_sets table
func (r *Repository) ListSets() (*sql.Rows, error) {
	const listStmt = "select * from exercise_sets"
	rows, err := r.db.Query(listStmt)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
