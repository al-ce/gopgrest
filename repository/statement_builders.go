package repository

import (
	"fmt"
	"slices"
	"strings"

	"gopgrest/rsql"
)

// buildWhereConditions builds a SQL WHERE clause from `conditions`.
// Placeholders values begin at `start`+1, e.g. if `start` == 5, a
// WHERE clause would begin with `WHERE x = $6 AND ...`
func buildWhereConditions(conditions []rsql.Condition, start int) (string, []any, error) {
	// If no params were passed, there should not be a WHERE clause
	if len(conditions) == 0 {
		return "", []any{}, nil
	}

	// values holds the order of column values that matches the placeholders
	values := []any{}
	// sqlConditions is an array of `col IN ([values])` statements joined by AND
	sqlConditions := []string{}
	// n is the number of the placeholder in the statement e.g. $1
	n := start + 1

	for _, cond := range conditions {
		var condition string
		columnName := makeQualifiedColumn(cond)

		// Null checks do not require placeholders or appending values array
		if slices.Contains([]string{"IS NULL", "IS NOT NULL"}, cond.SQLOperator) {
			condition = fmt.Sprintf("%s %s", columnName, cond.SQLOperator)
			sqlConditions = append(sqlConditions, condition)
			continue
		}

		// Check for empty condition values
		if len(cond.Values) == 0 {
			return "", []any{}, fmt.Errorf("Condition for col %s with no values", columnName)
		}

		// Append to values array in the same order we add conditions
		placeholders := []string{}
		for _, v := range cond.Values {
			// Placeholder value +1 should match values index
			placeholders = append(placeholders, fmt.Sprintf("$%d", n))
			values = append(values, v)
			n++
		}

		// Add `col {keyword} (...placeholders)` e.g.
		// `forename IN ($1,$2)`
		condition = fmt.Sprintf(
			"%s %s (%s)",
			columnName,
			cond.SQLOperator,
			strings.Join(placeholders, ","),
		)

		sqlConditions = append(sqlConditions, condition)
	}

	conditional := fmt.Sprintf("WHERE %s", strings.Join(sqlConditions, " AND "))
	return conditional, values, nil
}

func buildSelectColumns(query rsql.QueryParams) string {
	// If no columns were specified, the SELECT statement should be `SELECT *`
	if len(query.Columns) == 0 {
		return "*"
	}

	cols := []string{}
	for _, f := range query.Columns {
		if f.Alias != "" {
			cols = append(cols, fmt.Sprintf("%s AS %s", f.Name, f.Alias))
		} else {
			cols = append(cols, f.Name)
		}
	}
	return strings.Join(cols, ", ")
}

// buildJoinRelations builds SQL JOIN clauses
func buildJoinRelations(query rsql.QueryParams) string {
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

// makeQualifiedColumn returns a column name as a string with its qualifier if
// one was given, otherwise returns just the column name
func makeQualifiedColumn(cond rsql.Condition) string {
	var column string
	if cond.Column.Qualifier != "" {
		column = fmt.Sprintf("%s.%s", cond.Column.Qualifier, cond.Column.Name)
	} else {
		column = cond.Column.Name
	}
	return column
}
