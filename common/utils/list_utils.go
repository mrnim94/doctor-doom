package utils

func ListToUnique(list []string) []string {
	uniqueList := make([]string, 0)

	for _, item := range list {
		if !Contains(uniqueList, item) {
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}

func Contains(list []string, item string) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}

	return false
}
