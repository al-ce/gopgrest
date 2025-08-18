package service_test

import (
	"testing"

	"gopgrest/tests"
	"gopgrest/types"
)

func Test_ServiceUpdateRowsByRSQL(t *testing.T) {
	service := tests.NewTestService(t)
	expAuthors, err := tests.SelectRows(service.Repo, "SELECT * FROM authors WHERE forename = 'Anne'")
	if err != nil {
		t.Fatal(err)
	}

	// Set forename column in expected author values
	for _, expAuth := range expAuthors {
		expAuth["forename"] = "Beatrice"
	}

	// Update forename for each author named 'Anne'
	update := types.RowData{"forename": "Beatrice"}
	url := "/authors?forename==Anne"
	rowsAffected, err := service.UpdateRowsByRSQL("authors", url, &update)
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

	// Confirm rows were updated
	query := "SELECT * FROM authors WHERE forename != 'Virginia' ORDER BY id"
	gotRows, err := tests.SelectRows(service.Repo, query)
	if err != nil {
		t.Fatal(err)
	}
	if err := tests.CheckMapEquality(expAuthors, gotRows); err != nil {
		t.Fatal(err)
	}
}
