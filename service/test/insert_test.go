package service_test

import (
	"fmt"
	"testing"

	"ftrack/test_utils"
)

func TestService_InsertRows(t *testing.T) {
	for _, tt := range test_utils.GetInsertTests() {
		t.Run(tt.Name, func(t *testing.T) {
			serv, _ := test_utils.NewTestService(t)
			_, err := serv.InsertRow(&tt.NewRow, test_utils.TABLE1)
			if fmt.Sprintf("%s", tt.CustomErr) != fmt.Sprintf("%s", err) {
				t.Errorf("\nExp: %v\nGot: %v", tt.CustomErr, err)
			}
		})
	}
}
