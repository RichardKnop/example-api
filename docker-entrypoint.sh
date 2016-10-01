#!/bin/bash
set -e

if [ "$1" = 'example-api' ] && [ "$2" = 'runserver' ]; then
  $1 migrate
  $1 loaddata oauth/fixtures/scopes.yml
  $1 loaddata oauth/fixtures/roles.yml
fi

exec "$@"
