package service_test

import (
	"fmt"
	"testing"

	"gopgrest/apperrors"
	"gopgrest/test_utils"
)

func TestService_DeleteRow(t *testing.T) {
	serv, sampleRows := test_utils.NewTestService(t)

	t.Run("delete with valid ids", func(t *testing.T) {
		scannedRow := test_utils.ExerciseSet{}
		for id := range sampleRows {
			err := serv.DeleteRow(test_utils.TABLE1, fmt.Sprintf("%d", id))
			if err != nil {
				t.Errorf("pick err: %s", err)
			}
			// Try to get deleted row by id,
			row := serv.Repo.GetRowByID(test_utils.TABLE1, id)
			err = test_utils.ScanExerciseSetRow(&scannedRow, row)
			if err.Error() != "sql: no rows in result set" {
				t.Errorf(
					"Expected to delete row %d but found it:\n\t%v",
					id,
					scannedRow,
				)
			}
		}
	})

	t.Run("delete row with non-existent id", func(t *testing.T) {
		err := serv.DeleteRow(test_utils.TABLE1, fmt.Sprintf("%d", -1))

		expErr := apperrors.NewDeleteInvalidIDErr(test_utils.TABLE1, -1)
		if err.Error() != expErr.Error() {
			t.Errorf("\nExp: %s\nGot: %s", expErr, err)
		}
	})
}
