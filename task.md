# Assignment

Update the GET /api/chirps endpoint. It should accept an optional query parameter called sort. It can have 2 possible values:

asc - Sort the chirps in the response by created_at in ascending order
desc - Sort the chirps in the response by created_at in descending order
asc is the default if no sort query parameter is provided.

Keep it simple! You can just sort the chirps in-memory using sort.Slice.

Run and submit the CLI tests.

Examples of Valid URLs
GET <http://localhost:8080/api/chirps?sort=asc>
GET <http://localhost:8080/api/chirps?sort=desc>
GET <http://localhost:8080/api/chirps>
