package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/example-api/response"
)

// Handles requests to create a new user
// POST /v1/users
func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated client from the request context
	authenticatedAccount, err := GetAuthenticatedAccount(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Request body cannot be nil
	if r.Body == nil {
		response.Error(w, "Request body cannot be nil", http.StatusBadRequest)
		return
	}

	// Read the request body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into the request prototype
	userRequest := new(UserRequest)
	if err = json.Unmarshal(payload, userRequest); err != nil {
		logger.Errorf("Failed to unmarshal user request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new user account
	user, err := s.CreateUser(authenticatedAccount, userRequest)
	if err != nil {
		logger.Errorf("Create user error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// Create response
	userResponse, err := NewUserResponse(user)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Location header to the newly created resource
	w.Header().Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))
	// Write JSON response
	response.WriteJSON(w, userResponse, http.StatusCreated)
}
