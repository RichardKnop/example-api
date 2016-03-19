package web

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/recall/accounts"
	"github.com/RichardKnop/recall/oauth"
	"github.com/RichardKnop/recall/session"
	"github.com/gorilla/context"
)

type contextKey int

const (
	sessionServiceKey contextKey = 0
	clientKey         contextKey = 1
	confirmationKey   contextKey = 2
	passwordResetKey  contextKey = 3
)

var (
	// ErrSessionServiceNotPresent ...
	ErrSessionServiceNotPresent = errors.New("Session service not present in the request context")
	// ErrClientNotPresent ...
	ErrClientNotPresent = errors.New("Client not present in the request context")
	// ErrConfirmationNotPresent ...
	ErrConfirmationNotPresent = errors.New("Confirmation not present in the request context")
	// ErrPasswordResetNotPresent ...
	ErrPasswordResetNotPresent = errors.New("Password reset not present in the request context")
)

// Returns *session.Service from the request context
func getSessionService(r *http.Request) (session.ServiceInterface, error) {
	val, ok := context.GetOk(r, sessionServiceKey)
	if !ok {
		return nil, ErrSessionServiceNotPresent
	}

	sessionService, ok := val.(session.ServiceInterface)
	if !ok {
		return nil, ErrSessionServiceNotPresent
	}

	return sessionService, nil
}

// Returns *oauth.Client from the request context
func getClient(r *http.Request) (*oauth.Client, error) {
	val, ok := context.GetOk(r, clientKey)
	if !ok {
		return nil, ErrClientNotPresent
	}

	client, ok := val.(*oauth.Client)
	if !ok {
		return nil, ErrClientNotPresent
	}

	return client, nil
}

// Returns *accounts.Confirmation from the request context
func getConfirmation(r *http.Request) (*accounts.Confirmation, error) {
	val, ok := context.GetOk(r, confirmationKey)
	if !ok {
		return nil, ErrConfirmationNotPresent
	}

	confirmation, ok := val.(*accounts.Confirmation)
	if !ok {
		return nil, ErrConfirmationNotPresent
	}

	return confirmation, nil
}

// Returns *accounts.PasswordReset from the request context
func getPasswordReset(r *http.Request) (*accounts.PasswordReset, error) {
	val, ok := context.GetOk(r, passwordResetKey)
	if !ok {
		return nil, ErrPasswordResetNotPresent
	}

	passwordReset, ok := val.(*accounts.PasswordReset)
	if !ok {
		return nil, ErrPasswordResetNotPresent
	}

	return passwordReset, nil
}
