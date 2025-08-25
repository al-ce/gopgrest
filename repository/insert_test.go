package repository_test

import (
	"fmt"
	"strings"
	"testing"

	"gopgrest/apperrors"
	"gopgrest/assert"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_RepoInsertRow_Single(t *testing.T) {
	repo := tests.NewTestRepo(t)

	newRows := []types.RowData{
		{
			"forename": "N.K.",
			"surname":  "Jemisin",
			"born":     int64(1972),
			"died":     nil,
		},
		{
			"forename": "Martha",
			"surname":  "Nussbaum",
			"born":     int64(1947),
			"died":     nil,
		},
	}

	ids, err := repo.InsertRows("authors", newRows)
	assert.Try(t, err)

	// Turn got ids into str to retrieve from db in one query
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = fmt.Sprintf("%d", id)
	}
	gotRows, err := tests.SelectRows(
		repo,
		fmt.Sprintf(
			"SELECT * FROM authors WHERE id IN (%s) ORDER BY id",
			strings.Join(idStrs, ","),
		),
	)
	assert.Try(t, err)

	for idx, expAuthor := range newRows {
		gotAuthor := gotRows[idx]
		for k, v := range expAuthor {
			assert.IsEq(t, v, gotAuthor[k])
		}
	}
}

func Test_RepoInsertRow_NoRows(t *testing.T) {
	repo := tests.NewTestRepo(t)
	_, err := repo.InsertRows("authors", []types.RowData{})
	assert.ErrorsIs(t, err, apperrors.InsertWithNoRows)
}
