package util

import (
	"time"
)

// FormatTime formats a time object to RFC3339
func FormatTime(timestamp time.Time) string {
	return timestamp.UTC().Format(time.RFC3339)
}

// ParseTimestamp parses a string representation of a timestamp in RFC3339
// format and returns a time.Time instance
func ParseTimestamp(timestamp string) (time.Time, error) {
	// RFC3339 = "2006-01-02T15:04:05Z07:00"
	return time.Parse(time.RFC3339, timestamp)
}
