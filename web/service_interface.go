package web

import (
	"net/http"

	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/config"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	GetAccountsService() accounts.ServiceInterface

	// Needed for the newRoutes to be able to register handlers
	authorizeForm(w http.ResponseWriter, r *http.Request)
	authorize(w http.ResponseWriter, r *http.Request)
	loginForm(w http.ResponseWriter, r *http.Request)
	login(w http.ResponseWriter, r *http.Request)
	logout(w http.ResponseWriter, r *http.Request)
	registerForm(w http.ResponseWriter, r *http.Request)
	register(w http.ResponseWriter, r *http.Request)
}
