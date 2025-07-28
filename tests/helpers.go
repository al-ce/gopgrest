package tests

import (
	"database/sql"
	"fmt"
	"slices"

	"ftrack/repository"
	"ftrack/types"
)

// InsertSampleRows inserts sample rows into a repo
func InsertSampleRows(repo repository.Repository) {
	for _, sample := range SampleRows {
		result := repo.InsertRow(TABLE1, &sample)
		if result.Error != nil {
			panic("Failed to insert row, update insert tests")
		}
	}
}

func CheckExpectedErr(expectedErr any, err error) bool {
	return (expectedErr == nil && err != nil) ||
		(err != nil && err.Error() != expectedErr)
}

// ScanExerciseSetRow scans a result row into an ExerciseSet struct
func ScanExerciseSetRow(toScan *ExerciseSet, rows *sql.Rows) error {
	err := rows.Scan(
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
	return err
}

// FilterSampleRows filters the sample rows by a map of params
func FilterSampleRows(qf types.QueryFilters) []types.RowDataMap {
	m := []types.RowDataMap{}
	for _, row := range SampleRows {
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
