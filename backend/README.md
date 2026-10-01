# Forge Abacus — Backend

A calculator REST API built with Go. Supports seven arithmetic operations with explicit handling of floating-point edge cases.

## Architecture

```
backend/
├── main.go          # server setup, routing, CORS, graceful shutdown
├── handler/
│   ├── calculate.go # HTTP concerns: decode, validate, respond
│   └── model.go     # request/response types
└── calc/
    └── calc.go      # pure arithmetic, zero HTTP awareness
```

Three packages with a single responsibility each:

- **`calc`** — Pure computation. Takes an operator and two `float64` operands, returns a result or a named error. No HTTP imports.
- **`handler`** — HTTP layer. Unmarshals JSON, validates input, calls `calc.Compute`, marshals the response.
- **`main`** — Wiring. Routes requests, applies CORS, manages graceful shutdown.

## Setup

### Prerequisites

- Go 1.27+

### Local development

```bash
cd backend
go run .
```

The server starts on port 8080 by default. Override with the `PORT` environment variable:

```bash
PORT=3000 go run .
```

### Docker

```bash
cd backend
docker build -t forge-abacus .
docker run -p 8080:8080 forge-abacus
```

## API

### Endpoint

`POST /api/calculate`

### Request

```json
{
  "operation": "add",
  "a": 10,
  "b": 3
}
```

| Field       | Type   | Required    | Description                                              |
|-------------|--------|-------------|----------------------------------------------------------|
| `operation` | string | yes         | One of: `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`, `percentage` |
| `a`         | number | yes         | First operand                                            |
| `b`         | number | conditional | Required for all operations except `sqrt`                |

### Response

Success (200):
```json
{"result": 13}
```

Error (400):
```json
{"error": "division by zero is undefined"}
```

### Examples

```bash
# Addition
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"add","a":2,"b":3}'
# {"result":5}

# Subtraction
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"subtract","a":10,"b":4}'
# {"result":6}

# Multiplication
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"multiply","a":7,"b":6}'
# {"result":42}

# Division
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"divide","a":10,"b":3}'
# {"result":3.3333333333333335}

# Power
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"power","a":2,"b":10}'
# {"result":1024}

# Square root
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"sqrt","a":9}'
# {"result":3}

# Percentage (15% of 200)
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"percentage","a":15,"b":200}'
# {"result":30}
```

### Edge cases

```bash
# Division by zero
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"divide","a":1,"b":0}'
# {"error":"division by zero is undefined"}

# Negative square root
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"sqrt","a":-4}'
# {"error":"square root of negative number is undefined"}

# Missing required field
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"add","a":1}'
# {"error":"field 'b' is required"}

# Invalid JSON
curl -X POST http://localhost:8080/api/calculate \
  -d 'not json'
# {"error":"invalid JSON in request body"}

# Unsupported operation
curl -X POST http://localhost:8080/api/calculate \
  -d '{"operation":"modulo","a":10,"b":3}'
# {"error":"unsupported operation: modulo"}

# Wrong HTTP method
curl -X GET http://localhost:8080/api/calculate
# {"error":"method not allowed"}
```

## Design decisions

- **Zero external dependencies.** The standard library is sufficient for a calculator API. No third-party packages means no supply chain risk and no version management overhead.
- **`*float64` operands.** Pointer types distinguish "field absent" (`nil`) from "field is zero" (`0`). This prevents `{"operation":"add","a":5}` from silently treating the missing `b` as zero.
- **Layered error checking.** Domain errors (division by zero, negative square root) are caught before computation. A post-computation guard catches overflow to infinity and NaN from any operation, providing a safety net without per-operation special cases.
- **CORS wildcard.** `Access-Control-Allow-Origin: *` is appropriate for a stateless calculator with no authentication or credentials.
- **Graceful shutdown.** The server drains active requests on SIGINT/SIGTERM before exiting, preventing interrupted responses.
- **Minimal status codes.** The handler returns 400 for all client errors; routing returns 405 for wrong methods and 204 for CORS preflight. The error message body carries the specificity.

## Tests

Run all tests:

```bash
cd backend
go test ./...
```

With verbose output:

```bash
go test ./... -v
```

With coverage:

```bash
go test ./... -cover
```

Generate a coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
