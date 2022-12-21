package utils

import "strconv"

func MsToReadable(ms int64) string {
	seconds := ms / 1000
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	if days > 0 {
		return strconv.Itoa(int(days)) + "d"
	}

	if hours > 0 {
		return strconv.Itoa(int(hours)) + "h"
	}

	if minutes > 0 {
		return strconv.Itoa(int(minutes)) + "m"
	}

	if seconds > 0 {
		return strconv.Itoa(int(seconds)) + "s"
	}

	return strconv.Itoa(int(ms)) + "ms"
}
