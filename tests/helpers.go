package tests

import "ftrack/repository"

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
