package web

import (
	"github.com/RichardKnop/recall/routes"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the web package
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/web").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the web package
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "register_form",
			Method:      "GET",
			Pattern:     "/register",
			HandlerFunc: service.registerForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "register",
			Method:      "POST",
			Pattern:     "/register",
			HandlerFunc: service.register,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "confirm_email",
			Method:      "GET",
			Pattern:     "/confirm-email/{reference}",
			HandlerFunc: service.confirmEmail,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newConfirmationMiddleware(service),
			},
		},
		routes.Route{
			Name:        "login_form",
			Method:      "GET",
			Pattern:     "/login",
			HandlerFunc: service.loginForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "login",
			Method:      "POST",
			Pattern:     "/login",
			HandlerFunc: service.login,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "logout",
			Method:      "GET",
			Pattern:     "/logout",
			HandlerFunc: service.logout,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(service),
			},
		},
		routes.Route{
			Name:        "password_reset_form",
			Method:      "GET",
			Pattern:     "/password-reset/{reference}",
			HandlerFunc: service.passwordResetForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newPasswordResetMiddleware(service),
			},
		},
		routes.Route{
			Name:        "password_reset",
			Method:      "POST",
			Pattern:     "/password-reset/{reference}",
			HandlerFunc: service.passwordReset,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newPasswordResetMiddleware(service),
			},
		},
		routes.Route{
			Name:        "password_reset_success",
			Method:      "GET",
			Pattern:     "/password-reset-success",
			HandlerFunc: service.passwordResetSuccess,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
			},
		},
		routes.Route{
			Name:        "authorize_form",
			Method:      "GET",
			Pattern:     "/authorize",
			HandlerFunc: service.authorizeForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "authorize",
			Method:      "POST",
			Pattern:     "/authorize",
			HandlerFunc: service.authorize,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(service),
				newClientMiddleware(service),
			},
		},
	}
}
