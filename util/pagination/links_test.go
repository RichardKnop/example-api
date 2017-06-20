package pagination_test

import (
	"log"
	"net/url"
	"testing"

	"github.com/RichardKnop/example-api/util/pagination"
	"github.com/stretchr/testify/assert"
)

func TestGetLinks(t *testing.T) {
	// Test with both absolute and relative URI
	var testURLs = []string{
		"https://foo.bar/foobar?hello=world",
		"/foobar?hello=world",
	}
	for _, testURL := range testURLs {
		testGetLinks(t, testURL)
	}
}

func testGetLinks(t *testing.T, testURL string) {
	var (
		urlObject *url.URL
		first     string
		last      string
		previous  string
		next      string
		err       error
	)

	// Test URL object
	urlObject, err = url.Parse(testURL)
	if err != nil {
		log.Fatal(err)
	}

	// Test with zero results
	first, last, previous, next, err = pagination.GetLinks(
		urlObject,
		0, // count
		1, // page
		2, // limit
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "/foobar?hello=world&page=1", first)
		assert.Equal(t, "/foobar?hello=world&page=1", last)
		assert.Equal(t, "", previous)
		assert.Equal(t, "", next)
	}

	// Test first page
	first, last, previous, next, err = pagination.GetLinks(
		urlObject,
		10, // count
		1,  // page
		2,  // limit
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "/foobar?hello=world&page=1", first)
		assert.Equal(t, "/foobar?hello=world&page=5", last)
		assert.Equal(t, "", previous)
		assert.Equal(t, "/foobar?hello=world&page=2", next)
	}

	// Test middle page
	first, last, previous, next, err = pagination.GetLinks(
		urlObject,
		10, // count
		2,  // page
		2,  // limit
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "/foobar?hello=world&page=1", first)
		assert.Equal(t, "/foobar?hello=world&page=5", last)
		assert.Equal(t, "/foobar?hello=world&page=1", previous)
		assert.Equal(t, "/foobar?hello=world&page=3", next)
	}

	// Test last page
	first, last, previous, next, err = pagination.GetLinks(
		urlObject,
		10, // count
		5,  // page
		2,  // limit
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "/foobar?hello=world&page=1", first)
		assert.Equal(t, "/foobar?hello=world&page=5", last)
		assert.Equal(t, "/foobar?hello=world&page=4", previous)
		assert.Equal(t, "", next)
	}

	// Test page too big
	_, _, _, _, err = pagination.GetLinks(
		urlObject,
		10, // count
		6,  // page
		2,  // limit
	)
	if assert.NotNil(t, err) {
		assert.Equal(t, pagination.ErrPageTooBig, err)
	}

	// Test when limit > count
	first, last, previous, next, err = pagination.GetLinks(
		urlObject,
		10, // count
		1,  // page
		12, // limit
	)
	if assert.Nil(t, err) {
		assert.Equal(t, "/foobar?hello=world&page=1", first)
		assert.Equal(t, "/foobar?hello=world&page=1", last)
		assert.Equal(t, "", previous)
		assert.Equal(t, "", next)
	}
}
