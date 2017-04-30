#!/bin/sh

set -e

# to make sure etcd is ready (election ended and leader elected)
while ! etcdctl endpoint health &>/dev/null; do :; done

exec etcdctl put /config/example_api.json '{
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
  "Mailgun": {
    "Domain": "example.com",
    "APIKey": "mailgun_api_key",
    "PublicAPIKey": "mailgun_public_api_key"
  },
  "Web": {
    "Scheme": "http",
    "Host": "localhost:8080",
    "AppScheme": "http",
    "AppHost": "localhost:8000"
  },
  "AppSpecific": {
    "ConfirmationLifetime": 604800,
    "InvitationLifetime": 604800,
    "PasswordResetLifetime": 604800,
    "CompanyName": "Example Ltd",
    "CompanyNoreplyEmail": "noreply@example.com",
    "ConfirmationURLFormat": "%s://%s/confirm-email/%s",
    "InvitationURLFormat": "%s://%s/confirm-invitation/%s",
    "PasswordResetURLFormat": "%s://%s/reset-password/%s"
  },
  "IsDevelopment": true
}'
