package repository_test

import (
	"fmt"
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
			if fmt.Sprintf("%s", tt.PqErr) != fmt.Sprintf("%s", result.Error) {
				t.Errorf("\nExp: %s\nGot %s", tt.PqErr, result.Error)
			}
		})
	}
}
