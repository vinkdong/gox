package vtime

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParserTimestampMs(timestampMs int64) time.Time {
	return time.Unix(0, timestampMs*10e5)
}

func ParserTimestampNs(timestampNs int64) time.Time {
	return time.Unix(0, timestampNs)
}

func ParserTimestampS(timestampS int64) time.Time {
	return time.Unix(timestampS, 0)
}

var relativeRE = regexp.MustCompile(`^(now)?([+-]?\d+)([smhd])( ago)?$`)
var timeNow = time.Now

func ParseFlexibleTime(input string) (time.Time, error) {
	input = strings.TrimSpace(input)

	// 1. 相对时间：支持 "1h ago", "30m ago", "now+1h", "now-2d", "+5m", "-10s" 等
	if match := relativeRE.FindStringSubmatch(input); match != nil {
		var num int
		var err error
		numStr := match[2]
		if numStr == "" {
			num = 0
		} else {
			num, err = strconv.Atoi(numStr)
			if err != nil {
				return time.Time{}, err
			}
		}
		unit := match[3]

		var duration time.Duration
		switch unit {
		case "s":
			duration = time.Second
		case "m":
			duration = time.Minute
		case "h":
			duration = time.Hour
		case "d":
			duration = time.Hour * 24
		default:
			return time.Time{}, fmt.Errorf("unsupported unit: %s", unit)
		}

		// 处理方向：有 "ago" 或负号表示过去，否则是未来
		if strings.Contains(input, "ago") || (numStr != "" && numStr[0] == '-') {
			num = -int(math.Abs(float64(num))) // 确保是负数
		}

		return timeNow().Add(time.Duration(num) * duration), nil
	}

	// 2. Unix 时间戳（秒 或 毫秒）
	if ts, err := strconv.ParseInt(input, 10, 64); err == nil {
		if ts > 1000000000000 { // 毫秒
			return time.UnixMilli(ts), nil
		} else if ts > 1000000000 { // 秒
			return time.Unix(ts, 0), nil
		}
	}

	// 3. 标准时间格式
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
		"2006/01/02 15:04:05",
		"01/02/2006 15:04:05",
		"2006-01-02 15:04",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", input)
}
