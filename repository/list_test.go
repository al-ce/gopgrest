package repository_test

import (
	"testing"

	"gopgrest/repository"
	"gopgrest/rsql"
	"gopgrest/service"
	"gopgrest/test_utils"
	"gopgrest/types"
)

func Test_RepoListRows(t *testing.T) {
	repo := test_utils.NewTestRepo(t)

	// Test no RSQL query
	t.Run("GET /authors", func(t *testing.T) {
		expRows := []types.RowData{
			map[string]any{
				"born":     int64(1950),
				"died":     nil,
				"forename": "Anne",
				"id":       int64(1),
				"surname":  "Carson",
			},
			map[string]any{
				"born":     int64(1820),
				"died":     int64(1849),
				"forename": "Anne",
				"id":       int64(2),
				"surname":  "Brontë",
			},
			map[string]any{
				"born":     int64(1882),
				"died":     int64(1941),
				"forename": "Virginia",
				"id":       int64(3),
				"surname":  "Woolf",
			},
		}

		listRowsTester(t, repo, "authors", &rsql.Query{}, expRows)
	})

	// Test single filter: equality (rsql `==`, SQL `=`)
	t.Run("GET /authors?filter=forname==Anne", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
		}
		query := rsql.Query{Filters: filters}

		expRows := []types.RowData{
			map[string]any{
				"born":     int64(1950),
				"died":     nil,
				"forename": "Anne",
				"id":       int64(1),
				"surname":  "Carson",
			},
			map[string]any{
				"born":     int64(1820),
				"died":     int64(1849),
				"forename": "Anne",
				"id":       int64(2),
				"surname":  "Brontë",
			},
		}

		listRowsTester(t, repo, "authors", &query, expRows)
	})

	// Test single filter: inequality (rsql `!=`, SQL `!=`)
	t.Run("GET /authors?filter=surname!=Carson", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "surname", Values: []string{"Carson"}, SQLOperator: "!="},
		}
		query := rsql.Query{Filters: filters}

		expRows := []types.RowData{
			map[string]any{
				"born":     int64(1820),
				"died":     int64(1849),
				"forename": "Anne",
				"id":       int64(2),
				"surname":  "Brontë",
			},
			map[string]any{
				"born":     int64(1882),
				"died":     int64(1941),
				"forename": "Virginia",
				"id":       int64(3),
				"surname":  "Woolf",
			},
		}

		listRowsTester(t, repo, "authors", &query, expRows)
	})

	// Test multiple filter values (`;` separated)
	t.Run("GET /authors?filter=forname==Anne;surname==Carson", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
			{Column: "surname", Values: []string{"Carson"}, SQLOperator: "="},
		}
		query := rsql.Query{Filters: filters}

		expRows := []types.RowData{
			map[string]any{
				"born":     int64(1950),
				"died":     nil,
				"forename": "Anne",
				"id":       int64(1),
				"surname":  "Carson",
			},
		}

		listRowsTester(t, repo, "authors", &query, expRows)
	})
}

func listRowsTester(
	t *testing.T,
	repo repository.Repository,
	tableName string,
	query *rsql.Query,
	expRows []types.RowData,
) {
	rows, err := repo.ListRows(tableName, query)
	if err != nil {
		t.Fatalf("List err: %s", err)
	}
	defer rows.Close()

	gotRows, err := service.ScanRows(rows)
	if err != nil {
		t.Errorf("Scan err: %s", err)
	}
	checkMapEquality(t, expRows, gotRows)
}

func checkMapEquality(t *testing.T, expRows, gotRows []types.RowData) {
	if len(gotRows) != len(expRows) {
		t.Fatalf(
			"gotRows length %d does not match expRows length %d\nExp:\n%v\nGot:\n%v",
			len(gotRows),
			len(expRows),
			expRows,
			gotRows,
		)
	}
	for idx, expRow := range expRows {
		for k, expVal := range expRow {
			gotRow := gotRows[idx]
			gotVal, ok := gotRow[k]
			if !ok {
				t.Errorf("Expected key %s in row %v", k, gotRow)
			}
			if gotVal != expVal {
				t.Errorf(
					"Expected %s: %v (type %T)\nGot: %v (type %T)",
					k,
					expVal,
					expVal,
					gotVal,
					gotVal,
				)
			}
		}
	}
}
