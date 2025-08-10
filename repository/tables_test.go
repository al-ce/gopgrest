package repository_test

import (
	"slices"
	"testing"

	"gopgrest/repository"
	"gopgrest/test_utils"
)

func Test_GetPublicTables(t *testing.T) {
	tdb := test_utils.NewTestDB(t)

	// Make an array of all the tables in the test db
	tables, err := repository.GetPublicTables(tdb.DB)
	if err != nil {
		t.Errorf("%s", err)
	}

	expectedTables := []string{"authors", "books"}
	foundTables := []string{}
	for _, table := range tables {
		// Check for extraneous tables
		if !slices.Contains(expectedTables, table.Name) {
			t.Errorf("Found unexpected table %s in %v", table, expectedTables)
		}
		foundTables = append(foundTables, table.Name)
	}
	for _, table := range expectedTables {
		if !slices.Contains(foundTables, table) {
			t.Errorf("Expected to find table %s in %v", table, foundTables)
		}
	}
}
