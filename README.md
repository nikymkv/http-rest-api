# http-rest-api


Authorization service using access and refresh tokens, type jwt.

Command for use:
Get new Access and Refresh tokens:
`curl -v POST <your url/session-create> -d '{"user_id":"<your user_id>","fingerprint":"<your fingerprint>"}'`

Update pair tokens:
`curl -v POST <your url/session-refresh> -d '{"refresh_token":"<your refresh_token(base64 encoded)>"}'`

Delete current session:
`curl -v POST <your url/session-delete> -d '{"refresh_token":"<your refresh_token(base64 encoded)>"}'`

Delete all user sessions:
`curl -v POST <your url/delete-all-sessions> -d '{"refresh_token":"<your refresh_token(base64 encoded)>"}'`
