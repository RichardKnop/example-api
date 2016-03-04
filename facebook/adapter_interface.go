package facebook

import (
	fb "github.com/huandu/facebook"
	"golang.org/x/oauth2"
)

// AdapterInterface defines exported methods
type AdapterInterface interface {
	// Exported methods
	AuthCodeURL(state string) string
	Exchange(code string) (*oauth2.Token, error)
	GetMe(accessToken string) (fb.Result, error)
}
