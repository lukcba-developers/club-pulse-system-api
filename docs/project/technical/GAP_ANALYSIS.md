# Análisis de Paridad de Funcionalidades (Gap Analysis)
*Comparativa entre `club-management-system-api` (Legacy) y `club-pulse-system-api` (MVP Actual)*

Este documento detalla las funcionalidades que han sido migradas y aquellas que aún faltan para alcanzar la paridad completa con el sistema legado.

> [!IMPORTANT]
> La hoja de ruta para implementar las funcionalidades faltantes se encuentra en [POST_MVP_MIGRATION_PLAN.md](./POST_MVP_MIGRATION_PLAN.md).

## 1. Módulo de Autenticación (Auth)
**Estado: MVP Básico Migrado**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado | Fase Planificada |
| :--- | :---: | :---: | :---: | :---: |
| Login (Email/Pass) | ✅ | ✅ | Implementado | - |
| Registro | ✅ | ✅ | Implementado | - |
| **OAuth (Google)** | ✅ | ❌ | **Faltante** | **Fase 1** |
| Refresh Token | ✅ | ✅ | Implementado | - |
| Validación de Token | ✅ | ✅ | Implementado | - |
| Logout | ✅ | ✅ | Implementado | - |

## 2. Módulo de Usuarios (User)
**Estado: CRUD Básico Migrado**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado | Fase Planificada |
| :--- | :---: | :---: | :---: | :---: |
| CRUD Básico | ✅ | ✅ | Implementado | - |
| Búsqueda Avanzada | ✅ | ✅ | Implementado | - |
| **Perfil Deportivo** | ✅ | ✅ | **Implementado** (JSONB Schema) | **Fase 1 (Completada)** |
| **Estadísticas** | ✅ | ❌ | **Faltante** | **Fase 1** |
| **Gamification** | ✅ (Complex) | ❌ | **Faltante** (Wallet, Challenges) | **Fase 3** |

## 3. Módulo de Instalaciones (Facilities)
**Estado: Núcleo Migrado**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado | Fase Planificada |
| :--- | :---: | :---: | :---: | :---: |
| CRUD Canchas | ✅ | ✅ | Implementado | - |
| Equipamiento | ✅ | ✅ | Implementado | - |
| **Mantenimiento** | ✅ | ⚠️ | Parcial (Falta validación reservas) | **Fase 1 (Finalizar)** |
| Slots/Subscription | ✅ | ❌ | **Faltante** | **Fase 3** |

## 4. Módulo de Membresías (Membership)
**Estado: Simplificado para MVP**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado | Fase Planificada |
| :--- | :---: | :---: | :---: | :---: |
| Planes (Tiers) | ✅ | ✅ | Implementado | - |
| **Integración Pagos** | ✅ (Payments API) | ✅ | **Implementado** (Entidad Payment + Mock MP) | **Fase 2 (Completada)** |
| **POS Device** | ✅ | ❌ | **Faltante** | **Fase 3** |
| **Wallet/Virtual** | ✅ | ❌ | **Faltante** | **Fase 3** |
| **Facturación** | ✅ | ✅ | **Implementado** (BillingCycle/Fees) | **Fase 2 (Completada)** |
| **Reportes** | ✅ | ❌ | **Faltante** | **Fase 2** |
| Multi-Club | ✅ | ❌ | Simplificado | **Fase 3** |

## 5. Módulo de Reservas (Booking)
**Estado: Funcionalidad Core Migrada**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado | Fase Planificada |
| :--- | :---: | :---: | :---: | :---: |
| Crear Reserva | ✅ | ✅ | Implementado | - |
| Cancelar Reserva | ✅ | ✅ | Implementado | - |
| Pago de Reserva | ✅ | ✅ | **Implementado** (Via Payment Module) | **Fase 2 (Completada)** |

## Resumen de Prioridades
1.  **Fase 1**: Auth (OAuth) + User (Profile/Stats) + Facilities (Maintenance).
2.  **Fase 2**: Membership (Payments/Billing) - **CRÍTICO PARA NEGOCIO**.
3.  **Fase 3**: Facilities (Recurring) + Membership (Access/Multi-club).
