package tests

import (
	"database/sql"

	"ftrack/repository"
)

// InsertSampleRows inserts sample rows into a repo
func InsertSampleRows(repo repository.Repository, sampleRows []map[string]any) {
	for _, sample := range sampleRows {
		repo.InsertRow(TABLE1, &sample)
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
