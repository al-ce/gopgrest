package repository_test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	_ "github.com/lib/pq"

	"ftrack/tests"
	"ftrack/types"
)

func TestRepo_ListRows_InvalidFilters(t *testing.T) {
	invalidQueryTests := tests.GetInvalidQueryTests()
	for _, tt := range invalidQueryTests {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, _ := tests.NewTestRepo(t)
			_, err := repo.ListRows(tests.TABLE1, tt.Filters)

			if tests.CheckExpectedErr(tt.PqErr, err) {
				t.Errorf("Expected error: %v\nGot %v", tt.PqErr, err)
			}
		})
	}
}

func TestRepo_ListRows_NoFilters(t *testing.T) {
	t.Run("list all", func(t *testing.T) {
		repo, sampleRows := tests.NewTestRepo(t)
		tagMap := tests.GetTagMap(tests.ExerciseSet{})

		// List all rows in the table
		rows, err := repo.ListRows(tests.TABLE1, types.QueryFilter{})
		if tests.CheckExpectedErr(nil, err) {
			t.Errorf("Expected error: %v\nGot %v", nil, err)
		}

		// Track how many rows we got
		var gotCount int64 = 0
		scannedRow := tests.ExerciseSet{}
		for rows.Next() {

			// Scan rows into struct
			err := tests.ScanNextExerciseSetRow(&scannedRow, rows)
			if err != nil {
				t.Errorf("Scan err: %v", err)
			}

			// Confirm each column in the row matches the sample we inserted
			rowVal := reflect.ValueOf(scannedRow)
			sampleRow := sampleRows[gotCount]
			for idx := range rowVal.Type().NumField() {
				fieldName := tagMap[rowVal.Type().Field(idx).Name]
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
		expectedCount := int64(len(sampleRows))
		if expectedCount != gotCount {
			t.Errorf(
				"Expected %d rows\nGot %d\n",
				len(sampleRows),
				gotCount,
			)
		}
	})
}

func TestRepo_ListRows_ValidFilters(t *testing.T) {
	repo, sampleRows := tests.NewTestRepo(t)
	tagMap := tests.GetTagMap(tests.ExerciseSet{})
	filterTests := tests.GetValidFilterTests(sampleRows)

	// doNotFilter contains filters we will never look for, but also values
	// that we used in our sample rows. This allows us to test that our query
	// params exclude rows we didn't filter for
	doNotFilter := types.QueryFilter{
		"name":   {"bench press"},
		"weight": {"300"},
	}

	for _, tt := range filterTests {
		t.Run(tt.TestName, func(t *testing.T) {
			rows, err := repo.ListRows(tests.TABLE1, tt.Filters)
			if tests.CheckExpectedErr(tt.PqErr, err) {
				t.Errorf("Expected error: %v\nGot %v", tt.PqErr, err)
			}

			scannedRow := tests.ExerciseSet{}
			allScanned := []tests.ExerciseSet{}

			for rows.Next() {
				err := tests.ScanNextExerciseSetRow(&scannedRow, rows)
				if err != nil {
					t.Errorf("Scan err: %v", err)
				}

				rowVal := reflect.ValueOf(scannedRow)
				for colName, filterMap := range tt.Filters {
					fieldName := tests.GetFieldNameByColName(tagMap, colName, tests.ExerciseSet{})
					gotVal := fmt.Sprintf("%v", rowVal.FieldByName(fieldName))

					// Rows should include values we filtered for we know exist
					if !slices.Contains(filterMap, gotVal) {
						t.Errorf("Expected %s: %s in %s", colName, gotVal, filterMap)
					}
					// Rows should exclude values we never look for
					if slices.Contains(doNotFilter[colName], gotVal) {
						t.Errorf("False positive %s: %s", colName, gotVal)
					}
				}

				allScanned = append(allScanned, scannedRow)
			}

			if tt.RowCount != len(allScanned) {
				t.Errorf("Expected %d rows\nGot %d\n%v", tt.RowCount, len(allScanned), allScanned)
			}
		})
	}
}
