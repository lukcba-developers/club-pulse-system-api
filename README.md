# Club Pulse System API

## Introducción
Club Pulse es un monolito modular y robusto diseñado para gestionar las operaciones de un club deportivo. Sirve como un reemplazo moderno a la arquitectura de microservicios heredada, consolidando la lógica en una única API de backend eficiente y fácil de desplegar, junto con un frontend en Next.js.

## 🚀 Características
- **Arquitectura Modular**: Clara separación de responsabilidades (Autenticación, Usuarios, Instalaciones, Membresías).
- **Alto Rendimiento**: Construido con Go 1.23+ y Gin.
- **Frontend Moderno**: Next.js 14 con App Router y Tailwind CSS.
- **Despliegue Sencillo**: Stack dockerizado (API + Frontend + Postgres) gestionado a través de Docker Compose.

## 🛠️ Tecnologías
- **Backend**: Go, Gin, GORM, PostgreSQL.
- **Frontend**: Next.js (TypeScript), Tailwind CSS, Lucide Icons.
- **Infraestructura**: Docker, Docker Compose.

## 🏁 Primeros Pasos

### Prerrequisitos
- Docker & Docker Compose
- Go 1.23+ (opcional, para desarrollo local sin Docker)
- Node.js 20+ (opcional, para desarrollo local del frontend)

### Inicio Rápido (Recomendado)
1. **Clona el repositorio**:
   ```bash
   git clone https://github.com/lukcba-developers/club-pulse-system-api.git
   cd club-pulse-system-api
   ```

2. **Ejecuta con Docker Compose**:
   ```bash
   docker-compose up --build
   ```
   - El comando `docker-compose up --build` levanta tanto el backend como el frontend.
   - **Frontend**: `http://localhost:3000`
   - **Backend API**: `http://localhost:8080`

3. **Verifica la Instalación**:
   - Visita `http://localhost:3000` para ver la página de inicio de sesión.
   - Usuario de prueba por defecto (creado automáticamente si se ejecuta el seeder): `testuser@example.com` / `password123`.

## 📂 Estructura del Proyecto
```
.
├── backend/                # API Monolítica en Go
│   ├── cmd/api/            # Punto de entrada
│   └── internal/modules/   # Lógica de dominio (Auth, User, Facilities, Membership)
├── frontend/               # Aplicación Next.js
├── docs/                   # Documentación y planes
├── scripts/                # Scripts de utilidad (pruebas de integración, etc.)
└── docker-compose.yml      # Orquestación de servicios
```

## 🧪 Pruebas
Ejecuta la suite de pruebas de integración para verificar el estado y los flujos de trabajo del backend:
```bash
./scripts/integration_test.sh
```

## 📚 Documentación
- [Análisis y Plan del MVP](docs/MVP_ANALYSIS_AND_PLAN.md)
- [Contexto de la Migración](docs/MIGRATION_CONTEXT.md)
- [Documentación de la API](docs/API_DOCUMENTATION.md) (Ver para los endpoints)