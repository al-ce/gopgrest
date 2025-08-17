package tests

import (
	"testing"

	"gopgrest/service"
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
	rows, err := repo.DB.Query("SELECT surname FROM authors WHERE id = $1", result.ID)
	if err != nil {
		t.Fatalf("Could not pick author id %d: %s", result.ID, err)
	}
	defer rows.Close()
	gotRows, err := service.ScanRows(rows)
	if err != nil {
		t.Fatalf("Could not scan author id %d: %s", result.ID, err)
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
