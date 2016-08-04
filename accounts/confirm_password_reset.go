package accounts

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/RichardKnop/example-api/response"

	"github.com/gorilla/mux"
)

// ConfirmPasswordResetHandler - requests to complete a password reset by setting new password
func (s *Service) ConfirmPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated account from the request context
	_, err := GetAuthenticatedAccount(r)
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
	confirmPasswordResetRequest := new(ConfirmPasswordResetRequest)
	if err := json.Unmarshal(payload, confirmPasswordResetRequest); err != nil {
		logger.Errorf("Failed to unmarshal confirm password reset request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the reference from request URI
	vars := mux.Vars(r)
	reference := vars["reference"]

	// Fetch the password reset we want to get
	passwordReset, err := s.FindPasswordResetByReference(reference)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Reset the password
	if err := s.ResetPassword(passwordReset, confirmPasswordResetRequest.Password); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 204 no content response
	response.NoContent(w)
}
