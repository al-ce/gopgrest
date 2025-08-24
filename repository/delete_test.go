package repository_test

import (
	"errors"
	"testing"

	"gopgrest/apperrors"
	"gopgrest/rsql"
	"gopgrest/tests"
)

func Test_DeleteRowsByRSQL(t *testing.T) {
	// DELETE /authors?...
	t.Run("No query", func(t *testing.T) {
		repo := tests.NewTestRepo(t)
		deletedIDs, err := repo.DeleteRowsByRSQL("authors", []rsql.Condition{})
		if !errors.Is(err, apperrors.DeleteWithNoConditions) {
			t.Errorf("Expected error '%s', got '%s'", apperrors.DeleteWithNoConditions, err)
		}
		if len(deletedIDs) > 0 {
			t.Errorf("Expected 0 rows deleted, got %d", deletedIDs)
		}
	})

	// DELETE /authors?where=forname==Anne
	t.Run("Delete with single where condition", func(t *testing.T) {
		repo := tests.NewTestRepo(t)
		expCount, err := tests.CountRows(repo, "authors", "WHERE forename='Anne'")
		if err != nil {
			t.Fatal(err)
		}
		conditions := []rsql.Condition{
			{Column: rsql.Column{Name: "forename"}, Values: []string{"Anne"}, SQLOperator: "="},
		}
		deletedIDs, err := repo.DeleteRowsByRSQL("authors", conditions)
		if err != nil {
			t.Fatalf("Could not delete with conditions %v: %s", conditions, err)
		}
		if len(deletedIDs) != expCount {
			t.Errorf("Expected %d rows deleted, got %d", expCount, deletedIDs)
		}
		// Confirm authors no longer in DB
		gotRows, err := tests.SelectRows(repo, "SELECT * FROM authors WHERE forename = 'Anne'")
		if len(gotRows) != 0 {
			t.Errorf("Expected to not find authors, found %d", len(gotRows))
		}
	})

	// DELETE /authors?where=forname==Anne;born<1900
	t.Run("Delete with multiple conditions", func(t *testing.T) {
		repo := tests.NewTestRepo(t)
		expCount, err := tests.CountRows(repo, "authors", "WHERE forename='Anne' and born<1900")
		if err != nil {
			t.Fatal(err)
		}
		conditions := []rsql.Condition{
			{Column: rsql.Column{Name: "forename"}, Values: []string{"Anne"}, SQLOperator: "="},
			{Column: rsql.Column{Name: "born"}, Values: []string{"1900"}, SQLOperator: "<"},
		}
		deletedIDs, err := repo.DeleteRowsByRSQL("authors", conditions)
		if err != nil {
			t.Fatalf("Could not delete with conditions %v: %s", conditions, err)
		}
		if len(deletedIDs) != expCount {
			t.Errorf("Expected %d rows deleted, got %d", expCount, deletedIDs)
		}
		// Confirm authors no longer in DB
		gotRows, err := tests.SelectRows(
			repo,
			"SELECT * FROM authors WHERE forename = 'Anne' AND born < 1900",
		)
		if len(gotRows) != 0 {
			t.Errorf("Expected to not find authors, found %d", len(gotRows))
		}
	})
}
