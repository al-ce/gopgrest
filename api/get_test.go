package api_test

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"gopgrest/api"
	"gopgrest/tests"
	"gopgrest/types"
)

// Test_GET_RowByID tests requests with an id as a resource in the URL, e.g.
// `/authors/2`
func Test_GET_RowByID(t *testing.T) {
	repo := tests.NewTestRepo(t)
	expCount, err := tests.CountRows(repo, "authors", "")
	tests.Try(t, err)

	t.Run(fmt.Sprintf("ids 1-%d", expCount), func(t *testing.T) {
		for i := 1; i <= int(expCount); i++ {
			rawQuery := fmt.Sprintf("SELECT * FROM authors WHERE id = %d", i)
			apiGetRowsTester(t, rawQuery, fmt.Sprintf("/authors/%d", i))
		}
	})
}

// Test_GET_Rows_NoRSQL tests requests with no RSQL query params in the URL
func Test_GET_Rows_NoRSQL(t *testing.T) {
	t.Run("No RSQL", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors"
		apiGetRowsTester(t, rawQuery, "/authors")
	})
	t.Run("No RSQL, trailing `/`", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors"
		apiGetRowsTester(t, rawQuery, "/authors/")
	})
	t.Run("No RSQL, trailing `?`", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors"
		apiGetRowsTester(t, rawQuery, "/authors?")
	})
	t.Run("No RSQL, trailing `/?`", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors"
		apiGetRowsTester(t, rawQuery, "/authors/?")
	})
}

func Test_GET_Rows_RSQL_Select(t *testing.T) {
	t.Run("select one field", func(t *testing.T) {
		rawQuery := "SELECT surname FROM authors"
		apiGetRowsTester(t, rawQuery, "/authors?select=surname")
	})
	t.Run("select multiple fields", func(t *testing.T) {
		rawQuery := "SELECT forename, surname FROM authors"
		apiGetRowsTester(t, rawQuery, "/authors?select=forename,surname")
	})
}

func Test_GET_Rows_RSQL_Where(t *testing.T) {
	t.Run("single condition", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname='Woolf'"
		url := "/authors?where=surname==Woolf"
		apiGetRowsTester(t, rawQuery, url)
	})
	t.Run("multiple conditions", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname='Woolf' AND forename='Virginia'"
		url := "/authors?where=surname==Woolf;forename==Virginia"
		apiGetRowsTester(t, rawQuery, url)
	})
}

func Test_QueryParams_TrailingChars(t *testing.T) {
	t.Run("trailing `/`", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname='Woolf' AND forename='Virginia'"
		url := "/authors?where=surname==Woolf;forename==Virginia/"
		apiGetRowsTester(t, rawQuery, url)
	})
	t.Run("trailing `?`", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname='Woolf' AND forename='Virginia'"
		url := "/authors?where=surname==Woolf;forename==Virginia?"
		apiGetRowsTester(t, rawQuery, url)
	})
	t.Run("trailing `/?`", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname='Woolf' AND forename='Virginia'"
		url := "/authors?where=surname==Woolf;forename==Virginia/?"
		apiGetRowsTester(t, rawQuery, url)
	})
}

func Test_GET_Rows_RSQL_INNER_JOIN(t *testing.T) {
	t.Run("single JOIN", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books INNER JOIN genres on books.genre_id = genres.id"
		url := "/books?select=title,name:genre&inner_join=genres:books.genre_id==genres.id"
		apiGetRowsTester(t, rawQuery, url)
	})

	t.Run("multiple JOINs", func(t *testing.T) {
		rawQuery := "SELECT * FROM books JOIN authors ON books.author_id=authors.id JOIN genres ON books.genre_id=genres.id"
		url := "/books?join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id"
		apiGetRowsTester(t, rawQuery, url)
	})

	t.Run("single INNER JOIN", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books INNER JOIN genres on books.genre_id = genres.id"
		url := "/books?select=title,name:genre&inner_join=genres:books.genre_id==genres.id"
		apiGetRowsTester(t, rawQuery, url)
	})

	t.Run("single LEFT JOIN", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books LEFT JOIN genres ON books.genre_id = genres.id"
		url := "/books?select=title,name:genre&left_join=genres:books.genre_id==genres.id"
		apiGetRowsTester(t, rawQuery, url)
	})

	t.Run("single RIGHT JOIN", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books RIGHT JOIN genres ON books.genre_id = genres.id"
		url := "/books?select=title,name:genre&right_join=genres:books.genre_id==genres.id"
		apiGetRowsTester(t, rawQuery, url)
	})
}

// Test_GET_Tables checks if the "/" route returns a map of tables and their
// columns. For now, just tests that the request returns 200 and that the map
// has a key for each table
func Test_GET_Tables(t *testing.T) {
	_, rr := getRespBoilerplate(t, "/")
	tables := map[string]any{}
	err := json.Unmarshal(rr.Body.Bytes(), &tables)
	tests.Try(t, err)

	gotTables := slices.Collect(maps.Keys(tables))

	expectTables := []string{"authors", "books", "genres"}
	for _, table := range expectTables {
		if !slices.Contains(gotTables, table) {
			t.Errorf("Expected %s in got tables", table)
		}
	}
}

func apiGetRowsTester(
	t *testing.T,
	rawQuery string,
	path string,
) {
	ah, rr := getRespBoilerplate(t, path)

	gotRows := []types.RowData{}
	unmarshal(t, rr.Body.Bytes(), &gotRows)

	expRows, err := tests.SelectRows(ah.Repo, rawQuery)
	tests.Try(t, err)

	err = tests.CheckMapEquality(expRows, gotRows)
	if err != nil {
		t.Error(err)
	}
}

func getRespBoilerplate(t *testing.T, path string) (api.APIHandler, *httptest.ResponseRecorder) {
	ah := tests.NewTestAPIHandler(t)
	fullPath := fmt.Sprintf("%s", path)
	rr, err := tests.MakeHttpRequest(ah, http.MethodGet, fullPath, nil)
	tests.Try(t, err)

	if http.StatusOK != rr.Code {
		t.Fatalf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}
	return ah, rr
}

// a custom unmarshaller to convert float64 vals to int64. Fine for our tests
// since the test Tables don't have float64 values. Need this because
// json.Unmarshal takes all Number values to float64, but when we query the DB
// for rows via the Repo, the column types for e.g. id are int64
func unmarshal(t *testing.T, body []byte, dest *[]types.RowData) {
	err := json.Unmarshal(body, &dest)
	tests.Try(t, err)

	for _, row := range *dest {
		for k, v := range row {
			if fmt.Sprintf("%T", v) == "float64" {
				row[k] = int64(v.(float64))
			}
		}
	}
}
