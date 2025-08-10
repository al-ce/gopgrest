package repository_test

import (
	"reflect"
	"testing"

	"gopgrest/test_utils"
	"gopgrest/types"
)

func Test_UpdateRow(t *testing.T) {
	repo, sampleRows := test_utils.NewTestRepo(t)

	// Get first sample author from sample rows
	author := test_utils.SampleAuthor{}
	for _, v := range sampleRows.Authors {
		author = v
		break
	}

	// Update author attrs in DB
	updatedRow := types.RowData{"surname": "Sappho", "forename": "", "id": author.ID}
	err := repo.UpdateRowCol("authors", author.ID, &updatedRow)
	if err != nil {
		t.Errorf("Update err: %s\n", err)
	}

	// Get author we just updated
	row := repo.DB.QueryRow("SELECT * FROM authors WHERE id = $1", author.ID)
	if row.Err() != nil {
		t.Errorf("Could not get author %v by id %d", updatedRow, author.ID)
	}
	got, err := test_utils.ScanAuthor(row)
	if err != nil {
		t.Errorf("Insert ScanAuthor err %s\n", err)
	}
	exp := test_utils.AuthorRowDataToStruct(updatedRow)
	// Got should equal Expected
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("\nExp %v\nGot %v\n", exp, got)
	}
}
