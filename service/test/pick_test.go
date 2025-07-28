package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
)

func Test_PickRow(t *testing.T) {
	serv, sampleRows := tests.NewTestService(t)

	t.Run("pick with valid ids", func(t *testing.T) {
		for id, sample := range sampleRows {
			rowDataMap, err := serv.PickRow(tests.TABLE1, fmt.Sprintf("%d", id))
			if err != nil {
				t.Errorf("pick err: %s", err)
			}
			for k, v := range sample {
				if v != rowDataMap[k] {
					t.Errorf("Expected %s: %v\nGot: %v", k, v, rowDataMap[k])
				}
			}
		}
	})
}
