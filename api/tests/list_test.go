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
		t.Errorf("Unmarshal err: %s", err)
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
