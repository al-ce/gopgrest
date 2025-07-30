package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
)

func TestService_PickRow(t *testing.T) {
	serv, sampleRows := tests.NewTestService(t)

	t.Run("pick with valid ids", func(t *testing.T) {
		for id, sample := range sampleRows {
			rowDataMap, err := serv.PickRow(tests.TABLE1, fmt.Sprintf("%d", id))
			if err != nil {
				t.Errorf("pick err: %s", err)
			}
			for k, v := range sample {
				colVal := rowDataMap[k]
				if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", colVal) {
					t.Errorf("Expected %s: %v %T\nGot: %v %T", k, v, v, colVal, colVal)
				}
			}
		}
	})
}
