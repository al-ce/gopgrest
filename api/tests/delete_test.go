package api_test

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"gopgrest/apperrors"
	"gopgrest/test_utils"
)

func TestAPI_Delete_ValidID(t *testing.T) {
	ah, sampleRows := test_utils.NewTestAPIHandler(t)

	scannedRow := test_utils.ExerciseSet{}
	for id := range sampleRows {
		// Make request
		path := fmt.Sprintf("/%s/%d", test_utils.TABLE1, id)
		rr, err := test_utils.MakeHttpRequest(ah, http.MethodDelete, path, struct{}{})
		if err != nil {
			t.Error(err.Error())
		}
		if http.StatusOK != rr.Code {
			t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
		}

		// Try to get deleted row by id
		row := ah.Repo.GetRowByID(test_utils.TABLE1, id)
		err = test_utils.ScanExerciseSetRow(&scannedRow, row)
		if err.Error() != "sql: no rows in result set" {
			t.Errorf(
				"Expected to delete row %d but found it:\n\t%v",
				id,
				scannedRow,
			)
		}

	}
}

func TestAPI_Delete_NonexistentID(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)
	id := int64(math.MaxInt32) // psql integer value is 4 bytes signed, assuming that type for id
	path := fmt.Sprintf("/%s/%d", test_utils.TABLE1, id)
	rr, err := test_utils.MakeHttpRequest(ah, http.MethodDelete, path, struct{}{})
	if err != nil {
		t.Error(err.Error())
	}
	if http.StatusInternalServerError != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}

	expErr := apperrors.NewDeleteInvalidIDErr(test_utils.TABLE1, id)
	if rr.Body.String() != expErr.Error() {
		t.Errorf("\nExp: %s\nGot: %s", expErr, rr.Body)
	}
}
