package api_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"gopgrest/test_utils"
	"gopgrest/types"
)

func TestAPI_Update(t *testing.T) {
	ah, _ := test_utils.NewTestAPIHandler(t)

	sampleRow := types.RowData{
		"name":   "romanian deadlift",
		"weight": 309,
	}

	insertResult := ah.Repo.InsertRow(test_utils.TABLE1, &sampleRow)
	if insertResult.Error != nil {
		t.Errorf("Insert err %s", insertResult.Error)
	}

	tagMap := test_utils.GetTagMap(test_utils.ExerciseSet{})

	updateTests := test_utils.GetUpdateTests(insertResult)

	for _, tt := range updateTests {
		t.Run(tt.TestName, func(t *testing.T) {
			updateData := types.RowData{tt.Col: tt.Value}

			// Make request
			path := fmt.Sprintf("/%s/%d", test_utils.TABLE1, insertResult.ID)
			rr, err := test_utils.MakeHttpRequest(ah, http.MethodPut, path, updateData)
			if err != nil {
				t.Error(err.Error())
			}

			// Go to next test if this is an invalid update query
			if rr.Code == http.StatusInternalServerError {
				// Confirm expected error message
				if rr.Body.String() != tt.CustomErr {
					t.Errorf("\nExp: %s\nGot: %s", tt.CustomErr, rr.Body)
				} else {
					return
				}
			}

			// Get updated row
			updatedRow := test_utils.ExerciseSet{}
			err = test_utils.ScanExerciseSetRow(
				&updatedRow,
				ah.Repo.GetRowByID(
					test_utils.TABLE1,
					insertResult.ID,
				),
			)
			if err != nil {
				t.Errorf("Scan err: %s", err)
			}

			// Confirm update
			rowVal := reflect.ValueOf(updatedRow)
			fieldName := test_utils.GetFieldNameByColName(tagMap, tt.Col, test_utils.ExerciseSet{})
			gotVal := rowVal.FieldByName(fieldName).Interface()
			if tt.Value != gotVal {
				t.Errorf("Expected %v: %s\nGot %v", tt.Col, tt.Value, gotVal)
			}
		})
	}
}
