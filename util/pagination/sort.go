package pagination

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// Ascending sorting direction
	Ascending = "ASC"
	// Descending sorting direction
	Descending = "DESC"
)

var (
	// ErrEmptySortField ...
	ErrEmptySortField = errors.New("Empty sort field")
)

// GetSort parses and validates a sort parameter and returns a sorting map
func GetSort(sortParam string, sortableBy []string) (map[string]string, error) {
	var sorts = map[string]string{}

	if strings.Trim(sortParam, " ") == "" {
		return sorts, nil
	}

	for _, f := range strings.Split(sortParam, ",") {
		// Empty values not allowed
		f = strings.Trim(f, " ")
		if f == "" {
			return sorts, ErrEmptySortField
		}

		// If the field start with - (to indicate descended sort), shift it before
		// validator lookup
		var sortDirection = Ascending // assume ascending direction
		var i = 0
		if f[0] == '-' {
			i = 1
			sortDirection = Descending // switch to descending direction
		}

		// Make sure the field is sortable
		if !isSortableBy(f[i:], sortableBy) {
			return sorts, fmt.Errorf("Invalid sort field: %s", f)
		}

		sorts[f[i:]] = sortDirection
	}

	return sorts, nil
}

func isSortableBy(field string, sortableBy []string) bool {
	for _, s := range sortableBy {
		if field == s {
			return true
		}
	}
	return false
}
