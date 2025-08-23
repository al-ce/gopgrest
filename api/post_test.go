package api_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"gopgrest/tests"
	"gopgrest/types"
)

func Test_POST(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)

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

	rr, err := tests.MakeHttpRequest(ah, http.MethodPost, "/authors", newRows)
	if err != nil {
		t.Error(err.Error())
	}

	if http.StatusOK != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}

	// Expect that we got back an array of ids, like `[4, 5]`
	reInsertedIds := regexp.MustCompile(`^\[(.*)\]$`)
	match := reInsertedIds.FindStringSubmatch(rr.Body.String())
	if len(match) != 2 {
		t.Errorf("Expected resp body match on pattern: %v\nGot: %v", reInsertedIds, match)
	}
	ids := match[1]

	// Check for inserted rows in DB
	gotRows, err := tests.SelectRows(ah.Repo, fmt.Sprintf("SELECT * FROM authors WHERE id IN (%s)", ids))
	if err != nil {
		t.Error(err.Error())
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
