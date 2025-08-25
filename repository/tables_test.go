package repository_test

import (
	"slices"
	"testing"

	"gopgrest/assert"
	"gopgrest/repository"
	"gopgrest/tests"
)

func Test_GetPublicTables(t *testing.T) {
	tdb := tests.NewTestDB(t)

	// Make an array of all the tables in the test db
	tables, err := repository.GetPublicTables(tdb.DB)
	assert.Try(t, err)

	expectedTables := []string{"authors", "books", "genres"}
	foundTables := []string{}
	for _, table := range tables {
		// Check for extraneous tables
		assert.IsTrue(t, slices.Contains(expectedTables, table.Name))
		foundTables = append(foundTables, table.Name)
	}
	for _, table := range expectedTables {
		assert.IsTrue(t, slices.Contains(foundTables, table))
	}
}
