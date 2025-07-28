package repository_test

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	"ftrack/tests"
	"ftrack/types"
)

func TestInsertRow(t *testing.T) {
	insertTests := []struct {
		name       string
		newRow     types.RowDataMap
		expectRows int64
		expectErr  any
	}{
		{
			"ins row with valid col names/values",
			types.RowDataMap{
				"name":   "deadlift",
				"weight": 200,
				"reps":   10,
			},
			1,
			nil,
		},
		{
			"ins row with missing req cols",
			types.RowDataMap{
				"weight": 200,
				"reps":   10,
			},
			0,
			"pq: null value in column \"name\" of relation \"exercise_sets\" violates not-null constraint",
		},
		{
			"ins row with invalid values",
			types.RowDataMap{
				"weight": "not int",
			},
			0,
			"pq: invalid input syntax for type smallint: \"not int\"",
		},
		{
			"ins with invalid col names",
			types.RowDataMap{
				"not_a_col": 10,
			},
			0,
			fmt.Sprintf(
				"pq: column \"not_a_col\" of relation \"%s\" does not exist",
				tests.TABLE1),
		},
	}

	for _, tt := range insertTests {
		t.Run(tt.name, func(t *testing.T) {
			// Need new transaction for each subtest since some will be aborted
			// when they fail
			repo, _ := tests.NewTestRepo(t)

			result := repo.InsertRow(tests.TABLE1, &tt.newRow)
			if tests.CheckExpectedErr(tt.expectErr, result.Error) {
				t.Errorf("Expected error: %v\nGot %v", tt.expectErr, result.Error)
			}
		})
	}
}
