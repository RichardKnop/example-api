package web

import (
	"net/http"
)

func (s *Service) confirmEmail(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the confirmation from the request context
	confirmation, err := getConfirmation(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Confirm the email
	if err := s.accountsService.ConfirmUser(confirmation.User); err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Render the template
	errMsg, _ := sessionService.GetFlashMessage()
	renderTemplate(w, "confirm-email.html", map[string]interface{}{
		"error": errMsg,
	})
}
