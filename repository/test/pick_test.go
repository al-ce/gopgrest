package repository_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/tests"
)

func TestRepo_GetRowByID(t *testing.T) {
	repo, sampleRows := tests.NewTestRepo(t)

	t.Run("get row with valid id", func(t *testing.T) {
		scannedRow := tests.ExerciseSet{}
		tagMap := tests.GetTagMap(tests.ExerciseSet{})
		for id, sampleRow := range sampleRows {
			row := repo.GetRowByID(tests.TABLE1, fmt.Sprintf("%d", id))
			err := tests.ScanExerciseSetRow(&scannedRow, row)
			if err != nil {
				t.Errorf("Scan err: %v", err)
			}

			val := reflect.ValueOf(scannedRow)
			for colName, sampleValue := range sampleRow {
				fieldName := tests.GetFieldNameByColName(tagMap, colName, tests.ExerciseSet{})
				gotVal := val.FieldByName(fieldName).Interface()
				msg := fmt.Sprintf("\nExpected %s: %v %T\nGot %v %T",
					colName, sampleValue, sampleValue, gotVal, gotVal)
				if sampleValue != gotVal {
					t.Error(msg)
				}
			}

		}
	})
}
