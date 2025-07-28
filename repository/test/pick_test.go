package repository_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/tests"
)

func TestGetRowByID(t *testing.T) {
	repo, sampleRows := tests.NewTestRepo(t)

	t.Run("get row with valid id", func(t *testing.T) {
		scannedRow := tests.ExerciseSet{}
		for id, sampleRow := range sampleRows {
			row := repo.GetRowByID(tests.TABLE1, fmt.Sprintf("%d", id))
			err := tests.ScanExerciseSetRow(&scannedRow, row)
			if err != nil {
				t.Errorf("Scan err: %v", err)
			}

			val := reflect.ValueOf(scannedRow)
			for fieldName, sampleValue := range sampleRow {
				gotVal := val.FieldByName(fieldName).Interface()
				if sampleValue != gotVal {
					t.Errorf("Expected %s: %v\nGot %v", fieldName, sampleValue, gotVal)
				}
			}

		}
	})
}
