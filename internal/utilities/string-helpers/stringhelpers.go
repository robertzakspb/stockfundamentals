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

// Given a slice of 15 elements and a size of 4, the method would return 3 slices (3 with 4 elements and the last one with 3 elements)
func SplitInBatchesOf(size int, strings []string) [][]string {
	if len(strings) <= size {
		return [][]string{strings}
	}

	batches := [][]string{}

	filledBatchCount := len(strings) / size //All but the last batch will be filled completely
	lastBatchSize := len(strings) % size    //The last batch may have a different size

	for i := range filledBatchCount {
		batch := strings[size*i : size*(i+1)]
		batches = append(batches, batch)
	}

	if lastBatchSize == 0 {
		return batches
	}

	lastBatch := strings[filledBatchCount*size:]
	batches = append(batches, lastBatch)

	return batches
}
