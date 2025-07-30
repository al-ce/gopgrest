package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
	"ftrack/types"
)

func Test_ListRows(t *testing.T) {
	serv, sampleRows := tests.NewTestService(t)

	tagMap := tests.GetTagMap(tests.ExerciseSet{})

	t.Run("pick with no filter", func(t *testing.T) {
		rowDataMapSlice, err := serv.ListRows(tests.TABLE1, types.QueryFilters{})
		if err != nil {
			t.Errorf("ListRows err: %v\n", err)
		}
		for idx, rdm := range rowDataMapSlice {
			// index with offset of 1 should match id based on insert order
			sample := sampleRows[int64(idx+1)]
			for k, v := range sample {
				colName := tagMap[k]
				colVal := rdm[colName]
				if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", colVal) {
					t.Errorf("Expected %s: %v %T\nGot: %v %T", k, v, v, colVal, colVal)
				}
			}
		}
	})
}
