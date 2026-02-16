# UniSwap Frontend

React + TypeScript frontend for the campus marketplace.

## Setup

```bash
cd frontend
npm install
```

## Run

```bash
npm run dev
```

Runs at [http://localhost:5173](http://localhost:5173) by default.

## Backend URL

Set `VITE_API_URL` to your backend base URL (e.g. `http://localhost:8080/api`). Default is `http://localhost:8080/api`.

Create `frontend/.env`:

```
VITE_API_URL=http://localhost:8080/api
```

## Sprint 1 (Samarth)

- FE-1: React + TypeScript (Vite), folder structure
- FE-2: App routing (/, /login, /register, /listing/:id, /create, 404)
- FE-3: Login page | FE-4: Register page
- FE-5: JWT token storage (localStorage)
- FE-6: Protected route component
- FE-7: API service layer (axios, auth header, 401 → login)
