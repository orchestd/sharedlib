package dates

import (
	"testing"
	"time"
)
import . "github.com/smartystreets/goconvey/convey"

func TestGetShiftDate(t *testing.T) {
	Convey("get shift same date", t, func() {
		t, _ := time.Parse(DateTimeFormat, "2022-03-01 20:00:00")
		sd, err := getShiftDate(t, "05:00:00")
		So(err, ShouldBeNil)
		So(sd.Format(DateFormat), ShouldEqual, "2022-03-01")
	})
	Convey("get shift after midnight", t, func() {
		t, _ := time.Parse(DateTimeFormat, "2022-03-02 01:00:00")
		sd, err := getShiftDate(t, "05:00:00")
		So(err, ShouldBeNil)
		So(sd.Format(DateFormat), ShouldEqual, "2022-03-01")
	})
	Convey("get shift next day", t, func() {
		t, _ := time.Parse(DateTimeFormat, "2022-03-02 10:00:00")
		sd, err := getShiftDate(t, "05:00:00")
		So(err, ShouldBeNil)
		So(sd.Format(DateFormat), ShouldEqual, "2022-03-02")
	})
}

func TestOnlyDateWithinRange(t *testing.T) {
	Convey("within range", t, func() {
		sd, err := OnlyDateWithinRange("2022-03-04", "2022-03-03", "2022-03-10", )
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("outside range", t, func() {
		sd, err := OnlyDateWithinRange("2022-03-01", "2022-03-03", "2022-03-10")
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})
	Convey("within range start", t, func() {
		sd, err := OnlyDateWithinRange("2022-03-03", "2022-03-03", "2022-03-10", )
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("within range end", t, func() {
		sd, err := OnlyDateWithinRange("2022-03-10", "2022-03-03", "2022-03-10", )
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
}

func TestIsDateInTimeGroup(t *testing.T) {
	//hours
	Convey("test only hours true", t, func() {
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "calendar",
			TimeRange:        []string{"15:00", "18:00"},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-03 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("test only hours false", t, func() {
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "calendar",
			TimeRange:        []string{"16:00", "18:00"},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-03 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})

	//day of week calendar
	Convey("test only day of week by calendar date true", t, func() {
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "calendar",
			DaysOfWeek:       []string{"Thursday"},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-03 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("test only day of week by calendar date false", t, func() {
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "calendar",
			DaysOfWeek:       []string{"Wednesday"},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-03 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})

	//day of week dueDate
	Convey("test only day of week by dueDate date true", t, func() {
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "dueDate",
			DaysOfWeek:       []string{"Thursday"},
			DueDateEndTime:   "05:00:00",
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-04 01:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("test only day of week by dueDate date false", t, func() {
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "dueDate",
			DaysOfWeek:       []string{"Thursday"},
			DueDateEndTime:   "05:00:00",
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-04 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})

	//date range calendar
	Convey("test only date range true", t, func() {
		r1, _ := time.Parse(DateFormat, "2022-03-01")
		r2, _ := time.Parse(DateFormat, "2022-03-05")
		timeGroup := TimeGroup{
			DatesFilterType:  "range",
			DatesFilterFrame: "calendar",
			Dates:            []time.Time{r1, r2},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-03 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("test only date range false", t, func() {
		r1, _ := time.Parse(DateTimeFormat, "2022-03-01 15:15:15")
		r2, _ := time.Parse(DateTimeFormat, "2022-03-05 15:15:15")
		timeGroup := TimeGroup{
			DatesFilterType:  "range",
			DatesFilterFrame: "calendar",
			Dates:            []time.Time{r1, r2},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-06 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})

	//date range dueDate
	Convey("test only date range true", t, func() {
		r1, _ := time.Parse(DateFormat, "2022-03-01")
		r2, _ := time.Parse(DateFormat, "2022-03-05")
		timeGroup := TimeGroup{
			DatesFilterType:  "range",
			DatesFilterFrame: "dueDate",
			Dates:            []time.Time{r1, r2},
			DueDateEndTime:   "05:00:00",
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-03 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("test only date range false", t, func() {
		r1, _ := time.Parse(DateFormat, "2022-03-01")
		r2, _ := time.Parse(DateFormat, "2022-03-05")
		timeGroup := TimeGroup{
			DatesFilterType:  "range",
			DatesFilterFrame: "dueDate",
			Dates:            []time.Time{r1, r2},
			DueDateEndTime:   "05:00:00",
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-06 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})

	//date list calendar
	Convey("test only date list true", t, func() {
		r1, _ := time.Parse(DateFormat, "2022-03-01")
		r2, _ := time.Parse(DateFormat, "2022-03-05")
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "calendar",
			Dates:            []time.Time{r1, r2},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-01 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("test only date list false", t, func() {
		r1, _ := time.Parse(DateFormat, "2022-03-01")
		r2, _ := time.Parse(DateFormat, "2022-03-05")
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "calendar",
			Dates:            []time.Time{r1, r2},
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-03 15:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})

	//date list dueDate
	Convey("test only date list true", t, func() {
		r1, _ := time.Parse(DateFormat, "2022-03-01")
		r2, _ := time.Parse(DateFormat, "2022-03-05")
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "dueDate",
			Dates:            []time.Time{r1, r2},
			DueDateEndTime:   "05:00:00",
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-02 01:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeTrue)
	})
	Convey("test only date list false", t, func() {
		r1, _ := time.Parse(DateFormat, "2022-03-01")
		r2, _ := time.Parse(DateFormat, "2022-03-05")
		timeGroup := TimeGroup{
			DatesFilterType:  "list",
			DatesFilterFrame: "dueDate",
			Dates:            []time.Time{r1, r2},
			DueDateEndTime:   "05:00:00",
		}
		d, _ := time.Parse(DateTimeFormat, "2022-03-02 11:15:15")
		sd, err := IsDateInTimeGroup(d, timeGroup)
		So(err, ShouldBeNil)
		So(sd, ShouldBeFalse)
	})
}
