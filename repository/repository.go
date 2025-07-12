package repository

import (
	"database/sql"

	"ftrack/models"
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

// CreateSet inserts an exercise set in the exercise_sets table
func (r *Repository) CreateSet(setData *models.ExerciseSet) error {
	// Build create query
	const createStmnt = `
		insert into exercise_sets
		(
			name,
			performed_at,
			weight,
			unit,
			reps,
			set_count,
			notes,
			split_day,
			program,
			tags
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Execute create query
	result, err := r.db.Exec(createStmnt,
		setData.Name,
		setData.PerformedAt,
		setData.Weight,
		setData.Unit,
		setData.Reps,
		setData.SetCount,
		setData.Notes,
		setData.SplitDay,
		setData.Program,
		setData.Tags,
	)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// DeleteSet removes a set from the exercise_sets table by id
func (r *Repository) DeleteSet(id string) error {
	const deleteStmt = "delete from exercise_sets where id = $1"

	// Execute delete query
	result, err := r.db.Exec(deleteStmt, id)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}
