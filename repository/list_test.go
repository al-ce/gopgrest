package repository_test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	_ "github.com/lib/pq"

	"ftrack/repository"
	"ftrack/tests"
)

func TestListRows_ValidQueries(t *testing.T) {
	tdb := tests.GetTestDB(t)

	tx := tdb.BeginTX(t)
	repo := repository.NewRepository(tx)

	tests.InsertSampleRows(repo, []map[string]any{
		{
			"name":   "deadlift",
			"weight": 300,
		},
		{
			"name":   "deadlift",
			"weight": 200,
		},
		{
			"name":   "deadlift",
			"weight": 100,
		},
		{
			"name":   "squat",
			"weight": 300,
		},
		{
			"name":   "squat",
			"weight": 200,
		},
		{
			"name":   "squat",
			"weight": 100,
		},
		// Entries we will NOT filter for
		{
			"name":   "bench press",
			"weight": 300,
		},
	})

	filterTests := []struct {
		testName  string
		params    map[string][]string
		rowCount  int
		expectErr any
	}{
		{
			"list deadlifts",
			map[string][]string{
				"Name": {"deadlift"},
			},
			3,
			nil,
		},
		{
			"list deadlifts or squats",
			map[string][]string{
				"Name": {"deadlift", "squat"},
			},
			6,
			nil,
		},
		{
			"list weights of 100",
			map[string][]string{
				"Weight": {"100"},
			},
			2,
			nil,
		},
		{
			"list weights of 100 or 200",
			map[string][]string{
				"Weight": {"100", "200"},
			},
			4,
			nil,
		},
		{
			"list squats of weight 200",
			map[string][]string{
				"Name":   {"squat"},
				"Weight": {"200"},
			},
			1,
			nil,
		},
		{
			"list squats of weight 101 or 201",
			map[string][]string{
				"Name":   {"squat"},
				"Weight": {"100", "200"},
			},
			2,
			nil,
		},

		// Queries that should return 0 results
		{
			// non-existent exercise name
			"list presses",
			map[string][]string{
				"Name": {"press"},
			},
			0,
			nil,
		},
		{
			// valid exercise with no matching weight
			"list squats of weight 50",
			map[string][]string{
				"Name":   {"squat"},
				"Weight": {"50"},
			},
			0,
			nil,
		},
	}

	// doNotFilter contains filters we will never look for, but also values
	// that we used in our sample rows. This allows us to test that our query
	// params exclude rows we didn't filter for
	doNotFilter := map[string][]string{
		"Name":   {"bench press"},
		"Weight": {"300"},
	}

	for _, tt := range filterTests {
		t.Run(tt.testName, func(t *testing.T) {
			rows, err := repo.ListRows(tests.TABLE1, tt.params)
			if tests.CheckExpectedErr(tt.expectErr, err) {
				t.Errorf("Expected error: %v\nGot %v", tt.expectErr, err)
			}

			scannedRow := tests.ExerciseSet{}
			allScanned := []tests.ExerciseSet{}

			for rows.Next() {

				err := tests.ScanExerciseSetRow(&scannedRow, rows)
				if err != nil {
					t.Errorf("Scan err: %v", err)
				}

				val := reflect.ValueOf(scannedRow)
				for fieldName, fieldFilters := range tt.params {
					gotVal := fmt.Sprintf("%v", val.FieldByName(fieldName))
					// Rows should include values we filtered for we know exist
					if !slices.Contains(fieldFilters, gotVal) {
						t.Errorf("Expected %s: %s in %s", fieldName, gotVal, fieldFilters)
					}
					// Rows should exclude values we never look for
					if slices.Contains(doNotFilter[fieldName], gotVal) {
						t.Errorf("False positive %s: %s", fieldName, gotVal)
					}
				}

				allScanned = append(allScanned, scannedRow)
			}

			if tt.rowCount != len(allScanned) {
				t.Errorf("Expected %d rows\nGot %d\n%v", tt.rowCount, len(allScanned), allScanned)
			}
		})
	}
}
