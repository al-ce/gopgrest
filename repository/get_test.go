package repository_test

import (
	"testing"

	"gopgrest/rsql"
	"gopgrest/service"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_RepoGetRowById(t *testing.T) {
	repo := tests.NewTestRepo(t)
	expAuthors, err := tests.SelectRows(repo, "SELECT * FROM authors ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}

	for index, auth := range expAuthors {
		id := index + 1
		rows, err := repo.GetRowByID("authors", int64(id))
		if err != nil {
			t.Fatalf("Could not pick author id %d: %s", id, err)
		}
		defer rows.Close()
		gotRows, err := service.ScanRows(rows)
		if err != nil {
			t.Fatalf("Could not scan author id %d: %s", id, err)
		}
		if err := tests.CheckMapEquality([]types.RowData{auth}, gotRows); err != nil {
			t.Error(err)
		}
	}
}

func Test_RepoGetRows_NoQuery(t *testing.T) {
	// GET /authors
	t.Run("No RSQL query", func(t *testing.T) {
		rawQuery := "SELECT * from authors"
		repoGetRowsTester(t, rawQuery, "authors", rsql.QueryParams{})
	})
}

func Test_RepoGetRows_Where(t *testing.T) {
	// GET /authors?where=forname==Anne
	t.Run("Single condition: equality (rsql `==`, SQL `=`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
		}
		rawQuery := "SELECT * FROM authors WHERE forename = 'Anne'"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=surname!=Carson
	t.Run("Single condition: inequality (rsql `!=`, SQL `!=`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "surname", Values: []string{"Carson"}, SQLOperator: "!="},
		}
		rawQuery := "SELECT * FROM authors WHERE surname != 'Carson'"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=surname=in=Carson,Woolf
	t.Run("Single condition: in (rsql `=in=`, SQL `IN`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "surname", Values: []string{"Carson", "Woolf"}, SQLOperator: "IN"},
		}
		rawQuery := "SELECT * FROM authors WHERE surname IN ('Carson', 'Woolf')"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=surname=out=Carson,Woolf
	t.Run("Single condition: not in (rsql `=out=`, SQL `NOT IN`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "surname", Values: []string{"Carson", "Woolf"}, SQLOperator: "NOT IN"},
		}
		rawQuery := "SELECT * FROM authors WHERE surname NOT IN ('Carson', 'Woolf')"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=forename=like=Ann%
	t.Run("Single condition: like (rsql `=like=`, SQL `LIKE`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "forename", Values: []string{"Ann%"}, SQLOperator: "LIKE"},
		}
		rawQuery := "SELECT * FROM authors WHERE forename LIKE 'Ann%'"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=forename=notlike=Ann%
	t.Run("Single condition: not like (rsql `=notlike=`, SQL `NOT LIKE`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "forename", Values: []string{"Ann%"}, SQLOperator: "NOT LIKE"},
		}
		rawQuery := "SELECT * FROM authors WHERE forename NOT LIKE 'Ann%'"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=died=isnull=
	t.Run("Single condition: is null (rsql `=isnull=`, SQL `IS NULL`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "died", Values: []string{}, SQLOperator: "IS NULL"},
		}
		rawQuery := "SELECT * FROM authors WHERE died IS NULL"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=died=isnotnull=
	t.Run("Single condition: is not null (rsql `=isnotnull=`, SQL `IS NOT NULL`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "died", Values: []string{}, SQLOperator: "IS NOT NULL"},
		}
		rawQuery := "SELECT * FROM authors WHERE died IS NOT NULL"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=born=le=1882
	t.Run("Single condition: less than or equal to (rsql `=le=`, SQL `<=`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "born", Values: []string{"1882"}, SQLOperator: "<="},
		}
		rawQuery := "SELECT * FROM authors WHERE born <= 1882"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=born=le=1882
	t.Run("Single condition: less than (rsql `=lt=`, SQL `<`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "born", Values: []string{"1882"}, SQLOperator: "<"},
		}
		rawQuery := "SELECT * FROM authors WHERE born < 1882"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=born=ge=1882
	t.Run("Single condition: greater than or equal to (rsql `=ge=`, SQL `>=`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "born", Values: []string{"1882"}, SQLOperator: ">="},
		}
		rawQuery := "SELECT * FROM authors WHERE born >= 1882"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=born=gt=1882
	t.Run("Single condition: greater than (rsql `=gt=`, SQL `>`)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "born", Values: []string{"1882"}, SQLOperator: ">"},
		}
		rawQuery := "SELECT * FROM authors WHERE born > 1882"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=forename==Anne;surname==Carson
	t.Run("Multiple condition values (`;` separated)", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
			{Column: "surname", Values: []string{"Carson"}, SQLOperator: "="},
		}
		rawQuery := "SELECT * FROM authors WHERE forename = 'Anne' AND surname = 'Carson'"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?where=authors.forname==Anne
	t.Run("condition with qualifier", func(t *testing.T) {
		conditions := []rsql.Condition{
			{Column: "forename", Values: []string{"Anne"}, SQLOperator: "="},
		}
		rawQuery := "SELECT * FROM authors WHERE authors.forename = 'Anne'"
		rsqlQuery := rsql.QueryParams{Conditions: conditions}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})
}

func Test_RepoGetRows_Select(t *testing.T) {
	// GET /authors?select=forename,surname
	t.Run("Select, no qualifiers or aliases", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "forename"},
			{Name: "surname"},
		}
		rawQuery := "SELECT forename, surname FROM authors"
		rsqlQuery := rsql.QueryParams{Columns: columns}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?select=forename:first_name,surname:last_name
	t.Run("Columns with aliases", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "forename", Alias: "first_name"},
			{Name: "surname", Alias: "last_name"},
		}
		rawQuery := "SELECT forename AS first_name, surname AS last_name FROM authors"
		rsqlQuery := rsql.QueryParams{Columns: columns}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?select=authors.forename,authors.surname
	t.Run("Columns with qualifiers", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "forename", Qualifier: "authors"},
			{Name: "surname", Qualifier: "authors"},
		}
		rawQuery := "SELECT authors.forename, authors.surname FROM authors"
		rsqlQuery := rsql.QueryParams{Columns: columns}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})

	// GET /authors?select=authors.forename:first_name,authors.surname:last_name
	t.Run("Columns with aliases and qualifiers", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "forename", Alias: "first_name", Qualifier: "authors"},
			{Name: "surname", Alias: "last_name", Qualifier: "authors"},
		}
		rawQuery := "SELECT authors.forename AS first_name, authors.surname AS last_name FROM authors"
		rsqlQuery := rsql.QueryParams{Columns: columns}
		repoGetRowsTester(t, rawQuery, "authors", rsqlQuery)
	})
}

