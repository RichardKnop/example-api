package pagination_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/RichardKnop/example-api/util/pagination"
)

func TestGetParams(t *testing.T) {
	var (
		page       int
		limit      int
		sorts      map[string]string
		r          *http.Request
		err        error
		sortableBy []string
	)

	// Test default values
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, sorts, err = pagination.GetParams(r, sortableBy)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, page)
		assert.Equal(t, 25, limit)
		assert.Equal(t, 0, len(sorts))
	}

	// Test page <= 0
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, sorts, err = pagination.GetParams(r, sortableBy)
	if assert.Error(t, err) {
		assert.Equal(t, pagination.ErrPageTooSmall, err)
	}

	// Test limit too small
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=1&limit=0", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, sorts, err = pagination.GetParams(r, sortableBy)
	if assert.Error(t, err) {
		assert.Equal(t, pagination.ErrLimitTooSmall, err)
	}

	// Test limit too big
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=1&limit=1000", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, sorts, err = pagination.GetParams(r, sortableBy)
	if assert.Error(t, err) {
		assert.Equal(t, pagination.ErrLimitTooBig, err)
	}

	// Test valid page and limit
	r, err = http.NewRequest("GET", "http://1.2.3.4/v1/foo/bar?page=10&limit=50", nil)
	if err != nil {
		log.Fatal(err)
	}
	page, limit, sorts, err = pagination.GetParams(r, sortableBy)
	if assert.NoError(t, err) {
		assert.Equal(t, 10, page)
		assert.Equal(t, 50, limit)
		assert.Equal(t, 0, len(sorts))
	}
}
