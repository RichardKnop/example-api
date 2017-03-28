package facebook

import (
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/util/routes"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	loginResource = "login"
	loginPath     = "/" + loginResource
)

// RegisterRoutes registers route handlers for the facebook service
func (s *Service) RegisterRoutes(router *mux.Router, prefix string) {
	subRouter := router.PathPrefix(prefix).Subrouter()
	routes.AddRoutes(s.GetRoutes(), subRouter)
}

// GetRoutes returns []routes.Route slice for the facebook service
func (s *Service) GetRoutes() []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "facebook_login",
			Method:      "POST",
			Pattern:     loginPath,
			HandlerFunc: s.loginHandler,
			Middlewares: []negroni.Handler{
				accounts.NewClientAuthMiddleware(s.GetAccountsService()),
			},
		},
	}
}
