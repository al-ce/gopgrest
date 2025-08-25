package service_test

import (
	"testing"

	"gopgrest/assert"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_ServiceUpdateRowsByRSQL(t *testing.T) {
	service := tests.NewTestService(t)
	expAuthors, err := tests.SelectRows(service.Repo, "SELECT * FROM authors WHERE forename = 'Anne'")
	assert.Try(t, err)

	// Set forename column in expected author values
	for _, expAuth := range expAuthors {
		expAuth["forename"] = "Beatrice"
	}

	// Update forename for each author named 'Anne'
	update := types.RowData{"forename": "Beatrice"}
	url := "/authors?forename==Anne"
	ids, err := service.UpdateRowsByRSQL("authors", url, &update)
	assert.Try(t, err)
	assert.IsTrue(t, len(ids) == len(expAuthors))

	// Confirm rows were updated
	query := "SELECT * FROM authors WHERE forename != 'Virginia' ORDER BY id"
	gotRows, err := tests.SelectRows(service.Repo, query)
	assert.Try(t, err)

	err = tests.CheckMapEquality(expAuthors, gotRows)
	assert.Try(t, err)
}
