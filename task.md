## Step 1: Login

It would be pretty annoying if you had to enter your username and password every time you wanted to make a request to an API. Instead, after a user enters a username and password, our server should respond with a token (JWT) that's saved in the client's device.

The token remains valid until it expires, at which point the user will need to log in again.

Step 2: Using the Token
When the user wants to make a request to the API, they send the token along with the request in the HTTP headers. The server can then verify that the token is valid, which means the user is who they say they are.

Assignment
Add a GetBearerToken function to your auth package:
func GetBearerToken(headers http.Header) (string, error)

Auth information will come into our server in the Authorization header. Its value will look like this:

Bearer TOKEN_STRING

This function should look for the Authorization header in the headers parameter and return the TOKEN_STRING if it exists (stripping off the Bearer prefix and whitespace). If the header doesn't exist, return an error. This is an easy one to write a unit test for, and I'd recommend doing so.

Create a secret for your server and store it in your .env file. This is the secret used to sign and verify JWTs. By keeping it safe, no other servers will be able to create valid JWTs for your server. We will yet again use environment variables. You can generate a nice long random string on the command line like this:
openssl rand -base64 64

Secrets should NOT be stored in Git, just in case anyone malicious gains access to your repository.

Load the JWT secret from your .env file in your main() function and store it in your apiConfig struct.
Update your POST /api/login endpoint. It should accept a new, optional expires_in_seconds field in the request body:
{
"password": "04234",
"email": "<lane@example.com>",
"expires_in_seconds": 2
}

If it's specified by the client, use it as the expiration time. If it's not specified, use a default expiration time of 1 hour. If the client specified a number over 1 hour, use 1 hour as the expiration time.

Once you have the token created with the new params, respond to the request with a 200 code and this body shape:

{
"id": "5a47789c-a617-444a-8a80-b50359247804",
"created_at": "2021-07-01T00:00:00Z",
"updated_at": "2021-07-01T00:00:00Z",
"email": "<lane@example.com>",
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}

Update the POST /api/chirps endpoint. It is not an authenticated endpoint yet. To post a chirp, a user needs to have a valid JWT. The JWT will determine which user is posting the chirp. Use your GetBearerToken and ValidateJWT functions. If the JWT is invalid, return a 401 Unauthorized response.
