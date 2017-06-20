package pagination_test

import (
	"testing"

	"github.com/RichardKnop/example-api/util/pagination"
	"github.com/stretchr/testify/assert"
)

func TestGetSort(t *testing.T) {
	var (
		sortableBy = []string{"id", "created_at"}
		sorts      map[string]string
		err        error
	)

	// Empty sort param
	sorts, err = pagination.GetSort("", sortableBy)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(sorts))
	}

	// Empty sort param after trimming spaces
	sorts, err = pagination.GetSort(" ", sortableBy)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, len(sorts))
	}

	// Empty sort field
	sorts, err = pagination.GetSort(" ,", sortableBy)
	if assert.Error(t, err) {
		assert.Equal(t, pagination.ErrEmptySortField, err)
		assert.Equal(t, 0, len(sorts))
	}

	// Bogus sort field
	sorts, err = pagination.GetSort("id,bogus", sortableBy)
	if assert.Error(t, err) {
		assert.Equal(t, "Invalid sort field: bogus", err.Error())
		assert.Equal(t, 1, len(sorts))
		assert.Equal(t, pagination.Ascending, sorts["id"])
	}

	// Bogus sort field with minus sign
	sorts, err = pagination.GetSort("id,-bogus", sortableBy)
	if assert.Error(t, err) {
		assert.Equal(t, "Invalid sort field: -bogus", err.Error())
		assert.Equal(t, 1, len(sorts))
		assert.Equal(t, pagination.Ascending, sorts["id"])
	}

	// One field ascending
	sorts, err = pagination.GetSort("id", sortableBy)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(sorts))
		assert.Equal(t, pagination.Ascending, sorts["id"])
	}

	// One field descending
	sorts, err = pagination.GetSort("-id", sortableBy)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(sorts))
		assert.Equal(t, pagination.Descending, sorts["id"])
	}

	// Multiple fields
	sorts, err = pagination.GetSort("id,-created_at", sortableBy)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(sorts))
		assert.Equal(t, pagination.Ascending, sorts["id"])
		assert.Equal(t, pagination.Descending, sorts["created_at"])
	}
}
