package service_test

import (
	"fmt"
	"testing"

	"ftrack/tests"
)

func TestService_ListRows_InvalidFilters(t *testing.T) {
	for _, tt := range tests.GetInvalidQueryTests() {
		t.Run(tt.TestName, func(t *testing.T) {
			serv, _ := tests.NewTestService(t)
			_, err := serv.ListRows(tests.TABLE1, tt.Filters)
			if tests.CheckExpectedErr(tt.CustomErr, err) {
				t.Errorf("\nExp: %v\nGot: %v", tt.CustomErr, err)
			}
		})
	}
}

func TestService_ListRows_ValidFilters(t *testing.T) {
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
						t.Errorf("\nExp %s: %v %T\nGot: %v %T", k, v, v, colVal, colVal)
					}
				}
			}
		})
	}
}
