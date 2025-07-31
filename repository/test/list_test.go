package repository_test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	_ "github.com/lib/pq"

	"ftrack/test_utils"
	"ftrack/types"
)

func TestRepo_ListRows_InvalidFilters(t *testing.T) {
	invalidQueryTests := test_utils.GetInvalidQueryTests()
	for _, tt := range invalidQueryTests {
		t.Run(tt.TestName, func(t *testing.T) {
			repo, _ := test_utils.NewTestRepo(t)
			_, err := repo.ListRows(test_utils.TABLE1, tt.Filters)

			if test_utils.CheckExpectedErr(tt.PqErr, err) {
				t.Errorf("Expected error: %v\nGot %v", tt.PqErr, err)
			}
		})
	}
}

func TestRepo_ListRows_NoFilters(t *testing.T) {
	t.Run("list all", func(t *testing.T) {
		repo, sampleRows := test_utils.NewTestRepo(t)
		tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})

		// List all rows in the table
		rows, err := repo.ListRows(test_utils.TABLE1, types.QueryFilter{})
		if test_utils.CheckExpectedErr(nil, err) {
			t.Errorf("Expected error: %v\nGot %v", nil, err)
		}

		// Track how many rows we got
		var gotCount int64 = 0
		scannedRow := test_utils.ExerciseSet{}
		for rows.Next() {

			// Scan rows into struct
			err := test_utils.ScanNextExerciseSetRow(&scannedRow, rows)
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
	repo, sampleRows := test_utils.NewTestRepo(t)
	tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})
	filterTests := test_utils.GetValidFilterTests(sampleRows)

	// doNotFilter contains filters we will never look for, but also values
	// that we used in our sample rows. This allows us to test that our query
	// params exclude rows we didn't filter for
	doNotFilter := types.QueryFilter{
		"name":   {"bench press"},
		"weight": {"300"},
	}

	for _, tt := range filterTests {
		t.Run(tt.TestName, func(t *testing.T) {
			rows, err := repo.ListRows(test_utils.TABLE1, tt.Filters)
			if test_utils.CheckExpectedErr(tt.PqErr, err) {
				t.Errorf("Expected error: %v\nGot %v", tt.PqErr, err)
			}

			scannedRow := test_utils.ExerciseSet{}
			allScanned := []test_utils.ExerciseSet{}

			for rows.Next() {
				err := test_utils.ScanNextExerciseSetRow(&scannedRow, rows)
				if err != nil {
					t.Errorf("Scan err: %v", err)
				}

				rowVal := reflect.ValueOf(scannedRow)
				for colName, filterMap := range tt.Filters {
					fieldName := test_utils.GetFieldNameByColName(tagMap, colName, test_utils.ExerciseSet{})
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
