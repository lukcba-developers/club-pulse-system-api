# Club Pulse System

*Tu sistema de gestión de clubes deportivos, todo en uno.*

[![CI](https://github.com/lukcba-developers/club-pulse-system-api/actions/workflows/ci.yml/badge.svg)](https://github.com/lukcba-developers/club-pulse-system-api/actions/workflows/ci.yml)
[![Licencia: MIT](https://img.shields.io/badge/Licencia-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Club Pulse** es un sistema de gestión robusto y modular diseñado para centralizar todas las operaciones de un club deportivo. El proyecto consolida la lógica de negocio en una única API de backend de alto rendimiento construida con Go, junto con un frontend moderno y reactivo en Next.js.

---

## 🚀 Características Principales

El sistema está organizado en módulos cohesivos que cubren todas las necesidades de un club moderno:

-   **Gestión de Autenticación y Roles:** Sistema seguro de inicio de sesión (credenciales y OAuth) con un robusto control de acceso basado en roles (`MEMBER`, `ADMIN`, `SUPER_ADMIN`).
-   **Gestión de Usuarios y Familias:** Administración de perfiles de socios y sus grupos familiares.
-   **Membresías y Suscripciones:** Manejo de diferentes planes de membresía, suscripciones de socios y facturación.
-   **Instalaciones y Reservas:** Gestión del catálogo de instalaciones, consulta de disponibilidad en tiempo real y un completo sistema de reservas (incluyendo listas de espera).
-   **Disciplinas y Grupos:** Administración de las disciplinas deportivas del club y los grupos de entrenamiento.
-   **Campeonatos y Equipos:** Un sistema completo para crear y gestionar torneos, incluyendo inscripción de equipos, programación de partidos y tablas de posiciones.
-   **Control de Asistencia:** Herramientas para que los entrenadores registren la asistencia a las clases.
-   **Control de Acceso Físico:** Lógica para validar la entrada a las instalaciones a través de dispositivos como lectores QR.
-   **Tienda y Punto de Venta:** Una tienda integrada para vender merchandising, productos del buffet y más.
-   **Pagos y Billetera Virtual:** Integración con pasarelas de pago y gestión de una billetera virtual para cada socio.
-   **Gamificación (En desarrollo):** Sistema de Puntos de Experiencia (XP) y niveles para incentivar la participación.
-   **Notificaciones:** Servicio centralizado para enviar comunicaciones a los socios (Email, SMS).

## 🛠️ Stack Tecnológico

-   **Backend**: Go, Gin, GORM, PostgreSQL con PgVector.
-   **Frontend**: Next.js (TypeScript), React, Tailwind CSS.
-   **Infraestructura**: Docker, Docker Compose, Redis.

## 🏁 Guía de Inicio Rápido

Esta guía te permitirá levantar un entorno de desarrollo completo en tu máquina local.

### Prerrequisitos
-   Docker y Docker Compose.

### Pasos para la Instalación

1.  **Clona el repositorio**:
    ```bash
    git clone https://github.com/lukcba-developers/club-pulse-system-api.git
    cd club-pulse-system-api
    ```

2.  **Levanta los servicios con Docker Compose**:
    Este comando construirá y levantará los contenedores para el backend, frontend, base de datos y Redis.
    ```bash
    docker-compose up --build
    ```
    -   **Frontend**: Accesible en `http://localhost:3000`
    -   **Backend API**: Accesible en `http://localhost:8081`

3.  **Ejecuta las migraciones de la base de datos**:
    En una nueva terminal, ejecuta el siguiente comando para crear todas las tablas necesarias.
    ```bash
    docker-compose exec api go run ./cmd/migrate
    ```

4.  **(Opcional) Puebla la base de datos con datos de prueba**:
    Para tener datos iniciales (usuarios, clubes, etc.) y poder probar la aplicación, ejecuta el "seeder".
    ```bash
    docker-compose exec api go run ./cmd/seeder
    ```

Ahora, el entorno está listo. Puedes acceder a `http://localhost:3000` y usar las credenciales de prueba que pueda crear el seeder.

## 📂 Estructura del Proyecto
```
.
├── backend/                # API Monolítica en Go
│   ├── cmd/                # Puntos de entrada (api, migrate, seeder)
│   └── internal/
│       ├── core/           # Arquitectura central (puertos, errores)
│       ├── modules/        # Módulos de negocio (booking, user, store, etc.)
│       └── platform/       # Implementaciones de servicios (DB, logger, etc.)
├── frontend/               # Aplicación Next.js (App Router)
├── docs/
│   └── wiki/user/          # Wiki de usuario detallada por módulo
├── scripts/                # Scripts de utilidad
├── CONTRIBUTING.md         # Guía técnica para nuevos desarrolladores
├── tasks.md                # Lista de deuda técnica y mejoras pendientes
└── docker-compose.yml      # Orquestación de servicios Docker
```

## 📚 Documentación

La documentación es una pieza clave de este proyecto. Está diseñada para ser clara, completa y útil tanto para desarrolladores como para usuarios finales.

-   **[Wiki de Usuario](docs/wiki/user/README.md):** **(Lectura recomendada)** Es la fuente central de conocimiento sobre la funcionalidad del sistema. Detalla cada módulo de negocio, explicando su propósito, características y flujos de trabajo desde la perspectiva del usuario.
-   **Documentación Técnica por Módulo:** Cada módulo en `backend/internal/modules/` cuenta con su propio `README.md` detallando arquitectura, reglas de negocio y snippets de uso para desarrolladores.
-   **[Diagramas de Arquitectura](docs/architecture/diagrams.md):** Visualización de flujos críticos (Auth, Reservas, Semáforo del Jugador).

-   **[Guía para Contribuidores (`CONTRIBUTING.md`)](CONTRIBUTING.md):** **(Lectura obligatoria para desarrolladores)** Contiene la guía de arquitectura, configuración del entorno y el flujo de trabajo para añadir nuevas funcionalidades.

-   **[Lista de Tareas (`tasks.md`)](tasks.md):** Un listado de deuda técnica, funcionalidades incompletas y sugerencias de mejora para guiar el futuro desarrollo.

## 🧪 Pruebas

Para ejecutar la suite de pruebas de integración del backend, utiliza el siguiente script:
```bash
./scripts/integration_test.sh
```
