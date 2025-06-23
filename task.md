Add a MakeJWT function to your auth package:
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error)

Create and return a JWT using this JWT library, which you can import into your code by running:

go get -u github.com/golang-jwt/jwt/v5

Create a new token.
Use jwt.NewWithClaims
Use jwt.SigningMethodHS256 as the signing method.
Use jwt.RegisteredClaims as the claims.
Set the Issuer to "chirpy"
Set IssuedAt to the current time in UTC
Set ExpiresAt to the current time plus the expiration time (expiresIn)
Set the Subject to a stringified version of the user's id
Use token.SignedString to sign the token with the secret key. Refer to here for an overview of the different signing methods and their respective key types.
Add a ValidateJWT function to your auth package:
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error)

Use the jwt.ParseWithClaims function to validate the signature of the JWT and extract the claims into a \*jwt.Token struct. An error will be returned if the token is invalid or has expired.
If all is well with the token, use the token.Claims interface to get access to the user's id from the claims (which should be stored in the Subject field). Return the id as a uuid.UUID.

Add some more unit tests to the auth package. Make sure that you can create and validate JWTs, and that expired tokens are rejected and JWTs signed with the wrong secret are rejected.
Run and submit the CLI tests.

B
