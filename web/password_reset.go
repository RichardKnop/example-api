package web

import (
	"net/http"
)

func (s *Service) passwordResetForm(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the password reset from the request context
	_, err = getPasswordReset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Render the template
	errMsg, _ := sessionService.GetFlashMessage()
	renderTemplate(w, "password-reset.html", map[string]interface{}{
		"error":       errMsg,
		"queryString": getQueryString(r.URL.Query()),
	})
}

func (s *Service) passwordReset(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the password reset from the request context
	passwordReset, err := getPasswordReset(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check that both submitted passwords are the same
	if r.Form.Get("password") != r.Form.Get("password2") {
		sessionService.SetFlashMessage("Passwords are not the same")
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Reset the password
	err = s.GetAccountsService().ResetPassword(passwordReset, r.Form.Get("password"))
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Redirect to the success page
	redirectWithQueryString("/web/password-reset-success", r.URL.Query(), w, r)
}
