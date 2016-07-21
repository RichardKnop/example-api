package facebook

import (
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/routes"
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the accounts service
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/v1/facebook").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the accounts service
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "facebook_login",
			Method:      "POST",
			Pattern:     "/login",
			HandlerFunc: service.loginHandler,
			Middlewares: []negroni.Handler{
				accounts.NewAccountAuthMiddleware(service.GetAccountsService()),
			},
		},
	}
}
