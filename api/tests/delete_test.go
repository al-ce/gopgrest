package api_test

import (
	"fmt"
	"math"
	"net/http"
	"testing"

	"ftrack/test_utils"
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
		row := ah.Repo.GetRowByID(test_utils.TABLE1, fmt.Sprintf("%d", id))
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
	id := math.MaxInt32 // psql integer value is 4 bytes signed, assuming that type for id
	path := fmt.Sprintf("/%s/%d", test_utils.TABLE1, id)
	rr, err := test_utils.MakeHttpRequest(ah, http.MethodDelete, path, struct{}{})
	if err != nil {
		t.Error(err.Error())
	}
	if http.StatusInternalServerError != rr.Code {
		t.Errorf("\nExp StatusCode: %d\nGot: %d", http.StatusOK, rr.Code)
	}

	expErr := fmt.Sprintf(
		"row %d in table exercise_sets does not exist, did not attempt delete", id,
	)
	if rr.Body.String() != expErr {
		t.Errorf("\nExp: %s\nGot: %s", expErr, rr.Body)
	}
}
