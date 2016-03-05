package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/RichardKnop/recall/response"
	"github.com/RichardKnop/recall/util"
)

// Handles requests to create a new user (POST /v1/accounts/users)
func (s *Service) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated client from the request context
	authenticatedAccount, err := GetAuthenticatedAccount(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
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
	if err := json.Unmarshal(payload, userRequest); err != nil {
		logger.Errorf("Failed to unmarshal user request: %s", payload)
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if oauth user exists
	if s.GetOauthService().UserExists(userRequest.Email) {
		response.Error(w, "Email taken", http.StatusBadRequest)
		return
	}

	// Begin transaction
	tx := s.db.Begin()

	// Create a new user account
	user, err := s.CreateUserTx(tx, authenticatedAccount, userRequest)
	if err != nil {
		tx.Rollback() // rollback the transaction
		logger.Errorf("Create user error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new invitation
	confirmation := newConfirmation(user)
	if err := tx.Create(confirmation).Error; err != nil {
		tx.Rollback() // rollback the transaction
		logger.Errorf("Create confirmation error: %s", err)
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // rollback the transaction
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send confirmation email
	go func() {
		confirmationEmail := s.emailFactory.NewConfirmationEmail(confirmation)

		// Attemtp to send the confirmation email
		if err := s.emailService.Send(confirmationEmail); err != nil {
			logger.Errorf("Send email error: %s", err)
			return
		}

		// If the email was sent successfully, update the email_sent flag
		now := time.Now()
		s.db.Model(&confirmation).UpdateColumns(Confirmation{
			EmailSent:   true,
			EmailSentAt: util.TimeOrNull(&now),
		})
	}()

	// Create response
	userResponse, err := NewUserResponse(user)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set Location header to the newly created resource
	w.Header().Set("Location", fmt.Sprintf("/v1/accounts/users/%d", user.ID))
	// Write JSON response
	response.WriteJSON(w, userResponse, http.StatusCreated)
}
