package api_test

import (
	"net/http"
	"slices"
	"testing"

	"gopgrest/assert"
	"gopgrest/tests"
	"gopgrest/types"
)

func Test_PUT_NoConditions(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)
	updateData := types.RowData{"forename": "Emily"}
	rr, err := tests.MakeHttpRequest(ah, http.MethodPut, "/authors", updateData)
	assert.Try(t, err)
	// Do NOT want 200 (should be 400 but currently getting 500, need error
	// unwrapping?)
	assert.IsNotEq(t, rr.Code, http.StatusOK)
}

func Test_PUT_SingleCondition(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)
	updateData := types.RowData{"forename": "Emily"}
	rr, err := tests.MakeHttpRequest(ah, http.MethodPut, "/authors?surname==Brontë", updateData)
	assert.Try(t, err)
	assert.IsEq(t, rr.Code, http.StatusOK)

	ids := tests.ParseIDArrayResponse(t, rr.Body.String())

	// Confirm only row we wanted to update was updated
	selectedRows, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	assert.Try(t, err)

	for _, updatedRow := range selectedRows {

		forename := updatedRow["forename"]
		surname := updatedRow["surname"]
		thisID := updatedRow["id"].(int64)

		if slices.Contains(ids, thisID) {
			assert.IsEq(t, forename, "Emily")
			assert.IsEq(t, surname, "Brontë")
		} else {
			assert.IsNotEq(t, forename, "Emily")
			assert.IsNotEq(t, surname, "Brontë")
		}
	}
}

func Test_PUT_MultipleConditions(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)
	updateData := types.RowData{"forename": "Rachel"}
	rr, err := tests.MakeHttpRequest(ah, http.MethodPut, "/authors?forename==Anne;surname==Carson", updateData)
	assert.Try(t, err)
	assert.IsEq(t, rr.Code, http.StatusOK)

	ids := tests.ParseIDArrayResponse(t, rr.Body.String())

	// Confirm only row we wanted to update was updated
	selectedRows, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	assert.Try(t, err)

	for _, updatedRow := range selectedRows {

		forename := updatedRow["forename"]
		surname := updatedRow["surname"]
		thisID := updatedRow["id"].(int64)

		if slices.Contains(ids, thisID) {
			assert.IsEq(t, forename, "Rachel")
			assert.IsEq(t, surname, "Carson")
		} else {
			assert.IsNotEq(t, surname, "Carson")
		}
	}
}
