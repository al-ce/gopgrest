package repository_test

import (
	"testing"

	_ "github.com/lib/pq"

	"ftrack/tests"
)

func TestRepo_InsertRow(t *testing.T) {
	for _, tt := range tests.GetInsertTests() {
		t.Run(tt.Name, func(t *testing.T) {
			// Need new transaction for each subtest since some will be aborted
			// when they fail
			repo, _ := tests.NewTestRepo(t)
			result := repo.InsertRow(tests.TABLE1, &tt.NewRow)
			if tests.CheckExpectedErr(tt.ExpectErr, result.Error) {
				t.Errorf("Expected error: %v\nGot %v", tt.ExpectErr, result.Error)
			}
		})
	}
}
