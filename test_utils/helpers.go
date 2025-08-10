package test_utils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"

	"gopgrest/api"
	"gopgrest/types"
)

func MakeHttpRequest(ah api.APIHandler, method, path string, reqData any) (*httptest.ResponseRecorder, error) {
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		method,
		path,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}
	rr := httptest.NewRecorder()
	ah.ServeHTTP(rr, req)
	return rr, nil
}

// Get an author in the SampleAuthorsMap by a field name and value
func GetSampleAuthorByFieldValue(authors SampleAuthorsMap, fieldName string, value any) *SampleAuthor {
	for _, author := range authors {
		v := reflect.ValueOf(author)
		f := v.FieldByName(fieldName)
		if f.Interface() == value {
			return &author
		}
	}
	return nil
}

// Scan an author from sql.Rows (plural) into a SampleAuthor struct
func ScanAuthorFromRows(rows *sql.Rows) (SampleAuthor, error) {
	a := SampleAuthor{}
	err := rows.Scan(&a.ID, &a.Surname, &a.Forename)
	return a, err
}

// Scan an author from sql.Row (singular) into a SampleAuthor struct
func ScanAuthor(row *sql.Row) (SampleAuthor, error) {
	a := SampleAuthor{}
	err := row.Scan(&a.ID, &a.Surname, &a.Forename)
	return a, err
}

// Convert a RowData of author values into a SampleAuthor struct
func AuthorRowDataToStruct(rd types.RowData) SampleAuthor {
	return SampleAuthor{
		ID:       rd["id"].(int64),
		Surname:  rd["surname"].(string),
		Forename: rd["forename"].(string),
	}
}

// Convert a RowData of book values into a SampleBook struct
func BookRowDataToStruct(rd types.RowData) SampleBook {
	return SampleBook{
		ID:       rd["id"].(int64),
		Title:    rd["title"].(string),
		AuthorID: rd["author_id"].(int64),
	}
}
