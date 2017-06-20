package accounts

import (
	"github.com/RichardKnop/example-api/util/routes"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	usersResource          = "users"
	usersPath              = "/" + usersResource
	meResource             = "me"
	mePath                 = "/" + meResource
	invitationsResource    = "invitations"
	invitationsPath        = "/" + invitationsResource
	confirmationsResource  = "confirmations"
	confirmationsPath      = "/" + confirmationsResource
	passwordResetsResource = "password-resets"
	passwordResetsPath     = "/" + passwordResetsResource
)

// RegisterRoutes registers route handlers for the accounts service
func (s *Service) RegisterRoutes(router *mux.Router, prefix string) {
	subRouter := router.PathPrefix(prefix).Subrouter()
	routes.AddRoutes(s.GetRoutes(), subRouter)
}

// GetRoutes returns []routes.Route slice for the accounts service
func (s *Service) GetRoutes() []routes.Route {
	return []routes.Route{
		{
			Name:        "create_user",
			Method:      "POST",
			Pattern:     usersPath,
			HandlerFunc: s.createUserHandler,
			Middlewares: []negroni.Handler{
				NewClientAuthMiddleware(s),
			},
		},
		{
			Name:        "get_my_user",
			Method:      "GET",
			Pattern:     mePath,
			HandlerFunc: s.getMyUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(s),
			},
		},
		{
			Name:        "get_user",
			Method:      "GET",
			Pattern:     usersPath + "/{id:[0-9]+}",
			HandlerFunc: s.getUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(s),
			},
		},
		{
			Name:        "update_user",
			Method:      "PUT",
			Pattern:     usersPath + "/{id:[0-9]+}",
			HandlerFunc: s.updateUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(s),
			},
		},
		{
			Name:        "list_users",
			Method:      "GET",
			Pattern:     usersPath,
			HandlerFunc: s.listUsersHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(s),
			},
		},
		{
			Name:        "invite_user",
			Method:      "POST",
			Pattern:     invitationsPath,
			HandlerFunc: s.inviteUserHandler,
			Middlewares: []negroni.Handler{
				NewUserAuthMiddleware(s),
			},
		},
		{
			Name:        "confirm_email",
			Method:      "GET",
			Pattern:     confirmationsPath + "/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}",
			HandlerFunc: s.confirmEmailHandler,
			Middlewares: []negroni.Handler{
				NewClientAuthMiddleware(s),
			},
		},
		{
			Name:        "confirm_invitation",
			Method:      "POST",
			Pattern:     invitationsPath + "/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}",
			HandlerFunc: s.confirmInvitationHandler,
			Middlewares: []negroni.Handler{
				NewClientAuthMiddleware(s),
			},
		},
		{
			Name:        "create_password_reset",
			Method:      "POST",
			Pattern:     passwordResetsPath,
			HandlerFunc: s.createPasswordResetHandler,
			Middlewares: []negroni.Handler{
				NewClientAuthMiddleware(s),
			},
		},
		{
			Name:        "confirm_password_reset",
			Method:      "POST",
			Pattern:     passwordResetsPath + "/{reference:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}",
			HandlerFunc: s.confirmPasswordResetHandler,
			Middlewares: []negroni.Handler{
				NewClientAuthMiddleware(s),
			},
		},
	}
}
