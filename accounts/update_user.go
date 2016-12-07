package accounts

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/RichardKnop/example-api/logger"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/util/response"
	"github.com/gorilla/mux"
)

var (
	// ErrUpdateUserPermission ...
	ErrUpdateUserPermission = errors.New("Need permission to update user")
	// ErrUpdateUserMustBeSuperuser ...
	ErrUpdateUserMustBeSuperuser = errors.New("You must be a superuser to update these fields")
)

// Handles requests to update a user
// PUT /v1/users/{id:[0-9]+}
func (s *Service) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user from the request context
	authenticatedUser, err := GetAuthenticatedUser(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// Get the id from request URI and type assert it
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user we want to update
	user, err := s.FindUserByID(uint(userID))
	if err != nil {
		response.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Request body cannot be nil
	if r.Body == nil {
		response.Error(w, "Request body cannot be nil", http.StatusBadRequest)
		return
	}

	// Read the request body
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into the request prototype
	userRequest := new(UserRequest)
	if err = json.Unmarshal(payload, userRequest); err != nil {
		logger.ERROR.Printf("Failed to unmarshal user request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check permissions
	if err = checkUpdateUserPermissions(authenticatedUser, user, userRequest); err != nil {
		response.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Update the user
	if err = s.UpdateUser(user, userRequest); err != nil {
		logger.ERROR.Printf("Update user error: %s", err)
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// Create response
	userResponse, err := NewUserResponse(user)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write JSON response
	response.WriteJSON(w, userResponse, http.StatusOK)
}

func checkUpdateUserPermissions(authenticatedUser, user *models.User, req *UserRequest) error {
	// Superusers can update any users
	if authenticatedUser.OauthUser.RoleID.String == roles.Superuser {
		return nil
	}

	// Users can update their own accounts
	if authenticatedUser.ID == user.ID {
		return nil
	}

	// The user doesn't have the permission
	return ErrUpdateUserPermission
}
