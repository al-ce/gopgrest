package repository_test

import (
	"testing"

	"gopgrest/test_utils"
)

func Test_DeleteRow(t *testing.T) {
	repo, sampleRows := test_utils.NewTestRepo(t)

	for _, author := range sampleRows.Authors {
		// Delete author from db
		rowsAffected, err := repo.DeleteRow("authors", author.ID)
		if err != nil {
			t.Errorf("Delete err: %s\n", err)
		}
		// Expect 1 deleted row
		if rowsAffected != 1 {
			t.Errorf("Deleted 0 rows")
		}
		// Try to get author we just deleted
		row := repo.DB.QueryRow("SELECT * FROM authors WHERE id = $1", author.ID)
		author, err := test_utils.ScanAuthor(row)
		if err !=  nil {
			t.Errorf("ScanAuthor err: %s\n", err)
		}
		// Should have zero values for each field in SampleAuthor
		if author.ID != 0 || author.Surname != "" || author.Forename != "" {
			t.Errorf("Expected zero-values for scanned author: %v", author)
		}

	}
}
