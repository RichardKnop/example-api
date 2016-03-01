# Facebook

* [Login](#login)

## Login

Example request:

```
curl --compressed -v localhost:8080/v1/facebook/login \
	-H "Content-Type: application/json" \
	-u test_client_1:test_secret \
	-d "access_token=facebook_access_token"
```

Example response:

```json
{
	"id": 1,
	"access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
	"expires_in": 3600,
	"token_type": "Bearer",
	"scope": "read_write",
	"refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```
