package service_test

import (
	"fmt"
	"testing"

	"gopgrest/test_utils"
	"gopgrest/types"
)

func TestService_Update(t *testing.T) {
	serv, _ := test_utils.NewTestService(t)

	sampleRow := types.RowData{
		"name":   "romanian deadlift",
		"weight": 309,
	}

	insertResult := serv.Repo.InsertRow(test_utils.TABLE1, &sampleRow)
	if insertResult.Error != nil {
		t.Errorf("Insert err %s", insertResult.Error)
	}

	updateTests := test_utils.GetUpdateTests(insertResult)

	for _, tt := range updateTests {
		t.Run(tt.TestName, func(t *testing.T) {
			updateData := types.RowData{tt.Col: tt.Value}

			// Exec update query
			updateResult, err := serv.UpdateRow(
				test_utils.TABLE1,
				fmt.Sprintf("%d", tt.ID),
				&updateData,
			)
			if test_utils.CheckExpectedErr(tt.CustomErr, err) {
				t.Errorf("\nExp: %s\nGot: %s", tt.CustomErr, err)
			}
			// Go to next test if this is an invalid update query
			if err != nil {
				return
			}

			// Confirm updated field
			if updateResult[tt.Col] != tt.Value {
				t.Errorf(
					"\nExp %s: %s\nGot: %s",
					tt.Col,
					tt.Value,
					updateResult[tt.Col],
				)
			}
		})
	}
}
