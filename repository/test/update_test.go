package repository_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/tests"
	"ftrack/types"
)

type updateTest struct {
	testName  string
	id        string
	field     string
	value     any
	expectErr any
}

func TestUpdateRowCol(t *testing.T) {
	repo, _ := tests.NewTestRepo(t)

	sampleRow := types.RowData{
		"Name":   "romanian deadlift",
		"Weight": 309,
	}

	insertResult := repo.InsertRow(tests.TABLE1, &sampleRow)
	if insertResult.Error != nil {
		t.Errorf("Insert err %s", insertResult.Error)
	}

	updateTests := []updateTest{
		{
			"update valid string field",
			fmt.Sprintf("%d", insertResult.ID),
			"Name",
			"hack squat",
			nil,
		},
		{
			"update valid int field",
			fmt.Sprintf("%d", insertResult.ID),
			"Weight",
			299,
			nil,
		},
		{
			"update invalid field",
			fmt.Sprintf("%d", insertResult.ID),
			"not_a_col",
			"hack squat",
			"pq: column \"not_a_col\" of relation \"exercise_sets\" does not exist",
		},
	}

	for _, tt := range updateTests {
		t.Run(tt.testName, func(t *testing.T) {
			// Exec update query
			err := repo.UpdateRowCol(tests.TABLE1, tt.id, tt.field, tt.value)
			if tests.CheckExpectedErr(tt.expectErr, err) {
				t.Errorf("\nExp: %s\nGot: %s", tt.expectErr, err)
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
			gotVal := rowVal.FieldByName(tt.field).Interface()
			if tt.value != gotVal {
				t.Errorf("Expected %v: %s\nGot %v", tt.field, tt.value, gotVal)
			}
		})
	}
}
