package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	timeCheckerData = []struct {
		name          string
		startDate     string
		endDate       string
		checkDate     string
		tz            string
		layout        string
		errorExpected bool
		errorMessage  string
	}{
		{
			name:          "in time",
			startDate:     "2022-01-02",
			endDate:       "2022-02-01",
			checkDate:     "2022-01-28",
			tz:            TimezoneHavana,
			layout:        LayoutYyyyMmDd,
			errorExpected: false,
			errorMessage:  "error checking date",
		},
		{
			name:          "not in time",
			startDate:     "2022-01-02",
			endDate:       "2022-02-01",
			checkDate:     "2022-03-28",
			tz:            TimezoneHavana,
			layout:        LayoutYyyyMmDd,
			errorExpected: false,
			errorMessage:  "error checking date",
		},
		{
			name:          "invalid_fields_end_date",
			startDate:     "2022-01-02",
			endDate:       "aa-01-28",
			checkDate:     "2022-02-01",
			tz:            TimezoneHavana,
			layout:        LayoutYyyyMmDd,
			errorExpected: true,
			errorMessage:  "error invalid endDate",
		},
		{
			name:          "invalid_fields_tz",
			startDate:     "2022-01-02",
			endDate:       "2022-01-28",
			checkDate:     "2022-01-04",
			tz:            "TimezoneHavana",
			layout:        LayoutYyyyMmDd,
			errorExpected: true,
			errorMessage:  "error invalid tz value",
		},
	}
)

func TestNewTimeChecker(t *testing.T) {
	assertions := assert.New(t)
	for _, ch := range timeCheckerData {
		tch := NewTimeChecker(
			ch.layout, ch.tz, ch.startDate, ch.endDate, ch.checkDate,
		)
		if tch.InTime(tch.Check) {
			assertions.Equal("in time", ch.name)
		}
		if ch.errorExpected {
			t.Logf(
				"Not In time %s because has an error: %s",
				ch.name,
				ch.errorMessage,
			)
		}
		assertions.NotEqual("Not In time", ch.name)
	}

}
