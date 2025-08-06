package api_test

import (
	"fmt"
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
