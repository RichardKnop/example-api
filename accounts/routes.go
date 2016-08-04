package accounts

import (
	"github.com/RichardKnop/example-api/routes"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// RegisterRoutes registers route handlers for the accounts service
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/v1/accounts").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the accounts service
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "create_user",
			Method:      "POST",
			Pattern:     "/users",
			HandlerFunc: service.CreateUserHandler,
			Middlewares: []negroni.Handler{
				NewAccountAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "get_my_user",
			Method:      "GET",
			Pattern:     "/me",
			HandlerFunc: service.GetMyUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "get_user",
			Method:      "GET",
			Pattern:     "/users/{id:[0-9]+}",
			HandlerFunc: service.GetUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "update_user",
			Method:      "PUT",
			Pattern:     "/users/{id:[0-9]+}",
			HandlerFunc: service.UpdateUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "invite_user",
			Method:      "POST",
			Pattern:     "/invitations",
			HandlerFunc: service.InviteUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "create_password_reset",
			Method:      "POST",
			Pattern:     "/password-reset",
			HandlerFunc: service.CreatePasswordResetHandler,
			Middlewares: []negroni.Handler{
				NewAccountAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "confirm_email",
			Method:      "GET",
			Pattern:     "/confirmations/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}",
			HandlerFunc: service.ConfirmEmailHandler,
			Middlewares: []negroni.Handler{
				NewAccountAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "confirm_invitation",
			Method:      "POST",
			Pattern:     "/invitations/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}",
			HandlerFunc: service.ConfirmInvitationHandler,
			Middlewares: []negroni.Handler{
				NewAccountAuthMiddleware(service),
			},
		},
		routes.Route{
			Name:        "confirm_password_reset",
			Method:      "POST",
			Pattern:     "/password-resets/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}",
			HandlerFunc: service.ConfirmPasswordResetHandler,
			Middlewares: []negroni.Handler{
				NewAccountAuthMiddleware(service),
			},
		},
	}
}
