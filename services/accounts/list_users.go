package accounts

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/oauth/roles"
	"github.com/RichardKnop/example-api/util"
	"github.com/RichardKnop/example-api/util/pagination"
	"github.com/RichardKnop/example-api/util/response"
)

var (
	// ErrListUsersPermission ...
	ErrListUsersPermission = errors.New("Need permission to list users")
)

// Handles calls to list user accounts
// GET /v1/users
func (s *Service) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Check permissions
	if err = checkListUsersPermissions(authenticatedUser); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Get pagination params
	page, limit, sorts, err := pagination.GetParams(r, []string{"id"})
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Count total number of results
	count, err := s.PaginatedUsersCount()
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get pagination links
	first, last, previous, next, err := pagination.GetLinks(r.URL, count, page, limit)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get paginated results
	users, err := s.FindPaginatedUsers(
		pagination.GetOffsetForPage(count, page, limit),
		limit,
		sorts,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create response
	self := util.GetCurrentURL(r)
	listUsersResponse, err := NewListUsersResponse(
		count, page,
		self, first, last, previous, next,
		users,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, listUsersResponse, http.StatusOK)
}

func checkListUsersPermissions(authenticatedUser *models.User) error {
	// Superusers can list users
	if authenticatedUser.OauthUser.RoleID.String == roles.Superuser {
		return nil
	}

	// The user doesn't have the permission
	return ErrListUsersPermission
}
