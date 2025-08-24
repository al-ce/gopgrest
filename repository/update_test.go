package repository_test

import (
	"testing"

	"gopgrest/rsql"
	"gopgrest/tests"
	"gopgrest/types"
)

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
	conditions := []rsql.Condition{
		{Column: rsql.Column{Name: "forename"}, Values: []string{"Anne"}, SQLOperator: "="},
	}
	ids, err := repo.UpdateRowsByRSQL("authors", conditions, &update)
	if err != nil {
		t.Errorf("Update by RSQL err: %s", err)
	}
	if len(ids) != len(expAuthors) {
		t.Fatalf(
			"Expected to update %d columns, instead updated %d",
			len(expAuthors),
			len(ids),
		)
	}

	// Confirm rows were updated
	query := "SELECT * FROM authors WHERE forename != 'Virginia' ORDER BY id"
	gotRows, err := tests.SelectRows(repo, query)
	if err != nil {
		t.Fatal(err)
	}
	if err := tests.CheckMapEquality(expAuthors, gotRows); err != nil {
		t.Fatal(err)
	}
}
