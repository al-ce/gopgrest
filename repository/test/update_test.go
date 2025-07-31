package repository_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/test_utils"
	"ftrack/types"
)

func TestRepo_UpdateRowCol(t *testing.T) {
	repo, _ := test_utils.NewTestRepo(t)

	sampleRow := types.RowData{
		"name":   "romanian deadlift",
		"weight": 309,
	}

	insertResult := repo.InsertRow(test_utils.TABLE1, &sampleRow)
	if insertResult.Error != nil {
		t.Errorf("Insert err %s", insertResult.Error)
	}

	tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})

	updateTests := test_utils.GetUpdateTests(insertResult)

	for _, tt := range updateTests {
		t.Run(tt.TestName, func(t *testing.T) {
			// Exec update query
			err := repo.UpdateRowCol(test_utils.TABLE1, tt.ID, tt.Col, tt.Value)
			if test_utils.CheckExpectedErr(tt.PqErr, err) {
				t.Errorf("\nExp: %s\nGot: %s", tt.PqErr, err)
			}
			// Go to next test if this is an invalid update query
			if err != nil {
				return
			}

			// Get updated row
			updatedRow := test_utils.ExerciseSet{}
			err = test_utils.ScanExerciseSetRow(
				&updatedRow,
				repo.GetRowByID(
					test_utils.TABLE1,
					fmt.Sprintf("%d", insertResult.ID),
				),
			)
			if err != nil {
				t.Errorf("Scan err: %s", err)
			}
			// Confirm update
			rowVal := reflect.ValueOf(updatedRow)
			fieldName := test_utils.GetFieldNameByColName(tagMap, tt.Col, test_utils.ExerciseSet{})
			gotVal := rowVal.FieldByName(fieldName).Interface()
			if tt.Value != gotVal {
				t.Errorf("Expected %v: %s\nGot %v", tt.Col, tt.Value, gotVal)
			}
		})
	}
}
