package utils

import (
	"github.com/AlexanderKolesnkov/golang-utils-stuff/consts"
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

func TimeNowFormat() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

func CalcTimeframesCount(timeframe, start, end string) (int, error) {
	startTime, err := time.Parse(consts.TimeLayout, start)
	if err != nil {
		return 0, err
	}
	endTime, err := time.Parse(consts.TimeLayout, end)
	if err != nil {
		return 0, err
	}

	duration := endTime.Sub(startTime)

	switch timeframe {
	case "1":
		return int(duration.Minutes()), nil
	case "5":
		return int(duration.Minutes()) / 5, nil
	case "15":
		return int(duration.Minutes()) / 15, nil
	case "30":
		return int(duration.Minutes()) / 30, nil
	case "60":
		return int(duration.Hours()), nil
	case "240":
		return int(duration.Hours()) / 4, nil
	case "D":
		return int(duration.Hours()) / 24, nil
	}

	return -1, nil
}
