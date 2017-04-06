[![Codeship Status for RichardKnop/example-api](https://codeship.com/projects/eb1ee0d0-ac8c-0133-aaf6-0af8633f2b2a/status?branch=master)](https://codeship.com/projects/131678)

[![Travis Status for RichardKnop/example-api](https://travis-ci.org/RichardKnop/example-api.svg?branch=master)](https://travis-ci.org/RichardKnop/example-api)

# example-api

This is a base project to bootstrap and prototype quickly. It is useful as a starting point for REST APIs and includes full OAuth 2.0 implementation as well as basic endpoints to create and update a user, health check endpoint, Facebook integration, migrations and a ready to rumble Dockerfile.

It relies on `Postgres` for database and `etcd` for configuration but both are easily customizable. An [ORM library](https://github.com/jinzhu/gorm) is used for database communication.

# Index

* [example-api](#example-api)
* [Index](#index)
* [API Design](#api-design)
* [API Docs](../../../example-api/blob/master/docs/)
* [Dependencies](#dependencies)
* [Setup](#setup)
* [Test Data](#test-data)
* [Testing](#testing)
* [Docker](#docker)
* [Docker Compose](#docker-compose)
* [Building](../../../example-api/blob/master/BUILDING.md)

# API Design

The API is using REST architectural style which means resources are access and modified via HTTP methods:

- `GET` (to get a resource or paginate over resource sets)
- `POST` (to create new resources)
- `PUT` (to update resources)
- `DELETE` (to delete resources)

`PATCH` might be implemented later in order to enable partial updates of resources (see the [RFC](https://tools.ietf.org/html/rfc7396)).

The REST API objects are formatted according to JSON [HAL](http://stateless.co/hal_specification.html) specification. This means that each object has its own hyperlink clients can use to access it. Other related objects can be embedded into the response as well.

Simple example of JSON HAL resource:

```json
{
  "_links": {
    "self": {
      "href": "/v1/users/1"
    }
  },
  "id": 1,
  "email": "john@reese",
  "created_at": "2015-12-17T06:17:54Z",
  "updated_at": "2015-12-17T06:17:54Z"
}
```

Let's take a look at how a related object would be represented. The example bellow shows a file resource with embedded user object.

```json
{
  "_links": {
    "self": {
      "href":"/v1/files/1"
    }
  },
  "_embedded": {
    "user": {
      "_links": {
        "self": {
          "href":"/v1/users/1"
        }
      },
      "id":1
    }
  },
  "id":1,
  "user_id":1
}
```

Pagination example:

```json
{
  "_links":{
    "first":{
      "href":"/v1/files?page=1"
    },
    "last":{
      "href":"/v1/files?page=2"
    },
    "next":{
      "href":"/v1/files?page=2"
    },
    "prev":{
      "href":""
    },
    "self":{
      "href":"/v1/files?page=1"
    }
  },
  "_embedded":{
    "files":[
      {
        "_links":{
          "self":{
            "href":"/v1/files/1"
          }
        },
        "id":1
      },
      {
        "_links":{
          "self":{
            "href":"/v1/files/2"
          }
        },
        "id":2
      }
    ]
  },
  "count":2,
  "page":1
}
```

# Dependencies

According to [Go 1.5 Vendor experiment](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo), all dependencies are stored in the vendor directory. This approach is called `vendoring` and is the best practice for Go projects to lock versions of dependencies in order to achieve reproducible builds.

To update dependencies during development:

```sh
make update-deps
```

To install dependencies:

```sh
make install-deps
```

# Setup

If you are developing on OSX, install `etcd` and `Postgres``:

## etcd

```sh
brew install etcd
```

Load a development configuration into `etcd`:

```sh
etcdctl put /config/example_api.json '{
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
  "Mailgun": {
    "Domain": "example.com",
    "APIKey": "mailgun_api_key",
    "PublicAPIKey": "mailgun_private_api_key"
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
```

Check the config was loaded properly:

```sh
etcdctl get /config/example_api.json
```

## Postgres

```sh
brew install postgres
```

You might want to create a `Postgres` database:

```sh
createuser --createdb example_api
createdb -U example_api example_api
```

## Compile & Run

Compile the app:

```sh
go install .
```

Run migrations:

```sh
example-api migrate
```

And finally, run the app:

```sh
example-api runserver
```

When deploying, you can set etcd related environment variables:

* `ETCD_ENDPOINTS`
* `ETCD_CERT_FILE`
* `ETCD_KEY_FILE`
* `ETCD_CA_FILE`
* `ETCD_CONFIG_PATH`

# Test Data

You might want to insert some test data if you are testing locally using `curl` examples from this README:

```sh
example-api loaddata \
  services/oauth/fixtures/scopes.yml \
  services/oauth/fixtures/roles.yml \
  services/oauth/fixtures/test_clients.yml \
  services/oauth/fixtures/test_users.yml \
  services/accounts/fixtures/test_users.yml
```

# Testing

I have used a mix of unit and functional tests so you need to have `sqlite` installed in order for the tests to run successfully as the suite creates an in-memory database.

To run tests:

```sh
make test
```

# Docker

Build a Docker image and run the app in a container:

```sh
docker build -t example-api:latest .
docker run -e ETCD_ENDPOINTs=localhost:2379 -p 8080:8080 --name example-api example-api:latest
```

You can load fixtures with `docker exec` command:

```sh
docker exec <container_id> /go/bin/example-api loaddata \
  services/oauth/fixtures/scopes.yml \
  services/oauth/fixtures/roles.yml \
  services/oauth/fixtures/test_clients.yml
```

You can also execute interactive commands by passing `-i` flag:

```sh
docker exec -i <container_id> /go/bin/example-api createoauthclient
docker exec -i <container_id> /go/bin/example-api createsuperuser
```

# Docker Compose

You can use [docker-compose](https://docs.docker.com/compose/) to start the app, postgres, etcd in separate linked containers:

```sh
docker-compose up
```

During up process all configuration and fixtures will be loaded. After successful up you can check, that app is running using for example the health check request:

```sh
curl --compressed -v localhost:8080/v1/health
```
