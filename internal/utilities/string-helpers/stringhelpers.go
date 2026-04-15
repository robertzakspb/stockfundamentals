package stringhelpers

func RemoveDuplicatesFrom(slice []string) []string {
	noDuplicateSlice := []string{}

	for i := range slice {
		foundDuplicate := false
		for j := range noDuplicateSlice {
			if slice[i] == noDuplicateSlice[j] {
				foundDuplicate = true
			}
		}
		if !foundDuplicate {
			noDuplicateSlice = append(noDuplicateSlice, slice[i])
		}
	}

	return noDuplicateSlice
}