func Test_RepoGetRows_Joins(t *testing.T) {
	// GET /books?select=title,name&join=genres:books.genre_id==genres.id
	t.Run("Single JOIN relation (rsql `join`, SQL `JOIN`) (inner join)", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "title"},
			{Name: "name", Alias: "genre"},
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
		rawQuery := "SELECT title, name AS genre FROM books JOIN genres ON books.genre_id = genres.id"
		rsqlQuery := rsql.QueryParams{Columns: columns, Joins: joins}
		repoGetRowsTester(t, rawQuery, "books", rsqlQuery)
	})

	// GET /books?select=title,name:genre,surname&join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id
	t.Run("Multiple join relations",
		func(t *testing.T) {
			columns := []rsql.Column{
				{Name: "title"},
				{Name: "name", Alias: "genre"},
				{Name: "surname"},
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
			rawQuery := "SELECT title, name AS genre, surname FROM books JOIN authors ON books.author_id = authors.id JOIN genres ON books.genre_id = genres.id"
			rsqlQuery := rsql.QueryParams{Columns: columns, Joins: joins}
			repoGetRowsTester(t, rawQuery, "books", rsqlQuery)
		})

	// GET /books?select=title,surname&inner_join=genres:books.genre_id==genres.id
	t.Run("INNER JOIN relation (rsql `inner_join`, SQL `INNER JOIN`)", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "title"},
			{Name: "name", Alias: "genre"},
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
		rawQuery := "SELECT title, name AS genre FROM books INNER JOIN genres on books.genre_id = genres.id"
		rsqlQuery := rsql.QueryParams{Columns: columns, Joins: joins}
		repoGetRowsTester(t, rawQuery, "books", rsqlQuery)
	})

	// GET /books?select=title,name:genre&left_join=genres:books.genre_id==genres.id
	t.Run("LEFT JOIN relation (rsql `left_join`, SQL `LEFT JOIN`", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "title"},
			{Name: "name", Alias: "genre"},
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
		rawQuery := "SELECT title, name AS genre FROM books LEFT JOIN genres ON books.genre_id = genres.id"
		rsqlQuery := rsql.QueryParams{Columns: columns, Joins: joins}
		repoGetRowsTester(t, rawQuery, "books", rsqlQuery)
	})

	// GET /books?select=title,name:genre&right_join=genres:books.genre_id==genres.id
	t.Run("RIGHT JOIN relation (rsql `right_join`, SQL `RIGHT JOIN`", func(t *testing.T) {
		columns := []rsql.Column{
			{Name: "title"},
			{Name: "name", Alias: "genre"},
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
		rawQuery := "SELECT title, name AS genre FROM books RIGHT JOIN genres ON books.genre_id = genres.id"
		rsqlQuery := rsql.QueryParams{Columns: columns, Joins: joins}
		repoGetRowsTester(t, rawQuery, "books", rsqlQuery)
	})
}

func repoGetRowsTester(
	t *testing.T,
	rawQuery string,
	tableName string,
	rsqlQuery rsql.QueryParams,
) {
	repo := tests.NewTestRepo(t)

	expRows, err := tests.SelectRows(repo, rawQuery)
	if err != nil {
		t.Fatal(err)
	}

	// Testing GetRowsByRSQL result
	rows, err := repo.GetRowsByRSQL(tableName, rsqlQuery)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	gotRows, err := service.ScanRows(rows)
	if err != nil {
		t.Fatal(err)
	}

	if err := tests.CheckMapEquality(expRows, gotRows); err != nil {
		t.Error(err)
	}
}
