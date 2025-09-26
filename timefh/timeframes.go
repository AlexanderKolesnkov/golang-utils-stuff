package timefh

import (
	"fmt"
	"github.com/AlexanderKolesnkov/golang-utils-stuff/consts"
	"time"
)

const (
	Minute = 60
	Hour   = Minute * 60
	Day    = Hour * 24
)

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

func IsUnixEnds(unix int, duration int, counter *int) (bool, error) {
	switch duration {
	case Day * 3:
		if unix%86400 == 0 {
			*counter++
			if *counter == 3 {
				*counter = 0
				return true, nil
			}
		}

		return false, nil
	case Day:
		if unix%86400 == 0 {
			return true, nil
		}

		return false, nil
	case Hour * 12:
		if unix%86400 == 0 || unix%86400 == 43200 {
			return true, nil
		}

		return false, nil

	case Hour * 6:
		if unix%3600 == 0 && (unix/3600)%6 == 0 {
			return true, nil
		}

		return false, nil

	case Hour * 5:
		if unix%3600 == 0 && (unix/3600)%5 == 0 {
			return true, nil
		}

		return false, nil
	case Hour * 4:
		if unix%3600 == 0 && (unix/3600)%4 == 0 {
			return true, nil
		}

		return false, nil
	case Hour * 3:
		if unix%3600 == 0 && (unix/3600)%3 == 0 {
			return true, nil
		}

		return false, nil
	case Hour * 2:
		if unix%3600 == 0 && (unix/3600)%2 == 0 {
			return true, nil
		}

		return false, nil
	case Hour:
		if unix%3600 == 0 {
			return true, nil
		}

		return false, nil
	case Minute * 30:
		if ((unix/60)%60)%30 == 0 {
			return true, nil
		}

		return false, nil
	case Minute * 15:
		if ((unix/60)%60)%15 == 0 {
			return true, nil
		}

		return false, nil
	case Minute * 5:
		if ((unix/60)%60)%5 == 0 {
			return true, nil
		}

		return false, nil
	case Minute:
		if ((unix/60)%60)%1 == 0 {
			return true, nil
		}

		return false, nil
	default:
		return false, fmt.Errorf("not allowed duration: %v", duration)
	}
}

func ConvertTimeframeInUnix(timeframe string) int {
	switch timeframe {
	case "1":
		return 60
	case "5":
		return 300
	case "15":
		return 900
	case "30":
		return 1800
	case "60":
		return 3600
	case "240":
		return 14400
	case "D":
		return 86400
	default:
		return -1
	}
}
