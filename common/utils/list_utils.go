package utils

// Return unique items in a list
//
// @param list []string
//
// Example: ListToUnique([]string{"a", "b", "a", "c"}) => []string{"a", "b", "c"}
func ListToUnique(list []string) []string {
	uniqueList := make([]string, 0)

	for _, item := range list {
		if !Contains(uniqueList, item) {
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}

// Check if an item is in a list
//
// @param list []string
//
// @param item string
//
// Example: Contains([]string{"a", "b", "c"}, "b") => true
func Contains(list []string, item string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}

	return false
}
