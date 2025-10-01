package utils

import (
	"fmt"
	"time"
)

func ConvertTimeMilli(ms int) string {
	t := time.UnixMilli(int64(ms))
	return t.UTC().Format("2006-01-02 15:04:05")
}

func ConvertTimeSeconds(seconds int) string {
	t := time.Unix(int64(seconds), 0)
	return t.UTC().Format("2006-01-02 15:04:05")
}

func ConvertTimeFormat(ms int, format string) string {
	t := time.UnixMilli(int64(ms))
	return t.UTC().Format(format)
}

func TimeDeflect(start, tf string, deflect int) (string, int, error) {
	dt, err := time.Parse(time.DateTime, start)
	if err != nil {
		return "", 0, err
	}

	timestamp := AddTime(tf, dt, time.Duration(-deflect))

	return time.Unix(timestamp/1000, 0).UTC().Format(time.DateTime), int(dt.UnixMilli()), nil
}

func TimeSince(start time.Time) string {
	duration := time.Since(start)

	if duration < time.Second {
		return fmt.Sprintf("%.1fms\n", float64(duration.Microseconds())/1000)
	}

	return fmt.Sprintf("%.1fs\n", duration.Seconds())
}

func TimeNowFormat() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}
