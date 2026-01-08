# 🏗️ Diagramas de Arquitectura

Este documento visualiza los flujos de datos y lógica de negocio más críticos del sistema **Club Pulse**.

---

## 1. Registro y Login con Google
Este flujo describe la integración con Google OAuth2 y la creación automática de usuarios asociados a un club.

```mermaid
sequenceDiagram
    participant User as Socio (Frontend)
    participant Auth as Auth Module (Backend)
    participant Google as Google API
    participant DB as Postgres

    User->>Google: 1. Click "Login con Google"
    Google-->>User: 2. Retorna Auth Code
    User->>Auth: 3. POST /auth/google {code}
    Auth->>Google: 4. Intercambia Code por Token & Profiler
    Google-->>Auth: 5. Retorna Email & ID
    Auth->>DB: 6. Busca usuario por Email + ClubID
    alt El usuario no existe
        Auth->>DB: 7. Crea nuevo Usuario (Auto-registro)
    end
    Auth->>Auth: 8. Genera JWT & Refresh Token
    Auth-->>User: 9. Set HttpOnly Cookies & 200 OK
```

---

## 2. Ciclo de Vida de una Reserva
Flujo simplificado de una reserva desde la selección hasta la confirmación de pago por Webhook.

```mermaid
graph TD
    A[Socio: Selecciona Pista y Hora] --> B{¿Es de Pago?}
    B -- Sí --> C[Backend: Crea Reserva PENDING_PAYMENT]
    C --> D[Frontend: Redirige a Mercado Pago]
    D --> E[Socio: Completa Pago]
    E -.-> F[Mercado Pago: Notifica Webhook]
    F --> G[Payment Module: Valida Pago]
    G --> H[Booking Module: Cambia Estado a CONFIRMED]
    H --> I[Socio: Recibe Notificación de Éxito]
    B -- No --> J[Backend: Crea Reserva CONFIRMED]
    J --> I
```

---

## 3. El Semáforo del Jugador
Visualización de cómo el `TeamModule` orquesta datos de otros módulos para habilitar a un jugador.

```mermaid
graph LR
    subgraph "Socio Habilitado"
        A[Finanzas: Balance <= 0]
        B[Médico: EMMAC Válido]
    end

    A & B --> C{Regla AND}
    C -- TRUE --> D(JUGADOR HABILITADO)
    C -- FALSE --> E(JUGADOR INHABILITADO)

    subgraph "Fuentes de Datos"
        M[Membership Module] --> A
        U[User Documents] --> B
    end

    E --> F[Inconvocable en Torneos]
    E --> G[Bloqueado en Reservas]
```

---
⚠️ *Estos diagramas representan la lógica implementada hasta la fecha. Para detalles de implementación técnica, consulte los README.md de cada módulo.*
