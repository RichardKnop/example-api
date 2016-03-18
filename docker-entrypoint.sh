#!/bin/bash

# 1. Run database migrations
/go/bin/recall migrate

# 2. Load fixtures
/go/bin/recall loaddata \
  oauth/fixtures/scopes.yml \
  accounts/fixtures/roles.yml

# Finally, run the server
/go/bin/recall runserver
