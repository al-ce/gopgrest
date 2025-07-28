package repository_test

import (
	"fmt"
	"testing"

	"ftrack/repository"
	"ftrack/tests"
)

func TestDeleteRow(t *testing.T) {
	tdb := tests.GetTestDB(t)
	tx := tdb.BeginTX(t)
	repo := repository.NewRepository(tx)
	insertedSampleRows := tests.InsertSampleRows(repo)

	t.Run("delete row with valid id", func(t *testing.T) {
		scannedRow := tests.ExerciseSet{}
		for id := range insertedSampleRows {
			err := repo.DeleteRow(tests.TABLE1, fmt.Sprintf("%d", id))
			if err != nil {
				t.Errorf("Delete exec err: %v", err)
			}

			// Try to get deleted row by id
			row := repo.GetRowByID(tests.TABLE1, fmt.Sprintf("%d", id))
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
}
