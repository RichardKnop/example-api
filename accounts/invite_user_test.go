package accounts_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/accounts/roles"
	"github.com/RichardKnop/example-api/response"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/jsonhal"
	"github.com/stretchr/testify/assert"
)

func (suite *AccountsTestSuite) TestInviteUser() {
	// Prepare a request
	payload, err := json.Marshal(&accounts.InvitationRequest{
		Email: "john@reese",
	})
	assert.NoError(suite.T(), err, "JSON marshalling failed")
	r, err := http.NewRequest(
		"POST",
		"http://1.2.3.4/v1/accounts/invitations",
		bytes.NewBuffer(payload),
	)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Header.Set("Authorization", "Bearer test_user_token")

	// Mock invitation email
	suite.mockInvitationEmail()

	// Count before
	var (
		countBefore            int
		invitationsCountBefore int
	)
	suite.db.Model(new(accounts.User)).Count(&countBefore)
	suite.db.Model(new(accounts.Invitation)).Count(&invitationsCountBefore)

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	if !assert.Equal(suite.T(), 201, w.Code) {
		log.Print(w.Body.String())
	}

	// Count after
	var (
		countAfter            int
		invitationsCountAfter int
	)
	suite.db.Model(new(accounts.User)).Count(&countAfter)
	suite.db.Model(new(accounts.Invitation)).Count(&invitationsCountAfter)
	assert.Equal(suite.T(), countBefore+1, countAfter)
	assert.Equal(suite.T(), invitationsCountBefore+1, invitationsCountAfter)

	// Fetch the created invitation
	invitation := new(accounts.Invitation)
	assert.False(suite.T(), suite.db.
		Preload("InvitedUser.OauthUser").Preload("InvitedUser.Role").
		Preload("InvitedByUser.OauthUser").Preload("InvitedByUser.Role").
		Last(invitation).RecordNotFound())

	// And correct data was saved
	assert.Equal(suite.T(), invitation.InvitedUser.ID, invitation.InvitedUser.OauthUser.MetaUserID)
	assert.Equal(suite.T(), "john@reese", invitation.InvitedUser.OauthUser.Username)
	assert.False(suite.T(), invitation.InvitedUser.OauthUser.Password.Valid)
	assert.Equal(suite.T(), roles.User, invitation.InvitedUser.RoleID.String)
	assert.Equal(suite.T(), "test@user", invitation.InvitedByUser.OauthUser.Username)

	// Check the response body
	expected := &accounts.InvitationResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": &jsonhal.Link{
					Href: fmt.Sprintf("/v1/accounts/invitations/%d", invitation.ID),
				},
			},
		},
		ID:              invitation.ID,
		Reference:       invitation.Reference,
		InvitedUserID:   invitation.InvitedUser.ID,
		InvitedByUserID: invitation.InvitedByUser.ID,
		CreatedAt:       util.FormatTime(invitation.CreatedAt),
		UpdatedAt:       util.FormatTime(invitation.UpdatedAt),
	}
	expectedJSON, err := json.Marshal(expected)
	if assert.NoError(suite.T(), err) {
		response.TestResponseBody(suite.T(), w, string(expectedJSON))
	}

	// Wait for the email goroutine to finish
	<-time.After(5 * time.Millisecond)

	// Check that the mock object expectations were met
	suite.assertMockExpectations()
}
