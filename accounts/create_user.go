package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/RichardKnop/recall/response"
)

// Handles requests to create a new user (POST /v1/accounts/users)
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
	if err := json.Unmarshal(payload, userRequest); err != nil {
		log.Printf("Failed to unmarshal user request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if oauth user exists
	if s.GetOauthService().UserExists(userRequest.Email) {
		response.Error(w, "Email taken", http.StatusBadRequest)
		return
	}

	// Create a new user account
	user, err := s.CreateUser(authenticatedAccount, userRequest)
	if err != nil {
		log.Printf("Create user error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	userResponse, err := NewUserResponse(user)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Location header to the newly created resource
	w.Header().Set("Location", fmt.Sprintf("/v1/accounts/users/%d", user.ID))
	// Write JSON response
	response.WriteJSON(w, userResponse, http.StatusCreated)
}
