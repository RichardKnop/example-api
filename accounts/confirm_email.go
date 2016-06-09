package accounts

import (
	"net/http"

	"github.com/RichardKnop/recall/response"

	"github.com/gorilla/mux"
)

// Handles requests to confirm user's email based on a reference string
func (s *Service) confirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated account from the request context
	_, err := GetAuthenticatedAccount(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the reference from request URI
	vars := mux.Vars(r)
	reference := vars["reference"]

	// Fetch the confirmation we want to get
	confirmation, err := s.FindConfirmationByReference(reference)
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Confirm the user
	if !confirmation.User.Confirmed {
		if err := s.ConfirmUser(confirmation.User); err != nil {
			response.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 204 no content response
	response.NoContent(w)
}
