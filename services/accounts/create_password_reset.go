package accounts

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/example-api/log"
	"github.com/RichardKnop/example-api/util/response"
)

// Handles requests to reset a password
// POST /v1/password-resets
func (s *Service) createPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated client from the request context
	_, err := GetAuthenticatedClient(r)
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
	passwordResetRequest := new(PasswordResetRequest)
	if err = json.Unmarshal(payload, passwordResetRequest); err != nil {
		log.ERROR.Printf("Failed to unmarshal password reset request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user who wants to reset his/her password based on the email
	user, err := s.FindUserByEmail(passwordResetRequest.Email)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create a new password reset
	passwordReset, err := s.createPasswordReset(user)
	if err != nil {
		log.ERROR.Printf("Create password reset error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create password reset response
	passwordResetResponse, err := NewPasswordResetResponse(passwordReset)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	response.WriteJSON(w, passwordResetResponse, 201)
}
