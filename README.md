# ODS Quiz API

Questionnaire microservice for ODS Quiz. It serves the individual quiz,
calculates scores using the legacy scoring rules, and stores normalized
submission data for later BI analysis.

## Run locally

```bash
cp .env.example .env.local
```

Set local values:

```dotenv
PORT=8082
JWTSecret=local-development-secret
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=your_local_user
DB_PASSWORD=your_local_password
DB_NAME=odsquiz
DB_SSLMODE=disable
```

Then run:

```bash
go run ./cmd/api
```

## API

- `GET /api/quiz/questions`
- `POST /api/quiz/submissions`
- `GET /api/quiz/submissions`

Both submission routes need a bearer token from the auth API. A user whose
`users.suspended_at` is set gets `403 {"code":"account_suspended"}` on
`POST /api/quiz/submissions`; reading the own history stays allowed.

Submit payload:

```json
{
  "answers": [
    { "question_id": "uuid", "answer_option_id": "uuid" }
  ]
}
```

The service also accepts a compact legacy-friendly payload:

```json
{ "answers": [1, 2, 3, 4] }
```
