package tests

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"

	"ftrack/repository"
	"ftrack/types"
)

// InsertSampleRows inserts sample rows into a repo
func InsertSampleRows(repo repository.Repository) SampleRowsIdMap {
	// sampleRows are used to populate the test database
	sampleRows := []types.RowDataMap{
		{
			"Name":   "deadlift",
			"Weight": 300,
		},
		{
			"Name":   "deadlift",
			"Weight": 200,
		},
		{
			"Name":   "deadlift",
			"Weight": 100,
		},
		{
			"Name":   "squat",
			"Weight": 300,
		},
		{
			"Name":   "squat",
			"Weight": 200,
		},
		{
			"Name":   "squat",
			"Weight": 100,
		},
		// Entries we will NOT filter for
		{
			"Name":   "bench press",
			"Weight": 300,
		},
	}
	insertedRows := make(SampleRowsIdMap)
	for _, sample := range sampleRows {
		result := repo.InsertRow(TABLE1, &sample)
		if result.Error != nil {
			panic("Failed to insert row, update insert tests")
		}
		insertedRows[result.ID] = sample
	}
	return insertedRows
}

func CheckExpectedErr(expectedErr any, err error) bool {
	return (expectedErr == nil && err != nil) ||
		(err != nil && err.Error() != expectedErr)
}

// ScanExerciseSetRow scans a row into an ExerciseSet struct
func ScanExerciseSetRow(toScan *ExerciseSet, row *sql.Row) error {
	return row.Scan(
		&toScan.ID,
		&toScan.Name,
		&toScan.PerformedAt,
		&toScan.Weight,
		&toScan.Unit,
		&toScan.Reps,
		&toScan.SetCount,
		&toScan.Notes,
		&toScan.SplitDay,
		&toScan.Program,
		&toScan.Tags,
	)
}

// ScanNextExerciseSetRow scans the next row into an ExerciseSet struct
func ScanNextExerciseSetRow(toScan *ExerciseSet, rows *sql.Rows) error {
	return rows.Scan(
		&toScan.ID,
		&toScan.Name,
		&toScan.PerformedAt,
		&toScan.Weight,
		&toScan.Unit,
		&toScan.Reps,
		&toScan.SetCount,
		&toScan.Notes,
		&toScan.SplitDay,
		&toScan.Program,
		&toScan.Tags,
	)
}

// FilterSampleRows filters the sample rows by a map of params
func FilterSampleRows(qf types.QueryFilter, sampleRows SampleRowsIdMap) []types.RowDataMap {
	m := []types.RowDataMap{}
	for _, row := range sampleRows {
		match := true
		for k := range row {
			filterValue, exists := qf[k]
			if exists && !slices.Contains(filterValue, fmt.Sprintf("%v", row[k])) {
				match = false
				break
			}
		}
		if match {
			m = append(m, row)
		}
	}
	return m
}

// GetTagMap returns a map of json tags from a struct, assuming it has any
func GetTagMap(s any) TagMap {
	val := reflect.ValueOf(s)
	tagMap := make(TagMap)
	t := val.Type()

	for i := range val.NumField() {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		fieldName := t.Field(i).Name
		tagMap[fieldName] = jsonTag
	}
	return tagMap
}
