#!/bin/bash

# 1. Run database migrations
/go/bin/example-api migrate

# 2. Load fixtures
/go/bin/example-api loaddata \
  oauth/fixtures/scopes.yml \
  accounts/fixtures/roles.yml

# Finally, run the server
/go/bin/example-api runserver
