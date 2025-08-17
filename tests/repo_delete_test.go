package tests

import (
	"testing"

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
		rows, err := repo.GetRowByID("authors", int64(id))
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
