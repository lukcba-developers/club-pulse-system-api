# Plan de Migración e Implementación Post-MVP

Este documento detalla la hoja de ruta para evolucionar `club-pulse-system-api` desde su estado actual de MVP hacia un sistema completo de gestión de clubes, alcanzando la paridad con el sistema legado (`club-management-system-api`) e introduciendo nuevas capacidades.

**Fuente de Verdad para Migración:**
La lógica de negocio compleja, modelos de datos extendidos y flujos de trabajo administrativos deben ser migrados desde el repositorio legado:
`/Users/lukcba-macbook-pro/go/src/club-management-system-api/`

---

## Estrategia de Fases

Para garantizar estabilidad y entrega continua de valor, la implementación se dividirá en 3 fases estratégicas.

### Fase 1: Enriquecimiento de Experiencia de Usuario y "Core" (Semanas 1-2)
**Objetivo:** Mejorar la retención de usuarios, facilitar el acceso y completar el modelo de datos del usuario **(Enfoque Club Social)**.

| Módulo | Funcionalidad | Descripción Técnica | Fuente Legacy (Referencia) | Estado | Fase |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Auth** | **OAuth 2.0 (Google)** | Implementar login social para reducir fricción. <br>- Integrar `golang.org/x/oauth2`. <br>- Crear `link_social_account` usecase. | `legacy/auth/service/google_service.go` | | |
| **User** | **Categorización por Edad** | ✅ **[COMPLETADO]** Lógica automática: `DateOfBirth` -> Categoría (ej. "2012"). <br>- Implementado `CalculateCategory()` y tests. | `legacy/user/domain/model/user.go` | | |
| **User** | **Perfil Deportivo** | ✅ **[COMPLETADO]** Extendido `User` con `SportsPreferences` (JSONB). <br>- Permite flexibilidad sin tablas extra por ahora. | `legacy/user/dto/sport_profile_dto.go` | | |
| **Facilities** | **Mantenimiento** | ✅ | ✅ | **Implementado** (Blocked Slots) | **Fase 1 (Completada)** |
| Slots/Subscription | ✅ | ⚠️ | **En Progreso** (Entidad Prep) | **Fase 4** |

### Fase 2: Monetización y Operaciones Financieras (Semanas 3-4)
**Objetivo:** Habilitar el cobro real de **Cuotas Sociales** y reservas. Es la fase más crítica para la viabilidad del negocio.

| Módulo | Funcionalidad | Descripción Técnica | Fuente Legacy (Referencia) |
| :--- | :--- | :--- | :--- |
| **Membership** | **Gestión de Cuotas** | ✅ **[COMPLETADO]** Ciclos de facturación (`Monthly/Annual`) y `OutstandingBalance`. <br>- Implementado `CalculateLateFee`. | `legacy/membership/domain/entity/membership.go` |
| **Membership** | **Integración de Pagos** | ✅ **[COMPLETADO]** Entidad `Payment` y Mock de MercadoPago. <br>- Webhooks listos para procesar notificaciones. | `legacy/payments-api` |
| **Disciplines** | **Grupos de Entrenamiento** | **[NUEVO]** Crear entidad `TrainingGroup` (o adaptar `Championship`). <br>- Asignar usuarios a grupos por categoría. | `legacy/championship-api` (Adaptar) |

### Fase 3: Operaciones Avanzadas y Escalamiento (Semanas 5-6)
**Objetivo:** Soportar operaciones complejas de clubes grandes y múltiples sedes.

| Módulo | Funcionalidad | Descripción Técnica | Fuente Legacy (Referencia) |
| :--- | :--- | :--- | :--- |
| **Facilities** | **Slots Recurrentes** | Lógica para "Reservas Fijas" (ej. Escuela de fútbol, torneos). <br>- Conflict checking recurrente. | `legacy/booking/recurring` |
| **Membership** | **Control de Acceso** | ✅ **[COMPLETADO]** Sistema para validar entrada al club (QR Code dinámico). <br>- Endpoint `POST /access/entry`. | `legacy/access_control` |
| **Membership** | **Multi-Sede** | Adaptar DB para `club_id` en todas las tablas principales. <br>- Refactor masivo de queries (Tenant Isolation). | `legacy/multi_tenant` |
| **General** | **Notificaciones** | Sistema de email/push para recordatorios de reserva y pagos. <br>- Integrar con SendGrid/Firebase. | `legacy/notifications` |

---

## Detalle de Implementación por Módulo

### 1. Módulo de Autenticación (`/internal/modules/auth`)
- **Estado Actual:** Login/Registro básico (JWT).
- **Faltante:** OAuth, Recuperación de contraseña.
- **Plan:**
    1. Instalar `golang.org/x/oauth2`.
    2. Crear `auth/delivery/http/oauth_handler.go`.
    3. Migrar lógica de mapeo de usuario de Google a User interno.

### 2. Módulo de Usuarios (`/internal/modules/user`)
- **Estado Actual:** CRUD básico.
- **Faltante:** Perfil rico, Stats, Gamification basics.
- **Plan:**
    1. **Entidad Sport Profile**: Migrar DTO `SportProfileRequest` (Legacy: `user-api/internal/interfaces/api/dto`).
       - Campos: `Sports` ([]string), `PreferredPositions` (map), `SkillLevel` (map).
    2. **Gamification (Fase 1.5)**: Preparar esquema para `ChampionshipRecord` y `Wallet` (Legacy: `gamification.go`, `wallet.go`).
    3. Crear migraciones SQL para tablas `user_details`, `user_stats`.
    4. Exponer nuevos campos en `UserResponse`.

### 3. Módulo de Instalaciones (`/internal/modules/facilities`)
- **Estado Actual:** CRUD, Equipamiento (parcial).
- **Faltante:** Slots complejos, Bloqueo por mantenimiento (lógica de negocio).
- **Plan:**
    1. Asegurar que `CreateBooking` consulte tabla de `Maintenance`.
    2. Implementar `RecurringBookingStrategy`.

### 4. Módulo de Membresías (`/internal/modules/membership`) & Pagos
- **Estado Actual:** Planes estáticos.
- **Faltante:** **Motor de Pagos Completo** (Legacy: `payments-api`).
- **Plan (Basado en `payments-api` Legacy):**
    1. **Portar Arquitectura Clean**: `payments-api` tiene una estructura robusta (Ports & Adapters).
    2. **Entidad Payment**:
       - Status: `pending`, `approved`, `rejected`.
       - Metadatos: `ExternalID`, `POSDeviceID` (para integración futura).
    3. **Implementar Procesadores**:
       - `MercadoPagoProcessor` (Ref: `legacy/user-api/internal/infrastructure/payment`).
       - `StripeProcessor`.
    4. **Webhooks**: Endpoint unificado `POST /payments/webhook` que despacha a handlers específicos.

---

## Próximos Pasos Inmediatos (Para el Desarrollador)

1. **Aprobar este plan.**
2. Crear los tickets/issues correspondientes a la **Fase 1**.
3. Ejecutar migración del esquema de base de datos para `User Profile`.
