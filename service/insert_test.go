package service_test

import (
	"fmt"
	"strings"
	"testing"

	"gopgrest/apperrors"
	"gopgrest/assert"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_ServiceInsertRows(t *testing.T) {
	service := tests.NewTestService(t)

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

	ids, err := service.InsertRows(newRows, "authors")
	assert.Try(t, err)

	// Turn got ids into str to retrieve from db in one query
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = fmt.Sprintf("%d", id)
	}
	gotRows, err := tests.SelectRows(
		service.Repo,
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

func Test_ServiceInsertRow_NoRows(t *testing.T) {
	repo := tests.NewTestService(t)
	_, err := repo.InsertRows([]types.RowData{}, "authors")
	assert.ErrorsIs(t, err, apperrors.InsertWithNoRows)
}

func Test_ServiceInsertRow_BadTable(t *testing.T) {
	service := tests.NewTestService(t)

	_, gotErr := service.InsertRows([]types.RowData{{"dummy": "value"}}, "doesnotexist")
	expErr := apperrors.TableDoesNotExist
	assert.ErrorsIs(t, gotErr, expErr)
}

func Test_ServiceInsertRow_BadColumn(t *testing.T) {
	service := tests.NewTestService(t)
	badCol := "specialty"
	badRow := []types.RowData{{"surname": "Sappho", badCol: "lyric poetry"}}

	_, gotErr := service.InsertRows(badRow, "authors")
	expErr := apperrors.ColDoesNotExist
	assert.ErrorsIs(t, gotErr, expErr)
}

func Test_ServiceInsertRow_MismatchedColumns(t *testing.T) {
	service := tests.NewTestService(t)
	mismatchedCols := []types.RowData{
		{
			"surname": "Jemisin",
			"born":    int64(1972),
			"died":    nil,
		},
		{
			"forename": "Martha",
			"surname":  "Nussbaum",
		},
	}

	_, gotErr := service.InsertRows(mismatchedCols, "authors")
	expErr := apperrors.InsertColsDoNotMatch
	assert.ErrorsIs(t, gotErr, expErr)
}
