package tests

import (
	"testing"

	"gopgrest/apperrors"
	"gopgrest/rsql"
	"gopgrest/service"
	"gopgrest/types"
)

func Test_DeleteRowByID(t *testing.T) {
	repo := NewTestRepo(t)
	expAuthors := []types.RowData{AnneCarson, AnneBrontë, VirginiaWoolf}
	for index := range expAuthors {
		id := index + 1
		rowsAffected, err := repo.DeleteRowByID("authors", int64(id))
		if err != nil {
			t.Fatalf("Could not delete row %d: %s", id, err)
		}
		if rowsAffected != 1 {
			t.Fatalf("Expected to delete 1 row, deleted %d", rowsAffected)
		}
		// Confirm author no longer in DB
		rows, err := repo.DB.Query("SELECT * FROM authors WHERE id=$1", id)
		if err != nil {
			t.Fatalf("Could not pick author id %d: %s", id, err)
		}
		defer rows.Close()
		gotRows, err := service.ScanRows(rows)
		if err != nil {
			t.Fatalf("Could not scan author id %d: %s", id, err)
		}
		if len(gotRows) != 0 {
			t.Errorf("Expected to not find author w/ id %d, but found it", id)
		}
	}
}

func Test_DeleteRowsByRSQL(t *testing.T) {
	// DELETE /authors?...
	t.Run("No query", func(t *testing.T) {
		repo := NewTestRepo(t)
		rowsAffected, err := repo.DeleteRowsByRSQL("authors", []rsql.Filter{})
		if err != apperrors.DeleteWithNoFilters {
			t.Errorf("Expected error '%s', got '%s'", apperrors.DeleteWithNoFilters, err)
		}
		if rowsAffected != -1 {
			t.Errorf("Expected -1 rows affected (error), got %d", rowsAffected)
		}
	})

	// DELETE /authors?filter=forname==Anne
	t.Run("Delete with single parameter", func(t *testing.T) {
		repo := NewTestRepo(t)
		expCount, err := countRows(repo, "authors", "WHERE forename='Anne'")
		if err != nil {
			t.Fatal(err)
		}
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
		}
		rowsAffected, err := repo.DeleteRowsByRSQL("authors", filters)
		if err != nil {
			t.Errorf("Expected error '%s', got '%s'", apperrors.DeleteWithNoFilters, err)
		}
		if rowsAffected != expCount {
			t.Errorf("Expected %d rows deleted, got %d", expCount, rowsAffected)
		}
	})
}
