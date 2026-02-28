# Subtle SQL Injection (Go)

Intentionally vulnerable Go service for **testing security and static analysis tools**. The SQL injection is not in a single file: user input flows through several isolated layers and is concatenated into SQL only in a dedicated query-builder package.

## Why “subtle”?

- **No raw SQL in the HTTP handler** – the handler only reads query parameters and calls the service.
- **No SQL in the service layer** – the service forwards arguments to the repository.
- **Repository does not build strings** – it calls into `internal/querybuilder` and uses the returned clause in a query.
- **Vulnerability is only in `internal/querybuilder`** – `BuildCondition` and `OrderBy` concatenate user-controlled input into SQL.

A tool that only scans for `db.Query(...)` or string concatenation in the same file as the handler will miss it. Detecting this requires **data-flow / taint analysis** from HTTP input → handler → service → repository → querybuilder.

## Layout

```
sql-injection-subtle/
├── main.go                          # wiring only
├── internal/
│   ├── handler/user_handler.go      # reads q, username, sort, order from request
│   ├── service/user_service.go      # passes through to repository
│   ├── repository/user_repository.go # builds query using querybuilder
│   ├── querybuilder/
│   │   ├── where.go                 # BuildCondition(column, value) – concatenates value
│   │   └── order.go                 # OrderBy(column, direction) – concatenates column
│   └── database/schema.go           # safe DDL/seed only
├── README.md
└── poc.py                           # proof-of-concept requests
```

## Vulnerable endpoints

| Endpoint | Source of taint | Sink |
|----------|-----------------|------|
| `GET /users/search?q=` or `GET /users?username=` | `q` or `username` | `querybuilder.BuildCondition("username", value)` |
| `GET /users?sort=&order=` | `sort`, `order` | `querybuilder.OrderBy(sort, order)` |

## Build and run

```bash
cd sql-injection-subtle
go mod tidy
go run .
```

Server listens on `http://localhost:8080`.

## Proof of concept (for tool testing)

```bash
# Normal
curl "http://localhost:8080/users/search?q=alice"
curl "http://localhost:8080/users?username=admin"

# SQLi via search (WHERE clause)
curl "http://localhost:8080/users/search?q=alice' OR '1'='1"
curl "http://localhost:8080/users?username=admin'--"

# SQLi via sort (ORDER BY / second-order)
curl "http://localhost:8080/users?sort=id;(SELECT 1)--&order=asc"
```

Run the PoC script (starts server, sends requests, reports success/failure):

```bash
python3 poc.py
```

## What a good security/SAST tool should do

1. **Taint from HTTP** – treat query/form parameters as tainted.
2. **Follow calls** – handler → service → repository → querybuilder.
3. **Flag sinks** – use of tainted data in:
   - `querybuilder.BuildCondition(_, value)` (second argument).
   - `querybuilder.OrderBy(column, _)` (first argument).
   - Or any string that is concatenated into `db.Query(...)` / `db.QueryRow(...)`.

This project is for security tooling evaluation only. Do not use in production.
