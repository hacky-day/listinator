# Frontend

Vue frontend of the listinator.

---

## Setup (Production)

The frontend is embedded in the backend.

## Setup (Development)

### Requirements

- Node.js 18+ and npm
- `make`

### Quick Start

Simply run:

```bash
make clean run
```

This will:

1. Clean old artifacts
2. Install all dependencies
3. Run the frontend with a development server which also proxies all api
   requests to the backend for easier development
