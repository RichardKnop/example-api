package pagination

import (
	"net/http"
	"strconv"
)

// GetParams parses querystring and returns pagination params:
// - page
// - limit
// - sort map
func GetParams(r *http.Request, sortableBy []string) (int, int, map[string]string, error) {
	var (
		page    = 1            // default page
		limit   = DefaultLimit // default limit
		sortMap map[string]string
		err     error
	)

	// Get page from the querystring
	if r.URL.Query().Get("page") != "" {
		// String to int conversion
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			return 0, 0, sortMap, err
		}

		if page < 1 {
			return 0, 0, sortMap, ErrPageTooSmall
		}
	}

	// Get limit from the querystring
	if r.URL.Query().Get("limit") != "" {
		// String to int conversion
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			return 0, 0, sortMap, err
		}

		if limit < MinLimit {
			return 0, 0, sortMap, ErrLimitTooSmall
		}

		if limit > MaxLimit {
			return 0, 0, sortMap, ErrLimitTooBig
		}
	}

	// Get sort from the querystring
	if r.URL.Query().Get("sort") != "" {
		sortMap, err = GetSort(r.URL.Query().Get("sort"), sortableBy)
		if err != nil {
			return 0, 0, sortMap, err
		}
	}

	return page, limit, sortMap, nil
}
