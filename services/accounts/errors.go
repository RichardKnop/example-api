package accounts

import (
	"net/http"

	"github.com/RichardKnop/example-api/services/oauth"
)

var (
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
