package accounts

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/example-api/oauth"
)

var (
	// ErrAccountAuthenticationRequired ...
	ErrAccountAuthenticationRequired = errors.New("Account authentication required")
	// ErrUserAuthenticationRequired ...
	ErrUserAuthenticationRequired = errors.New("User authentication required")
	// ErrAccountOrUserAuthenticationRequired ...
	ErrAccountOrUserAuthenticationRequired = errors.New("Account or user authentication required")

	errStatusCodeMap = map[error]int{
		ErrSuperuserOnlyManually:     http.StatusBadRequest,
		oauth.ErrRoleNotFound:        http.StatusBadRequest,
		oauth.ErrUsernameTaken:       http.StatusBadRequest,
		oauth.ErrUserPasswordNotSet:  http.StatusBadRequest,
		oauth.ErrInvalidUserPassword: http.StatusBadRequest,
	}
)

func getErrStatusCode(err error) int {
	code, ok := errStatusCodeMap[err]
	if ok {
		return code
	}

	return http.StatusInternalServerError
}
