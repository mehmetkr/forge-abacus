# Forge Abacus — Frontend

A calculator UI built with React and TypeScript. Consumes the backend REST API to perform seven arithmetic operations.

## Architecture

```
frontend/
├── src/
│   ├── main.tsx                  # React DOM mount
│   ├── index.css                 # global reset
│   ├── App.tsx                   # root component, layout
│   ├── App.module.css
│   ├── calculate.ts              # API call, Operation type
│   └── components/
│       ├── Calculator.tsx        # state, orchestration, form
│       ├── Calculator.module.css
│       ├── NumberInput.tsx       # labeled numeric input (used twice)
│       └── NumberInput.module.css
├── package.json
├── index.html
├── vite.config.ts
├── vitest.config.ts
├── Dockerfile
└── nginx.conf
```

Three layers:

- **`calculate.ts`** — API function and `Operation` type. Calls `POST /api/calculate`, returns the result or throws with the backend's error message. Network failures throw "Could not reach the server".
- **`Calculator`** — Owns all form state. Validates inputs with `Number()` + `Number.isFinite()` before sending. Renders the operation dropdown, number inputs, submit button, and result/error display.
- **`NumberInput`** — Reusable labeled input with `inputMode="decimal"` for mobile keyboards. Instantiated twice (A and B), with B hidden for square root.

## Setup

### Prerequisites

- Node.js 22+
- Backend running on port 8080 (see `backend/README.md`)

### Local development

```bash
cd frontend
npm install
npm run dev
```

The dev server starts on `http://localhost:5173` and proxies `/api` requests to `localhost:8080`.

### Docker Compose (full stack)

From the repository root:

```bash
docker compose up
```

This starts the backend and frontend together. The calculator is available at `http://localhost:3000`. nginx serves the frontend and proxies API requests to the backend.

## Tests

Run all tests:

```bash
cd frontend
npm test
```

With verbose output:

```bash
npx vitest run --reporter=verbose
```

With coverage (requires `@vitest/coverage-v8`):

```bash
npm install -D @vitest/coverage-v8
npx vitest run --coverage
```

## Design decisions

- **Zero runtime dependencies beyond React.** No state management libraries, no CSS frameworks, no HTTP clients. React state, CSS Modules, and the Fetch API are sufficient for a single-page calculator.
- **`Operation` union type.** A TypeScript union (`"add" | "subtract" | ... | "percentage"`) is the single source of truth for valid operations, catching invalid values at compile time.
- **Client-side validation complements, not duplicates.** The frontend validates that inputs are present and finite (rejecting values JSON cannot represent). Domain validation (division by zero, negative square root) is left to the backend.
- **Backend error messages shown directly.** The backend already returns human-readable error messages. No client-side mapping layer.
- **`inputMode="decimal"`** on number inputs triggers a numeric soft keyboard on mobile without restricting input to digits only (negative numbers and decimals are valid).
- **`aria-live="polite"`** on the result area so screen readers announce computation results and errors.
- **CSS Modules for scoping.** Scoped styles without runtime cost or external dependencies. Each component has a co-located `.module.css` file.
