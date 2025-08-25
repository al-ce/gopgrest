package repository_test

import (
	"testing"

	"gopgrest/apperrors"
	"gopgrest/assert"
	"gopgrest/rsql"
	"gopgrest/tests"
)

func Test_DeleteRowsByRSQL(t *testing.T) {
	// DELETE /authors?...
	t.Run("No query", func(t *testing.T) {
		repo := tests.NewTestRepo(t)
		deletedIDs, err := repo.DeleteRowsByRSQL("authors", []rsql.Condition{})
		assert.ErrorsIs(t, err, apperrors.DeleteWithNoConditions)
		assert.IsTrue(t, len(deletedIDs) == 0)
	})

	// DELETE /authors?where=forname==Anne
	t.Run("Delete with single where condition", func(t *testing.T) {
		repo := tests.NewTestRepo(t)
		expCount, err := tests.CountRows(repo, "authors", "WHERE forename='Anne'")
		assert.Try(t, err)

		conditions := []rsql.Condition{
			{Column: rsql.Column{Name: "forename"}, Values: []string{"Anne"}, SQLOperator: "="},
		}
		deletedIDs, err := repo.DeleteRowsByRSQL("authors", conditions)
		assert.Try(t, err)
		assert.IsEq(t, len(deletedIDs), expCount)

		// Confirm authors no longer in DB
		gotRows, err := tests.SelectRows(repo, "SELECT * FROM authors WHERE forename = 'Anne'")
		assert.IsTrue(t, len(gotRows) == 0)
	})

	// DELETE /authors?where=forname==Anne;born<1900
	t.Run("Delete with multiple conditions", func(t *testing.T) {
		repo := tests.NewTestRepo(t)
		expCount, err := tests.CountRows(repo, "authors", "WHERE forename='Anne' and born<1900")
		assert.Try(t, err)

		conditions := []rsql.Condition{
			{Column: rsql.Column{Name: "forename"}, Values: []string{"Anne"}, SQLOperator: "="},
			{Column: rsql.Column{Name: "born"}, Values: []string{"1900"}, SQLOperator: "<"},
		}
		deletedIDs, err := repo.DeleteRowsByRSQL("authors", conditions)
		assert.Try(t, err)

		assert.IsEq(t, len(deletedIDs), expCount)

		// Confirm authors no longer in DB
		gotRows, err := tests.SelectRows(
			repo,
			"SELECT * FROM authors WHERE forename = 'Anne' AND born < 1900",
		)
		assert.IsTrue(t, len(gotRows) == 0)
	})
}
