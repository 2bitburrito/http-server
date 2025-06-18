# Assignment

- Create an internal/auth package and expose two functions:
- func HashPassword(password string) (string, error): Hash the password using the bcrypt.GenerateFromPassword function. Bcrypt is a secure hash function that is intended for use with passwords.
  func CheckPasswordHash(hash, password string) error: Use the bcrypt.CompareHashAndPassword function to compare the password that the user entered in the HTTP request with the password that is stored in the database.
  I wrote a couple of simple unit tests to ensure the package is working as expected.

Update the POST /api/users endpoint. The body parameters should now require a new password field:

```json{
"password": "04234",
"email": "<lane@example.com>"
}
```

As long as your server uses HTTPS in production, it's safe to send raw passwords in HTTP requests, because the entire request is encrypted.

Use your internal package's HashPassword function to hash the password before storing it in the database. Do NOT return the hashed password in the response. Again, that would be a security risk.

Add a POST /api/login endpoint. This endpoint should allow a user to login. In a future exercise, this endpoint will be used to give the user a token that they can use to make authenticated requests. For now, let's just make sure password validation is working. It should accept this body:
{
"password": "04234",
"email": "<lane@example.com>"
}

You'll need a new query to look up a user by their email address (you don't have access to an ID here). Once you have the user, check to see if their password matches the stored hash using your internal package. If either the user lookup or the password comparison errors, just return a 401 Unauthorized response with the message "Incorrect email or password".

If the passwords match, return a 200 OK response and a copy of the user resource (without the password of course):

{
"id": "f0f87ec2-a8b5-48cc-b66a-a85ce7c7b862",
"created_at": "2021-07-07T00:00:00Z",
"updated_at": "2021-07-07T00:00:00Z",
"email": "<lane@example.com>"
}

Run and submit the CLI tests.
