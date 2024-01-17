package helper

import (
	"fmt"
	"strconv"
)

// Convert age to minutes
//
// Example:
//
//	1s -> 1000
//	1m -> 60000
//	1h -> 3600000
//	1d -> 86400000
//	1w -> 604800000
//	1M -> 2592000000
func AgeToMs(age string) int {
	unit := age[len(age)-1:]
	ageInt, err := strconv.Atoi(age[:len(age)-1])
	if err != nil {
		fmt.Println(err)
		return 0
	}

	switch unit {
	//case "s":
	//	return ageInt * 1000
	case "m":
		return ageInt
	case "h":
		return ageInt * 60
	case "d":
		return ageInt * 60 * 24
	case "w":
		return ageInt * 60 * 24 * 7
	case "M":
		return ageInt * 60 * 24 * 30
	case "y":
		return ageInt * 60 * 24 * 365
	default:
		return 0
	}
}
