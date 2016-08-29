[![Codeship Status for RichardKnop/example-api](https://codeship.com/projects/eb1ee0d0-ac8c-0133-aaf6-0af8633f2b2a/status?branch=master)](https://codeship.com/projects/131678)

[![Travis Status for RichardKnop/example-api](https://travis-ci.org/RichardKnop/example-api.svg?branch=master)](https://travis-ci.org/RichardKnop/example-api)

# Example API

This is a base project to bootstrap and prototype quickly. It is useful as a starting point for REST APIs and includes full OAuth 2.0 implementation as well as basic endpoints to create and update a user, health check endpoint, Facebook integration, migrations and a ready to rumble Dockerfile.

It relies on `Postgres` for database and `etcd` for configuration but both are easily customizable. An [ORM library](https://github.com/jinzhu/gorm) is used for database communication.

# Index

* [Example API](#example-api)
* [Index](#index)
* [Docs](../../../example-api/blob/master/docs/)
* [Dependencies](#dependencies)
* [Setup](#setup)
* [Test Data](#test-data)
* [Testing](#testing)
* [Docker](#docker)
* [Docker Compose](#docker-compose)

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
createuser --createdb example_api
createdb -U example_api example_api
```

Load a development configuration into `etcd`:

```
curl -L http://localhost:2379/v2/keys/config/example_api.json -XPUT -d value='{
    "Database": {
        "Type": "postgres",
        "Host": "localhost",
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
        "APIKey": "sendgrid_api_key"
    },
    "Web": {
        "Scheme": "http",
        "Host": "localhost:8080",
        "AppScheme": "http",
        "AppHost": "localhost:8000"
    },
    "AppSpecific": {
        "PasswordResetLifetime": 604800,
        "CompanyName": "Your Company Name",
        "CompanyEmail": "contact@example.com"
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
	oauth/fixtures/test_users.yml \
	accounts/fixtures/roles.yml \
	accounts/fixtures/test_accounts.yml \
	accounts/fixtures/test_users.yml
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
docker build -t example-api:latest .
docker run -e ETCD_ENDPOINT=localhost:2379 -p 8080:8080 --name example-api example-api:latest
```

You can load fixtures with `docker exec` command:

```
docker exec <container_id> /go/bin/example-api loaddata \
	oauth/fixtures/scopes.yml \
	oauth/fixtures/test_clients.yml
```

You can also execute interactive commands by passing `-i` flag:

```
docker exec -i <container_id> /go/bin/example-api createaccount
docker exec -i <container_id> /go/bin/example-api createsuperuser
```

# Docker Compose

You can use [docker-compose](https://docs.docker.com/compose/) to start the app, postgres, etcd in separate linked containers:

```
docker-compose up
```

During up process all configuration and fixtures will be loaded. After successful up you can check, that app is running using for example the health check request:
```
curl --compressed -v localhost:8080/v1/health
```
