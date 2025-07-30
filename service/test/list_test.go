package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
)

func Test_Service_ListRows(t *testing.T) {
	serv, sampleRows := tests.NewTestService(t)
	filterTests := tests.GetValidFilterTests(sampleRows)

	for _, tt := range filterTests {
		t.Run(tt.TestName, func(t *testing.T) {
			rowDataMap, err := serv.ListRows(tests.TABLE1, tt.Filters)
			if err != nil {
				t.Errorf("ListRows err: %v\n", err)
			}
			filteredSampleRows := tests.FilterSampleRows(tt.Filters, sampleRows)
			for idx, rdm := range rowDataMap {
				sample := filteredSampleRows[idx]
				for k, v := range sample {
					colVal := rdm[k]
					if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", colVal) {
						t.Errorf("Expected %s: %v %T\nGot: %v %T", k, v, v, colVal, colVal)
					}
				}
			}
		})
	}
}
