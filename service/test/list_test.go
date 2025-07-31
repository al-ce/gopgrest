package service_test

import (
	"fmt"
	"testing"

	"ftrack/test_utils"
)

func TestService_ListRows_InvalidFilters(t *testing.T) {
	for _, tt := range test_utils.GetInvalidQueryTests() {
		t.Run(tt.TestName, func(t *testing.T) {
			serv, _ := test_utils.NewTestService(t)
			_, err := serv.ListRows(test_utils.TABLE1, tt.Filters)
			if test_utils.CheckExpectedErr(tt.CustomErr, err) {
				t.Errorf("\nExp: %v\nGot: %v", tt.CustomErr, err)
			}
		})
	}
}

func TestService_ListRows_ValidFilters(t *testing.T) {
	serv, sampleRows := test_utils.NewTestService(t)
	filterTests := test_utils.GetValidFilterTests(sampleRows)

	for _, tt := range filterTests {
		t.Run(tt.TestName, func(t *testing.T) {
			rowDataMap, err := serv.ListRows(test_utils.TABLE1, tt.Filters)
			if err != nil {
				t.Errorf("ListRows err: %v\n", err)
			}
			filteredSampleRows := test_utils.FilterSampleRows(tt.Filters, sampleRows)
			for idx, rdm := range *rowDataMap {
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
