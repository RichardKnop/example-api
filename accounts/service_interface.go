package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/util/routes"
	"github.com/gorilla/mux"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	GetOauthService() oauth.ServiceInterface
	GetRoutes() []routes.Route
	RegisterRoutes(router *mux.Router, prefix string)
	GetUserCredentialsFromToken(token string) (*models.User, error)
	GetClientCredentialsFromBaseAuth(r *http.Request) (*models.OauthClient, error)
	GetClientCredentialsFromToken(token string) (*models.OauthClient, error)
	GetMixedCredentialsFromToken(token string) (*models.OauthClient, *models.User, error)
	FindUserByOauthUserID(oauthUserID uint) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(userID uint) (*models.User, error)
	FindUserByFacebookID(facebookID string) (*models.User, error)
	CreateUser(oauthClient *models.OauthClient, userRequest *UserRequest) (*models.User, error)
	UpdateUser(user *models.User, userRequest *UserRequest) error
	PaginatedUsersCount() (int, error)
	FindPaginatedUsers(offset, limit int, sorts map[string]string) ([]*models.User, error)
	FindConfirmationByReference(reference string) (*models.Confirmation, error)
	ConfirmUser(confirmation *models.Confirmation) error
	FindPasswordResetByReference(reference string) (*models.PasswordReset, error)
	ResetPassword(passwordReset *models.PasswordReset, password string) error
	GetOrCreateFacebookUser(oauthClient *models.OauthClient, facebookID string, userRequest *UserRequest) (*models.User, error)
	CreateSuperuser(oauthClient *models.OauthClient, email, password string) (*models.User, error)
	FindInvitationByReference(reference string) (*models.Invitation, error)
	InviteUser(invitedByUser *models.User, invitationRequest *InvitationRequest) (*models.Invitation, error)
	ConfirmInvitation(invitation *models.Invitation, password string) error
	GetUserFromQueryString(r *http.Request) (*models.User, error)
}
