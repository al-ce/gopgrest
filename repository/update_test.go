package repository_test

import (
	"testing"

	"gopgrest/assert"
	"gopgrest/rsql"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_RepoUpdateRowByRSQL(t *testing.T) {
	repo := tests.NewTestRepo(t)
	expAuthors, err := tests.SelectRows(repo, "SELECT * FROM authors WHERE forename = 'Anne'")
	assert.Try(t, err)

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
	assert.Try(t, err)
	assert.IsTrue(t, len(ids) == len(expAuthors))

	// Confirm rows were updated
	query := "SELECT * FROM authors WHERE forename != 'Virginia' ORDER BY id"
	gotRows, err := tests.SelectRows(repo, query)
	assert.Try(t, err)

	err = tests.CheckMapEquality(expAuthors, gotRows)
	if err != nil {
		t.Fatal(err)
	}
}
