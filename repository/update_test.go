package repository_test

import (
	"testing"

	"gopgrest/repository"
	"gopgrest/rsql"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_RepoUpdateRowById(t *testing.T) {
	repo := tests.NewTestRepo(t)
	expAuthors, err := tests.SelectRows(repo, "SELECT * FROM authors")
	if err != nil {
		t.Fatal(err)
	}

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

	err = verifyUpdatedColumns(
		repo,
		expAuthors,
		// Get all authors initially named 'Anne', i.e. all but Virginia
		"SELECT * FROM authors",
	)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_RepoUpdateRowByRSQL(t *testing.T) {
	repo := tests.NewTestRepo(t)
	expAuthors, err := tests.SelectRows(repo, "SELECT * FROM authors WHERE forename = 'Anne'")
	if err != nil {
		t.Fatal(err)
	}

	// Set forename column in expected author values
	for _, expAuth := range expAuthors {
		expAuth["forename"] = "Beatrice"
	}

	// Update forename for each author named 'Anne'
	update := types.RowData{"forename": "Beatrice"}
	filters := []rsql.Filter{
		{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
	}
	err = repo.UpdateRowByRSQL("authors", filters, &update)
	if err != nil {
		t.Errorf("Update by RSQL err: %s", err)
	}

	err = verifyUpdatedColumns(
		repo,
		expAuthors,
		// Get all authors initially named 'Anne', i.e. all but Virginia
		"SELECT * FROM authors WHERE forename != 'Virginia' ORDER BY id",
	)
	if err != nil {
		t.Fatal(err)
	}
}

func verifyUpdatedColumns(repo repository.Repository, expAuthors []types.RowData, query string) error {
	gotRows, err := tests.SelectRows(repo, query)
	if err != nil {
		return err
	}
	if err := tests.CheckMapEquality(expAuthors, gotRows); err != nil {
		return err
	}
	return nil
}
