package stringutils

func ListContains(haystack []string, needle string) bool {
	found := false

	for _, elm := range haystack {
		if elm == needle {
			found = true
			break
		}
	}

	return found
}
