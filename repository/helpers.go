package repository

import (
	"fmt"

	"gopgrest/types"
)

// GetTable gets a table from the tables slice by name
func (r *Repository) GetTable(tableName string) (*Table, error) {
	for _, t := range r.Tables {
		if t.Name == tableName {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("table does not exist")
}

func (r *Repository) IsValidColumn(table Table, col string) bool {
	_, ok := table.ColumnMap[col]
	return ok
}

// buildConditionalClause builds a SQL WHERE clause to select a row. Ex: when
// params == `[name:[bob alice] age:[45]]`, the name in the row must be either
// bob or alice and the age must be 45.
func buildConditionalClause(qf types.QueryFilter) (string, []any, error) {
	// If no params were passed, there should not be a WHERE clause
	if len(qf) == 0 {
		return "", []any{}, nil
	}

	clause := " WHERE ("
	values := []any{}

	// n is the number of the placeholder in the statement e.g. $1
	n := 0
	for k, vals := range qf {
		// Check for empty filter values
		if len(vals) == 0 {
			return "", []any{}, fmt.Errorf("attempt to filter on key %s with no values", k)
		}

		// Join all values for this key with OR, allow any match
		for _, v := range vals {
			n += 1
			clause += fmt.Sprintf("%s = $%d OR ", k, n)
			values = append(values, v)
		}
		// Strip final " OR "
		clause = clause[:len(clause)-4]
		// Join all cols with AND, require at least one match from each key
		clause += ") AND ("
	}
	// Strip final " AND ("
	clause = clause[:len(clause)-6]

	return clause, values, nil
}
