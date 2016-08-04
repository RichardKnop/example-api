package facebook

import (
	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/routes"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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
			HandlerFunc: service.LoginHandler,
			Middlewares: []negroni.Handler{
				accounts.NewAccountAuthMiddleware(service.GetAccountsService()),
			},
		},
	}
}
