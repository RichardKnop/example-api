package facebook

import (
	"fmt"

	"github.com/RichardKnop/recall/config"
	fb "github.com/huandu/facebook"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

// Adapter ...
type Adapter struct {
	oauth2Cnf *oauth2.Config
}

// NewAdapter starts a new Adapter instance
func NewAdapter(cnf *config.Config) *Adapter {
	return &Adapter{
		oauth2Cnf: &oauth2.Config{
			ClientID:     cnf.Facebook.AppID,
			ClientSecret: cnf.Facebook.AppSecret,
			RedirectURL: fmt.Sprintf(
				"%s://%s/v1/facebook/redirect",
				cnf.Web.Scheme,
				cnf.Web.Host,
			),
			Scopes: []string{
				"public_profile",
				"email",
			},
			Endpoint: facebook.Endpoint, // https://godoc.org/golang.org/x/oauth2/facebook
		},
	}
}

// AuthCodeURL generates an authorisation URL
func (a *Adapter) AuthCodeURL(state string) string {
	return a.oauth2Cnf.AuthCodeURL(state)
}

// Exchange exchanges auth code for an access token
func (a *Adapter) Exchange(code string) (*oauth2.Token, error) {
	return a.oauth2Cnf.Exchange(nil, code)
}

// GetMe returns user profile data from facebook
func (a *Adapter) GetMe(accessToken string) (fb.Result, error) {
	return fb.Get("/me", fb.Params{
		"fields": []string{
			"id",
			"first_name",
			"last_name",
			"email",
		},
		"access_token": accessToken,
	})
}
