package test_utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"

	"gopgrest/api"
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
