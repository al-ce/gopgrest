package tests

import (
	"fmt"
	"testing"

	"gopgrest/types"
)

func Test_RepoInsertRow(t *testing.T) {
	repo := NewTestRepo(t)
	author := types.RowData{
		"surname": "Sappho",
	}
	result := repo.InsertRow("authors", &author)
	if result.Error != nil {
		t.Fatalf("Could not insert author %v: %s", author, result.Error)
	}
	gotRows, err := selectRows(
		repo,
		fmt.Sprintf("SELECT surname FROM authors WHERE id = %d", result.ID),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotRows) != 1 {
		t.Fatalf("Expected 1 result for pick, got %d", len(gotRows))
	}
	if gotSurname, ok := gotRows[0]["surname"]; ok {
		if gotSurname != author["surname"] {
			t.Errorf("Expected '%s' Got '%s'", author["surname"], gotSurname)
		}
	} else {
		t.Fatalf("Picked row does not have 'surname' column")
	}
}
