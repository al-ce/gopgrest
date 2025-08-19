package repository_test

import (
	"errors"
	"fmt"
	"testing"

	"gopgrest/apperrors"
	"gopgrest/rsql"
	"gopgrest/tests"
)

func Test_DeleteRowByID(t *testing.T) {
	repo := tests.NewTestRepo(t)
	sampleAuthors, err := tests.SelectRows(repo, "SELECT * FROM authors")
	if err != nil {
		t.Fatal(err)
	}
	for index := range sampleAuthors {
		id := index + 1
		rowsAffected, err := repo.DeleteRowByID("authors", int64(id))
		if err != nil {
			t.Fatalf("Could not delete row %d: %s", id, err)
		}
		if rowsAffected != 1 {
			t.Fatalf("Expected to delete 1 row, deleted %d", rowsAffected)
		}
		// Confirm author no longer in DB
		gotRows, err := tests.SelectRows(repo, fmt.Sprintf("SELECT * FROM AUTHORS WHERE id=%d", id))
		if len(gotRows) != 0 {
			t.Errorf("Expected to not find author w/ id %d, but found it", id)
		}
	}
}

func Test_DeleteRowsByRSQL(t *testing.T) {
	// DELETE /authors?...
	t.Run("No query", func(t *testing.T) {
		repo := tests.NewTestRepo(t)
		rowsAffected, err := repo.DeleteRowsByRSQL("authors", []rsql.Condition{})
		if !errors.Is(err, apperrors.DeleteWithNoConditions) {
			t.Errorf("Expected error '%s', got '%s'", apperrors.DeleteWithNoConditions, err)
		}
		if rowsAffected != -1 {
			t.Errorf("Expected -1 rows affected (error), got %d", rowsAffected)
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
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
		}
		rowsAffected, err := repo.DeleteRowsByRSQL("authors", conditions)
		if err != nil {
			t.Fatalf("Could not delete with conditions %v: %s", conditions, err)
		}
		if rowsAffected != expCount {
			t.Errorf("Expected %d rows deleted, got %d", expCount, rowsAffected)
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
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
			{Column: "born", Values: []string{"1900"}, SQLOperator: "<"},
		}
		rowsAffected, err := repo.DeleteRowsByRSQL("authors", conditions)
		if err != nil {
			t.Fatalf("Could not delete with conditions %v: %s", conditions, err)
		}
		if rowsAffected != expCount {
			t.Errorf("Expected %d rows deleted, got %d", expCount, rowsAffected)
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
