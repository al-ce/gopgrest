package service_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gopgrest/apperrors"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_ServiceInsertRow_Single(t *testing.T) {
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
	if err != nil {
		t.Fatalf("Could not insert authors\n%v:\n%s", newRows, err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

	for idx, expAuthor := range newRows {
		gotAuthor := gotRows[idx]
		for k, v := range expAuthor {
			if v != gotAuthor[k] {
				t.Errorf("Exp %s %v %T:\nGot %v %T", k, v, v, gotAuthor[k], gotAuthor[k])
			}
		}
	}
}

func Test_ServiceInsertRow_NoRows(t *testing.T) {
	repo := tests.NewTestService(t)
	_, err := repo.InsertRows([]types.RowData{}, "authors")
	if !errors.Is(err, apperrors.InsertWithNoRows) {
		t.Errorf("Expected err '%s' got '%s'", apperrors.InsertWithNoRows, err)
	}
}

func Test_ServiceInsertRow_BadTable(t *testing.T) {
	service := tests.NewTestService(t)

	_, gotErr := service.InsertRows([]types.RowData{{"dummy": "value"}}, "doesnotexist")
	expErr := apperrors.TableDoesNotExist
	validateError(t, expErr, gotErr)
}

func Test_ServiceInsertRow_BadColumn(t *testing.T) {
	service := tests.NewTestService(t)
	badCol := "specialty"
	badRow := []types.RowData{{"surname": "Sappho", badCol: "lyric poetry"}}

	_, gotErr := service.InsertRows(badRow, "authors")
	expErr := apperrors.ColDoesNotExist
	validateError(t, expErr, gotErr)
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
	validateError(t, expErr, gotErr)
}

func validateError(t *testing.T, expErr, gotErr error) {
	if !errors.Is(gotErr, expErr) {
		t.Errorf("\nExp err '%s'\nGot err '%s'", expErr, errors.Unwrap(gotErr))
	}
}
