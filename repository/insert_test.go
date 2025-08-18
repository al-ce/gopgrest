package repository_test

import (
	"fmt"
	"strings"
	"testing"

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
	if err != nil {
		t.Fatalf("Could not insert authors\n%v:\n%s", newRows, err)
	}

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
