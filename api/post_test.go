package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"gopgrest/assert"
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

	assert.Try(t, err)
	assert.IsEq(t, rr.Code, http.StatusOK)

	// Expect that we got back an array of ids, like `[4, 5]`
	tests.ParseIDArrayResponse(t, rr.Body.String())

	// Check for inserted rows in DB
	ids := strings.Trim(rr.Body.String(), "[]")
	gotRows, err := tests.SelectRows(
		ah.Repo,
		fmt.Sprintf("SELECT * FROM authors WHERE id IN (%s)", ids),
	)
	assert.Try(t, err)

	for idx, expAuthor := range newRows {
		gotAuthor := gotRows[idx]
		for k, v := range expAuthor {
			assert.IsEq(t, v, gotAuthor[k])
		}
	}
}
