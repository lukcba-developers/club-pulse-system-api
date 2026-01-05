# Guía para Contribuidores

¡Gracias por tu interés en contribuir a Club Pulse System! Esta guía proporciona toda la información necesaria para entender la arquitectura del proyecto y cómo puedes empezar a colaborar.

## 1. Configuración del Entorno de Desarrollo

El proyecto está completamente contenedorizado usando Docker, lo que simplifica enormemente la configuración inicial.

### Prerrequisitos

-   Tener instalado [Docker](https://www.docker.com/get-started) y Docker Compose.

### Pasos para la Instalación

1.  **Clonar el Repositorio:**
    ```bash
    git clone https://github.com/lukcba-developers/club-pulse-system-api.git
    cd club-pulse-system-api
    ```

2.  **Levantar los Servicios:**
    Todos los servicios (backend, frontend, base de datos y Redis) se levantan con un solo comando.
    ```bash
    docker-compose up --build
    ```
    -   El backend (`api`) estará disponible en `http://localhost:8081`.
    -   El frontend (`web`) estará disponible en `http://localhost:3000`.
    -   La base de datos (`postgres`) expondrá su puerto en `localhost:5432`.

3.  **Ejecutar las Migraciones:**
    Para crear la estructura de la base de datos, abre una nueva terminal y ejecuta el comando de migración dentro del contenedor del backend.
    ```bash
    docker-compose exec api go run ./cmd/migrate
    ```

4.  **(Opcional) Poblar la Base de Datos con Datos de Prueba:**
    Para tener datos iniciales con los que trabajar (usuarios, clubes, etc.), puedes ejecutar el "seeder".
    ```bash
    docker-compose exec api go run ./cmd/seeder
    ```

¡Y eso es todo! El entorno de desarrollo ya está completamente funcional.

## 2. Arquitectura del Proyecto

El proyecto es un monorepo que contiene dos aplicaciones principales: `backend` y `frontend`.

### Backend (Go)

El backend sigue una arquitectura hexagonal (también conocida como Puertos y Adaptadores), lo que permite un excelente desacoplamiento y testeabilidad.

-   `backend/cmd/`: Contiene los puntos de entrada de la aplicación (`api`, `migrate`, `seeder`).
-   `backend/internal/core/`: Define los elementos centrales de la arquitectura, como los `ports` (interfaces).
-   `backend/internal/platform/`: Implementaciones concretas de servicios transversales como `logger`, `database`, `redis`, `websocket`, etc.
-   `backend/internal/modules/`: **Aquí reside la lógica de negocio**. Cada módulo (ej: `booking`, `user`) está a su vez dividido en:
    -   `domain/`: Contiene los modelos de datos del negocio y las interfaces de los repositorios. Es el corazón del módulo.
    -   `application/`: Contiene los casos de uso, que orquestan la lógica de negocio.
    -   `infrastructure/`: Contiene las implementaciones concretas de las interfaces del dominio. Por ejemplo:
        -   `repository/`: Implementación del repositorio usando PostgreSQL.
        -   `http/`: `handler.go` que define los manejadores de la API y registra las rutas.

### Frontend (Next.js)

El frontend está construido con Next.js y sigue las convenciones modernas de React.

-   `frontend/app/`: Utiliza el App Router de Next.js para la estructura de rutas.
-   `frontend/components/`: Componentes de UI reutilizables.
-   `frontend/services/`: Contiene funciones que realizan las llamadas a la API del backend. Cada servicio (ej: `booking-service.ts`) se corresponde con un módulo del backend.
-   `frontend/lib/`: Utilidades generales, como la instancia configurada de Axios.
-   `frontend/context/`: Contextos de React para el manejo de estado global (ej: `auth-context.tsx`).

## 3. Cómo Añadir una Nueva Funcionalidad (Ejemplo)

Imaginemos que queremos añadir un endpoint `GET /facilities/{id}/maintenance-history`.

1.  **Dominio (`domain`):**
    -   Añadir la estructura `Maintenance` en un archivo `maintenance.go`.
    -   Añadir el método `ListMaintenanceHistory(facilityID string)` a la interfaz `FacilityRepository`.

2.  **Infraestructura (`infrastructure`):**
    -   Implementar el nuevo método en `repository/postgres.go`, escribiendo la consulta a la base de datos.

3.  **Aplicación (`application`):**
    -   Crear un nuevo caso de uso `ListMaintenanceHistory` que llame al método del repositorio.

4.  **Handler (`infrastructure/http`):**
    -   Crear una nueva función `ListMaintenanceHistory(c *gin.Context)` en `handler.go`.
    -   Esta función llamará al caso de uso y devolverá el resultado como JSON.
    -   Registrar la nueva ruta `GET /:id/maintenance-history` en la función `RegisterRoutes`.

5.  **Frontend (`services`):**
    -   Añadir una nueva función `getMaintenanceHistory(facilityId)` en `facility-service.ts` que llame al nuevo endpoint.

Este flujo de trabajo asegura que el código se mantenga organizado, desacoplado y fácil de mantener.
