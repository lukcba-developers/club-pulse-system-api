# Análisis de Paridad de Funcionalidades (Gap Analysis)
*Comparativa entre `club-management-system-api` (Legacy) y `club-pulse-system-api` (MVP Actual)*

Este documento detalla las funcionalidades que han sido migradas y aquellas que aún faltan para alcanzar la paridad completa con el sistema legado.

## 1. Módulo de Autenticación (Auth)
**Estado: MVP Básico Migrado**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado |
| :--- | :---: | :---: | :---: |
| Login (Email/Pass) | ✅ | ✅ | Implementado |
| Registro | ✅ | ✅ | Implementado |
| **OAuth (Google)** | ✅ | ❌ | **Faltante** |
| Refresh Token | ✅ | ✅ | Implementado |
| Validación de Token | ✅ | ✅ | Implementado (Middleware) |
| Logout | ✅ | ✅ | Implementado |
| Gestión de Sesiones | ✅ | ✅ | Implementado (List/Revoke) |
| Admin / Analytics | ✅ | ✅ | Implementado (Auth Logs) |

---

## 2. Módulo de Usuarios (User)
**Estado: CRUD Básico Migrado**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado |
| :--- | :---: | :---: | :---: |
| Crear Perfil | ✅ | ✅ | Implementado |
| Ver Perfil | ✅ | ✅ | Implementado |
| Actualizar Perfil | ✅ | ✅ | Implementado |
| Listar (Paginado) | ✅ | ✅ | Implementado |
| **Borrar Usuario** | ✅ | ✅ | Implementado |
| Búsqueda Avanzada | ✅ | ✅ | Implementado (Nombre/Email) |
| Perfil Deportivo | ✅ | ❌ | **Faltante** |
| Estadísticas | ✅ | ❌ | **Faltante** |

---

## 3. Módulo de Instalaciones (Facilities)
**Estado: Núcleo Migrado**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado |
| :--- | :---: | :---: | :---: |
| Crear/Ver/Listar | ✅ | ✅ | Implementado |
| **Gestión de Equipamiento** | ✅ | ✅ | Implementado (Entidad + CRUD) |
| **Mantenimiento** | ✅ | ✅ | Implementado (Tareas + Conflicto con Reservas) |
| Slots/Subscription | ✅ | ❌ | **Faltante** |

---

## 4. Módulo de Membresías (Membership)
**Estado: Simplificado para MVP**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado |
| :--- | :---: | :---: | :---: |
| Planes (Tiers) | ✅ | ✅ | Implementado |
| Membresía Usuario | ✅ | ✅ | Implementado |
| **Multi-Club** | ✅ | ❌ | Simplificado a Club Único |
| Facturación/Pagos | ✅ | ❌ | **Integración de Pagos Faltante** |
| Control de Acceso | ✅ | ❌ | Faltante (Suspend/Activate) |
| Reportes Financieros | ✅ | ❌ | **Faltante** |

---

## 5. Módulo de Reservas (Booking)
**Estado: Funcionalidad Core Migrada**

| Funcionalidad | Legacy | Club Pulse (Nuevo) | Estado |
| :--- | :---: | :---: | :---: |
| Crear Reserva | ✅ | ✅ | Implementado |
| Buscar Conflicto | ✅ | ✅ | Implementado |
| Listar Reservas | ✅ | ✅ | Implementado |
| Cancelar | ✅ | ✅ | Implementado |
| Completar (Flow) | ✅ | ✅ | Implementado (Lógica Backend) |
| Disponibilidad (API)| ✅ | ✅ | Implementado (Endpoint /availability) |

## Conclusión y Recomendación

Hemos migrado exitosamente el **"Core MVP"** que permite el flujo principal: **Registrarse -> Ver Instalaciones -> Contratar Membresía -> Reservar**.

Sin embargo, el sistema legado contiene lógica de negocio avanzada que no ha sido portada, principalmente:
1.  **Operaciones Administrativas**: Mantenimiento, Equipamiento, Analytics, Reportes.
2.  **Robustez de Auth**: Recuperación de contraseña, OAuth, Gestión de Sesiones.
3.  **Monetización Real**: Integración con pasarelas de pago (MercadoPago fue visto en el código legacy).

**Siguientes Pasos Recomendados:**
1.  Dependiendo del objetivo inmediato ("MVP vs Feature Parity"), decidir si migrar **Pagos** y **OAuth** es prioritario.
2.  Módulos como **Mantenimiento** y **Equipamiento** pueden posponerse para una fase 2.
