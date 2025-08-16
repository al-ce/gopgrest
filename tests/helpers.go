package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopgrest/api"
	"gopgrest/types"
)

func checkMapEquality(t *testing.T, expRows, gotRows []types.RowData) {
	if len(gotRows) != len(expRows) {
		t.Fatalf(
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
				t.Errorf("Expected key %s in row %v", k, gotRow)
			}
			if gotVal != expVal {
				t.Errorf(
					"Expected %s: %v (type %T)\nGot: %v (type %T)",
					k,
					expVal,
					expVal,
					gotVal,
					gotVal,
				)
			}
		}
	}
}

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
