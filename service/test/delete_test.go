package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
)

func TestService_DeleteRow(t *testing.T) {
	serv, sampleRows := tests.NewTestService(t)

	t.Run("delete with valid ids", func(t *testing.T) {
		scannedRow := tests.ExerciseSet{}
		for id := range sampleRows {
			err := serv.DeleteRow(tests.TABLE1, fmt.Sprintf("%d", id))
			if err != nil {
				t.Errorf("pick err: %s", err)
			}
			// Try to get deleted row by id,
			row := serv.Repo.GetRowByID(tests.TABLE1, fmt.Sprintf("%d", id))
			err = tests.ScanExerciseSetRow(&scannedRow, row)
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
		err := serv.DeleteRow(tests.TABLE1, fmt.Sprintf("%d", -1))
		if fmt.Sprintf("%v", err) != fmt.Sprintf(
			"row %d in table %s does not exist, did not attempt delete",
			-1, tests.TABLE1,
		) {
			t.Errorf("Expected non-existent id, but delete was successful: %v", err)
		}
	})
}


