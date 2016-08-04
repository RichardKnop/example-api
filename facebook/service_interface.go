package facebook

import (
	"net/http"

	"github.com/RichardKnop/example-api/accounts"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetAccountsService() accounts.ServiceInterface

	// Needed for the NewRoutes to be able to register handlers
	LoginHandler(w http.ResponseWriter, r *http.Request)
}
