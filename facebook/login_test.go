package facebook_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/RichardKnop/example-api/facebook"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/oauth/tokentypes"
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
	suite.db.Model(new(models.User)).Count(&countBefore)

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
	suite.db.Model(new(models.User)).Count(&countAfter)
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
		"id":         suite.users[1].FacebookID.String,
		"email":      suite.users[1].OauthUser.Username,
		"name":       suite.users[1].GetName(),
		"first_name": suite.users[1].FirstName.String,
		"last_name":  suite.users[1].LastName.String,
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(models.User)).Count(&countBefore)

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
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", facebook.ErrAccountMismatch.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *FacebookTestSuite) TestLoginExistingUser() {
	var (
		testOauthUser *models.OauthUser
		testUser      *models.User
		err           error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetAccountsService().GetOauthService().CreateUser(
		roles.User,
		"harold@finch",
		"", // empty password
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser, err = models.NewUser(
		suite.accounts[0],
		testOauthUser,
		"some_facebook_id", // facebook ID
		"Harold",
		"Finch",
		"some_picture",
		true, // confirmed
	)
	assert.NoError(suite.T(), err, "Failed to create a new test account user object")
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.OauthUser = testOauthUser

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
		"id":         "some_facebook_id",
		"email":      testUser.OauthUser.Username,
		"name":       testUser.GetName(),
		"first_name": testUser.FirstName.String,
		"last_name":  testUser.LastName.String,
		"picture": map[string]interface{}{
			"data": map[string]interface{}{
				"url": testUser.Picture.String,
			},
		},
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(models.User)).Count(&countBefore)

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
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the logged in user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// The user should not have changed
	assert.Equal(suite.T(), testUser.OauthUser.Username, user.OauthUser.Username)
	assert.Equal(suite.T(), testUser.FacebookID.String, user.FacebookID.String)
	assert.Equal(suite.T(), testUser.FirstName.String, user.FirstName.String)
	assert.Equal(suite.T(), testUser.LastName.String, user.LastName.String)
	assert.Equal(suite.T(), testUser.Picture.String, user.Picture.String)

	// Fetch oauth tokens
	accessToken := new(models.OauthAccessToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		First(accessToken).RecordNotFound())
	refreshToken := new(models.OauthRefreshToken)
	assert.False(suite.T(), models.OauthRefreshTokenPreload(suite.db).
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *FacebookTestSuite) TestLoginUpdatesExistingUser() {
	var (
		testOauthUser *models.OauthUser
		testUser      *models.User
		err           error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetAccountsService().GetOauthService().CreateUser(
		roles.User,
		"harold@finch",
		"", // empty password
	)
	assert.NoError(suite.T(), err, "Failed to insert a test oauth user")
	testUser, err = models.NewUser(
		suite.accounts[0],
		testOauthUser,
		"some_facebook_id", // facebook ID
		"Harold",
		"Finch",
		"some_picture",
		true, // confirmed
	)
	assert.NoError(suite.T(), err, "Failed to create a new test account user object")
	err = suite.db.Create(testUser).Error
	assert.NoError(suite.T(), err, "Failed to insert a test user")
	testUser.OauthUser = testOauthUser

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
		"id":         "new_facebook_id",
		"email":      testUser.OauthUser.Username,
		"name":       "New Name",
		"first_name": "New First Name",
		"last_name":  "New Last Name",
		"picture": map[string]interface{}{
			"data": map[string]interface{}{
				"url": "new_picture",
			},
		},
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(models.User)).Count(&countBefore)

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
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore, countAfter)

	// Fetch the updated user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).First(user, testUser.ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), testUser.OauthUser.Username, user.OauthUser.Username)
	assert.Equal(suite.T(), "New First Name", user.FirstName.String)
	assert.Equal(suite.T(), "New Last Name", user.LastName.String)
	assert.Equal(suite.T(), "new_facebook_id", user.FacebookID.String)
	assert.Equal(suite.T(), "new_picture", user.Picture.String)
	assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)

	// Fetch oauth tokens
	accessToken := new(models.OauthAccessToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		First(accessToken).RecordNotFound())
	refreshToken := new(models.OauthRefreshToken)
	assert.False(suite.T(), models.OauthRefreshTokenPreload(suite.db).
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
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
		"id":         "new_facebook_id",
		"email":      "new@user",
		"name":       "John Reese",
		"first_name": "John",
		"last_name":  "Reese",
		"picture": map[string]interface{}{
			"data": map[string]interface{}{
				"url": "johns_picture",
			},
		},
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(models.User)).Count(&countBefore)

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
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "new@user", user.OauthUser.Username)
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
	assert.Equal(suite.T(), "new_facebook_id", user.FacebookID.String)
	assert.Equal(suite.T(), "johns_picture", user.Picture.String)
	assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)

	// Fetch oauth tokens
	accessToken := new(models.OauthAccessToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		First(accessToken).RecordNotFound())
	refreshToken := new(models.OauthRefreshToken)
	assert.False(suite.T(), models.OauthRefreshTokenPreload(suite.db).
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
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
		"id":         "new_facebook_id",
		"email":      nil,
		"name":       "John Reese",
		"first_name": "John",
		"last_name":  "Reese",
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(models.User)).Count(&countBefore)

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
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "new_facebook_id@facebook.com", user.OauthUser.Username)
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
	assert.Equal(suite.T(), "new_facebook_id", user.FacebookID.String)
	assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)
	assert.True(suite.T(), user.Confirmed)

	// Fetch oauth tokens
	accessToken := new(models.OauthAccessToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		First(accessToken).RecordNotFound())
	refreshToken := new(models.OauthRefreshToken)
	assert.False(suite.T(), models.OauthRefreshTokenPreload(suite.db).
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *FacebookTestSuite) TestLoginCreatesNewUserNoPicture() {
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
		"id":         "new_facebook_id",
		"email":      "new@user",
		"name":       "John Reese",
		"first_name": "John",
		"last_name":  "Reese",
		"picture":    nil,
	}, nil)

	// Count before
	var countBefore int
	suite.db.Model(new(models.User)).Count(&countBefore)

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
	suite.db.Model(new(models.User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), user.ID, user.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "new@user", user.OauthUser.Username)
	assert.Equal(suite.T(), "John", user.FirstName.String)
	assert.Equal(suite.T(), "Reese", user.LastName.String)
	assert.Equal(suite.T(), "new_facebook_id", user.FacebookID.String)
	assert.False(suite.T(), user.Picture.Valid)
	assert.Equal(suite.T(), roles.User, user.OauthUser.RoleID.String)

	// Fetch oauth tokens
	accessToken := new(models.OauthAccessToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		First(accessToken).RecordNotFound())
	refreshToken := new(models.OauthRefreshToken)
	assert.False(suite.T(), models.OauthRefreshTokenPreload(suite.db).
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       user.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
