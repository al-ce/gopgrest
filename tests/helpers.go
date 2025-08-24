package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"gopgrest/api"
	"gopgrest/repository"
	"gopgrest/service"
	"gopgrest/types"
)

func MakeHttpRequest(ah api.APIHandler, method, path string, reqData any) (*httptest.ResponseRecorder, error) {
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		path,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	rr := httptest.NewRecorder()
	ah.ServeHTTP(rr, req)
	return rr, nil
}

func CheckMapEquality(expRows, gotRows []types.RowData) error {
	if len(gotRows) != len(expRows) {
		return fmt.Errorf(
			"gotRows length %d does not match expRows length %d\nExp:\n%v\nGot:\n%v",
			len(gotRows),
			len(expRows),
			expRows,
			gotRows,
		)
	}
	for idx, expRow := range expRows {
		for k, expVal := range expRow {
			gotRow := gotRows[idx]
			gotVal, ok := gotRow[k]
			if !ok {
				return fmt.Errorf("Expected key %s in row %v", k, gotRow)
			}
			if gotVal != expVal {
				return fmt.Errorf(
					"\nExpected %s: %v (type %T)\nGot: %v (type %T)",
					k,
					expVal,
					expVal,
					gotVal,
					gotVal,
				)
			}
		}
	}
	return nil
}

func CountRows(repo repository.Repository, tableName, condition string) (int, error) {
	var count int64
	row := repo.DB.QueryRow(fmt.Sprintf(
		"SELECT COUNT(*) FROM %s %s",
		tableName,
		condition,
	))
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return int(count), nil
}

func SelectRows(repo repository.Repository, query string) ([]types.RowData, error) {
	rows, err := repo.DB.Query(query)
	if err != nil {
		return []types.RowData{}, err
	}
	defer rows.Close()
	gotRows, err := service.ScanRows(rows)
	if err != nil {
		return []types.RowData{}, err
	}
	return gotRows, nil
}

func ParseIDArrayResponse(t *testing.T, resp string) []int64 {
	// Expect that we got back an array of ids, like `[4, 5]`
	reInsertedIds := regexp.MustCompile(`^\[(\d+,? ?)+\]$`)
	match := reInsertedIds.FindStringSubmatch(resp)
	if len(match) != 2 {
		t.Fatalf("Expected resp body match on pattern: %v\nGot: %v", reInsertedIds, match)
	}
	idStrs := strings.Split(match[1], ", ")
	ids := make([]int64, len(idStrs))
	for i, _id := range idStrs {
		gotID, err := strconv.ParseInt(_id, 10, 64)
		Try(t, err)
		ids[i] = gotID
	}
	return ids
}

// Try fails the test if err is not nil
func Try(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

// AssertEq fails the test if got != exp
func AssertEq(t *testing.T, got any, exp any)  {
	if got != exp {
		t.Fatalf("%v != %v", got, exp)
	}
}

// AssertNotEq fails the test if got == exp
func AssertNotEq(t *testing.T, got any, exp any)  {
	if got == exp {
		t.Fatalf("%v == %v", got, exp)
	}
}
