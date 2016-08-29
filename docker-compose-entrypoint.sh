#!/bin/bash
set -e

if [ "$1" = 'example-api' ] && [ "$2" = 'runserver' ]; then
  curl -L http://etcd:2379/v2/keys/config/example_api.json -XPUT -d value='{
    "Database": {
      "Type": "postgres",
      "Host": "postgres",
      "Port": 5432,
      "User": "example_api",
      "Password": "",
      "DatabaseName": "example_api",
      "MaxIdleConns": 5,
      "MaxOpenConns": 5
    },
    "Oauth": {
      "AccessTokenLifetime": 3600,
      "RefreshTokenLifetime": 1209600,
      "AuthCodeLifetime": 3600
    },
    "Facebook": {
  		"AppID": "facebook_app_id",
  		"AppSecret": "facebook_app_secret"
  	},
  	"Sendgrid": {
  		APIKey: "sendgrid_api_key"
  	},
    "Web": {
      "Scheme": "http",
      "Host": "localhost:8080",
      "AppScheme": "http",
      "AppHost": "localhost:8000"
    },
    "AppSpecific": {
      "PasswordResetLifetime": 604800,
      "CompanyName": "Your company",
      "CompanyNoreplyEmail": "noreply@example.com"
    },
    "IsDevelopment": true
  }'

  $1 migrate
  $1 loaddata oauth/fixtures/scopes.yml
  $1 loaddata oauth/fixtures/test_clients.yml
  $1 loaddata accounts/fixtures/roles.yml
fi

exec "$@"
