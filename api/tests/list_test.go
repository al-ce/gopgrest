package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"ftrack/test_utils"
	"ftrack/types"
)

func TestAPI_List_AllRows(t *testing.T) {
	ah, sampleRows := test_utils.NewTestAPIHandler(t)
	path := fmt.Sprintf("/%s", test_utils.TABLE1)
	rr, err := test_utils.MakeHttpRequest(ah, http.MethodGet, path, struct{}{})
	if err != nil {
		t.Error(err.Error())
	}
	if http.StatusOK != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}

	var respItems types.RowDataIdMap
	err = json.NewDecoder(rr.Body).Decode(&respItems)
	if err != nil {
		t.Errorf("Decode err: %s", err)
	}

	if len(respItems) != len(sampleRows) {
		t.Errorf(
			"Response row length (%d) does not match expected sample rows length (%d)",
			len(sampleRows),
			len(respItems),
		)
	}

	for id, sr := range sampleRows {
		for field, expVal := range sr {
			gotVal := respItems[id][field]
			if fmt.Sprintf("%v", gotVal) != fmt.Sprintf("%v", expVal) {
				t.Errorf("\nExpected %s: %s\nGot: %s", field, expVal, gotVal)
			}
		}
	}
}

func TestAPI_List_ValidFilters(t *testing.T) {
	ah, sampleRows := test_utils.NewTestAPIHandler(t)
	filterTests := test_utils.GetValidFilterTests(sampleRows)

	for _, tt := range filterTests {
		t.Run(tt.TestName, func(t *testing.T) {
			qp := makeQueryParams(tt.Filters)
			path := fmt.Sprintf("/%s?%s", test_utils.TABLE1, qp)
			rr, err := test_utils.MakeHttpRequest(ah, http.MethodGet, path, struct{}{})
			if err != nil {
				t.Error(err.Error())
			}
			if http.StatusOK != rr.Code {
				t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
			}

			// Decode response into a map of row data
			var respItems types.RowDataIdMap
			err = json.NewDecoder(rr.Body).Decode(&respItems)
			if err != nil {
				t.Errorf("Decode err: %s", err)
			}

			fsr := test_utils.FilterSampleRows(tt.Filters, sampleRows)

			if len(respItems) != len(fsr) {
				t.Errorf(
					"Response row length (%d) does not match expected filtered rows length (%d)",
					len(fsr),
					len(respItems),
				)
			}

			// Rows with matching ids should have same values we filtered for
			for idx, item := range respItems {
				for field, expVal := range fsr[idx] {
					gotVal := item[field]
					if fmt.Sprintf("%v", gotVal) != fmt.Sprintf("%v", expVal) {
						t.Errorf(
							"\nExp %s: %v %T\nGot: %v %T",
							field, expVal, expVal, gotVal, gotVal,
						)
					}
				}
			}
		})
	}
}

func makeQueryParams(qf types.QueryFilter) (queryParams string) {
	queryParams = ""
	for key, list := range qf {
		for _, item := range list {
			queryParams += fmt.Sprintf("&%s=%s", key, item)
		}
	}
	// Strip initial '&'
	return queryParams[1:]
}
