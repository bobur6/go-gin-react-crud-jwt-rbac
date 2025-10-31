# Assignment 3 – Full-Stack Application

This repository contains a Golang backend (Gin) secured with JWT + RBAC and a React frontend (Create React App) that consumes the API, handles authentication, and manages global state through React Context.

## Project Structure

```
assignment3/
├── backend/   # Go API server, in-memory data store, JWT auth, RBAC
└── frontend/  # React client, Context state management, Axios API layer
```

## Backend

### Prerequisites

- Go 1.24+

### Configuration

Environment variables (optional, defaults shown in brackets):

| Variable               | Default                   | Description                                        |
|------------------------|---------------------------|----------------------------------------------------|
| `PORT`                 | `8080`                    | HTTP port for the API                              |
| `JWT_SECRET`           | `change-me-in-production` | Secret used to sign JWT tokens                     |
| `JWT_ISSUER`           | `assignment3-backend`     | Issuer claim for JWT tokens                        |
| `JWT_EXPIRY_MINUTES`   | `60`                      | Token lifetime in minutes                          |
| `ADMIN_USERNAME`       | `admin`                   | Username for the seeded admin account              |
| `ADMIN_PASSWORD`       | `admin123`                | Password for the seeded admin account              |
| `FRONTEND_ORIGINS`     | *(empty)*                 | Extra allowed origins for CORS (comma-separated)   |

> **PowerShell note:** set variables per session using `$env:PORT = "8080"` (no `export`).  
> To see the current value run `Get-ChildItem Env:PORT`.

### Run the API

```powershell
cd backend
$env:PORT = "8080"
$env:JWT_SECRET = "supersecret"
# Optional: allow extra frontend origins (example with LAN IP)
$env:FRONTEND_ORIGINS = "http://localhost:3000,http://127.0.0.1:3000,http://10.0.120.40:3000"
# Wildcard for troubleshooting only:
# $env:FRONTEND_ORIGINS = "*"
go run ./cmd/server
```

Allowed origins always include:

- `http://localhost:3000`
- `http://127.0.0.1:3000`
- The API host itself (`http://localhost:<PORT>` and `http://127.0.0.1:<PORT>`)

Add extra hosts (e.g. LAN IPs) via `FRONTEND_ORIGINS`. Use `*` only if you need to temporary disable checks.

The server exposes routes prefixed with `/api` (register, login, items CRUD, health check). Default admin credentials: **admin / admin123**.

### Tests

```powershell
cd backend
go test ./...
```

## Frontend

### Prerequisites

- Node.js 18+
- npm 9+

### Setup & Run

```powershell
cd frontend
npm install
npm start
```

Open <http://localhost:3000>. The dev server proxies API calls to <http://localhost:8080>. If you prefer using a LAN IP (e.g. `http://10.0.120.40:3000`), ensure that origin is listed in `FRONTEND_ORIGINS` before starting the backend.

### Build

```powershell
cd frontend
npm run build
```

## Features

### Role-Based Access Control (RBAC)

The application supports two roles with different privileges:

**User Role:**
- Create items
- View all items
- Edit own items
- Cannot delete items

**Admin Role:**
- All user privileges
- Edit any item
- Delete any item
- **Manage users** (view all users, delete users)

### Admin User Management

Admins have access to a dedicated User Management panel where they can:
- View all registered users
- See user roles and registration dates
- Delete users (except their own account)
- Refresh the user list

## Development Tips

- Backend and frontend can run simultaneously (ports 8080 and 3000 by default).  
- JWT tokens are stored in `localStorage`. Use the **Sign Out** button to clear them during development.  
- The in-memory store is reset whenever the backend restarts.  
- To seed additional demo data, adjust `cmd/server/main.go`.
- Default admin credentials: **admin / admin123**
