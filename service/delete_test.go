package service_test

import (
	"testing"

	"gopgrest/tests"
)

func Test_ServiceDeleteRowsByRSQL(t *testing.T) {
	service := tests.NewTestService(t)

	expAuthors, err := tests.SelectRows(service.Repo, "SELECT * FROM authors WHERE forename = 'Anne'")
	if err != nil {
		t.Fatal(err)
	}

	// Delete rows with matching conditions
	url := "/authors?forename==Anne"
	rowsAffected, err := service.DeleteRowsByRSQL("authors", url)
	if err != nil {
		t.Errorf("Update by RSQL err: %s", err)
	}
	if rowsAffected != int64(len(expAuthors)) {
		t.Fatalf(
			"Expected to update %d columns, instead updated %d",
			len(expAuthors),
			rowsAffected,
		)
	}

	// Confirm rows were deleted
	query := "SELECT * FROM authors WHERE forename = 'Anne' ORDER BY id"
	gotRows, err := tests.SelectRows(service.Repo, query)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotRows) > 0 {
		t.Errorf("Expected 0 rows, got %d", len(gotRows))
	}
}
