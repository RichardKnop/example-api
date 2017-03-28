package accounts_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/oauth/roles"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestGetMyUserRequiresUserAuthentication() {
	// Prepare a request
	r, err := http.NewRequest("GET", "http://1.2.3.4/v1/me", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code, "This requires an authenticated user")
}

func (suite *AccountsTestSuite) TestGetMyUser() {
	// Prepare a request
	r, err := http.NewRequest("GET", "http://1.2.3.4/v1/me", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_user_token")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route) {
		assert.Equal(suite.T(), "get_my_user", match.Route.GetName())
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()

	// Check the status code
	if !assert.Equal(suite.T(), 200, w.Code) {
		log.Print(w.Body.String())
	}

	// Fetch the user
	user := new(models.User)
	notFound := models.UserPreload(suite.db).First(user, suite.users[1].ID).RecordNotFound()
	assert.False(suite.T(), notFound)

	// Check the response
	expected := &accounts.UserResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/users/%d", user.ID),
				},
			},
		},
		ID:        user.ID,
		Email:     user.OauthUser.Username,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Role:      roles.User,
		Confirmed: user.Confirmed,
		CreatedAt: util.FormatTime(&user.CreatedAt),
		UpdatedAt: util.FormatTime(&user.UpdatedAt),
	}
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
