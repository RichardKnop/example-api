package accounts

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrSuperuserOnlyManually: http.StatusBadRequest,
		ErrRoleNotFound:          http.StatusBadRequest,
		ErrEmailTaken:            http.StatusBadRequest,
	}
)
