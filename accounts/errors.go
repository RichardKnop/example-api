package accounts

import (
	"net/http"

	"github.com/RichardKnop/recall/oauth"
)

var (
	errStatusCodeMap = map[error]int{
		ErrSuperuserOnlyManually: http.StatusBadRequest,
		ErrRoleNotFound:          http.StatusBadRequest,
		oauth.ErrUsernameTaken:   http.StatusBadRequest,
	}
)
