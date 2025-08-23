package apperrors

import (
	"errors"
	"fmt"

	"gopgrest/rsql"
)

var (
	GetByIdNotUnique         = errors.New("GET by id returned multiple rows")
	InsertWithNoRows         = errors.New("Cannot insert with no rows")
	InsertColsDoNotMatch     = errors.New("Columns in rows to insert do not match")
	InsertValTypesDoNotMatch = errors.New("Value types in rows to insert do not match")
	DeleteWithNoConditions   = errors.New("Will not DELETE with no WHERE conditions")
	UpdateWithNoConditions   = errors.New("Will not UPDATE with no WHERE conditions")

	TableDoesNotExist = errors.New("Table does not exist")
	ColDoesNotExist   = errors.New("Column not found in given table")
)

func NewDeleteInvalidIDErr(tableName string, id int64) error {
	return fmt.Errorf("cannot delete, no row in table %s with id %d", tableName, id)
}

func NewUpdateNoMatchingConditionsErr(tableName string, conditions []rsql.Condition) error {
	return fmt.Errorf(
		"cannot update, no rows in table %s with conditions %v",
		tableName, conditions,
	)
}
