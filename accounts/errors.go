package accounts

import (
	"net/http"
)

var (
	errStatusCodeMap = map[error]int{
		ErrSuperuserOnlyManually: http.StatusBadRequest,
		ErrEmailTaken:            http.StatusBadRequest,
		ErrEmailCannotBeChanged:  http.StatusBadRequest,
	}
)
