package tests

import (
	"testing"

	"gopgrest/service"
	"gopgrest/types"
)

func Test_RepoGetRowById(t *testing.T) {
	repo := NewTestRepo(t)
	expAuthors := []types.RowData{AnneCarson, AnneBrontë, VirginiaWoolf}

	for index, auth := range expAuthors {
		id := index + 1
		rows, err := repo.GetRowByID("authors", int64(id))
		if err != nil {
			t.Fatalf("Could not pick author id %d: %s", id, err)
		}
		defer rows.Close()
		gotRows, err := service.ScanRows(rows)
		if err != nil {
			t.Fatalf("Could not scan author id %d: %s", id, err)
		}
		if err := checkMapEquality([]types.RowData{auth}, gotRows); err != nil {
			t.Error(err)
		}
	}
}
