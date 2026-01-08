# Documentación de la API del Sistema Club Pulse

## Visión General
Esta API impulsa el Sistema Club Pulse, gestionando la autenticación de usuarios, perfiles, instalaciones, membresías y reservas.
Está construida siguiendo una Arquitectura Limpia (Clean Architecture) y un patrón de Monolito Modular.

- **Base URL**: `http://localhost:8081/api/v1`
- **Autenticación**: El sistema utiliza **Cookies HttpOnly** (`access_token` y `refresh_token`). No es necesario incluir el header manual `Authorization` en el navegador después del login.

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

### Iniciar Sesión (Login)
Autentica a un usuario y establece cookies seguras.
- **Endpoint**: `POST /auth/login`
- **Público**
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "securePassword123"
  }
  ```
- **Respuesta**: `200 OK`
  ```json
  { "message": "Login successful" }
  ```
- **Cookies**: Se establecen `access_token` (24h) y `refresh_token` (7d).

### Login Social (Google)
Inicia sesión utilizando un código de intercambio de Google OAuth.
- **Endpoint**: `POST /auth/google`
- **Público**
- **Body**: `{ "code": "4/0A..." }`

### Ver Sesiones Activas
Lista todas las sesiones activas del usuario actual.
- **Endpoint**: `GET /auth/sessions`

### Revocar Sesión
Cierra una sesión específica.
- **Endpoint**: `DELETE /auth/sessions/:id`

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

### Becas (Scholarships) - Admin
Gestiona descuentos especiales para socios.
- **Endpoint**: `POST /memberships/scholarships`
- **Body**:
  ```json
  {
    "user_id": "uuid",
    "percentage": 50,
    "reason": "Deportista Federado",
    "valid_until": "2025-12-31T23:59:59Z"
  }
  ```

## 5. Módulo de Reservas (Bookings)

### Crear Reserva
Reserva una instalación para un horario específico.
- **Endpoint**: `POST /bookings`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "facility_id": "uuid-de-instalacion",
    "start_time": "2025-10-20T10:00:00Z",
    "end_time": "2025-10-20T11:00:00Z",
    "guest_details": [
      {
        "name": "Nombre del Invitado",
        "dni": "12345678"
      }
    ]
  }
  ```
- **Notas**: El `user_id` se obtiene automáticamente del token de autenticación. `guest_details` es opcional. Se requiere **Certificado Médico** válido para procesar la reserva.
- **Errores**:
  - `409 Conflict`: Si el horario ya está ocupado.

### Listar Mis Reservas
Obtiene todas las reservas del usuario actual.
- **Endpoint**: `GET /bookings`
- **Headers**: `Authorization: Bearer <token>`

### Listar Todas las Reservas (Admin)
Obtiene todas las reservas del club, con opción de filtrar por fechas.
- **Endpoint**: `GET /bookings/all`
- **Headers**: `Authorization: Bearer <token>`
- **Rol Requerido**: `ADMIN` o `SUPER_ADMIN`
- **Query Params**:
  - `from`: Fecha de inicio en formato `YYYY-MM-DD`.
  - `to`: Fecha de fin en formato `YYYY-MM-DD`.

### Cancelar Reserva
Cancela una reserva existente. Al cancelar, se puede disparar el proceso de lista de espera.
- **Endpoint**: `DELETE /bookings/:id`
- **Headers**: `Authorization: Bearer <token>`

### Consultar Disponibilidad
Obtiene los slots de tiempo y su estado para una instalación y fecha.
- **Endpoint**: `GET /bookings/availability`
- **Headers**: `Authorization: Bearer <token>`
- **Query Params**:
  - `facility_id`: UUID de la instalación.
  - `date`: Fecha en formato `YYYY-MM-DD`.
- **Respuesta Exitosa (200 OK)**:
  ```json
  {
    "data": [
      {
        "start_time": "2025-10-20T09:00:00Z",
        "end_time": "2025-10-20T10:00:00Z",
        "available": true
      },
      {
        "start_time": "2025-10-20T10:00:00Z",
        "end_time": "2025-10-20T11:00:00Z",
        "available": false
      }
    ]
  }
  ```

### Unirse a la Lista de Espera
Añade al usuario a la lista de espera para una instalación en un horario específico.
- **Endpoint**: `POST /bookings/waitlist`
- **Body**:
  ```json
  {
    "resource_id": "uuid-de-instalacion",
    "target_date": "2025-10-20T10:00:00Z"
  }
  ```

### Reservas Recurrentes (Admin)
Define una regla de reserva fija (ej: todos los lunes).
- **Endpoint**: `POST /bookings/recurring-rules`
- **Body**:
  ```json
  {
    "facility_id": "uuid",
    "type": "WEEKLY",
    "day_of_week": 1,
    "start_time": "18:00",
    "end_time": "19:00",
    "start_date": "2025-01-01",
    "end_date": "2025-12-31"
  }
  ```


## 6. Módulo de Pagos (Payments)

### Iniciar Pago (Checkout)
Genera una preferencia de pago y devuelve la URL para redirigir al usuario.
- **Endpoint**: `POST /payments/checkout`
### Registrar Pago Offline (Admin)
Registra un pago realizado en efectivo o mediante intercambio de mano de obra.
- **Endpoint**: `POST /payments`
- **Body**:
  ```json
  {
    "amount": 5000.00,
    "method": "CASH", // "CASH", "LABOR_EXCHANGE", "TRANSFER"
    "payer_id": "uuid-del-socio",
    "reference_id": "uuid-factura",
    "notes": "Detalle del pago o trabajo realizado"
  }
  ```

### Integración Webhook (MercadoPago)
Endpoint para recibir notificaciones de estado de pagos digitales.
- **Endpoint**: `POST /payments/webhook`
- **Público**
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
- **Body**:
  ```json
  {
    "user_id": "student-uuid",
    "status": "PRESENT", // PRESENT, ABSENT, LATE, EXCUSED
    "notes": "Llegó tarde por tráfico"
  }
  ```

## 9. Módulo de Grupos Familiares (Family)

### Crear Grupo Familiar
Crea un grupo para vincular múltiples socios bajo una misma unidad de facturación.
- **Endpoint**: `POST /family/groups`
- **Body**: `{ "name": "Familia García" }`

### Invitar a Miembro
Añade un miembro al grupo familiar.
- **Endpoint**: `POST /family/groups/:id/members`
- **Body**: `{ "user_id": "uuid-miembro", "role": "MEMBER" }`
- **Respuesta**: `200 OK`
