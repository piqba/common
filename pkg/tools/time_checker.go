package tools

import "time"

const (
	LayoutYyyyMmDd = "2006-01-02"
	TimezoneHavana = "America/Havana"
)

// NewTimeChecker ...
func NewTimeChecker(layout, tz, dateStart, endDate, dateCheck string) *TimeChecker {
	startToParse, _ := time.Parse(layout, dateStart)
	checkToParse, _ := time.Parse(layout, dateCheck)
	beginningOfMonth, _ := BeginningOfMonth(startToParse, tz)
	endOfMonth, _ := ConvertTimeADayBefore(endDate, layout, tz)

	return &TimeChecker{
		Start:  beginningOfMonth,
		End:    endOfMonth,
		Check:  checkToParse,
		Tz:     tz,
		Layout: layout,
	}
}

// InTime ...
func (c *TimeChecker) InTime(timeToCheck time.Time) bool {
	return timeToCheck.After(c.Start) && timeToCheck.Before(c.End)
}

// ConvertTimeADayBefore ...
func ConvertTimeADayBefore(dateEnd, layout, zone string) (time.Time, error) {
	var tz, err = time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, err
	}
	end, err := time.Parse(layout, dateEnd)
	if err != nil {
		return time.Time{}, err
	}
	return end.In(tz).Add(-time.Second), nil
}

// BeginningOfMonth ...
func BeginningOfMonth(t time.Time, zone string) (time.Time, error) {
	var tz, err = time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, tz), nil
}

// EndOfMonth ...
func EndOfMonth(end time.Time, zone string) (time.Time, error) {
	var tz, err = time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, err
	}
	beginningOfMonth, err := BeginningOfMonth(end, zone)
	if err != nil {
		return time.Time{}, err
	}
	return beginningOfMonth.In(tz).AddDate(0, 1, 0).Add(-time.Second), nil
}
