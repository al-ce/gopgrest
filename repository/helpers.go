package repository

import (
	"fmt"
)

// TableExists check if a table exists in the database
func (r *Repository) TableExists(table string) bool {
	for _, t := range r.tables {
		if t.name == table {
			return true
		}
	}
	return false
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
		// Join all cols with AND, require at least one match from each key
		clause += ") AND ("
	}
	// Strip final " AND ("
	clause = clause[:len(clause)-6]

	return clause, values
}
