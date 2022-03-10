package dates

import (
	"fmt"
	"time"
)

const DateFormat = "2006-01-02"
const TimeFormat = "15:04:05"
const DateTimeFormat = "2006-01-02 15:04:05"

type TimeGroup struct {
	Dates            []time.Time
	DatesFilterType  string
	DatesFilterFrame string
	TimeRange        []string
	DaysOfWeek       []string
	DueDateEndTime   string
}

func IsDateInTimeGroup(date time.Time, timeGroup TimeGroup) (bool, error) {
	if date.IsZero() {
		return false, fmt.Errorf("IsDateInTimeGroup: date must have value")
	}

	//test time
	if len(timeGroup.TimeRange) > 1 {
		fromTime := timeGroup.TimeRange[0] + ":00"
		toTime := timeGroup.TimeRange[1] + ":59"
		curTime := date.Format(TimeFormat)
		ok, err := OnlyTimeWithinRange(fromTime, toTime, curTime)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	dateFilter := date
	if timeGroup.DatesFilterFrame == "dueDate" {
		if timeGroup.DueDateEndTime == "" {
			return false, fmt.Errorf("IsDateInTimeGroup: dueDateEndTime must have value for DatesFilterFrame = dueDate")
		}
		shiftDate, err := getShiftDate(date, timeGroup.DueDateEndTime)
		if err != nil {
			return false, err
		}
		dateFilter = shiftDate
	}

	//test date range
	if timeGroup.DatesFilterType == "range" && len(timeGroup.Dates) > 1 {

		fromDate := timeGroup.Dates[0].Format(DateFormat)
		toDate := timeGroup.Dates[1].Format(DateFormat)
		dateF := dateFilter.Format(DateFormat)
		ok, err := OnlyDateWithinRange(dateF, fromDate, toDate)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	//test date list
	if timeGroup.DatesFilterType == "list" && len(timeGroup.Dates) > 0 {
		found := false
		for _, d := range timeGroup.Dates {
			if DateEqual(d, dateFilter) {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	//test day of week
	if len(timeGroup.DaysOfWeek) > 0 {
		found := false
		for _, w := range timeGroup.DaysOfWeek {
			if dateFilter.Weekday().String() == w {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	return true, nil
}

func getShiftDate(date time.Time, dueDateEndTime string) (time.Time, error) {
	curTime := date.Format(TimeFormat)
	ok, err := OnlyTimeWithinRange("00:00:00", dueDateEndTime, curTime)
	if err != nil {
		return time.Time{}, err
	}
	if ok {
		return date.Add(-24 * time.Hour), nil
	} else {
		return date, nil
	}
}

func OnlyDateWithinRange(curDate, fromDate, toDate string) (bool, error) {
	from, err := time.Parse(DateFormat, fromDate)
	if err != nil {
		return false, err
	}
	to, err := time.Parse(DateFormat, toDate)
	if err != nil {
		return false, err
	}
	date, err := time.Parse(DateFormat, curDate)
	if err != nil {
		return false, err
	}
	return TimeWithinRange(from, to, date), nil
}

//copy from helpers
func OnlyTimeWithinRange(fromTime, toTime, curTime string) (bool, error) {
	fromTimeParsed, err := time.Parse(TimeFormat, fromTime)
	if err != nil {
		return false, err
	}
	toTimeParsed, err := time.Parse(TimeFormat, toTime)
	if err != nil {
		return false, err
	}
	curTimeParsed, err := time.Parse(TimeFormat, curTime)
	if err != nil {
		return false, err
	}
	if fromTimeParsed.After(toTimeParsed) {
		toTimeParsed = toTimeParsed.Add(24 * time.Hour)
	}
	if curTimeParsed.Before(fromTimeParsed) {
		curTimeParsed = curTimeParsed.Add(24 * time.Hour)
	}
	return TimeWithinRange(fromTimeParsed, toTimeParsed, curTimeParsed), nil
}

//copy from helpers
func TimeWithinRange(from, to, date time.Time) bool {
	return (from == date || to == date) || (date.After(from) && date.Before(to))
}

//copy from helpers
func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
