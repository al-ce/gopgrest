package repository_test

import (
	"testing"

	"gopgrest/apperrors"
	"gopgrest/test_utils"
)

func TestRepo_DeleteRow(t *testing.T) {
	repo, sampleRows := test_utils.NewTestRepo(t)

	t.Run("delete row with valid id", func(t *testing.T) {
		scannedRow := test_utils.ExerciseSet{}
		for id := range sampleRows {
			err := repo.DeleteRow(test_utils.TABLE1, id)
			if err != nil {
				t.Errorf("Delete exec err: %v", err)
			}
			// Try to get deleted row by id
			row := repo.GetRowByID(test_utils.TABLE1, id)
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
		err := repo.DeleteRow(test_utils.TABLE1, -1)
		expErr := apperrors.NewDeleteInvalidIDErr(test_utils.TABLE1, -1)
		if err.Error() != expErr.Error() {
			t.Errorf("\nExp: %s\nGot: %s", expErr, err)
		}
	})
}
