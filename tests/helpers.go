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
func InsertSampleRows(repo repository.Repository) types.RowDataIdMap {
	// sampleRows are used to populate the test database
	sampleRows := []types.RowData{
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
	}
	insertedRows := make(types.RowDataIdMap)
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
func FilterSampleRows(qf types.QueryFilter, sampleRows types.RowDataIdMap) types.RowDataIdMap {
	m := make(types.RowDataIdMap)
	for id, row := range sampleRows {
		match := true
		for k := range row {
			filterValue, exists := qf[k]
			if exists && !slices.Contains(filterValue, fmt.Sprintf("%v", row[k])) {
				match = false
				break
			}
		}
		if match {
			m[id] = row
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

func GetFieldNameByColName(tm TagMap, colName string, s any) string {
	val := reflect.ValueOf(s)
	t := val.Type()
	for i := range val.NumField() {
		fieldName := t.Field(i).Name
		if tm[fieldName] == colName {
			return fieldName
		}
	}
	panic(fmt.Sprintf("col name %s not found in tagmap %v for %v", colName, tm, s))
}

// MakeFilterTest constructs test params for testing a filtered ListRows call
func MakeFilterTest(testName string, qf types.QueryFilter, sampleRows types.RowDataIdMap, expectErr any) FilterTest {
	return FilterTest{
		TestName:  testName,
		Filters:   qf,
		RowCount:  len(FilterSampleRows(qf, sampleRows)),
		ExpectErr: expectErr,
	}
}

func GetValidFilterTests(sampleRows types.RowDataIdMap) []FilterTest {
	return []FilterTest{
		MakeFilterTest(
			"list deadlifts",
			types.QueryFilter{
				"name": {"deadlift"},
			},
			sampleRows,
			nil,
		),
		MakeFilterTest(
			"list deadlifts or squats",
			types.QueryFilter{
				"name": {"deadlift", "squat"},
			},
			sampleRows,
			nil,
		),
		MakeFilterTest(
			"list weights of 100",
			types.QueryFilter{
				"weight": {"100"},
			},
			sampleRows,
			nil,
		),
		MakeFilterTest(
			"list weights of 100 or 200",
			types.QueryFilter{
				"weight": {"100", "200"},
			},
			sampleRows,
			nil,
		),
		MakeFilterTest(
			"list squats of weight 200",
			types.QueryFilter{
				"name":   {"squat"},
				"weight": {"200"},
			},
			sampleRows,
			nil,
		),
		MakeFilterTest(
			"list squats of weight 101 or 201",
			types.QueryFilter{
				"name":   {"squat"},
				"weight": {"100", "200"},
			},
			sampleRows,
			nil,
		),

		// Queries that should return 0 results
		MakeFilterTest(
			// non-existent exercise name
			"list presses",
			types.QueryFilter{
				"name": {"press"},
			},
			sampleRows,
			nil,
		),
		MakeFilterTest(
			// valid exercise with no matching weight
			"list squats of weight 50",
			types.QueryFilter{
				"name":   {"squat"},
				"weight": {"50"},
			},
			sampleRows,
			nil,
		),
	}
}

func GetInvalidQueryTests() []FilterTest {
	return []FilterTest{
		MakeFilterTest(
			"empty filter value",
			types.QueryFilter{
				"name": {},
			},
			types.RowDataIdMap{},
			"attempt to filter on key name with no values",
		),
		MakeFilterTest(
			"invalid column names",
			types.QueryFilter{
				"not_a_col": {"value"},
			},
			types.RowDataIdMap{},
			"pq: column \"not_a_col\" does not exist",
		),
		MakeFilterTest(
			"invalid column values",
			types.QueryFilter{
				"weight": {"not int"},
			},
			types.RowDataIdMap{},
			"pq: invalid input syntax for type smallint: \"not int\"",
		),
	}
}
