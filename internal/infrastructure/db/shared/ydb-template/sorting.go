package ydbtemplate

import "strings"

func addSortingToQuery(sb *strings.Builder, sortBy, direction string) {
	if sortBy == "" {
		return
	}

	sortByQuery := "\nORDER BY " + sortBy

	switch direction {
	case ">":
		sortByQuery += " ASC\n"

	case "<":
		sortByQuery += " DESC\n"
	default:
		sortByQuery += " DESC\n"
	}

	sb.WriteString(sortByQuery)
}
