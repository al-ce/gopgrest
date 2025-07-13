package repository

import (
	"database/sql"
	"fmt"

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
func (r *Repository) ListSets(params map[string][]string) (*sql.Rows, error) {
	// Build list query with optional conditional filters
	conditional, values := buildConditionalClause(params)
	listStmt := "select * from exercise_sets" + conditional

	// Execute list query
	rows, err := r.db.Query(listStmt, values...)
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

// UpdateSetField updates a field in an exercise set row
func (r *Repository) UpdateSetField(id, field string, value any) error {
	// Build update query
	updateStmt := fmt.Sprintf(
		"update exercise_sets set %s = $1 where id = $2",
		field,
	)

	if _, err := r.db.Exec(updateStmt, value, id); err != nil {
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

// buildConditionalClause builds a SQL WHERE clause to select a row. Ex: when
// params == `[name:[bob alice] age:[45]]`, the name in the row must be either
// bob or alice and the age must be 45.
func buildConditionalClause(params map[string][]string) (string, []any) {
	// If no params were passed, there should not be a WHERE clause
	if len(params) == 0 {
		return "", []any{}
	}

	clause := " WHERE ("
	values := []any{}

	// n is the number of the placeholder in the statement e.g. $1
	n := 1
	for k, vals := range params {
		// Join all values for this key with OR, allow any match
		for _, v := range vals {
			clause += fmt.Sprintf("%s = $%d OR ", k, n)
			n += 1
			values = append(values, v)
		}
		// Strip final " OR "
		clause = clause[:len(clause)-4]
		// Join all keys with AND, require at least one match from each key
		clause += ") AND ("
	}
	// Strip final " AND ("
	clause = clause[:len(clause)-6]

	return clause, values
}
