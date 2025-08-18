package service_test

import (
	"fmt"
	"testing"

	"gopgrest/tests"
	"gopgrest/types"
)

func Test_ServiceGetRowByID(t *testing.T) {
	service := tests.NewTestService(t)
	expAuthors, err := tests.SelectRows(service.Repo, "SELECT * FROM authors ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}

	for index, auth := range expAuthors {
		idAsStr := fmt.Sprintf("%d", index+1)
		gotRowData, err := service.GetRowByID("authors", idAsStr)
		if err != nil {
			t.Fatalf("Could not pick author id %s: %s", idAsStr, err)
		}
		if err := tests.CheckMapEquality([]types.RowData{auth}, []types.RowData{gotRowData}); err != nil {
			t.Error(err)
		}
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

func Test_ServiceGetRows_Filters(t *testing.T) {
	t.Run("Single filter: equality (rsql `==`, SQL `=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename = 'Anne'"
		url := "/authors?filter=forename==Anne"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: inequality (rsql `!=`, SQL `!=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname != 'Carson'"
		url := "/authors?filter=surname!=Carson"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: in (rsql `=in=`, SQL `IN`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname IN ('Carson', 'Woolf')"
		url := "/authors?filter=surname=in=Carson,Woolf"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: not in (rsql `=out=`, SQL `NOT IN`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE surname NOT IN ('Carson', 'Woolf')"
		url := "/authors?filter=surname=out=Carson,Woolf"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: like (rsql `=like=`, SQL `LIKE`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename LIKE 'Ann%'"
		url := "/authors?filter=forename=like=Ann%"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: not like (rsql `=notlike=`, SQL `NOT LIKE`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename NOT LIKE 'Ann%'"
		url := "/authors?filter=forename=notlike=Ann%"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: is null (rsql `=isnull=`, SQL `IS NULL`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE died IS NULL"
		url := "/authors?filter=died=isnull="
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: is not null (rsql `=isnotnull=`, SQL `IS NOT NULL`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE died IS NOT NULL"
		url := "/authors?filter=died=isnotnull="
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: less than or equal to (rsql `=le=`, SQL `<=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born <= 1882"
		url := "/authors?filter=born=le=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: less than (rsql `=lt=`, SQL `<`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born < 1882"
		url := "/authors?filter=born=lt=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: greater than or equal to (rsql `=ge=`, SQL `>=`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born >= 1882"
		url := "/authors?filter=born=ge=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Single filter: greater than (rsql `=gt=`, SQL `>`)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE born > 1882"
		url := "/authors?filter=born=gt=1882"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Multiple filter values (`;` separated)", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE forename = 'Anne' AND surname = 'Carson'"
		url := "/authors?filter=forename==Anne;surname==Carson"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Filter with qualifier", func(t *testing.T) {
		rawQuery := "SELECT * FROM authors WHERE authors.forename = 'Anne'"
		url := "/authors?filter=authors.forename==Anne"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})
}

func Test_ServiceGetRows_Fields(t *testing.T) {
	t.Run("Fields, no qualifiers or aliases", func(t *testing.T) {
		rawQuery := "SELECT forename, surname FROM authors"
		url := "/authors?fields=forename,surname"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Fields with aliases", func(t *testing.T) {
		rawQuery := "SELECT forename AS first_name, surname AS last_name FROM authors"
		url := "/authors?fields=forename:first_name,surname:last_name"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Fields with qualifiers", func(t *testing.T) {
		rawQuery := "SELECT authors.forename, authors.surname FROM authors"
		url := "/authors?fields=authors.forename,authors.surname"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})

	t.Run("Fields with aliases and qualifiers", func(t *testing.T) {
		rawQuery := "SELECT authors.forename AS first_name, authors.surname AS last_name FROM authors"
		url := "/authors?fields=authors.forename:first_name,authors.surname:last_name"
		serviceGetRowsTester(t, rawQuery, "authors", url)
	})
}

func Test_ServiceGetRows_Joins(t *testing.T) {
	t.Run("Single JOIN relation (rsql `join`, SQL `JOIN`) (inner join)", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books JOIN genres ON books.genre_id = genres.id"
		url := "/books?fields=title,name:genre&join==genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("Multiple join relations", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre, surname FROM books JOIN authors ON books.author_id = authors.id JOIN genres ON books.genre_id = genres.id"
		url := "/books?fields=title,name:genre,surname&join=authors:books.author_id==authors.id;genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("INNER JOIN relation (rsql `inner_join`, SQL `INNER JOIN`)", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books INNER JOIN genres on books.genre_id = genres.id"
		url := "/books?fields=title,name:genre&inner_join=genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("LEFT JOIN relation (rsql `left_join`, SQL `LEFT JOIN`", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books LEFT JOIN genres ON books.genre_id = genres.id"
		url := "/books?fields=title,name:genre&left_join=genres:books.genre_id==genres.id"
		serviceGetRowsTester(t, rawQuery, "books", url)
	})

	t.Run("RIGHT JOIN relation (rsql `right_join`, SQL `RIGHT JOIN`", func(t *testing.T) {
		rawQuery := "SELECT title, name AS genre FROM books RIGHT JOIN genres ON books.genre_id = genres.id"
		url := "/books?fields=title,name:genre&right_join=genres:books.genre_id==genres.id"
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
	if err != nil {
		t.Fatal(err)
	}

	// Testing GetRowsByRSQL result
	gotRows, err := service.GetRowsByRSQL(tableName, url)
	if err != nil {
		t.Fatal(err)
	}
	if err := tests.CheckMapEquality(expRows, gotRows); err != nil {
		t.Error(err)
	}
}
