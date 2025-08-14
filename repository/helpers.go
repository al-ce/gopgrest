package repository

import (
	"fmt"
	"slices"
	"strings"

	"gopgrest/rsql"
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

// buildWhereConditions builds a SQL WHERE clause
func buildWhereConditions(query *rsql.Query) (string, []any, error) {
	// If no params were passed, there should not be a WHERE clause
	if query == nil || len(query.Filters) == 0 {
		return "", []any{}, nil
	}

	// values holds the order of column values that matches the placeholders
	values := []any{}
	// conditions is an array of `col IN ([values])` statements joined by AND
	conditions := []string{}
	// n is the number of the placeholder in the statement e.g. $1
	n := 0

	for _, f := range query.Filters {
		var condition string

		// Null checks do not require placeholders or appending values array
		if slices.Contains([]string{"IS NULL", "IS NOT NULL"}, f.SQLOperator) {
			condition = fmt.Sprintf("%s %s", f.Column, f.SQLOperator)
			conditions = append(conditions, condition)
			continue
		}

		// Check for empty filter values
		if len(f.Values) == 0 {
			return "", []any{}, fmt.Errorf("attempt to filter on col %s with no values", f.Column)
		}

		// Append to values array in the same order we add conditions
		placeholders := []string{}
		for _, v := range f.Values {
			// Placeholder value +1 should match values index
			n += 1
			placeholders = append(placeholders, fmt.Sprintf("$%d", n))
			values = append(values, v)
		}

		// Add `col {keyword} (...placeholders)` e.g.
		// `forename IN ($1,$2)`
		condition = fmt.Sprintf(
			"%s %s (%s)",
			f.Column,
			f.SQLOperator,
			strings.Join(placeholders, ","),
		)

		conditions = append(conditions, condition)
	}

	conditional := fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	return conditional, values, nil
}
