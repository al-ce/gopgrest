package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
)

func Test_PickRow(t *testing.T) {
	serv, sampleRows := tests.NewTestService(t)

	tagMap := tests.GetTagMap(tests.ExerciseSet{})

	t.Run("pick with valid ids", func(t *testing.T) {
		for id, sample := range sampleRows {
			rowDataMap, err := serv.PickRow(tests.TABLE1, fmt.Sprintf("%d", id))
			if err != nil {
				t.Errorf("pick err: %s", err)
			}
			for k, v := range sample {
				colName := tagMap[k]
				colVal := rowDataMap[colName]
				if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", colVal) {
					t.Errorf("Expected %s: %v %T\nGot: %v %T", k, v, v, colVal, colVal)
				}
			}
		}
	})
}
