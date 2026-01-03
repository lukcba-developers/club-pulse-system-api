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
    "category": "1990 (Master)",
    "date_of_birth": "1990-01-01T00:00:00Z",
    "sports_preferences": {
      "main_sport": "tennis",
      "level": "intermediate"
    },
    "created_at": "..."
  }
  ```

### Actualizar Perfil
Actualiza la información personal del usuario.
- **Endpoint**: `PUT /users/me`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "name": "John Doe Updated",
    "date_of_birth": "1990-01-01T00:00:00Z",
    "sports_preferences": {
      "main_sport": "tennis",
      "hand": "right"
    }
  }
  ```
- **Respuesta**: `200 OK` (Devuelve el objeto User actualizado)

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

## 6. Módulo de Pagos (Payments)

### Iniciar Pago (Checkout)
Genera una preferencia de pago y devuelve la URL para redirigir al usuario.
- **Endpoint**: `POST /payments/checkout`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "amount": 5000.00,
    "currency": "ARS",
    "description": "Cuota Mensual - Enero",
    "payer_email": "user@example.com",
    "reference_id": "membership-uuid", // ID de lo que se paga
    "reference_type": "MEMBERSHIP"
  }
  ```
- **Respuesta**:
  ```json
  {
    "init_point": "https://www.mercadopago.com.ar/checkout/v1/redirect/...",
    "preference_id": "..."
  }
  ```

### Integración Webhook (MercadoPago)
Endpoint para recibir notificaciones de estado de pagos.
- **Endpoint**: `POST /payments/webhook`
- **Público** (Debe validar firma en producción)
- **Query Params**:
  - `type`: Tipo de evento (ej. `payment`).
- **Body**: Payload del proveedor de pagos.
- **Respuesta**: `200 OK`

## 7. Módulo de Control de Acceso (Access)

### Registrar Ingreso/Egreso
Valida y registra el acceso de un usuario (Simulación QR).
- **Endpoint**: `POST /access/entry`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": "uuid-del-usuario",
    "facility_id": "uuid-instalacion-opcional",
    "direction": "IN" // IN o OUT
  }
  ```
- **Respuesta Exitoso (200)**:
  ```json
  {
    "status": "GRANTED",
    "user_id": "...",
    "timestamp": "..."
  }
  ```
- **Respuesta Denegado (403)**:
  ```json
  {
    "status": "DENIED",
    "reason": "Membership Inactive / Debt"
  }
  ```

## 8. Módulo de Asistencia (Attendance)

### Obtener Lista de Asistencia (Vista Coach)
Obtiene o crea la lista de asistencia para un grupo/categoría en una fecha.
- **Endpoint**: `GET /attendance/groups/:group`
- **Headers**: `Authorization: Bearer <token>`
- **Query Params**:
  - `date`: YYYY-MM-DD (Opcional, default hoy)
- **Ejemplo**: `GET /attendance/groups/2012?date=2025-10-20`
- **Respuesta**:
  ```json
  {
    "id": "list-uuid",
    "group": "2012",
    "date": "2025-10-20T00:00:00Z",
    "records": [
      {
        "user_id": "student-uuid",
        "status": "ABSENT",
        "user": { "name": "..." } // Si se incluye
      }
    ]
  }
  ```

### Marcar Asistencia
Actualiza el estado de un alumno en la lista.
- **Endpoint**: `POST /attendance/:listID/records`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "user_id": "student-uuid",
    "status": "PRESENT", // PRESENT, ABSENT, LATE, EXCUSED
    "notes": "Llegó tarde por tráfico"
  }
  ```
- **Respuesta**: `200 OK`
