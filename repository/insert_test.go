package repository_test

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	"ftrack/repository"
	"ftrack/tests"
)

func TestInsertRow(t *testing.T) {
	tdb := tests.GetTestDB(t)

	repo := repository.NewRepository(tdb.DB)

	insertTests := []struct {
		name       string
		newRow     map[string]any
		expectRows int64
		expectErr  any
	}{
		{
			"ins row with valid col names/values",
			map[string]any{
				"name":   "deadlift",
				"weight": 200,
				"reps":   10,
			},
			1,
			nil,
		},
		{
			"ins row with missing req cols",
			map[string]any{
				"weight": 200,
				"reps":   10,
			},
			0,
			"pq: null value in column \"name\" of relation \"exercise_sets\" violates not-null constraint",
		},
		{
			"ins row with invalid values",
			map[string]any{
				"weight": "not int",
			},
			0,
			"pq: invalid input syntax for type smallint: \"not int\"",
		},
		{
			"ins with invalid col names",
			map[string]any{
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
			rowsCreated, err := repo.InsertRow(tests.TABLE1, &tt.newRow)
			if rowsCreated != tt.expectRows {
				t.Errorf("Expected rows: %d\nGot: %v", rowsCreated, tt.expectRows)
			}
			if (tt.expectErr == nil && err != nil) ||
				(err != nil && err.Error() != tt.expectErr) {
				t.Errorf("Expected error: %v\nGot %v", tt.expectErr, err)
			}
		})
	}
}
