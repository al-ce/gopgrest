package tests

import (
	"testing"

	"gopgrest/rsql"
	"gopgrest/service"
	"gopgrest/types"
)

func Test_RepoListRows_NoQuery(t *testing.T) {
	// GET /authors
	t.Run("No RSQL query", func(t *testing.T) {
		expRows := []types.RowData{AnneCarson, AnneBrontë, VirginiaWoolf}
		listRowsTester(t, "authors", &rsql.Query{}, expRows)
	})
}

func Test_RepoListRows_Filters(t *testing.T) {
	// GET /authors?filter=forname==Anne
	t.Run("Single filter: equality (rsql `==`, SQL `=`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson, AnneBrontë}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=surname!=Carson
	t.Run("Single filter: inequality (rsql `!=`, SQL `!=`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "surname", Values: []string{"Carson"}, SQLOperator: "!="},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneBrontë, VirginiaWoolf}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=surname=in=Carson,Woolf
	t.Run("Single filter: in (rsql `=in=`, SQL `IN`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "surname", Values: []string{"Carson", "Woolf"}, SQLOperator: "IN"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson, VirginiaWoolf}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=surname=out=Carson,Woolf
	t.Run("Single filter: not in (rsql `=out=`, SQL `NOT IN`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "surname", Values: []string{"Carson", "Woolf"}, SQLOperator: "NOT IN"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneBrontë}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=forename=like=Ann%
	t.Run("Single filter: like (rsql `=like=`, SQL `LIKE`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Ann%"}, SQLOperator: "LIKE"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson, AnneBrontë}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=forename=notlike=Ann%
	t.Run("Single filter: not like (rsql `=notlike=`, SQL `NOT LIKE`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Ann%"}, SQLOperator: "NOT LIKE"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{VirginiaWoolf}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=died=isnull=
	t.Run("Single filter: is null (rsql `=isnull=`, SQL `IS NULL`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "died", Values: []string{}, SQLOperator: "IS NULL"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=died=isnotnull=
	t.Run("Single filter: is not null (rsql `=isnotnull=`, SQL `IS NOT NULL`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "died", Values: []string{}, SQLOperator: "IS NOT NULL"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneBrontë, VirginiaWoolf}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=born=le=1882
	t.Run("Single filter: less than or equal to (rsql `=le=`, SQL `<=`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "born", Values: []string{"1882"}, SQLOperator: "<="},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneBrontë, VirginiaWoolf}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=born=le=1882
	t.Run("Single filter: less than (rsql `=lt=`, SQL `<`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "born", Values: []string{"1882"}, SQLOperator: "<"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneBrontë}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=born=ge=1882
	t.Run("Single filter: greater than or equal to (rsql `=ge=`, SQL `>=`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "born", Values: []string{"1882"}, SQLOperator: ">="},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson, VirginiaWoolf}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=born=gt=1882
	t.Run("Single filter: greater than (rsql `=gt=`, SQL `>`)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "born", Values: []string{"1882"}, SQLOperator: ">"},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=forename==Anne;surname==Carson
	t.Run("Multiple filter values (`;` separated)", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
			{Column: "surname", Values: []string{"Carson"}, SQLOperator: "="},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?filter=authors.forname==Anne
	t.Run("Filter with qualifier", func(t *testing.T) {
		filters := []rsql.Filter{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
		}
		query := rsql.Query{Filters: filters}
		expRows := []types.RowData{AnneCarson, AnneBrontë}
		listRowsTester(t, "authors", &query, expRows)
	})
}

func Test_RepoListRows_Fields(t *testing.T) {
	// GET /authors?fields=forename,surname
	t.Run("Fields, no qualifiers or aliases", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "forename"},
			{Column: "surname"},
		}
		query := rsql.Query{Fields: fields}
		expRows := []types.RowData{
			map[string]any{
				"forename": "Anne",
				"surname":  "Carson",
			},
			map[string]any{
				"forename": "Anne",
				"surname":  "Brontë",
			},
			map[string]any{
				"forename": "Virginia",
				"surname":  "Woolf",
			},
		}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?fields=forename:first_name,surname:last_name
	t.Run("Fields with aliases", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "forename", Alias: "first_name"},
			{Column: "surname", Alias: "last_name"},
		}
		query := rsql.Query{Fields: fields}
		expRows := []types.RowData{
			map[string]any{
				"first_name": "Anne",
				"last_name":  "Carson",
			},
			map[string]any{
				"first_name": "Anne",
				"last_name":  "Brontë",
			},
			map[string]any{
				"first_name": "Virginia",
				"last_name":  "Woolf",
			},
		}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?fields=authors.forename,authors.surname
	t.Run("Fields with qualifiers", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "forename", Qualifier: "authors"},
			{Column: "surname", Qualifier: "authors"},
		}
		query := rsql.Query{Fields: fields}
		expRows := []types.RowData{
			map[string]any{
				"forename": "Anne",
				"surname":  "Carson",
			},
			map[string]any{
				"forename": "Anne",
				"surname":  "Brontë",
			},
			map[string]any{
				"forename": "Virginia",
				"surname":  "Woolf",
			},
		}
		listRowsTester(t, "authors", &query, expRows)
	})

	// GET /authors?fields=authors.forename:first_name,authors.surname:last_name
	t.Run("Fields with aliases and qualifiers", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "forename", Alias: "first_name", Qualifier: "authors"},
			{Column: "surname", Alias: "last_name", Qualifier: "authors"},
		}
		query := rsql.Query{Fields: fields}
		expRows := []types.RowData{
			map[string]any{
				"first_name": "Anne",
				"last_name":  "Carson",
			},
			map[string]any{
				"first_name": "Anne",
				"last_name":  "Brontë",
			},
			map[string]any{
				"first_name": "Virginia",
				"last_name":  "Woolf",
			},
		}
		listRowsTester(t, "authors", &query, expRows)
	})
}

func Test_RepoListRows_Joins(t *testing.T) {
	// GET /books?fields=title,surname&join=authors:books.author_id==authors.id;genres:book.genre_id==genres.id
	t.Run("Single JOIN relation (rsql `join`, SQL `JOIN`) (inner join)", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "title"},
			{Column: "name", Alias: "genre"},
		}
		joins := []rsql.JoinRelation{
			{
				Type:           "JOIN",
				Table:          "genres",
				LeftQualifier:  "books",
				LeftCol:        "genre_id",
				RightQualifier: "genres",
				RightCol:       "id",
			},
		}
		query := rsql.Query{Fields: fields, Joins: joins}
		expRows := []types.RowData{
			map[string]any{
				"title": "Autobiography of Red",
				"genre": "Romance",
			},
			map[string]any{
				"title": "The Tenant of Wildfell Hall",
				"genre": "Epistolary",
			},
			map[string]any{
				"title": "To The Lighthouse",
				"genre": "Modernism",
			},
		}
		listRowsTester(t, "books", &query, expRows)
	})

	// GET /books?fields=title,name:genre,surname&join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id
	t.Run("Multiple join relations",
		func(t *testing.T) {
			fields := []rsql.Field{
				{Column: "title"},
				{Column: "name", Alias: "genre"},
				{Column: "surname"},
			}
			joins := []rsql.JoinRelation{
				{
					Type:           "JOIN",
					Table:          "authors",
					LeftQualifier:  "books",
					LeftCol:        "author_id",
					RightQualifier: "authors",
					RightCol:       "id",
				},
				{
					Type:           "JOIN",
					Table:          "genres",
					LeftQualifier:  "books",
					LeftCol:        "genre_id",
					RightQualifier: "genres",
					RightCol:       "id",
				},
			}
			query := rsql.Query{Fields: fields, Joins: joins}
			expRows := []types.RowData{
				map[string]any{
					"surname": "Carson",
					"title":   "Autobiography of Red",
					"genre":   "Romance",
				},
				map[string]any{
					"surname": "Brontë",
					"title":   "The Tenant of Wildfell Hall",
					"genre":   "Epistolary",
				},
				map[string]any{
					"surname": "Woolf",
					"title":   "To The Lighthouse",
					"genre":   "Modernism",
				},
			}
			listRowsTester(t, "books", &query, expRows)
		})

	// GET /books?fields=title,surname&inner_join=authors:books.author_id==authors.id
	t.Run("INNER JOIN relation (rsql `inner_join`, SQL `INNER JOIN`)", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "title"},
			{Column: "name", Alias: "genre"},
		}
		joins := []rsql.JoinRelation{
			{
				Type:           "JOIN",
				Table:          "genres",
				LeftQualifier:  "books",
				LeftCol:        "genre_id",
				RightQualifier: "genres",
				RightCol:       "id",
			},
		}
		query := rsql.Query{Fields: fields, Joins: joins}
		expRows := []types.RowData{
			map[string]any{
				"title": "Autobiography of Red",
				"genre": "Romance",
			},
			map[string]any{
				"title": "The Tenant of Wildfell Hall",
				"genre": "Epistolary",
			},
			map[string]any{
				"title": "To The Lighthouse",
				"genre": "Modernism",
			},
		}
		listRowsTester(t, "books", &query, expRows)
	})

	// GET /books?fields=title,name:genre&left_join=genres:books.genre_id==genres.id
	t.Run("LEFT JOIN relation (rsql `left_join`, SQL `LEFT JOIN`", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "title"},
			{Column: "name", Alias: "genre"},
		}
		joins := []rsql.JoinRelation{
			{
				Type:           "LEFT JOIN",
				Table:          "genres",
				LeftQualifier:  "books",
				LeftCol:        "genre_id",
				RightQualifier: "genres",
				RightCol:       "id",
			},
		}
		query := rsql.Query{Fields: fields, Joins: joins}
		expRows := []types.RowData{
			map[string]any{
				"title": "Autobiography of Red",
				"genre": "Romance",
			},
			map[string]any{
				"title": "The Tenant of Wildfell Hall",
				"genre": "Epistolary",
			},
			map[string]any{
				"title": "To The Lighthouse",
				"genre": "Modernism",
			},
			map[string]any{
				"title": "Mrs. Dalloway",
				"genre": nil,
			},
		}
		listRowsTester(t, "books", &query, expRows)
	})

	// GET /books?fields=title,name:genre&right_join=genres:books.genre_id==genres.id
	t.Run("RIGHT JOIN relation (rsql `right_join`, SQL `RIGHT JOIN`", func(t *testing.T) {
		fields := []rsql.Field{
			{Column: "title"},
			{Column: "name", Alias: "genre"},
		}
		joins := []rsql.JoinRelation{
			{
				Type:           "RIGHT JOIN",
				Table:          "genres",
				LeftQualifier:  "books",
				LeftCol:        "genre_id",
				RightQualifier: "genres",
				RightCol:       "id",
			},
		}
		query := rsql.Query{Fields: fields, Joins: joins}
		expRows := []types.RowData{
			map[string]any{
				"title": "Autobiography of Red",
				"genre": "Romance",
			},
			map[string]any{
				"title": "The Tenant of Wildfell Hall",
				"genre": "Epistolary",
			},
			map[string]any{
				"title": "To The Lighthouse",
				"genre": "Modernism",
			},
			map[string]any{
				"title": nil,
				"genre": "Dystopian",
			},
		}
		listRowsTester(t, "books", &query, expRows)
	})
}

func listRowsTester(
	t *testing.T,
	tableName string,
	query *rsql.Query,
	expRows []types.RowData,
) {
	repo := NewTestRepo(t)
	rows, err := repo.ListRowsByRSQL(tableName, query)
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
