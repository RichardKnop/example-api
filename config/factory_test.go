package config_test

import (
	"testing"

	"github.com/RichardKnop/example-api/config"
	"github.com/stretchr/testify/assert"
)

type Foo struct {
	cnf *config.Config
}

func TestConfigReloading(t *testing.T) {
	config.Cnf.Oauth.AuthCodeLifetime = 123
	foo := &Foo{cnf: config.Cnf}
	assert.Equal(t, 123, foo.cnf.Oauth.AuthCodeLifetime)
	newCnf := &config.Config{Oauth: config.OauthConfig{AuthCodeLifetime: 9999}}
	assert.Equal(t, 123, foo.cnf.Oauth.AuthCodeLifetime)
	config.RefreshConfig(newCnf)
	assert.Equal(t, 9999, foo.cnf.Oauth.AuthCodeLifetime)
}
