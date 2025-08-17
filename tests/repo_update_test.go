package tests

import (
	"testing"

	"gopgrest/service"
	"gopgrest/types"
)

func Test_RepoUpdateRowById(t *testing.T) {
	repo := NewTestRepo(t)
	expAuthors := []types.RowData{AnneCarson, AnneBrontë, VirginiaWoolf}
	// Set forename column in expected author values
	for _, expAuth := range expAuthors {
		expAuth["forename"] = "Beatrice"
	}

	// Update forename for each author in DB
	for index := range expAuthors {
		id := int64(index + 1)

		update := types.RowData{"forename": "Beatrice"}
		err := repo.UpdateRowByID("authors", id, &update)
		if err != nil {
			t.Fatalf("Could not update author id %d: %s", id, err)
		}
	}

	// Verify updated column
	rows, err := repo.DB.Query("SELECT * FROM authors")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	gotRows, err := service.ScanRows(rows)
	if err != nil {
		t.Fatalf("Could not scan author rows: %s", err)
	}
	if err := checkMapEquality(expAuthors, gotRows); err != nil {
		t.Errorf("%s\nExp %v\nGot %v", err, expAuthors, gotRows)
	}
}
