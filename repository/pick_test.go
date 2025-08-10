package repository_test

import (
	"reflect"
	"testing"

	"gopgrest/test_utils"
)

func Test_GetRowByID(t *testing.T) {
	repo, sampleRows := test_utils.NewTestRepo(t)

	for _, exp := range sampleRows.Authors {
		// Get author from db
		row := repo.GetRowByID("authors", exp.ID)
		got, err := test_utils.ScanAuthor(row)
		if err != nil {
			t.Errorf("ScanAuthor err: %s\n", err)
		}
		// Got should equal Expected
		if !reflect.DeepEqual(got, exp) {
			t.Errorf("\nExp %v\nGot %v\n", exp, got)
		}
	}
}
