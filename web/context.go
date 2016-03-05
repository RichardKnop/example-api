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
)

var (
	errSessionServiceNotPresent = errors.New("Session service not present in the request context")
	errClientNotPresent         = errors.New("Client not present in the request context")
	errConfirmationNotPresent   = errors.New("Confirmation not present in the request context")
)

// Returns *session.Service from the request context
func getSessionService(r *http.Request) (session.ServiceInterface, error) {
	val, ok := context.GetOk(r, sessionServiceKey)
	if !ok {
		return nil, errSessionServiceNotPresent
	}

	sessionService, ok := val.(session.ServiceInterface)
	if !ok {
		return nil, errSessionServiceNotPresent
	}

	return sessionService, nil
}

// Returns *oauth.Client from the request context
func getClient(r *http.Request) (*oauth.Client, error) {
	val, ok := context.GetOk(r, clientKey)
	if !ok {
		return nil, errClientNotPresent
	}

	client, ok := val.(*oauth.Client)
	if !ok {
		return nil, errClientNotPresent
	}

	return client, nil
}

// Returns *accounts.Confirmation from the request context
func getConfirmation(r *http.Request) (*accounts.Confirmation, error) {
	val, ok := context.GetOk(r, confirmationKey)
	if !ok {
		return nil, errConfirmationNotPresent
	}

	confirmation, ok := val.(*accounts.Confirmation)
	if !ok {
		return nil, errConfirmationNotPresent
	}

	return confirmation, nil
}
