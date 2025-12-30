# Club Pulse System API

## Introduction
Club Pulse is a robust, modular monolith designed to manage sports club operations. It serves as a modern replacement for the legacy microservices architecture, consolidating logic into a single, efficient, and easy-to-deploy backend API with a Next.js frontend.

## 🚀 Features
- **Modular Architecture**: Clean separation of concerns (Auth, User, Facilities, Membership).
- **High Performance**: Built with Go 1.23+ and Gin.
- **Modern Frontend**: Next.js 14 App Router with Tailwind CSS.
- **Easy Deployment**: Dockerized stack (API + Frontend + Postgres) managed via Docker Compose.

## 🛠 Tech Stack
- **Backend**: Go, Gin, GORM, PostgreSQL.
- **Frontend**: Next.js (TypeScript), Tailwind CSS, Lucide Icons.
- **Infrastructure**: Docker, Docker Compose.

## 🏁 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (optional, for local dev without Docker)
- Node.js 20+ (optional, for local frontend dev)

### Quick Start (Recommended)
1. **Clone the repository**:
   ```bash
   git clone <repo-url>
   cd club-pulse-system-api
   ```

2. **Run with Docker Compose**:
   ```bash
   docker-compose up --build
   ```
   - Legacy command: `docker-compose up --build` works for both backend and frontend.
   - **Frontend**: `http://localhost:3000`
   - **Backend API**: `http://localhost:8080`

3. **Verify Installation**:
   - Visit `http://localhost:3000` to see the login page.
   - Default test user (auto-created if seed runs): `testuser@example.com` / `password123`.

## 📂 Project Structure
```
.
├── backend/                # Go Monolith API
│   ├── cmd/api/            # Entrypoint
│   └── internal/modules/   # Domain logic (Auth, User, Facilities, Membership)
├── frontend/               # Next.js Application
├── docs/                   # Documentation & Plans
├── scripts/                # Utility scripts (Integration tests, etc.)
└── docker-compose.yml      # Orchestration
```

## 🧪 Testing
Run the integration test suite to verify the backend health and workflows:
```bash
./scripts/integration_test.sh
```

## 📚 Documentation
- [MVP Analysis & Plan](docs/MVP_ANALYSIS_AND_PLAN.md)
- [Migration Context](docs/MIGRATION_CONTEXT.md)
- [API Documentation](docs/API_DOCUMENTATION.md) (See for endpoints)
