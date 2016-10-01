package health

import (
	"github.com/gorilla/mux"
	"github.com/RichardKnop/example-api/routes"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetRoutes() []routes.Route
	RegisterRoutes(router *mux.Router, prefix string)
}
