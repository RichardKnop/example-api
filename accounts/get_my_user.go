package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/response"
)

// Handles requests to get own user data (GET /v1/accounts/me)
func (s *Service) getMyUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Create response
	userResponse, err := NewUserResponse(authenticatedUser)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, userResponse, http.StatusOK)
}
