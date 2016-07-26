package util_test

import (
	"testing"
	"time"

	"github.com/RichardKnop/example-api/util"
	"github.com/stretchr/testify/assert"
)

func TestFormatTime(t *testing.T) {
	var (
		timestamp        time.Time
		expected, actual string
	)

	// UTC
	timestamp = time.Date(2012, 12, 11, 8, 52, 31, 493729031, time.UTC)
	expected = "2012-12-11T08:52:31Z"
	actual = util.FormatTime(timestamp)
	assert.Equal(t, expected, actual)

	// UTC
	timestamp = time.Date(2012, 12, 11, 8, 52, 31, 493729031, time.FixedZone("HKT", 8*3600))
	expected = "2012-12-11T00:52:31Z"
	actual = util.FormatTime(timestamp)
	assert.Equal(t, expected, actual)
}

func TestParseTimestamp(t *testing.T) {
	var (
		parsedTimestamp time.Time
		err             error
	)

	parsedTimestamp, err = util.ParseTimestamp("bogus")
	assert.NotNil(t, err)

	parsedTimestamp, err = util.ParseTimestamp("2016-05-04T12:08:35Z")
	assert.Nil(t, err)
	assert.Equal(t, 2016, parsedTimestamp.UTC().Year())
	assert.Equal(t, time.May, parsedTimestamp.UTC().Month())
	assert.Equal(t, 4, parsedTimestamp.UTC().Day())
	assert.Equal(t, 12, parsedTimestamp.UTC().Hour())
	assert.Equal(t, 8, parsedTimestamp.UTC().Minute())
	assert.Equal(t, 35, parsedTimestamp.UTC().Second())

	parsedTimestamp, err = util.ParseTimestamp("2016-05-04T12:08:35+07:00")
	assert.Nil(t, err)
	assert.Equal(t, 2016, parsedTimestamp.UTC().Year())
	assert.Equal(t, time.May, parsedTimestamp.UTC().Month())
	assert.Equal(t, 4, parsedTimestamp.UTC().Day())
	assert.Equal(t, 5, parsedTimestamp.UTC().Hour())
	assert.Equal(t, 8, parsedTimestamp.UTC().Minute())
	assert.Equal(t, 35, parsedTimestamp.UTC().Second())
}
