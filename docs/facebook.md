# Facebook

* [Login](#login)

## Login

General flow:

1) Mobile app displays a Facebook login button.

2) The `oauth2` flow happens completely within the app. No interaction with the API.

3) Once app fetches an access token from Facebook, it POSTs it to `/v1/facebook/login` endpoint which does registration/login behind the scenes (creates a new account for the user if it doesn't exist yet) and returns access token.

4) The mobile app is now logged in.

Example request:

```
curl --compressed -v localhost:8080/v1/facebook/login \
	-u test_client_1:test_secret \
	-d "access_token=facebook_access_token" \
	-d "scope=read_write"
```

Example response:

```json
{
  "id": 1,
  "user_id": 1,
  "access_token": "00ccd40e-72ca-4e79-a4b6-67c95e2e3f1c",
  "expires_in": 3600,
  "token_type": "Bearer",
  "scope": "read_write",
  "refresh_token": "6fd8d272-375a-4d8a-8d0f-43367dc8b791"
}
```
