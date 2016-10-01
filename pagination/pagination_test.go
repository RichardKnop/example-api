package pagination_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/RichardKnop/example-api/pagination"
)

func TestGetOffsetForPage(t *testing.T) {
	var offset int

	// First page offset should be zero
	offset = pagination.GetOffsetForPage(
		10, // count
		1,  // page
		2,  // limit
	)
	assert.Equal(t, 0, offset)

	// Second page offset should be 2
	offset = pagination.GetOffsetForPage(
		10, // count
		2,  // page
		2,  // limit
	)
	assert.Equal(t, 2, offset)

	// Last page offset should be 8
	offset = pagination.GetOffsetForPage(
		10, // count
		5,  // page
		2,  // limit
	)
	assert.Equal(t, 8, offset)
}
