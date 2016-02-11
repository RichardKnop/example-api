[![Codeship Status for RichardKnop/recall](https://codeship.com/projects/eb1ee0d0-ac8c-0133-aaf6-0af8633f2b2a/status?branch=master)](https://codeship.com/projects/131678)

[![Travis Status for RichardKnop/recall](https://travis-ci.org/RichardKnop/recall.svg?branch=master)](https://travis-ci.org/RichardKnop/recall)

# Recall

Recall is a base project to bootstrap and prototype quickly. It is useful as a starting point for REST APIs and includes full OAuth 2.0 implementation as well as basic endpoints to create and update a user.

# Index

* [Recall](#recall)
* [Index](#index)
* [Docs](../../../recall/blob/master/docs/)
* [Dependencies](#dependencies)
* [Setup](#setup)
* [Test Data](#test-data)
* [Testing](#testing)
* [Docker](#docker)

# Dependencies

According to [Go 1.5 Vendor experiment](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo), all dependencies are stored in the vendor directory. This approach is called `vendoring` and is the best practice for Go projects to lock versions of dependencies in order to achieve reproducible builds.

To update dependencies during development:

```
make update-deps
```

To install dependencies:

```
make install-deps
```

# Setup

If you are developing on OSX, install `etcd`, `Postgres`:

```
brew install etcd
brew install postgres
```

You might want to create a `Postgres` database:

```
createuser --createdb recall
createdb -U area recall
```

Load a development configuration into `etcd`:

```
curl -L http://localhost:2379/v2/keys/config/recall.json -XPUT -d value='{
	"Database": {
		"Type": "postgres",
		"Host": "localhost",
		"Port": 5432,
		"User": "recall",
		"Password": "",
		"DatabaseName": "recall",
		"MaxIdleConns": 5,
		"MaxOpenConns": 5
	},
	"Oauth": {
		"AccessTokenLifetime": 3600,
		"RefreshTokenLifetime": 1209600,
		"AuthCodeLifetime": 3600
	},
	"Session": {
		"Secret": "test_secret",
		"Path": "/",
		"MaxAge": 604800,
		"HTTPOnly": true
	},
	"TrustedClient": {
		"ClientID": "test_client",
		"Secret": "test_secret"
	},
	"Sendgrid": {
		"APIKey": "sendgrid_api_key"
	},
	"Web": {
		"Scheme": "http",
		"Host": "localhost:8080"
	},
	"IsDevelopment": true
}'
```

Run migrations:

```
go run main.go migrate
```

And finally, run the app:

```
go run main.go runserver
```

When deploying, you can set `ETCD_HOST` and `ETCD_PORT` environment variables.

# Test Data

You might want to insert some test data if you are testing locally using `curl` examples from this README:

```
go run main.go loaddata \
	oauth/fixtures/scopes.yml \
	oauth/fixtures/test_clients.yml \
	oauth/fixtures/test_users.yml
```

# Testing

I have used a mix of unit and functional tests so you need to have `sqlite` installed in order for the tests to run successfully as the suite creates an in-memory database.

To run tests:

```
make test
```

# Docker

Build a Docker image and run the app in a container:

```
docker build -t recall .
docker run -e ETCD_HOST=localhost -e ETCD_PORT=2379 -p 6060:8080 recall
```

You can load fixtures with `docker exec` command:

```
docker exec <container_id> /go/bin/recall loaddata \
	oauth/fixtures/scopes.yml \
	oauth/fixtures/test_clients.yml
```

You can also execute interactive commands by passing `-i` flag:

```
docker exec -i <container_id> /go/bin/area-api createaccount
docker exec -i <container_id> /go/bin/area-api createsuperuser
```
