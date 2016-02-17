package web

import (
	"net/http"

	"github.com/RichardKnop/recall/accounts"
)

func (s *Service) registerForm(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	errMsg, _ := sessionService.GetFlashMessage()
	renderTemplate(w, "register.html", map[string]interface{}{
		"error":       errMsg,
		"queryString": getQueryString(r.URL.Query()),
	})
}

func (s *Service) register(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the client from the request context
	client, err := getClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check that the submitted email hasn't been registered already
	if s.GetAccountsService().GetOauthService().UserExists(r.Form.Get("email")) {
		sessionService.SetFlashMessage("Email taken")
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Fetch the account based on oauth client
	account, err := s.GetAccountsService().FindAccountByOauthClientID(client.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a user
	_, err = s.GetAccountsService().CreateUser(
		account,
		&accounts.UserRequest{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		},
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Redirect to the login page
	redirectWithQueryString("/web/login", r.URL.Query(), w, r)
}
