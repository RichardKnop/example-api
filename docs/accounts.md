# CURL Examples

# Index

* [CURL Examples](#curl-examples)
* [Index](#index)
* [Accounts](#accounts)
  * [Users](#users)
    * [Create User](#create-user)
    * [Get Me](#get-me)
    * [Update User](#update-user)
  * [Reset Password](#reset-password)

# Accounts

## Users

### Create User

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/users \
	-H "Content-Type: application/json" \
	-u test_client_1:test_secret \
	-d '{
		"email": "test@user",
		"password": "test_password"
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/accounts/users/1"
		}
	},
	"id": 1,
	"email": "test@user",
	"first_name": "",
	"last_name": "",
	"role": "user",
	"confirmed": false,
	"created_at": "2015-12-17T06:17:54Z",
	"updated_at": "2015-12-17T06:17:54Z"
}
```

### Get Me

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/me \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c"
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/accounts/users/1"
		}
	},
	"id": 1,
	"email": "test@user",
	"first_name": "",
	"last_name": "",
	"role": "user",
	"confirmed": true,
	"created_at": "2015-12-17T06:17:54Z",
	"updated_at": "2015-12-17T06:17:54Z"
}
```

### Update User

Example request:

```
curl -XPUT --compressed -v localhost:8080/v1/accounts/users/1 \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer 00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c" \
	-d '{
		"email": "test@user_updated",
		"first_name": "test_first_name_updated",
		"last_name": "test_last_name_updated",
	}'
```

Example response:

```json
{
	"_links": {
		"self": {
			"href": "/v1/accounts/users/1"
		}
	},
	"id": 1,
	"email": "test@user_updated",
	"first_name": "test_first_name_updated",
	"last_name": "test_last_name_updated",
	"role": "user",
	"confirmed": true,
	"created_at": "2015-12-17T06:17:54Z",
	"updated_at": "2015-12-18T07:09:15Z"
}
```

## Reset Password

Example request:

```
curl --compressed -v localhost:8080/v1/accounts/passwordreset \
	-H "Content-Type: application/json" \
	-u test_client_1:test_secret \
	-d '{
		"email": "test@user"
	}'
```

Returns `204` empty response on success.
