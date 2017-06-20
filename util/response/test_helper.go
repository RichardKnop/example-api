package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestHandlerFailsWithoutAuthenticatedUser ...
func TestHandlerFailsWithoutAuthenticatedUser(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) {
	r, err := http.NewRequest("", "", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()

	handler(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

// TestListBadRequests tests a list response for common bad request failures
func TestListBadRequests(t *testing.T, entity string, router *mux.Router) {
	code := http.StatusBadRequest
	TestListFailsBadPage(t, entity, router, code)
	TestListFailsPageTooBig(t, entity, router, code)
}

// TestListFailsBadPage tests a list endpoint for a bad page response
func TestListFailsBadPage(t *testing.T, entity string, router *mux.Router, code int) {
	url := fmt.Sprintf("http://1.2.3.4/v1/%s?page=bad_page", entity)
	msg := "strconv.ParseInt: parsing \\\"bad_page\\\": invalid syntax"

	TestGetErrorExpectedResponse(t, router, url, msg, code)
}

// TestListFailsPageTooBig a list endpoint for a too big page response
func TestListFailsPageTooBig(t *testing.T, entity string, router *mux.Router, code int) {
	url := fmt.Sprintf("http://1.2.3.4/v1/%s?page=9999", entity)
	msg := "Page too big"

	TestGetErrorExpectedResponse(t, router, url, msg, code)
}

// TestGetErrorExpectedResponse ...
func TestGetErrorExpectedResponse(t *testing.T, router *mux.Router, url, msg string, code int) {
	TestErrorExpectedResponse(t, router, "GET", url, msg, code, nil)
}

// TestPutErrorExpectedResponse ...
func TestPutErrorExpectedResponse(t *testing.T, router *mux.Router, url, msg string, code int, data io.Reader) {
	TestErrorExpectedResponse(t, router, "PUT", url, msg, code, data)
}

// TestPostErrorExpectedResponse ...
func TestPostErrorExpectedResponse(t *testing.T, router *mux.Router, url, msg string, code int, data io.Reader) {
	TestErrorExpectedResponse(t, router, "POST", url, msg, code, data)
}

// TestErrorExpectedResponse is the generic test code for testing for a bad response
func TestErrorExpectedResponse(t *testing.T, router *mux.Router, operation, url, msg string, code int, data io.Reader) {
	// Prepare a request
	r, err := http.NewRequest(
		operation,
		url,
		data,
	)
	assert.NoError(t, err)

	// Mock authentication
	r.Header.Set("Authorization", "Bearer test_token")

	// And serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	TestResponseForError(t, w, msg, code)
}

// TestResponseForError tests a response w to see if it returned an error msg with http code
func TestResponseForError(t *testing.T, w *httptest.ResponseRecorder, msg string, code int) {
	assert.NotNil(t, w)
	assert.Equal(
		t,
		code,
		w.Code,
		fmt.Sprintf("Expected a %d response but got %d", code, w.Code),
	)
	TestResponseBody(t, w, getErrorJSON(msg))
}

// TestEmptyResponse tests an empty 204 response
func TestEmptyResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, 204, w.Code)
	TestResponseBody(t, w, "")
}

// TestResponseObject tests response body is equal to expected object in JSON form
func TestResponseObject(t *testing.T, w *httptest.ResponseRecorder, expected interface{}, code int) {
	assert.Equal(
		t,
		code,
		w.Code,
		fmt.Sprintf("Expected a %d response but got %d", code, w.Code),
	)
	jsonBytes, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.Equal(
		t,
		string(jsonBytes),
		strings.TrimRight(w.Body.String(), "\n"),
		"Should have returned correct body text",
	)
}

// TestResponseBody tests response body is equal to expected string
func TestResponseBody(t *testing.T, w *httptest.ResponseRecorder, expected string) {
	assert.Equal(
		t,
		expected,
		strings.TrimRight(w.Body.String(), "\n"),
		"Should have returned correct body text",
	)
}

// TestListValidResponse ...
func TestListValidResponse(t *testing.T, router *mux.Router, entity string, items []interface{}, assertExpectations func()) {
	TestListValidResponseWithParams(t, router, entity, items, assertExpectations, nil)
}

// TestListValidResponseWithParams tests a list endpoint for a valid response with default settings
func TestListValidResponseWithParams(t *testing.T, router *mux.Router, entity string, items []interface{}, assertExpectations func(), params map[string]string) {
	u, err := url.Parse(fmt.Sprintf("http://1.2.3.4/v1/%s", entity))
	assert.NoError(t, err, "URL parsing failed")

	// add any params
	for k, v := range params {
		q := u.Query()
		q.Set(k, v)
		u.RawQuery = q.Encode()
	}

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)
	assert.NoError(t, err)

	// Mock authentication
	r.Header.Set("Authorization", "Bearer test_token")

	// And serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	assertExpectations()
	assert.Equal(t, http.StatusOK, w.Code)

	baseURI := u.RequestURI()

	q := u.Query()
	q.Set("page", "1")
	u.RawQuery = q.Encode()

	pagedURI := u.RequestURI()

	expected := &ListResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": {
					Href: baseURI,
				},
				"first": {
					Href: pagedURI,
				},
				"last": {
					Href: pagedURI,
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				entity: jsonhal.Embedded(items),
			},
		},
		Count: uint(len(items)),
		Page:  1,
	}
	expectedJSON, err := json.Marshal(expected)

	if assert.NoError(t, err, "JSON marshalling failed") {
		TestResponseBody(t, w, string(expectedJSON))
	}
}

func getErrorJSON(msg string) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", msg)
}
