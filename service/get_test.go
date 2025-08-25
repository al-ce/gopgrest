package service_test

import (
	"fmt"
	"testing"

	"gopgrest/assert"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_ServiceGetRowByID(t *testing.T) {
	service := tests.NewTestService(t)
	expAuthors, err := tests.SelectRows(service.Repo, "SELECT * FROM authors ORDER BY id")
	assert.Try(t, err)

	for index, auth := range expAuthors {
		idAsStr := fmt.Sprintf("%d", index+1)
		gotRowData, err := service.GetRowByID("authors", idAsStr)
		assert.Try(t, err)
		err = tests.CheckMapEquality([]types.RowData{auth}, []types.RowData{gotRowData})
		assert.Try(t, err)
	}
}

func Test_RepoGetRows_NoQuery(t *testing.T) {
	// GET /authors
	t.Run("No RSQL query", func(t *testing.T) {
		rawQuery := "SELECT * from authors"
		url := "/authors"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})
}

func Test_ServiceGetRows_Where(t *testing.T) {
	t.Run("Single condition: equality (rsql `==`, SQL `=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename = 'Anne'"
		url := "/authors?where=forename==Anne"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: inequality (rsql `!=`, SQL `!=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname != 'Carson'"
		url := "/authors?where=surname!=Carson"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: in (rsql `=in=`, SQL `IN`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname IN ('Carson', 'Woolf')"
		url := "/authors?where=surname=in=Carson,Woolf"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: not in (rsql `=out=`, SQL `NOT IN`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname NOT IN ('Carson', 'Woolf')"
		url := "/authors?where=surname=out=Carson,Woolf"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: like (rsql `=like=`, SQL `LIKE`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename LIKE 'Ann%'"
		url := "/authors?where=forename=like=Ann%"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: not like (rsql `=notlike=`, SQL `NOT LIKE`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename NOT LIKE 'Ann%'"
		url := "/authors?where=forename=notlike=Ann%"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: is null (rsql `=isnull=`, SQL `IS NULL`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE died IS NULL"
		url := "/authors?where=died=isnull="
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: is not null (rsql `=isnotnull=`, SQL `IS NOT NULL`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE died IS NOT NULL"
		url := "/authors?where=died=isnotnull="
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: less than or equal to (rsql `=le=`, SQL `<=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born <= 1882"
		url := "/authors?where=born=le=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: less than (rsql `=lt=`, SQL `<`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born < 1882"
		url := "/authors?where=born=lt=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: greater than or equal to (rsql `=ge=`, SQL `>=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born >= 1882"
		url := "/authors?where=born=ge=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single condition: greater than (rsql `=gt=`, SQL `>`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born > 1882"
		url := "/authors?where=born=gt=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Multiple conditions (`;` separated)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename = 'Anne' AND surname = 'Carson'"
		url := "/authors?where=forename==Anne;surname==Carson"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Condition with qualifier", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE authors.forename = 'Anne'"
		url := "/authors?where=authors.forename==Anne"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})
}

func Test_ServiceGetRows_Select(t *testing.T) {
	t.Run("Select, no qualifiers or aliases", func(t *testing.T) {
		rawQuery := "SELECT forename, surname FROM authors"
		url := "/authors?select=forename,surname"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Columns with aliases", func(t *testing.T) {
		rawQuery := "SELECT forename AS first_name, surname AS last_name FROM authors"
		url := "/authors?select=forename:first_name,surname:last_name"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Columns with qualifiers", func(t *testing.T) {
		rawQuery := "SELECT authors.forename, authors.surname FROM authors"
		url := "/authors?select=authors.forename,authors.surname"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Columns with aliases and qualifiers", func(t *testing.T) {
		rawQuery := "SELECT authors.forename AS first_name, authors.surname AS last_name FROM authors"
		url := "/authors?select=authors.forename:first_name,authors.surname:last_name"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})
}

func Test_ServiceGetRows_Joins(t *testing.T) {
	t.Run("Single JOIN relation (rsql `join`, SQL `JOIN`) (inner join)", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books JOIN genres ON books.genre_id = genres.id"
		url := "/books?select=title,name:genre&join==genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("Multiple join relations", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre, surname FROM books JOIN authors ON books.author_id = authors.id JOIN genres ON books.genre_id = genres.id"
		url := "/books?select=title,name:genre,surname&join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("INNER JOIN relation (rsql `inner_join`, SQL `INNER JOIN`)", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books INNER JOIN genres on books.genre_id = genres.id"
		url := "/books?select=title,name:genre&inner_join=genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("LEFT JOIN relation (rsql `left_join`, SQL `LEFT JOIN`", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books LEFT JOIN genres ON books.genre_id = genres.id"
		url := "/books?select=title,name:genre&left_join=genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("RIGHT JOIN relation (rsql `right_join`, SQL `RIGHT JOIN`", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books RIGHT JOIN genres ON books.genre_id = genres.id"
		url := "/books?select=title,name:genre&right_join=genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})
}

func serviceGetRowsTester(
	t *testing.T,
	rawQuery string,
	tableName string,
	url string,
) {
	service := tests.NewTestService(t)

	expRows, err := tests.SelectRows(service.Repo, rawQuery)
	assert.Try(t, err)

	// Testing GetRowsByRSQL result
	gotRows, err := service.GetRowsByRSQL(tableName, url)
	assert.Try(t, err)
	if err != nil {
		t.Fatal(err)
	}
	err = tests.CheckMapEquality(expRows, gotRows)
	assert.Try(t, err)
}
