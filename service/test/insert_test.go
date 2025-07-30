package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
)

func TestService_InsertRows(t *testing.T) {
	for _, tt := range tests.GetInsertTests() {
		t.Run(tt.Name, func(t *testing.T) {
			serv, _ := tests.NewTestService(t)
			_, err := serv.InsertRow(&tt.NewRow, tests.TABLE1)
			if fmt.Sprintf("%s", tt.PqErr) != fmt.Sprintf("%s", err) {
				t.Errorf("Got: %v", err)
			}
		})
	}
}
