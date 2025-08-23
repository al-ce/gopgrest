package api_test

import (
	"encoding/json"
	"maps"
	"net/http"
	"slices"
	"testing"

	"gopgrest/tests"
)

// Test_GET_Tables checks if the "/" route returns a map of tables and their
// columns. For now, just tests that the request returns 200 and that the map
// has a key for each table
func Test_GET_Tables(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)

	rr, err := tests.MakeHttpRequest(ah, http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	if http.StatusOK != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}

	tables := map[string]any{}
	err = json.Unmarshal(rr.Body.Bytes(), &tables)
	if err != nil {
		t.Error(err.Error())
	}

	gotTables := slices.Collect(maps.Keys(tables))

	expectTables := []string{"authors", "books", "genres"}
	for _, table := range expectTables {
		if !slices.Contains(gotTables, table) {
			t.Errorf("Expected %s in got tables", table)
		}
	}
}
