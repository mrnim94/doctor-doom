package helper

import (
	"fmt"
	"strconv"
)

// Constants for time unit conversions
const (
	MinutesInHour = 60
	HoursInDay    = 24
	DaysInWeek    = 7
	DaysInMonth   = 30 // Approximation
	DaysInYear    = 365
)

// DurationToMinutes converts a duration string to minutes.
// Supports minutes (m), hours (h), days (d), weeks (w), months (M), and years (y).
func DurationToMinutes(duration string) int64 {
	unit := duration[len(duration)-1:]
	value, err := strconv.Atoi(duration[:len(duration)-1])
	if err != nil {
		fmt.Println(err)
		return 0
	}
	switch unit {
	case "m":
		return int64(value)
	case "h":
		return int64(value) * MinutesInHour
	case "d":
		return int64(value) * MinutesInHour * HoursInDay
	case "w":
		return int64(value) * MinutesInHour * HoursInDay * DaysInWeek
	case "M":
		return int64(value) * MinutesInHour * HoursInDay * DaysInMonth
	case "y":
		return int64(value) * MinutesInHour * HoursInDay * DaysInYear
	default:
		return 0
	}
}
