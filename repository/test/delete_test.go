package repository_test

import (
	"fmt"
	"testing"

	"ftrack/test_utils"
)

func TestRepo_DeleteRow(t *testing.T) {
	repo, sampleRows := test_utils.NewTestRepo(t)

	t.Run("delete row with valid id", func(t *testing.T) {
		scannedRow := test_utils.ExerciseSet{}
		for id := range sampleRows {
			err := repo.DeleteRow(test_utils.TABLE1, fmt.Sprintf("%d", id))
			if err != nil {
				t.Errorf("Delete exec err: %v", err)
			}
			// Try to get deleted row by id
			row := repo.GetRowByID(test_utils.TABLE1, fmt.Sprintf("%d", id))
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
		err := repo.DeleteRow(test_utils.TABLE1, fmt.Sprintf("%d", -1))
		if fmt.Sprintf("%v", err) != fmt.Sprintf(
			"row %d in table %s does not exist, did not attempt delete",
			-1, test_utils.TABLE1,
		) {
			t.Errorf("Expected non-existent id, but delete was successful: %v", err)
		}
	})
}
