package pagination

import (
	"errors"
	"fmt"
)

var (
	// DefaultLimit defines number of items to return per page
	DefaultLimit = 25
	// MinLimit defines minimum limit
	MinLimit = 1
	// MaxLimit defines maximum limit of items to return per page
	MaxLimit = 100
	// ErrPageTooSmall - page is not a positive integer
	ErrPageTooSmall = errors.New("Page must be > 0")
	// ErrPageTooBig - page is greater number than there are actual pages
	ErrPageTooBig = errors.New("Page too big")
	// ErrLimitTooSmall - limit is less than MinLimit
	ErrLimitTooSmall = fmt.Errorf("Limit must be < %d", MinLimit)
	// ErrLimitTooBig - limit is more than MaxLimit
	ErrLimitTooBig = fmt.Errorf("Limit must be <= %d", MaxLimit)
)

// GetOffsetForPage returns an offset for a page
func GetOffsetForPage(count, page, limit int) int {
	return limit * (page - 1)
}
