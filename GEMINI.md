# Project Context: club-pulse-system-api

This file serves as a persistent memory for the Gemini agent to understand the project's context, structure, and conventions.

## 1. Project Overview
- **Type:** Monorepo-style project.
- **Backend:** A Go application located in the `backend/` directory.
- **Frontend:** A Next.js (TypeScript/React) application located in the `frontend/` directory.
- **Containerization:** Uses `docker-compose.yml` for service orchestration.

## 2. Key Technologies & Languages
- **Backend:** Go
- **Frontend:** Next.js, React, TypeScript, Tailwind CSS
- **Package Management:** Go Modules (`go.mod`), NPM (`package.json`)

## 3. Important Files & Directories
- `backend/cmd/api/main.go`: The main entry point for the Go API server.
- `frontend/app/`: The main application directory for the Next.js frontend, using the App Router.
- `package.json`: Located in `frontend/`, contains scripts for running, testing, and building the frontend.
- `go.mod`: Located in `backend/`, contains dependencies for the Go application.
- `docker-compose.yml`: Defines the services, networks, and volumes for local development.

## 4. Version Control
- **Provider:** GitHub
- **Repository URL:** https://github.com/lukcba-developers/club-pulse-system-api

## 5. User Preferences
- The user prefers all communication to be in Spanish.