#!/bin/bash
set -e

# Migrate database and load fixtures
/go/bin/example-api migrate
/go/bin/example-api loaddata \
  oauth/fixtures/scopes.yml \
  accounts/fixtures/roles.yml

# It is important that if your image uses such a script,
# that script should use exec so that the script’s process
# is replaced by your software. If you do not use exec,
# then signals sent by docker will go to your wrapper script
# instead of your software’s process.
exec /go/bin/example-api runserver

exec "$@"
