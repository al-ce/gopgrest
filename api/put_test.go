package api_test

import (
	"net/http"
	"testing"

	"gopgrest/tests"
	"gopgrest/types"
)

func Test_PUT_NoConditions(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)
	updateData := types.RowData{"forename": "Emily"}
	rr, err := tests.MakeHttpRequest(ah, http.MethodPut, "/authors", updateData)
	tests.Try(t, err)
	// Do NOT want 200 (should be 400 but currently getting 500, need error
	// unwrapping?)
	if rr.Code == http.StatusOK {
		t.Fatalf("\nExp StatusCode: %d\nGot: %d\nResp: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
}

func Test_PUT_SingleCondition(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)
	updateData := types.RowData{"forename": "Emily"}
	rr, err := tests.MakeHttpRequest(ah, http.MethodPut, "/authors?surname==Brontë", updateData)
	tests.Try(t, err)
	if rr.Code != http.StatusOK {
		t.Fatalf("\nExp StatusCode: %d\nGot: %d\nResp: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	// Confirm only row we wanted to update was updated
	selectedRows, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	tests.Try(t, err)

	for _, updatedRow := range selectedRows {
		forename, ok := updatedRow["forename"]
		if !ok {
			t.Errorf("Could not get forename from %v", updatedRow)
		}
		surname, ok := updatedRow["surname"]
		if !ok {
			t.Errorf("Could not get surname from %v", updatedRow)
		}
		if surname == "Brontë" && forename != "Emily" {
			t.Errorf("Row did not update: %v", updatedRow)
		} else if surname != "Brontë" && forename == "Emily" {
			t.Errorf("Row should not have updated: %v", updatedRow)
		}
	}
}

func Test_PUT_MultipleConditions(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)
	updateData := types.RowData{"forename": "Rachel"}
	rr, err := tests.MakeHttpRequest(ah, http.MethodPut, "/authors?forename==Anne;surname==Carson", updateData)
	tests.Try(t, err)
	if http.StatusOK != rr.Code {
		t.Fatalf("\nExp StatusCode: %d\nGot: %d\nResp: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	// Confirm only row we wanted to update was updated
	selectedRows, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	tests.Try(t, err)
	for _, updatedRow := range selectedRows {
		forename, ok := updatedRow["forename"]
		if !ok {
			t.Errorf("Could not get forename from %v", updatedRow)
		}
		surname, ok := updatedRow["surname"]
		if !ok {
			t.Errorf("Could not get surname from %v", updatedRow)
		}
		if surname == "Carson" && forename != "Rachel" {
			t.Errorf("Row did not update: %v", updatedRow)
		} else if surname != "Carson" && forename == "Rachel" {
			t.Errorf("Row should not have updated: %v", updatedRow)
		}
	}
}
