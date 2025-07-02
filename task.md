# Assignment

Add a new DELETE /api/chirps/{chirpID} route to your server that deletes a chirp from the database by its id.
This is an authenticated endpoint, so be sure to check the token in the header. Only allow the deletion of a chirp if the user is the author of the chirp.
If they are not, return a 403 status code.
If the chirp is deleted successfully, return a 204 status code.
If the chirp is not found, return a 404 status code.
Run and submit the CLI tests.
