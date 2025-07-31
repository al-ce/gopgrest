package repository_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/tests"
	"ftrack/types"
)

func TestRepo_UpdateRowCol(t *testing.T) {
	repo, _ := tests.NewTestRepo(t)

	sampleRow := types.RowData{
		"name":   "romanian deadlift",
		"weight": 309,
	}

	insertResult := repo.InsertRow(tests.TABLE1, &sampleRow)
	if insertResult.Error != nil {
		t.Errorf("Insert err %s", insertResult.Error)
	}

	tagMap := tests.GetTagMap(tests.ExerciseSet{})

	updateTests := tests.GetUpdateTests(insertResult)

	for _, tt := range updateTests {
		t.Run(tt.TestName, func(t *testing.T) {
			// Exec update query
			err := repo.UpdateRowCol(tests.TABLE1, tt.ID, tt.Col, tt.Value)
			if tests.CheckExpectedErr(tt.PqErr, err) {
				t.Errorf("\nExp: %s\nGot: %s", tt.PqErr, err)
			}
			// Go to next test if this is an invalid update query
			if err != nil {
				return
			}

			// Get updated row
			updatedRow := tests.ExerciseSet{}
			err = tests.ScanExerciseSetRow(
				&updatedRow,
				repo.GetRowByID(
					tests.TABLE1,
					fmt.Sprintf("%d", insertResult.ID),
				),
			)
			if err != nil {
				t.Errorf("Scan err: %s", err)
			}
			// Confirm update
			rowVal := reflect.ValueOf(updatedRow)
			fieldName := tests.GetFieldNameByColName(tagMap, tt.Col, tests.ExerciseSet{})
			gotVal := rowVal.FieldByName(fieldName).Interface()
			if tt.Value != gotVal {
				t.Errorf("Expected %v: %s\nGot %v", tt.Col, tt.Value, gotVal)
			}
		})
	}
}
