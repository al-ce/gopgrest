package repository_test

import (
	"reflect"
	"testing"

	"gopgrest/test_utils"
	"gopgrest/types"
)

func Test_InsertRow(t *testing.T) {
	repo, _ := test_utils.NewTestRepo(t)

	t.Run("Insert row into authors", func(t *testing.T) {
		// Insert new row into test db
		newRow := types.RowData{"surname": "Ahmed", "forename": "Sara"}
		result := repo.InsertRow("authors", &newRow)
		if result.Error != nil {
			t.Errorf("Insert err: %s\n", result.Error)
		}
		// Get row we just inserted
		row := repo.DB.QueryRow("SELECT * FROM authors WHERE id = $1", result.ID)
		if row.Err() != nil {
			t.Errorf("Could not get author %v by id %d", newRow, result.ID)
		}
		got, err := test_utils.ScanAuthor(row)
		if err != nil {
			t.Errorf("Insert ScanAuthor err %s\n", err)
		}
		// Update newRow to reflect id we just got
		newRow["id"] = result.ID
		// Make SampleAuthor struct from new RowData
		exp := test_utils.AuthorRowDataToStruct(newRow)
		// Got should equal Expected
		if !reflect.DeepEqual(got, exp) {
			t.Errorf("\nExp %v\nGot %v\n", exp, got)
		}
	})
}
