package utils

func MsToReadable(ms int64) string {
	seconds := ms / 1000
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	if days > 0 {
		return string(rune(days)) + "d"
	}

	if hours > 0 {
		return string(rune(hours)) + "h"
	}

	if minutes > 0 {
		return string(rune(minutes)) + "m"
	}

	if seconds > 0 {
		return string(rune(seconds)) + "s"
	}

	return string(rune(ms)) + "ms"
}
