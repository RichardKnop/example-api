package accounts

import (
	"fmt"
	"time"

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

// NewUserResponse creates new UserResponse instance
func NewUserResponse(user *User) (*UserResponse, error) {
	response := &UserResponse{
		ID:        user.ID,
		Email:     user.OauthUser.Username,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Role:      user.RoleID.String,
		Confirmed: user.Confirmed,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
	}

	// Set the self link
	response.SetLink(
		"self", // name
		fmt.Sprintf("/v1/accounts/users/%d", user.ID), // href
		"", // title
	)

	return response, nil
}
