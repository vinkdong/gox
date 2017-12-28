package vtime

import (
	"strconv"
	"time"
)

type Time struct {
	Format string
	Unit   string
	Value  string
	Time   time.Time
	TZ     string
}

func (t *Time) Transfer(to *Time) error {
	stdTime, err := t.Parser()
	if err != nil {
		return err
	}
	to.FromTime(stdTime)
	return nil
}

// Parser VTime Format to time
func (t *Time) Parser() (time.Time, error) {
	var err error
	if t.Format == "timestamp" {
		formTimeValue, err := strconv.ParseInt(t.Value, 0, 64)
		switch t.Unit {
		case "ms":
			t.Time = ParserTimestampMs(formTimeValue)
			break
		case "ns":
			t.Time = ParserTimestampNs(formTimeValue)
		case "s":
			t.Time = ParserTimestampS(formTimeValue)
		default:
			t.Time = time.Now()
		}
		return t.Time, err
	} else {
		t.Time, err = time.Parse(t.Format, t.Value)
	}
	return t.Time, err
}

func (t *Time) FromTime(stdTime time.Time) {
	if t.TZ != ""{
		loc, err := time.LoadLocation(t.TZ)
		if err == nil{
			stdTime.In(loc)
		}
	}
	if t.Format == "timestamp" {
		switch t.Unit {
		case "ms":
			ttime := stdTime.UnixNano()
			t.Value = strconv.FormatInt(ttime, 0)
			return
		case "ns":
			ttime := stdTime.UnixNano()
			t.Value = strconv.FormatInt(ttime, 0)
			return
		case "s":
			ttime := stdTime.Unix()
			t.Value = strconv.FormatInt(ttime, 0)
		default:
			return
		}
	} else {
		t.Value = stdTime.Format(t.Format)
	}
}
