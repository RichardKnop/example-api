package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/RichardKnop/jsonhal"
	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestCreateUser() {
	// Prepare a request
	payload, err := json.Marshal(&UserRequest{
		Email:    "test@user2",
		Password: "test_password",
	})
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/accounts/users",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Set("Authorization", "Basic dGVzdF9jbGllbnQ6dGVzdF9zZWNyZXQ=")

	// Check the routing
	match := new(mux.RouteMatch)
	suite.router.Match(r, match)
	if assert.NotNil(suite.T(), match.Route.GetName()) {
		assert.Equal(suite.T(), "create_user", match.Route.GetName())
	}

	// Count before
	var (
		countBefore int
	)
	suite.db.Model(new(User)).Count(&countBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter int
	)
	suite.db.Model(new(User)).Count(&countAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)

	// Fetch the created user
	user := new(User)
	notFound := suite.db.Preload("Account").Preload("OauthUser").
		Preload("Role").Last(user).RecordNotFound()
	assert.False(suite.T(), notFound)

	// And correct data was saved
	assert.Equal(suite.T(), "test@user2", user.OauthUser.Username)
	assert.Equal(suite.T(), "", user.FirstName.String)
	assert.Equal(suite.T(), "", user.LastName.String)
	assert.Equal(suite.T(), roles.User, user.Role.Name)

	// Check the Location header
	assert.Equal(
		suite.T(),
		fmt.Sprintf("/v1/accounts/users/%d", user.ID),
		w.Header().Get("Location"),
	)

	// Check the response body
	expected := &UserResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/users/%d", user.ID),
				},
			},
		},
		ID:        user.ID,
		Email:     "test@user2",
		Role:      roles.User,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(
		suite.T(),
		string(expectedJSON),
		strings.TrimRight(w.Body.String(), "\n"), // trim the trailing \n
	)
}
