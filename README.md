# Forge Abacus

A full-stack calculator application. Go REST API backend, React/TypeScript frontend.

## Quick start

```bash
docker compose up
```

Open http://localhost:3000 in your browser. To use a different port, change the `3000:80` mapping in `docker-compose.yml`.

## Project structure

| Directory  | Description                        | Details                          |
|------------|------------------------------------|----------------------------------|
| `backend/` | REST API — Go, zero dependencies   | [backend/README.md](backend/README.md) |
| `frontend/`| Calculator UI — React, TypeScript  | [frontend/README.md](frontend/README.md) |

## API

Single endpoint: `POST /api/calculate`

Supports seven operations: `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`, `percentage`.

```bash
curl -X POST http://localhost:3000/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":2,"b":3}'
# {"result":5}
```

See [backend/README.md](backend/README.md) for the full API reference.

## Local development

Start the backend and frontend separately:

```bash
cd backend && go run .
```

```bash
cd frontend && npm install && npm run dev
```

The frontend dev server (http://localhost:5173) proxies `/api` requests to the backend on port 8080.

## Tests

```bash
cd backend && go test ./...
cd frontend && npm test
```
