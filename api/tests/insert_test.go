package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"ftrack/test_utils"
)

func TestAPI_Insert(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)

	reqData := test_utils.ExerciseSet{Name: "deadlift", Weight: 325}

	path := fmt.Sprintf("/%s", test_utils.TABLE1)
	rr, err := test_utils.MakeRequest(ah, reqData, path)
	if err != nil {
		t.Errorf("MakeRequest err: %s", err)
	}
	if http.StatusOK != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}
}
