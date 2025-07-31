package repository_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/test_utils"
)

func TestRepo_GetRowByID(t *testing.T) {
	repo, sampleRows := test_utils.NewTestRepo(t)

	t.Run("get row with valid id", func(t *testing.T) {
		scannedRow := test_utils.ExerciseSet{}
		tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})
		for id, sampleRow := range sampleRows {
			row := repo.GetRowByID(test_utils.TABLE1, fmt.Sprintf("%d", id))
			err := test_utils.ScanExerciseSetRow(&scannedRow, row)
			if err != nil {
				t.Errorf("Scan err: %v", err)
			}

			val := reflect.ValueOf(scannedRow)
			for colName, sampleValue := range sampleRow {
				fieldName := test_utils.GetFieldNameByColName(tagMap, colName, test_utils.ExerciseSet{})
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
