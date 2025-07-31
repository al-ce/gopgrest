package service_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/tests"
	"ftrack/types"
)

func TestService_Update(t *testing.T) {
	serv, _ := tests.NewTestService(t)

	sampleRow := types.RowData{
		"name":   "romanian deadlift",
		"weight": 309,
	}

	insertResult := serv.Repo.InsertRow(tests.TABLE1, &sampleRow)
	if insertResult.Error != nil {
		t.Errorf("Insert err %s", insertResult.Error)
	}

	tagMap := tests.GetTagMap(tests.ExerciseSet{})

	updateTests := tests.GetUpdateTests(insertResult)

	for _, tt := range updateTests {
		t.Run(tt.TestName, func(t *testing.T) {
			updateData := types.RowData{tt.Col: tt.Value}

			// Exec update query
			err := serv.UpdateRow(tests.TABLE1, tt.ID, &updateData)
			if tests.CheckExpectedErr(tt.CustomErr, err) {
				t.Errorf("\nExp: %s\nGot: %s", tt.CustomErr, err)
			}
			// Go to next test if this is an invalid update query
			if err != nil {
				return
			}

			// Get updated row
			updatedRow := tests.ExerciseSet{}
			err = tests.ScanExerciseSetRow(
				&updatedRow,
				serv.Repo.GetRowByID(
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
