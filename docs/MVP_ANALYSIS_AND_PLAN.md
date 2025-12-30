# Análisis y Plan de MVP Robusto - Club Pulse

## 1. Análisis de Situación Actual
El proyecto actual (`club-management-system-api`) utiliza una arquitectura de microservicios distribuida con más de 12 servicios (Auth, User, Booking, Membership, etc.) y un BFF (Frontend Gateway).

**Problemas identificados para un MVP:**
- **Complejidad de Despliegue**: Coordinar 12+ contenedores es costoso y complejo para una fase inicial.
- **Sobrecarga de Infraestructura**: Requiere muchos recursos (CPU/RAM) ociosos.
- **Latencia**: Saltos de red entre BFF y servicios.
- **Dificultad de Desarrollo**: Mantener contratos entre múltiples servicios ralentiza la iteración rápida necesaria en un MVP.

## 2. Estrategia de Solución: MVP Robusto (Monolito Modular)
Para cumplir con los objetivos de "fácil de desplegar", "económico" y "robusto", adoptaremos una arquitectura de **Monolito Modular**.

### ¿Por qué Monolito Modular?
- **Despliegue Simple**: Un solo binario/contenedor para todo el backend.
- **Bajo Costo**: Se ejecuta perfectamente en una instancia pequeña (ej. Railway Starter o VPS de ).
- **Robusto**: Utilizamos Go (tipado fuerte, concurrencia) con Clean Architecture.
- **Escalable**: El código se organiza en módulos (`auth`, `booking`, `user`) que respetan límites estrictos. Si un módulo crece mucho, *ese* módulo se puede extraer a un microservicio en el futuro sin reescribir la lógica.

## 3. Arquitectura Propuesta

### Frontend
- **Tecnología**: Next.js (aprovechando `clubpulse-front-api` existente).
- **Ubicación**: `/frontend` dentro del nuevo proyecto.
- **Comunicación**: Llamadas REST directas al Backend Monolito (sin BFF complejo intermedio, el backend sirve como API Gateway lógico).

### Backend
- **Tecnología**: Go 1.24+ con Gin o Echo.
- **Ubicación**: `/backend` dentro del nuevo proyecto.
- **Estructura**:
  ```text
  backend/
  ├── cmd/api/          # Entrypoint (Main)
  ├── internal/
  │   ├── core/         # Shared Kernel (Domain events, errors, logger)
  │   ├── modules/      # Módulos aislados (Dominio, Capas)
  │   │   ├── auth/
  │   │   ├── user/
  │   │   ├── booking/
  │   │   └── membership/
  │   └── platform/     # Infraestructura (DB, Redis, Email)
  └── pkg/              # Código público/utils
  ```

## 4. Plan de Migración e Implementación
Trabajaremos migrando la lógica esencial de los microservicios actuales al nuevo monolito.

### Fase 1: Setup y Core (Días 1-2)
### Fase 1: Setup y Core (Días 1-2)
- [x] Crear estructura de carpetas.
- [x] Definir `go.mod` y dependencias (Gin, Gorm/Sqlc, Viper).
- [x] Configurar Logger, Middlewares (CORS, Recovery) y Database Connection.

### Fase 2: Módulos Críticos (Días 3-5)
- [x] **Auth Module**: Migrar lógica de login/registro (JWT).
- [x] **User Module**: Perfiles y roles.
- [x] **Facilities Module**: Definición de canchas/espacios.

### Fase 3: Lógica de Negocio (Días 6-10)
- [/] **Booking Module**: Motor de reservas (validación de conflictos) (En Progreso).
- [x] **Membership Module**: Planes y suscripciones básicas.

### Fase 4: Integración Frontend
### Fase 4: Integración Frontend
- [x] Setup inicial Next.js + Tailwind.
- [x] Implementación UI Auth (Login).
- [x] Dashboard Layout principal.
- [x] Páginas de Facilities y Membership.
- [ ] Consumo de API Booking.

## 5. Estrategia de Despliegue (DevOps)
- **Docker**: Un solo `Dockerfile` multistage para el backend (~20MB imagen final).
- **Base de Datos**: Una sola instancia Postgres con esquemas separados por módulo (`auth_schema`, `booking_schema`) para mantener aislamiento lógico.
- **CI/CD**: GitHub Actions simple (Build -> Test -> Deploy to Railway/Render).

---
**Resultado**: Un sistema que se siente como grandes ligas pero corre con el presupuesto de ligas menores.
