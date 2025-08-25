package service_test

import (
	"testing"

	"gopgrest/assert"
	"gopgrest/tests"
)

func Test_ServiceDeleteRowsByRSQL(t *testing.T) {
	service := tests.NewTestService(t)

	expAuthors, err := tests.SelectRows(service.Repo, "SELECT * FROM authors WHERE forename = 'Anne'")
	assert.Try(t, err)

	// Delete rows with matching conditions
	url := "/authors?forename==Anne"
	deletedIDs, err := service.DeleteRowsByRSQL("authors", url)
	assert.Try(t, err)
	assert.IsTrue(t, len(deletedIDs) == len(expAuthors))

	// Confirm rows were deleted
	query := "SELECT * FROM authors WHERE forename = 'Anne' ORDER BY id"
	gotRows, err := tests.SelectRows(service.Repo, query)
	assert.Try(t, err)
	assert.IsTrue(t, len(gotRows) == 0)
}
