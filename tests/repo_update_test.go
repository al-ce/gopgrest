package tests

import (
	"testing"

	"gopgrest/rsql"
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
		t.Fatal(err)
	}
	if err := checkMapEquality(expAuthors, gotRows); err != nil {
		t.Errorf("%s\nExp %v\nGot %v", err, expAuthors, gotRows)
	}
}

func Test_RepoUpdateRowByRSQL(t *testing.T) {
	repo := NewTestRepo(t)
	expAuthors := []types.RowData{AnneCarson, AnneBrontë}
	// Set forename column in expected author values
	for _, expAuth := range expAuthors {
		expAuth["forename"] = "Beatrice"
	}

	// Update forename for each author named 'Anne'
	update := types.RowData{"forename": "Beatrice"}
	filters := []rsql.Filter{
		{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
	}
	err := repo.UpdateRowByRSQL("authors", filters, &update)
	if err != nil {
		t.Errorf("Update by RSQL err: %s", err)
	}

	// Get all authors initially named 'Anne', i.e. all but Virginia
	rows, err := repo.DB.Query("SELECT * FROM authors WHERE forename != 'Virginia' ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	gotRows, err := service.ScanRows(rows)
	if err != nil {
		t.Fatal(err)
	}

	if err := checkMapEquality(expAuthors, gotRows); err != nil {
		t.Errorf("%s\nExp %v\nGot %v", err, expAuthors, gotRows)
	}
}
