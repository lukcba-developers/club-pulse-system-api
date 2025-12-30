# Documentación de la API del Sistema Club Pulse

## Visión General
Esta API impulsa el Sistema Club Pulse, gestionando la autenticación de usuarios, perfiles, instalaciones, membresías y reservas.
Está construida siguiendo una Arquitectura Limpia (Clean Architecture) y un patrón de Monolito Modular.

- **Base URL**: `http://localhost:8081/api/v1`
- **Autenticación**: Bearer Token (JWT)

## 1. Módulo de Autenticación (Auth)

### Registro de Usuario
Crea una nueva cuenta de usuario.
- **Endpoint**: `POST /auth/register`
- **Público**
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securePassword123",
    "name": "John Doe",
    "role": "MEMBER" // Opcional, por defecto MEMBER
  }
  ```

### Iniciar Sesión
Autentica a un usuario y devuelve un token JWT.
- **Endpoint**: `POST /auth/login`
- **Público**
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securePassword123"
  }
  ```
- **Respuesta**:
  ```json
  {
    "access_token": "eyJhbGciOi...",
    "expires_in": 86400
  }
  ```

## 2. Módulo de Usuarios (Users)

### Obtener Perfil Actual
Obtiene los detalles del usuario autenticado actual.
- **Endpoint**: `GET /users/me`
- **Headers**: `Authorization: Bearer <token>`
- **Respuesta**:
  ```json
  {
    "id": "uuid...",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "MEMBER",
    "created_at": "..."
  }
  ```

## 3. Módulo de Instalaciones (Facilities)

### Listar Instalaciones
Obtiene una lista de todas las instalaciones deportivas disponibles.
- **Endpoint**: `GET /facilities`
- **Público**
- **Respuesta**:
  ```json
  {
    "data": [
      {
        "id": "uuid...",
        "name": "Cancha de Tenis 1",
        "type": "court",
        "status": "active",
        "hourly_rate": 50.00,
        "location": { "name": "Principal" },
        "specifications": { "surface": "clay" }
      }
    ]
  }
  ```

## 4. Módulo de Membresías (Memberships)

### Listar Niveles (Tiers)
Obtiene todos los planes de membresía disponibles.
- **Endpoint**: `GET /memberships/tiers`
- **Headers**: `Authorization: Bearer <token>`
- **Respuesta**:
  ```json
  {
    "data": [
      {
        "id": "uuid...",
        "name": "Gold",
        "monthly_fee": 99.99,
        "benefits": ["Acceso Total", "Sauna"]
      }
    ]
  }
  ```

### Obtener Mi Membresía
Obtiene el estado de la suscripción del usuario actual.
- **Endpoint**: `GET /memberships`
- **Headers**: `Authorization: Bearer <token>`

## 5. Módulo de Reservas (Bookings)

### Crear Reserva
Reserva una instalación para un horario específico.
- **Endpoint**: `POST /bookings`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": "uuid-del-usuario",
    "facility_id": "uuid-de-instalacion",
    "start_time": "2025-10-20T10:00:00Z",
    "end_time": "2025-10-20T11:00:00Z"
  }
  ```
- **Errores**:
  - `409 Conflict`: Si el horario ya está ocupado por otra reserva o por mantenimiento programado.

### Listar Mis Reservas
Obtiene todas las reservas del usuario actual (o todas si es Admin).
- **Endpoint**: `GET /bookings`
- **Headers**: `Authorization: Bearer <token>`

### Cancelar Reserva
Cancela una reserva existente.
- **Endpoint**: `DELETE /bookings/:id`
- **Headers**: `Authorization: Bearer <token>`

### Consultar Disponibilidad
Obtiene los slots de tiempo bloqueados por reservas o mantenimiento para una fecha e instalación específica.
- **Endpoint**: `GET /bookings/availability`
- **Headers**: `Authorization: Bearer <token>`
- **Query Params**:
  - `facility_id`: UUID de la instalación.
  - `date`: Fecha en formato YYYY-MM-DD.
- **Respuesta**:
  ```json
  {
    "data": [
       // Lista de bloqueos (reservas o mantenimiento)
       // Nota: Formato específico dependerá de la implementación final
    ]
  }
  ```
