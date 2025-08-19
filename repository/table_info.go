package repository

import (
	"gopgrest/apperrors"
)

// GetTable gets a table from the tables slice by name
func (r *Repository) GetTable(tableName string) (*Table, error) {
	for _, t := range r.Tables {
		if t.Name == tableName {
			return &t, nil
		}
	}
	return nil, apperrors.TableDoesNotExist
}

func (r *Repository) IsValidColumn(table Table, col string) bool {
	_, ok := table.ColumnMap[col]
	return ok
}
