package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"ftrack/test_utils"
)

func TestAPI_Pick_ValidID(t *testing.T) {
	ah, sampleRows := test_utils.NewTestAPIHandler(t)
	tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})

	for id, sr := range sampleRows {
		// Make request
		path := fmt.Sprintf("/%s/%d", test_utils.TABLE1, id)
		rr, err := test_utils.MakeHttpRequest(ah, http.MethodGet, path, struct{}{})
		if err != nil {
			t.Error(err.Error())
		}
		if http.StatusOK != rr.Code {
			t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
		}

		// Decode response
		var resp test_utils.ExerciseSet
		err = json.NewDecoder(rr.Body).Decode(&resp)
		if err != nil {
			t.Errorf("Decode err")
		}

		// Check values of decoded response against sample row
		val := reflect.ValueOf(resp)
		for expCol, expVal := range sr {
			fieldName := test_utils.GetFieldNameByColName(tagMap, expCol, test_utils.ExerciseSet{})
			gotVal := val.FieldByName(fieldName).Interface()
			if expVal != gotVal {
				t.Errorf("\nExp %s: %v\nGot: %v", expCol, expVal, gotVal)
			}
		}
	}
}
