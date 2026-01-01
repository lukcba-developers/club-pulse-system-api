# Contexto de Migración y Arquitectura del Sistema

## 1. Objetivo del Proyecto
El objetivo es migrar el sistema legado `club-management-system-api` (microservicios) a una nueva arquitectura llamada `club-pulse-system-api` (Monolito Modular) para mejorar la mantenibilidad, reducir costos de infraestructura (MVP) y aplicar Clean Code y mejores prácticas de Go.

## 2. Estado Actual (Backend)
El backend se encuentra en `/Users/lukcba-macbook-pro/go/src/club-pulse-system-api/backend`.
Está construido con **Go**, **Gin**, **GORM** y **PostgreSQL**.

### Arquitectura (Clean Architecture)
El proyecto sigue una estructura de capas estricta:
- **Domain**: Entidades y definiciones de interfaces (Core Business Logic). Independiente de frameworks.
- **Application**: Casos de uso (Lógica de la aplicación). Orquesta el flujo de datos.
- **Infrastructure**: Implementaciones concretas (Repositorios DB, Handlers HTTP, Servicios Externos).

### Módulos Implementados
1.  **Auth Module**:
    - Gestión de Login/Registro.
    - JWT para seguridad.
    - Repositorio Postgres.
2.  **User Module**:
    - Gestión de perfil de usuario.
    - Integrado con Auth Middleware.

3.  **Facilities Module**:
    - CRUD completo de instalaciones.
    - Soporte para specs en JSONB.
    - Repositorio Postgres.
4.  **Membership Module**:
    - Gestión de Membresías y Niveles (Tiers).
    - Lógica de facturación y estado.
5.  **Booking Module**:
    - Motor de reservas completo.
    - Detección de conflictos (Overlap).
    - Estado de reserva (CONFIRMED/CANCELLED).
6.  **Payment Module**:
    - Procesamiento de pagos (Mock MercadoPago).
    - Webhooks y conciliación.
    - Repositorio independiente.

## 3. Análisis de Legacy (club-management-system-api)
Se ha analizado el código fuente original para extraer la lógica de negocio y replicarla con mejoras.

### Facilities API (Instalaciones)
- **Estado**: Migrado y Operativo.

## 4. Próximos Pasos (Completados)
- [x] Migración del módulo de **Booking** (Último bloque crítico).
- [x] **Fase 1**: Core de Usuarios (Categorías y Perfil Deportivo).
- [x] **Fase 2**: Integración de Pagos y Membresías Avanzadas.
- [x] Integración final del Frontend (PricingCards y Facilities funcionan con datos reales).
- [x] Pruebas E2E completas:
    - [x] Conflicto de reservas validado.
    - [x] Listado de instalaciones validado.
    - [x] Sistema de Membresías (Tiers) validado.
    - [x] Sistema de Membresías (Tiers) validado.
    - [x] Seeders creados.
- [x] Funcionalidades de Cierre (Gap Analysis):
    - [x] Búsqueda Avanzada de Usuarios.
    - [x] Gestión de Sesiones y Auditoría Auth.
    - [x] Frontend para Gestión de Usuarios y Perfil.

El sistema backend ahora es completamente funcional como un Monolito Modular.
El frontend puede consumir datos reales.
