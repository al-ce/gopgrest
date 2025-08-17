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

// buildWhereConditions builds a SQL WHERE clause from `filters`. Placeholders
// values begin at `start`+1, e.g. if `start` == 5, a WHERE clause would begin
// with `WHERE x = $6 AND ...`
func buildWhereConditions(filters []rsql.Filter, start int) (string, []any, error) {
	// If no params were passed, there should not be a WHERE clause
	if len(filters) == 0 {
		return "", []any{}, nil
	}

	// values holds the order of column values that matches the placeholders
	values := []any{}
	// conditions is an array of `col IN ([values])` statements joined by AND
	conditions := []string{}
	// n is the number of the placeholder in the statement e.g. $1
	n := start + 1

	for _, f := range filters {
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
			placeholders = append(placeholders, fmt.Sprintf("$%d", n))
			values = append(values, v)
			n++
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

	conditional := fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
	return conditional, values, nil
}

func buildColumnsToReturn(query rsql.Query) string {
	// If no columns were specified, the SELECT statement should be `SELECT *`
	if len(query.Fields) == 0 {
		return "*"
	}

	cols := []string{}
	for _, f := range query.Fields {
		if f.Alias != "" {
			cols = append(cols, fmt.Sprintf("%s AS %s", f.Column, f.Alias))
		} else {
			cols = append(cols, f.Column)
		}
	}
	return strings.Join(cols, ", ")
}

// buildJoinRelations builds SQL JOIN clauses
func buildJoinRelations(query rsql.Query) string {
	if len(query.Joins) == 0 {
		return ""
	}
	joins := []string{}
	for _, j := range query.Joins {
		joins = append(
			joins,
			fmt.Sprintf(
				"%s %s ON %s.%s = %s.%s",
				j.Type,
				j.Table,
				j.LeftQualifier,
				j.LeftCol,
				j.RightQualifier,
				j.RightCol,
			),
		)
	}
	return strings.Join(joins, " ")
}
