# Go Clean Architecture HTTP Server (User CRUD)

This project demonstrates a minimal HTTP server in Go following Clean Architecture principles with a `User` CRUD example. It uses only the standard library and Go 1.22's `http.ServeMux` path patterns.

## Requirements

- Go 1.22+

## Run

```bash
go build ./...
GO111MODULE=on go run ./cmd/server
```

The server starts on `:8080`.

## Endpoints

- POST `/api/v1/users`
- GET `/api/v1/users`
- GET `/api/v1/users/{id}`
- PUT `/api/v1/users/{id}`
- DELETE `/api/v1/users/{id}`
- GET `/healthz`

## Examples

Create:

```bash
curl -s -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}' | jq
```

List:

```bash
curl -s http://localhost:8080/api/v1/users | jq
```

Get by ID:

```bash
curl -s http://localhost:8080/api/v1/users/1 | jq
```

Update:

```bash
curl -s -X PUT http://localhost:8080/api/v1/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Cooper","email":"acooper@example.com"}' | jq
```

Delete:

```bash
curl -i -X DELETE http://localhost:8080/api/v1/users/1
```

## Structure

```
cmd/
  server/
    main.go
internal/
  app/
    router.go
    middleware.go
  adapter/
    http/
      user_handler.go
  domain/
    user.go
  repository/
    memory/
      user_repository.go
  usecase/
    user_service.go
```

## Notes

- The repository is in-memory for demo purposes. Swap it with a real implementation (e.g., Postgres) by implementing `internal/domain.UserRepository` and wiring it in `cmd/server/main.go`.
