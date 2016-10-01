package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/routes"
	"github.com/gorilla/mux"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	GetOauthService() oauth.ServiceInterface
	GetRoutes() []routes.Route
	RegisterRoutes(router *mux.Router, prefix string)
	GetUserCredentialsFromToken(token string) (*User, error)
	GetClientCredentialsFromBaseAuth(r *http.Request) (*Account, error)
	GetClientCredentialsFromToken(token string) (*Account, error)
	GetMixedCredentialsFromToken(token string) (*Account, *User, error)
	FindAccountByOauthClientID(oauthClientID uint) (*Account, error)
	FindAccountByID(accountID uint) (*Account, error)
	FindAccountByName(name string) (*Account, error)
	CreateAccount(name, description, key, secret, redirectURI string) (*Account, error)
	FindUserByOauthUserID(oauthUserID uint) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindUserByID(userID uint) (*User, error)
	FindUserByFacebookID(facebookID string) (*User, error)
	CreateUser(account *Account, userRequest *UserRequest) (*User, error)
	UpdateUser(user *User, userRequest *UserRequest) error
	PaginatedUsersCount() (int, error)
	FindPaginatedUsers(offset, limit int, sorts map[string]string) ([]*User, error)
	FindConfirmationByReference(reference string) (*Confirmation, error)
	ConfirmUser(confirmation *Confirmation) error
	FindPasswordResetByReference(reference string) (*PasswordReset, error)
	ResetPassword(passwordReset *PasswordReset, password string) error
	GetOrCreateFacebookUser(account *Account, facebookID string, userRequest *UserRequest) (*User, error)
	CreateSuperuser(account *Account, email, password string) (*User, error)
	FindInvitationByReference(reference string) (*Invitation, error)
	InviteUser(invitedByUser *User, invitationRequest *InvitationRequest) (*Invitation, error)
	ConfirmInvitation(invitation *Invitation, password string) error
	GetUserFromQueryString(r *http.Request) (*User, error)
}
