package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"ftrack/api"
	"ftrack/test_utils"
)

func TestAPI_Insert_ValidReq(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)
	reqData := map[string]any{"name": "deadlift", "weight": 325}
	tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})

	rr, err := makeInsertRequest(ah, reqData)
	if err != nil {
		t.Error(err.Error())
	}

	if http.StatusOK != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}

	// Check response body for id
	respBody := rr.Body.String()
	ReInsertedId := regexp.MustCompile(`^row\s([0-9])+.*$`)
	match := ReInsertedId.FindStringSubmatch(respBody)
	if len(match) != 2 {
		t.Errorf("Expected respBody match on pattern: %v\nGot: %v", ReInsertedId, match)
	}
	id := match[1]

	// Check db for inserted row by id
	scannedRow := test_utils.ExerciseSet{}
	row := ah.Repo.GetRowByID(test_utils.TABLE1, fmt.Sprintf("%s", id))
	err = test_utils.ScanExerciseSetRow(&scannedRow, row)
	if err != nil {
		t.Errorf("Scan err: %s", err)
	}

	// Compare inserted values to got values
	val := reflect.ValueOf(scannedRow)
	for reqColName, reqVal := range reqData {
		fieldName := test_utils.GetFieldNameByColName(tagMap, reqColName, test_utils.ExerciseSet{})
		gotVal := val.FieldByName(fieldName).Interface()
		if reqVal != gotVal {
			t.Errorf("\nExp %s: %v\nGot: %v", reqColName, reqVal, gotVal)
		}
	}
}

func TestAPI_Insert_MissingReqField(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)

	reqData := map[string]any{"weight": 325}
	rr, err := makeInsertRequest(ah, reqData)
	if err != nil {
		t.Errorf("MakeRequest err: %s", err)
	}

	if http.StatusInternalServerError != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestAPI_Insert_InvalidField(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)

	reqData := map[string]any{"page_count": 1000}
	rr, err := makeInsertRequest(ah, reqData)
	if err != nil {
		t.Errorf("MakeRequest err: %s", err)
	}

	if http.StatusInternalServerError != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusInternalServerError, rr.Code)
	}
}

func makeInsertRequest(ah api.APIHandler, reqData any) (*httptest.ResponseRecorder, error) {
	path := fmt.Sprintf("/%s", test_utils.TABLE1)
	rr, err := test_utils.MakeRequest(ah, reqData, http.MethodPost, path)
	if err != nil {
		return nil, err
	}
	return rr, nil
}
