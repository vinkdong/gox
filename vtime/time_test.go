package vtime

import (
	"fmt"
	"testing"
	"time"
)

const (
	TestTimeMS     = 1514461867000
	TestTimeNS     = 1514461867000000000
	TestTimeS      = 1514461867
	TestTimeYear   = 2017
	TestTimeSecond = 7
	TestTimeMinute = 51
	TestTimeString = "2017-12-28 19:51:07"
	TestTimeFormat = "2006-01-02 15:04:05"
	TestTimeTZ     = "Asia/Shanghai"
)

func checkTimeIn(t *testing.T, time time.Time) {
	if time.UTC().Year() != TestTimeYear {
		t.Errorf("Expect Time year is %d, But got %d", TestTimeYear, time.UTC().Year())
	}
	if time.Second() != TestTimeSecond {
		t.Errorf("Expect Time second is %d, But got %d", TestTimeSecond, time.UTC().Second())
	}
	if time.Minute() != TestTimeMinute {
		t.Errorf("Expect Time minute is %d, But got %d", TestTimeMinute, time.UTC().Minute())
	}
	if time.UnixNano() != TestTimeNS {
		t.Errorf("Expect Timestamp is %d, but got %d", TestTimeNS, time.UnixNano())
	}
}

func TestParserTimestampMs(t *testing.T) {
	var tsMs int64
	tsMs = TestTimeMS
	time := ParserTimestampMs(tsMs)
	checkTimeIn(t, time)
}

func TestParserTimestampNs(t *testing.T) {
	var tsNs int64
	tsNs = TestTimeNS
	time := ParserTimestampNs(tsNs)
	checkTimeIn(t, time)
}

func TestParserTimestampS(t *testing.T) {
	var tsS int64
	tsS = TestTimeS
	time := ParserTimestampS(tsS)
	checkTimeIn(t, time)
}

func TestParserVTime(t *testing.T) {
	vt := Time{
		Format: "2006-01-02 15:04:05",
		Value:  "2017-12-28 19:51:07",
		TZ:     TestTimeTZ,
	}
	time, err := vt.Parser()
	if err != nil {
		t.Fail()
	}
	checkTimeIn(t, time)
}

func TestFromTime(t *testing.T) {
	vt := Time{
		Format: TestTimeFormat,
		TZ:     TestTimeTZ,
	}
	var tsS int64
	tsS = TestTimeS
	time := ParserTimestampS(tsS)
	vt.FromTime(time)
	if vt.Value != TestTimeString {
		t.Errorf("expect time value is %s but got %s", TestTimeString, vt.Value)
	}
}

func TestTimeTransfer(t *testing.T) {
	to := &Time{
		Format: TestTimeFormat,
		TZ:     TestTimeTZ,
	}
	from := &Time{
		Format: "timestamp",
		Value:  fmt.Sprintf("%d", TestTimeMS),
		Unit:   "ms",
	}
	from.Transfer(to)
	if to.Value != TestTimeString {
		t.Errorf("expect time value is %s but got %s", TestTimeString, to.Value)
	}
}

func TestTime_FromRelativeTime(t *testing.T) {
	tm := &Time{
		Format: TestTimeFormat,
		TZ:     TestTimeTZ,
	}
	tm.FromRelativeTime("now+1h")
	if tm.Time.Hour() != timeNow().UTC().Hour()+1 {
		t.Errorf("expect time hour is %d but got %d", timeNow().UTC().Hour()+1, tm.Time.Hour())
	}
}

func TestTime_ParseFlexibleTime(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2025, 12, 17, 14, 30, 0, 0, time.UTC)
	}
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		// 相对时间 - 过去（ago）
		{"1h ago", "1h ago", timeNow().Add(-1 * time.Hour), false},
		{"30m ago", "30m ago", timeNow().Add(-30 * time.Minute), false},
		{"2d ago", "2d ago", timeNow().Add(-48 * time.Hour), false},
		{"45s ago", "45s ago", timeNow().Add(-45 * time.Second), false},

		// 相对时间 - now + 未来 / 过去
		{"now+1h", "now+1h", timeNow().Add(1 * time.Hour), false},
		{"now+90m", "now+90m", timeNow().Add(90 * time.Minute), false},
		{"now-2h", "now-2h", timeNow().Add(-2 * time.Hour), false},
		{"now-5d", "now-5d", timeNow().Add(-5 * 24 * time.Hour), false},

		// 简写正负号（等价于 now+ / now-）
		{"+10m", "+10m", timeNow().Add(10 * time.Minute), false},
		{"-5s", "-5s", timeNow().Add(-5 * time.Second), false},
		{"+3d", "+3d", timeNow().Add(3 * 24 * time.Hour), false},

		// Unix 时间戳
		{"unix seconds", "1734432000", time.Unix(1734432000, 0).UTC(), false},
		{"unix milliseconds", "1734432000000", time.UnixMilli(1734432000000).UTC(), false},

		// 标准时间格式
		{"RFC3339", "2025-12-17T15:30:00Z", time.Date(2025, 12, 17, 15, 30, 0, 0, time.UTC), false},
		{"RFC3339 with offset", "2025-12-17T22:30:00+08:00", time.Date(2025, 12, 17, 14, 30, 0, 0, time.UTC), false},
		{"common format", "2025-12-17 14:30:00", time.Date(2025, 12, 17, 14, 30, 0, 0, time.UTC), false},
		{"date only", "2025-12-17", time.Date(2025, 12, 17, 0, 0, 0, 0, time.UTC), false},

		// 错误案例
		{"invalid format", "invalid time", time.Time{}, true},
		{"unsupported unit", "5w ago", time.Time{}, true},
		{"empty string", "", time.Time{}, true},
		{"random text", "hello world", time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlexibleTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlexibleTime(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && !got.Equal(tt.want) {
				t.Errorf("ParseFlexibleTime(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
