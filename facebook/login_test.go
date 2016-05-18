package facebook

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/RichardKnop/recall/oauth"
	"github.com/gorilla/mux"
	fb "github.com/huandu/facebook"
	"github.com/stretchr/testify/assert"
)

func (suite *FacebookTestSuite) TestLoginFacebookCallFails() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/facebook/login", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{
		"access_token": {"facebook_token"},
		"scope":        {"read_write"},
	}
	r.SetBasicAuth("test_client_1", "test_secret")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "facebook_login", match.Route.GetName())
	}

	// Mock fetching profile data from facebook
	suite.mockFacebookGetMe(nil, errors.New("Some error from facebook"))

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.adapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 401, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"Some error from facebook\"}",
		strings.TrimSpace(w.Body.String()),
	)
}

// This checks that error is returned when an API key being used to make a call
// is different from the one used previously when the facebook user was created.
func (suite *FacebookTestSuite) TestLoginErrAccountMismatch() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/facebook/login", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{
		"access_token": {"facebook_token"},
		"scope":        {"read_write"},
	}
	r.SetBasicAuth("test_client_2", "test_secret")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "facebook_login", match.Route.GetName())
	}

	// Mock fetching profile data from facebook
	suite.mockFacebookGetMe(fb.Result{
		"id":         interface{}("facebook_id_2"),
		"email":      interface{}("test@user"),
		"first_name": interface{}("test_first_name_2"),
		"last_name":  interface{}("test_last_name_2"),
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.adapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 401, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", ErrAccountMismatch.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *FacebookTestSuite) TestLoginExistingUser() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/facebook/login", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{
		"access_token": {"facebook_token"},
		"scope":        {"read_write"},
	}
	r.SetBasicAuth("test_client_1", "test_secret")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "facebook_login", match.Route.GetName())
	}

	// Mock fetching profile data from facebook
	suite.mockFacebookGetMe(fb.Result{
		"id":         interface{}("facebook_id_2"),
		"email":      interface{}(suite.users[1].OauthUser.Username),
		"first_name": interface{}("test_first_name_2"),
		"last_name":  interface{}("test_last_name_2"),
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.adapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the logged in user
	user := new(accounts.User)
	notFound := suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role.Permissions").First(user, suite.users[1].ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// The user should not have changed
	assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
	assert.Equal(suite.T(), "facebook_id_2", user.FacebookID.String)
	assert.Equal(suite.T(), "test_first_name_2", user.FirstName.String)
	assert.Equal(suite.T(), "test_last_name_2", user.LastName.String)

	// Fetch oauth tokens
	accessToken := new(oauth.AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())
	refreshToken := new(oauth.RefreshToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    oauth.TokenType,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *FacebookTestSuite) TestLoginUpdatesExistingUser() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/facebook/login", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{
		"access_token": {"facebook_token"},
		"scope":        {"read_write"},
	}
	r.SetBasicAuth("test_client_1", "test_secret")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "facebook_login", match.Route.GetName())
	}

	// Mock fetching profile data from facebook
	suite.mockFacebookGetMe(fb.Result{
		"id":         interface{}("new_facebook_id"),
		"email":      interface{}(suite.users[1].OauthUser.Username),
		"first_name": interface{}("Harold"),
		"last_name":  interface{}("Finch"),
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.adapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(accounts.User)
	notFound := suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role.Permissions").First(user, suite.users[1].ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), "test@user", user.OauthUser.Username)
	assert.Equal(suite.T(), "Harold", user.FirstName.String)
	assert.Equal(suite.T(), "Finch", user.LastName.String)
	assert.Equal(suite.T(), "new_facebook_id", user.FacebookID.String)
	assert.Equal(suite.T(), roles.User, user.Role.ID)
	assert.True(suite.T(), user.Confirmed)

	// Fetch oauth tokens
	accessToken := new(oauth.AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())
	refreshToken := new(oauth.RefreshToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    oauth.TokenType,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *FacebookTestSuite) TestLoginCreatesNewUser() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/facebook/login", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{
		"access_token": {"facebook_token"},
		"scope":        {"read_write"},
	}
	r.SetBasicAuth("test_client_1", "test_secret")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "facebook_login", match.Route.GetName())
	}

	// Mock fetching profile data from facebook
	suite.mockFacebookGetMe(fb.Result{
		"id":         interface{}("new_facebook_id"),
		"email":      interface{}("new@user"),
		"first_name": interface{}("John"),
		"last_name":  interface{}("Reese"),
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.adapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created user
	user := new(accounts.User)
	notFound := suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role.Permissions").Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "new@user", user.OauthUser.Username)
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
	assert.Equal(suite.T(), "new_facebook_id", user.FacebookID.String)
	assert.Equal(suite.T(), roles.User, user.Role.ID)
	assert.True(suite.T(), user.Confirmed)

	// Fetch oauth tokens
	accessToken := new(oauth.AccessToken)
	assert.False(suite.T(), suite.db.First(accessToken).RecordNotFound())
	refreshToken := new(oauth.RefreshToken)
	assert.False(suite.T(), suite.db.First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    oauth.TokenType,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *FacebookTestSuite) TestLoginCreatesNewUserNilEmail() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/facebook/login", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{
		"access_token": {"facebook_token"},
		"scope":        {"read_write"},
	}
	r.SetBasicAuth("test_client_1", "test_secret")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "facebook_login", match.Route.GetName())
	}

	// Mock fetching profile data from facebook
	suite.mockFacebookGetMe(fb.Result{
		"id":         interface{}("new_facebook_id"),
		"email":      interface{}(nil),
		"first_name": interface{}("John"),
		"last_name":  interface{}("Reese"),
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(accounts.User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check mock expectations were met
	suite.adapterMock.AssertExpectations(suite.T())

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var countAfter int
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created user
	user := new(accounts.User)
	notFound := suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role.Permissions").Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "new_facebook_id@facebook.com", user.OauthUser.Username)
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
	assert.Equal(suite.T(), "new_facebook_id", user.FacebookID.String)
	assert.Equal(suite.T(), roles.User, user.Role.ID)
	assert.True(suite.T(), user.Confirmed)

	// Fetch oauth tokens
	accessToken := new(oauth.AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())
	refreshToken := new(oauth.RefreshToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    oauth.TokenType,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
