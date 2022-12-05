package dates

import (
	"fmt"
	"time"
)

const DateFormat = "2006-01-02"
const TimeFormat = "15:04:05"
const DateTimeFormat = "2006-01-02 15:04:05"
const DateTimeMsFormat = "2006-01-02 15:04:05.000"

type TimeGroup struct {
	Dates            []string
	DatesFilterType  string
	DatesFilterFrame string
	TimeRange        []string
	DaysOfWeek       []string
	DueDateEndTime   string
}

// supported timeOfDay = hh:mm  || hh:mm:ss
func SetTimeOfDay(date time.Time, timeOfDay string) (time.Time, error) {
	if len(timeOfDay) == 5 {
		timeOfDay = timeOfDay + ":00"
	} else if len(timeOfDay) != 7 {
		return time.Time{}, fmt.Errorf("SetTimeOfDay: unsupported format: %s,supported format hh:mm || hh:mm:ss", timeOfDay)
	}
	datePlusTime, err := time.Parse(DateTimeFormat, date.Format(DateFormat)+" "+timeOfDay)
	if err != nil {
		return time.Time{}, err
	}
	return datePlusTime, nil

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
		if len(timeGroup.Dates) == 1 {
			return false, fmt.Errorf("IsDateInTimeGroup: dates must have 2 values for DatesFilterType = range")
		}
		fromDate := timeGroup.Dates[0]
		toDate := timeGroup.Dates[1]
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
			dt, err := time.Parse(DateFormat, d)
			if err != nil {
				return false, err
			}
			if DateEqual(dt, dateFilter) {
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

func DayAndMonthEqual(date1, date2 time.Time) bool {
	_, m1, d1 := date1.Date()
	_, m2, d2 := date2.Date()
	return m1 == m2 && d1 == d2
}

func MonthEqual(date1, date2 time.Time) bool {
	_, m1, _ := date1.Date()
	_, m2, _ := date2.Date()
	return m1 == m2
}

func MinDate(d1 time.Time, d2 time.Time) time.Time {
	if d1.Before(d2) {
		return d1
	} else {
		return d2
	}
}

func MaxDate(d1 time.Time, d2 time.Time) time.Time {
	if d1.After(d2) {
		return d1
	} else {
		return d2
	}
}

func SetSameLocale(time1 time.Time, time2 *time.Time) {
	if time1.Location() != time2.Location() {
		*time2 = time.Date(time2.Year(), time2.Month(), time2.Day(), time2.Hour(), time2.Minute(), time2.Second(), time2.Nanosecond(), time1.Location())
	}
}
