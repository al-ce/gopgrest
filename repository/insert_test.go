package repository_test

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"

	"gopgrest/test_utils"
)

func TestRepo_InsertRow(t *testing.T) {
	for _, tt := range test_utils.GetInsertTests() {
		t.Run(tt.Name, func(t *testing.T) {
			// Need new transaction for each subtest since some will be aborted
			// when they fail
			repo, _ := test_utils.NewTestRepo(t)
			result := repo.InsertRow(test_utils.TABLE1, &tt.NewRow)
			if fmt.Sprintf("%s", tt.PqErr) != fmt.Sprintf("%s", result.Error) {
				t.Errorf("\nExp: %s\nGot %s", tt.PqErr, result.Error)
			}
		})
	}
}
