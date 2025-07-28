package repository_test

import (
	"fmt"
	"reflect"
	"testing"

	"ftrack/repository"
	"ftrack/tests"
	"ftrack/types"
)

type updateTest struct {
	testName  string
	id        string
	field     string
	value     string
	expectErr any
}

func TestUpdateRowCol(t *testing.T) {
	tdb := tests.GetTestDB(t)

	tx := tdb.BeginTX(t)
	repo := repository.NewRepository(tx)

	sampleRow := types.RowDataMap{
		"Name":   "romanian deadlift",
		"Weight": 309,
	}

	insertResult := repo.InsertRow(tests.TABLE1, &sampleRow)
	if insertResult.Error != nil {
		t.Errorf("Insert err %s", insertResult.Error)
	}

	updateTests := []updateTest{
		{
			"update valid field",
			fmt.Sprintf("%d", insertResult.ID),
			"Name",
			"hack squat",
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
				repo.GetRowByID(tests.TABLE1, int(insertResult.ID)),
			)
			if err != nil {
				t.Errorf("Scan err: %s", err)
			}
			// Confirm update
			rowVal := reflect.ValueOf(updatedRow)
			gotVal := fmt.Sprintf("%v", rowVal.FieldByName(tt.field))
			if tt.value != gotVal {
				t.Errorf("Expected %s: %s\nGot %s", tt.field, tt.value, gotVal)
			}
		})
	}
}
