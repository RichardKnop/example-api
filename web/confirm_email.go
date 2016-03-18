package web

import (
	"net/http"
)

func (s *Service) confirmEmail(w http.ResponseWriter, r *http.Request) {
	// Get the confirmation from the request context
	confirmation, err := getConfirmation(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Confirm the email
	var confirmationErr string
	if err := s.accountsService.ConfirmUser(confirmation.User); err != nil {
		confirmationErr = err.Error()
	}

	// Render the template
	renderTemplate(w, "confirm-email.html", map[string]interface{}{
		"error":      confirmationErr,
		"hideLogout": true,
	})
}
