package accounts

import (
	"fmt"

	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/example-api/util/response"
	"github.com/RichardKnop/jsonhal"
)

// UserResponse ...
type UserResponse struct {
	jsonhal.Hal
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Confirmed bool   `json:"confirmed"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// EmailTokenResponse is an abstract class for common responses to objects
// derived from EmailTokenModel
type EmailTokenResponse struct {
	jsonhal.Hal
	ID          uint   `json:"id"`
	Reference   string `json:"reference"`
	EmailSent   bool   `json:"email_sent"`
	EmailSentAt string `json:"email_sent_at"`
	ExpiresAt   string `json:"expires_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ConfirmationResponse ...
type ConfirmationResponse struct {
	EmailTokenResponse
	UserID uint `json:"user_id"`
}

// InvitationResponse ...
type InvitationResponse struct {
	EmailTokenResponse
	InvitedUserID   uint `json:"invited_user_id"`
	InvitedByUserID uint `json:"invited_by_user_id"`
}

// PasswordResetResponse ...
type PasswordResetResponse struct {
	EmailTokenResponse
	UserID uint `json:"user_id"`
}

// NewUserResponse creates new UserResponse instance
func NewUserResponse(o *models.User) (*UserResponse, error) {
	response := &UserResponse{
		ID:        o.ID,
		Email:     o.OauthUser.Username,
		FirstName: o.FirstName.String,
		LastName:  o.LastName.String,
		Role:      o.OauthUser.RoleID.String,
		Confirmed: o.Confirmed,
		CreatedAt: util.FormatTime(&o.CreatedAt),
		UpdatedAt: util.FormatTime(&o.UpdatedAt),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/users/%d", o.ID), // href
		"", // title
	)

	return response, nil
}

// NewConfirmationResponse creates new ConfirmationResponse instance
func NewConfirmationResponse(o *models.Confirmation) (*ConfirmationResponse, error) {
	response := &ConfirmationResponse{
		EmailTokenResponse: EmailTokenResponse{
			ID:          o.ID,
			Reference:   o.Reference,
			EmailSent:   o.EmailSent,
			EmailSentAt: util.FormatTime(o.EmailSentAt),
			ExpiresAt:   util.FormatTime(&o.ExpiresAt),
			CreatedAt:   util.FormatTime(&o.CreatedAt),
			UpdatedAt:   util.FormatTime(&o.UpdatedAt),
		},
		UserID: o.User.ID,
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/confirmations/%d", o.ID), // href
		"", // title
	)

	return response, nil
}

// NewInvitationResponse creates new InvitationResponse instance
func NewInvitationResponse(o *models.Invitation) (*InvitationResponse, error) {
	response := &InvitationResponse{
		EmailTokenResponse: EmailTokenResponse{
			ID:          o.ID,
			Reference:   o.Reference,
			EmailSent:   o.EmailSent,
			EmailSentAt: util.FormatTime(o.EmailSentAt),
			ExpiresAt:   util.FormatTime(&o.ExpiresAt),
			CreatedAt:   util.FormatTime(&o.CreatedAt),
			UpdatedAt:   util.FormatTime(&o.UpdatedAt),
		},
		InvitedUserID:   o.InvitedUser.ID,
		InvitedByUserID: o.InvitedByUser.ID,
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/invitations/%d", o.ID), // href
		"", // title
	)

	return response, nil
}

// NewPasswordResetResponse creates new PasswordResetResponse instance
func NewPasswordResetResponse(o *models.PasswordReset) (*PasswordResetResponse, error) {
	response := &PasswordResetResponse{
		EmailTokenResponse: EmailTokenResponse{
			ID:          o.ID,
			Reference:   o.Reference,
			EmailSent:   o.EmailSent,
			EmailSentAt: util.FormatTime(o.EmailSentAt),
			ExpiresAt:   util.FormatTime(&o.ExpiresAt),
			CreatedAt:   util.FormatTime(&o.CreatedAt),
			UpdatedAt:   util.FormatTime(&o.UpdatedAt),
		},
		UserID: o.User.ID,
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/password-resets/%d", o.ID), // href
		"", // title
	)

	return response, nil
}

// NewListUsersResponse creates new *response.ListResponse instance
func NewListUsersResponse(count, page int, self, first, last, previous, next string, users []*models.User) (*response.ListResponse, error) {
	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponse, err := NewUserResponse(user)
		if err != nil {
			return nil, err
		}
		userResponses[i] = userResponse
	}

	listResponse := response.NewListResponse(
		count, page, self, first, last, previous, next, "users", userResponses,
	)

	return listResponse, nil
}
