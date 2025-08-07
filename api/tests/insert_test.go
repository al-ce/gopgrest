package api_test

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"gopgrest/test_utils"
)

func TestAPI_Insert_ValidReq(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)
	tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})

	path := fmt.Sprintf("/%s", test_utils.TABLE1)
	reqData := map[string]any{"name": "deadlift", "weight": 325}
	rr, err := test_utils.MakeHttpRequest(ah, http.MethodPost, path, reqData)
	if err != nil {
		t.Error(err.Error())
	}

	if http.StatusOK != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}

	// Check response body for id
	respBody := rr.Body.String()
	ReInsertedId := regexp.MustCompile(`^.*?([0-9]+).*$`)
	match := ReInsertedId.FindStringSubmatch(respBody)
	if len(match) != 2 {
		t.Errorf("Expected respBody match on pattern: %v\nGot: %v", ReInsertedId, match)
	}
	id := match[1]
	idInt, err := strconv.ParseInt(id, 10, 64)

	// Check db for inserted row by id
	scannedRow := test_utils.ExerciseSet{}
	row := ah.Repo.GetRowByID(test_utils.TABLE1, idInt)
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

	path := fmt.Sprintf("/%s", test_utils.TABLE1)
	reqData := map[string]any{"weight": 325}
	rr, err := test_utils.MakeHttpRequest(ah, http.MethodPost, path, reqData)
	if err != nil {
		t.Errorf("MakeRequest err: %s", err)
	}

	if http.StatusInternalServerError != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestAPI_Insert_InvalidField(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)

	path := fmt.Sprintf("/%s", test_utils.TABLE1)
	reqData := map[string]any{"page_count": 1000}
	rr, err := test_utils.MakeHttpRequest(ah, http.MethodPost, path, reqData)
	if err != nil {
		t.Errorf("MakeRequest err: %s", err)
	}

	if http.StatusInternalServerError != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusInternalServerError, rr.Code)
	}
}
