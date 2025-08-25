package api_test

import (
	"net/http"
	"slices"
	"testing"

	"gopgrest/assert"
	"gopgrest/tests"
)

func Test_DELETE_NoConditions(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)
	rr, err := tests.MakeHttpRequest(ah, http.MethodDelete, "/authors", nil)
	assert.Try(t, err)
	// Do NOT want 200 (should be 400 but currently getting 500, need error
	// unwrapping?)
	assert.IsNotEq(t, rr.Code, http.StatusOK)
}

func Test_DELETE_SingleCondition(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)

	// Get all authors before deletions
	allAuthors, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	assert.Try(t, err)

	rr, err := tests.MakeHttpRequest(ah, http.MethodDelete, "/authors?surname==Brontë", nil)
	assert.Try(t, err)
	assert.IsEq(t, rr.Code, http.StatusOK)

	ids := tests.ParseIDArrayResponse(t, rr.Body.String())

	// Track author IDs we expect to find
	keptIDs := []int64{}
	// Confirm deleted IDs are ones we expect
	for _, row := range allAuthors {
		thisID := row["id"].(int64)
		surname := row["surname"]
		if surname == "Brontë" {
			assert.IsTrue(t, slices.Contains(ids, thisID))
		} else {
			assert.IsTrue(t, !slices.Contains(ids, thisID))
			keptIDs = append(keptIDs, thisID)
		}
	}

	currentAuthors, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	assert.Try(t, err)
	// Confirm that current authors don't include deleted ones
	for _, row := range currentAuthors {
		thisID := row["id"].(int64)
		assert.IsTrue(t, slices.Contains(keptIDs, thisID))
		surname := row["surname"]
		assert.IsNotEq(t, surname, "Brontë")
	}
}

func Test_DELETE_MultipleConditions(t *testing.T) {
	ah := tests.NewTestAPIHandler(t)

	// Get all authors before deletions
	allAuthors, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	assert.Try(t, err)

	rr, err := tests.MakeHttpRequest(ah, http.MethodDelete, "/authors?forename==Anne;surname==Carson", nil)
	assert.Try(t, err)
	assert.IsEq(t, rr.Code, http.StatusOK)

	ids := tests.ParseIDArrayResponse(t, rr.Body.String())

	// Track author IDs we expect to find
	keptIDs := []int64{}
	// Confirm deleted IDs are ones we expect
	for _, row := range allAuthors {
		thisID := row["id"].(int64)
		forename := row["forename"]
		surname := row["surname"]
		if surname == "Carson" && forename == "Anne" {
			assert.IsTrue(t, slices.Contains(ids, thisID))
		} else {
			assert.IsTrue(t, !slices.Contains(ids, thisID))
			keptIDs = append(keptIDs, thisID)
		}
	}

	currentAuthors, err := tests.SelectRows(ah.Repo, "SELECT * FROM authors")
	assert.Try(t, err)
	// Confirm that current authors don't include deleted ones
	for _, row := range currentAuthors {
		t.Logf("🪚 row: %v\n", row)
		thisID := row["id"].(int64)
		assert.IsTrue(t, slices.Contains(keptIDs, thisID))
		surname := row["surname"]
		forename := row["forename"]
		assert.IsTrue(t, !(surname == "Carson" && forename == "Anne"))
	}
}
