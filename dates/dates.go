package dates

import (
	"fmt"
	"time"
)

type TimeGroup struct {
	Dates            []time.Time
	DatesFilterType  string
	DatesFilterFrame string
	TimeRange        []string
	DaysOfWeek       []int
	DueDateEndTime   string
}

func IsDateInTimeGroup(date time.Time, timeGroup TimeGroup) (bool, error) {
	if date.IsZero() {
		return false, fmt.Errorf("IsDateInTimeGroup: date must have value")
	}

	//test time
	if len(timeGroup.TimeRange) > 1 {
		fromTime := timeGroup.TimeRange[0] + ":00"
		toTime := timeGroup.TimeRange[1] + ":00"
		curTime := date.Format("15:04:05")
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
		fromDate := timeGroup.Dates[0].Format("2006-01-02")
		toDate := timeGroup.Dates[1].Format("2006-01-02")
		dateF := dateFilter.Format("2006-01-02")
		ok, err := OnlyDateWithinRange(fromDate, toDate, dateF)
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
			if dateFilter.Weekday() == time.Weekday(w) {
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
	curTime := date.Format("15:04:05")
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

func OnlyDateWithinRange(fromDate, toDate, curDate string) (bool, error) {
	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return false, err
	}
	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return false, err
	}
	date, err := time.Parse("2006-01-02", curDate)
	if err != nil {
		return false, err
	}
	return (date.After(from) && date.Before(to)) || (from == date || to == date), nil
}

//copy from helpers
func OnlyTimeWithinRange(fromTime, toTime, curTime string) (bool, error) {
	fromTimeParsed, err := time.Parse("15:04:05", fromTime)
	if err != nil {
		return false, err
	}
	toTimeParsed, err := time.Parse("15:04:05", toTime)
	if err != nil {
		return false, err
	}
	curTimeParsed, err := time.Parse("15:04:05", curTime)
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
	return (date.After(from) && date.Before(to)) || (from == date || to == date)
}

//copy from helpers
func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
