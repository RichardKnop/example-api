package accounts

import (
	"net/http"
	"strconv"

	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/oauth/tokentypes"
	"github.com/RichardKnop/example-api/util/response"
	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
)

// Handles requests to confirm user's email
// GET /v1/confirmations/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}
func (s *Service) confirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated client from the request context
	authenticatedClient, err := GetAuthenticatedClient(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the reference from request URI
	vars := mux.Vars(r)
	reference := vars["reference"]

	// Fetch the confirmation we want to work with (by reference from email link)
	confirmation, err := s.FindConfirmationByReference(reference)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Confirm the user
	if err = s.ConfirmUser(confirmation); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create confirmation response
	confirmationResponse, err := NewConfirmationResponse(confirmation)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Was autologin flag passed in the query string
	autoLogin, _ := strconv.ParseBool(r.URL.Query().Get("autologin"))

	// If autologin == true, login the user and embed access token in the response object
	if autoLogin {
		// Login the user
		accessToken, refreshToken, err := s.GetOauthService().Login(
			authenticatedClient,
			confirmation.User.OauthUser,
			"read_write",
		)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create access token response
		accessTokenResponse, err := oauth.NewAccessTokenResponse(
			accessToken,
			refreshToken,
			s.cnf.Oauth.AccessTokenLifetime,
			tokentypes.Bearer,
		)
		if err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set embedded access token
		confirmationResponse.SetEmbedded(
			"access-token",
			jsonhal.Embedded(accessTokenResponse),
		)
	}

	// Write the response
	response.WriteJSON(w, confirmationResponse, 200)
}
