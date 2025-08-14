package repository_test

import (
	"reflect"
	"testing"

	"gopgrest/test_utils"
	"gopgrest/types"
)

func Test_ListRows(t *testing.T) {
	repo, sampleRows := test_utils.NewTestRepo(t)

	t.Run("List all authors, no filter", func(t *testing.T) {
		// Expected authors are the same as sample authors (no filter)
		expAuthors := sampleRows.Authors
		// Query DB for all authors, no filters
		rows, err := repo.ListRows("authors", types.RSQLFilters{})
		if err != nil {
			t.Fatalf("List err: %s", err)
		}
		defer rows.Close()

		// Make a SampleAuthorsMap that should match the expected one
		gotAuthors := test_utils.SampleAuthorsMap{}
		for rows.Next() {
			got, err := test_utils.ScanAuthorFromRows(rows)
			if err != nil {
				t.Fatalf("List all authors scan err: %s\n", err)
			}
			gotAuthors[got.ID] = got
		}

		// Got should equal Expected
		for _, got := range gotAuthors {
			exp := expAuthors[got.ID]
			if !reflect.DeepEqual(got, exp) {
				t.Errorf("\nExp %v\nGot %v\n", exp, got)
			}
		}
	})

	t.Run("List authors, filter for forename 'Anne'", func(t *testing.T) {
		// Filter sample rows for authors with forname Anne
		expAuthors := test_utils.SampleAuthorsMap{}
		for id, author := range sampleRows.Authors {
			if author.Forename == "Anne" {
				expAuthors[id] = author
			}
		}

		// Query test db for authors with forename Anne
		rows, err := repo.ListRows("authors", types.RSQLFilters{"forename": []string{"Anne"}})
		if err != nil {
			t.Fatalf("List err: %s", err)
		}
		defer rows.Close()

		// Make a SampleAuthorsMap that should match the expected one
		gotAuthors := test_utils.SampleAuthorsMap{}
		for rows.Next() {
			got, err := test_utils.ScanAuthorFromRows(rows)
			if err != nil {
				t.Fatalf("List all authors scan err: %s\n", err)
			}
			gotAuthors[got.ID] = got
		}

		// Got should equal Expected
		for _, got := range gotAuthors {
			exp := expAuthors[got.ID]
			if !reflect.DeepEqual(got, exp) {
				t.Errorf("\nExp %v\nGot %v\n", exp, got)
			}
		}
	})
}
