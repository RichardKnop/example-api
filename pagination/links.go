package pagination

import (
	"fmt"
	"math"
	"net/url"
)

// GetLinks returns links for first, last, previous and next page
func GetLinks(urlObject *url.URL, count, page, limit int) (string, string, string, string, error) {
	var (
		first    string
		last     string
		previous string
		next     string
		q        url.Values
	)

	// Number of pages
	nuPages := int(math.Ceil(float64(count) / float64(limit)))
	if nuPages < 1 {
		nuPages = 1
	}

	// Page too big
	if page > nuPages {
		return first, last, previous, next, ErrPageTooBig
	}

	// First page
	q = urlObject.Query()
	q.Set("page", fmt.Sprintf("%d", 1))
	first = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())

	// Last page
	q = urlObject.Query()
	q.Set("page", fmt.Sprintf("%d", nuPages))
	last = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())

	// Previous page
	if page > 1 {
		q := urlObject.Query()
		q.Set("page", fmt.Sprintf("%d", page-1))
		previous = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())
	}

	// Next page
	if page < nuPages {
		q := urlObject.Query()
		q.Set("page", fmt.Sprintf("%d", page+1))
		next = fmt.Sprintf("%s?%s", urlObject.Path, q.Encode())
	}

	return first, last, previous, next, nil
}
