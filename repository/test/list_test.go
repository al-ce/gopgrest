package repository_test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	_ "github.com/lib/pq"

	"ftrack/repository"
	"ftrack/tests"
	"ftrack/types"
)

type filterTest struct {
	testName  string
	filters   types.QueryFilters
	rowCount  int
	expectErr any
}

func makeFilterTest(testName string, qf types.QueryFilters, expectErr any) filterTest {
	return filterTest{
		testName,
		qf,
		len(tests.FilterSampleRows(qf)),
		expectErr,
	}
}

func TestListRows_InvalidFilters(t *testing.T) {
	tdb := tests.NewTestDB(t)

	invalidQueryTests := []filterTest{
		{
			"empty filter value",
			types.QueryFilters{
				"name": {},
			},
			0,
			"attempt to filter on key name with no values",
		},
		{
			"invalid column names",
			types.QueryFilters{
				"not_a_col": {"value"},
			},
			0,
			"pq: column \"not_a_col\" does not exist",
		},
		{
			"invalid column values",
			types.QueryFilters{
				"weight": {"not int"},
			},
			0,
			"pq: invalid input syntax for type smallint: \"not int\"",
		},
	}

	for _, tt := range invalidQueryTests {
		t.Run(tt.testName, func(t *testing.T) {
			tx := tdb.BeginTX(t)
			repo := repository.NewRepository(tx)
			_, err := repo.ListRows(tests.TABLE1, tt.filters)
			if tests.CheckExpectedErr(tt.expectErr, err) {
				t.Errorf("Expected error: %v\nGot %v", tt.expectErr, err)
			}
		})
	}
}

func TestListRows_NoFilters(t *testing.T) {
	t.Run("list all", func(t *testing.T) {
		tdb := tests.NewTestDB(t)
		tx := tdb.BeginTX(t)
		repo := repository.NewRepository(tx)
		tests.InsertSampleRows(repo)

		// List all rows in the table
		rows, err := repo.ListRows(tests.TABLE1, types.QueryFilters{})
		if tests.CheckExpectedErr(nil, err) {
			t.Errorf("Expected error: %v\nGot %v", nil, err)
		}

		// Track how many rows we got
		gotCount := 0
		scannedRow := tests.ExerciseSet{}
		for rows.Next() {

			// Scan rows into struct
			err := tests.ScanNextExerciseSetRow(&scannedRow, rows)
			if err != nil {
				t.Errorf("Scan err: %v", err)
			}

			// Confirm each column in the row matches the sample we inserted
			rowVal := reflect.ValueOf(scannedRow)
			sampleRow := tests.SampleRows[gotCount]
			for idx := range rowVal.Type().NumField() {
				fieldName := rowVal.Type().Field(idx).Name
				expectedVal := fmt.Sprintf("%v", sampleRow[fieldName])
				gotVal := fmt.Sprintf("%v", rowVal.Field(idx))
				if expectedVal != "<nil>" && expectedVal != gotVal {
					t.Errorf(
						"Expected %s: %v\nGot %v",
						fieldName,
						expectedVal,
						gotVal,
					)
				}
			}
			gotCount++
		}

		// Confirm we got the same amount of rows we inserted
		expectedCount := len(tests.SampleRows)
		if expectedCount != gotCount {
			t.Errorf(
				"Expected %d rows\nGot %d\n",
				len(tests.SampleRows),
				gotCount,
			)
		}
	})
}

func TestListRows_ValidFilters(t *testing.T) {
	tdb := tests.NewTestDB(t)

	tx := tdb.BeginTX(t)
	repo := repository.NewRepository(tx)

	tests.InsertSampleRows(repo)

	filterTests := []struct {
		testName  string
		filters   types.QueryFilters
		rowCount  int
		expectErr any
	}{
		makeFilterTest(
			"list deadlifts",
			types.QueryFilters{
				"Name": {"deadlift"},
			},
			nil,
		),
		makeFilterTest(
			"list deadlifts or squats",
			types.QueryFilters{
				"Name": {"deadlift", "squat"},
			},
			nil,
		),
		makeFilterTest(
			"list weights of 100",
			types.QueryFilters{
				"Weight": {"100"},
			},
			nil,
		),
		makeFilterTest(
			"list weights of 100 or 200",
			types.QueryFilters{
				"Weight": {"100", "200"},
			},
			nil,
		),
		makeFilterTest(
			"list squats of weight 200",
			types.QueryFilters{
				"Name":   {"squat"},
				"Weight": {"200"},
			},
			nil,
		),
		makeFilterTest(
			"list squats of weight 101 or 201",
			types.QueryFilters{
				"Name":   {"squat"},
				"Weight": {"100", "200"},
			},
			nil,
		),

		// Queries that should return 0 results
		makeFilterTest(
			// non-existent exercise name
			"list presses",
			types.QueryFilters{
				"Name": {"press"},
			},
			nil,
		),
		makeFilterTest(
			// valid exercise with no matching weight
			"list squats of weight 50",
			types.QueryFilters{
				"Name":   {"squat"},
				"Weight": {"50"},
			},
			nil,
		),
	}

	// doNotFilter contains filters we will never look for, but also values
	// that we used in our sample rows. This allows us to test that our query
	// params exclude rows we didn't filter for
	doNotFilter := types.QueryFilters{
		"Name":   {"bench press"},
		"Weight": {"300"},
	}

	for _, tt := range filterTests {
		t.Run(tt.testName, func(t *testing.T) {
			rows, err := repo.ListRows(tests.TABLE1, tt.filters)
			if tests.CheckExpectedErr(tt.expectErr, err) {
				t.Errorf("Expected error: %v\nGot %v", tt.expectErr, err)
			}

			scannedRow := tests.ExerciseSet{}
			allScanned := []tests.ExerciseSet{}

			for rows.Next() {

				err := tests.ScanNextExerciseSetRow(&scannedRow, rows)
				if err != nil {
					t.Errorf("Scan err: %v", err)
				}

				val := reflect.ValueOf(scannedRow)
				for fieldName, fieldFilters := range tt.filters {
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
