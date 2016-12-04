package accounts

import (
	"net/http"
	"strconv"

	"github.com/RichardKnop/example-api/models"
)

// GetUserFromQueryString parses user_id from the query string and
// returns a matching *User instance or an error
func (s *Service) GetUserFromQueryString(r *http.Request) (*models.User, error) {
	// If no user_id query string parameter found, just return
	if r.URL.Query().Get("user_id") == "" {
		return nil, nil
	}

	// If the user_id is present in the query string, try to convert it to int
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return nil, err
	}

	// Get the authenticated user from the request context if present
	authenticatedUser, _ := GetAuthenticatedUser(r)

	// If the user ID matches the authenticated user, just return it
	if authenticatedUser != nil && uint(userID) == authenticatedUser.ID {
		return authenticatedUser, nil
	}

	// Fetch the user from the database
	user, err := s.FindUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
