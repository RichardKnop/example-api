package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/jinzhu/gorm"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	GetOauthService() oauth.ServiceInterface
	FindAccountByOauthClientID(oauthClientID uint) (*Account, error)
	FindAccountByID(accountID uint) (*Account, error)
	FindAccountByName(name string) (*Account, error)
	CreateAccount(name, description, key, secret, redirectURI string) (*Account, error)
	FindUserByOauthUserID(oauthUserID uint) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindUserByID(userID uint) (*User, error)
	FindUserByFacebookID(facebookID string) (*User, error)
	CreateUser(account *Account, userRequest *UserRequest) (*User, error)
	CreateUserTx(tx *gorm.DB, account *Account, userRequest *UserRequest) (*User, error)
	UpdateUser(user *User, userRequest *UserRequest) error
	FindConfirmationByReference(reference string) (*Confirmation, error)
	ConfirmUser(user *User) error
	FindPasswordResetByReference(reference string) (*PasswordReset, error)
	ResetPassword(passwordReset *PasswordReset, password string) error
	GetOrCreateFacebookUser(account *Account, facebookID string, userRequest *UserRequest) (*User, error)
	CreateSuperuser(account *Account, email, password string) (*User, error)
	FindInvitationByID(invitationID uint) (*Invitation, error)
	FindInvitationByReference(reference string) (*Invitation, error)
	InviteUser(invitedByUser *User, invitationRequest *InvitationRequest) (*Invitation, error)
	InviteUserTx(tx *gorm.DB, invitedByUser *User, invitationRequest *InvitationRequest) (*Invitation, error)
	ConfirmInvitation(invitation *Invitation, password string) error

	// Needed for the newRoutes to be able to register handlers
	getUserCredentials(token string) (*Account, *User, error)
	getClientCredentials(r *http.Request) (*Account, *User, error)
	createUserHandler(w http.ResponseWriter, r *http.Request)
	getMyUserHandler(w http.ResponseWriter, r *http.Request)
	getUserHandler(w http.ResponseWriter, r *http.Request)
	updateUserHandler(w http.ResponseWriter, r *http.Request)
	inviteUserHandler(w http.ResponseWriter, r *http.Request)
	createPasswordResetHandler(w http.ResponseWriter, r *http.Request)
	confirmEmailHandler(w http.ResponseWriter, r *http.Request)
	confirmInvitationHandler(w http.ResponseWriter, r *http.Request)
	confirmPasswordResetHandler(w http.ResponseWriter, r *http.Request)
}
