package tests

import (
	"fmt"
	"testing"

	"gopgrest/types"
)

func Test_ServiceGetRowByID(t *testing.T) {
	service := NewTestService(t)
	expAuthors, err := selectRows(service.Repo, "SELECT * FROM authors ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}

	for index, auth := range expAuthors {
		idAsStr := fmt.Sprintf("%d", index+1)
		gotRowData, err := service.GetRowByID("authors", idAsStr)
		if err != nil {
			t.Fatalf("Could not pick author id %s: %s", idAsStr, err)
		}
		if err := checkMapEquality([]types.RowData{auth}, []types.RowData{gotRowData}); err != nil {
			t.Error(err)
		}
	}
}
