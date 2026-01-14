# Plan de Pruebas de Navegabilidad E2E (End-to-End)

Este documento define la suite completa de pruebas para verificar la navegabilidad y funcionalidad de la aplicación `club-pulse-system-api`. Las pruebas están organizadas por **Roles de Usuario** y **Módulos Funcionales**, basándose en la estructura real del proyecto.

---

## 1. Matriz de Roles y Permisos

| Rol | Descripción | Alcance Principal |
| :--- | :--- | :--- |
| **`SUPER_ADMIN`** | Administrador de plataforma | Gestión multi-tenant, clubes, configuración global |
| **`ADMIN`** | Administrador de club | Instalaciones, usuarios, reservas, pagos, campeonatos, config del club |
| **`COACH`** | Entrenador | Gestión de equipos, asistencia, viajes, jugadores |
| **`MEMBER`** | Socio/Jugador | Reservas, perfil, membresía, tienda, campeonatos |
| **`MEDICAL_STAFF`** | Personal médico | Acceso a datos sensibles (GDPR Art. 9) |
| **`GUEST`** | Invitado/Público | Registro, páginas legales, landing de club |

---

## 2. Suite de Pruebas por Rol

### 2.1. Navegación Pública (Sin Autenticación)

**Objetivo:** Verificar accesibilidad de páginas públicas y flujos de entrada.

#### Módulo: Legal y Estático
| ID | Caso de Prueba | Rutas | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-PUB-01 | Ver Términos y Condiciones | `/legal/terms` | Texto legal visible, título "Términos y Condiciones" | Alta |
| TC-PUB-02 | Ver Política de Privacidad | `/legal/privacy` | Texto de privacidad visible, secciones GDPR | Alta |
| TC-PUB-03 | Landing Default (Root) | `/` | Redirección a login o landing genérica (si existe) | Baja |

#### Módulo: Registro de Jugador
| ID | Caso de Prueba | Rutas | Flujo Detallado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-PUB-04 | Formulario de Registro | `/register-player` | 1. Acceder a url<br>2. Verificar campos: Nombre, Email, Password, Fecha Nacimiento | Alta |
| TC-PUB-05 | Registro Exitoso | `/register-player` | 1. Llenar datos válidos<br>2. Submit<br>3. Verificar toast éxito<br>4. Redirección a `/login` | Crítica |
| TC-PUB-06 | Validación Email Duplicado | `/register-player` | 1. Usar email existente (`member@clubpulse.com`)<br>2. Submit<br>3. Verificar error "Email en uso" | Alta |
| TC-PUB-07 | Registro Menor de Edad | `/register-player` | 1. Seleccionar fecha nacimiento < 18 años<br>2. Verificar aparición campo "Tutor Legal" | Media |

#### Módulo: Landing Pública por Club `[clubSlug]`
| ID | Caso de Prueba | Rutas | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-PUB-08 | Ver landing de club | `/{clubSlug}` | Hero section, logo del club, enlaces a store/campeonatos | Media |
| TC-PUB-09 | Navegación Tienda Pública | `/{clubSlug}/store` | Catálogo de productos público (si aplica) | Baja |
| TC-PUB-10 | Navegación Torneos | `/{clubSlug}/championships` | Lista pública de campeonatos | Media |
| TC-PUB-11 | Ver Detalle Torneo | `/{clubSlug}/championships/[id]` | Fixture y tabla de posiciones pública | Media |

---

### 2.2. Rol: MEMBER (Socio/Jugador)

**Objetivo:** Verificar flujos transaccionales del socio (reservas, pagos, perfil).

#### Módulo: Autenticación
| ID | Caso de Prueba | Rutas | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-01 | Login Exitoso | `/login` | Credenciales válidas -> Redirección a `/` (Dashboard) | Crítica |
| TC-MEM-02 | Login Fallido | `/login` | Credenciales erróneas -> Mensaje error visible | Alta |
| TC-MEM-03 | Google Auth | `/login` -> `/google/callback` | Flujo OAuth completo -> Dashboard | Alta |

#### Módulo: Dashboard y Perfil
| ID | Caso de Prueba | Rutas | Acción | Resultado Esperado |
|:---|:---|:---|:---|:---|
| TC-MEM-04 | Ver Dashboard | `/` | Carga inicial | Widgets: Próxima reserva, Balance, Nivel (Gamification) |
| TC-MEM-05 | Ver Perfil | `/profile` | Click Avatar | Formulario edición datos, cambio password |
| TC-MEM-06 | Ver Membresía | `/membership` | Click Nav | Estado suscripción, plan actual, historial pagos |

#### Módulo: Reservas (Bookings) - **CRÍTICO**
| ID | Caso de Prueba | Rutas | Flujo Detallado |
|:---|:---|:---|:---|
| TC-MEM-07 | Iniciar Reserva | `/bookings/new` | 1. Seleccionar Deporte/Cancha<br>2. Ver disponibilidad (grid/calendar) |
| TC-MEM-08 | Confirmar Reserva | `/bookings/new` | 1. Click slot disponible<br>2. Confirmar modal<br>3. Redirección a `/bookings` |
| TC-MEM-09 | Ver Mis Reservas | `/bookings` | Lista de reservas futuras y pasadas |
| TC-MEM-10 | Cancelar Reserva | `/bookings` | 1. Click en reserva futura<br>2. Cancelar<br>3. Verificar estado "Cancelled" |

#### Módulo: Tienda (Store)
| ID | Caso de Prueba | Rutas | Acción |
|:---|:---|:---|:---|
| TC-MEM-11 | Agregar al Carrito | `/store` | Click "Agregar" en producto -> Contador carrito ++ |
| TC-MEM-12 | Checkout Real | `/store` | Abrir carrito -> "Checkout" -> Integración MercadoPago simulada |

---

### 2.3. Rol: COACH (Entrenador)

**Objetivo:** Verificar herramientas de gestión de equipo.

#### Módulo: Dashboard & Asistencia
| ID | Caso de Prueba | Rutas | Flujo Detallado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-COA-01 | Dashboard Coach | `/coach` | Ver lista resumen de equipos asignados | Alta |
| TC-COA-02 | Ver Asistencia | `/coach/attendance` | Ver grid de jugadores vs fechas | Crítica |
| TC-COA-03 | **Tomar Asistencia** | `/coach/attendance` | 1. Click en celda (Jugador/Fecha)<br>2. Toggle Presente/Ausente/Tarde<br>3. Verificar guardado automático (Toast) | Crítica |
| TC-COA-04 | Nota de Asistencia | `/coach/attendance` | 1. Click derecho/largo en celda<br>2. Agregar observación ("Lesionado") | Media |

---

### 2.4. Rol: ADMIN (Administrador de Club)

**Objetivo:** Gestión integral del club (Instalaciones, Usuarios, Configuración).

#### Módulo: Gestión de Instalaciones
| ID | Caso de Prueba | Rutas | Resultado Esperado |
|:---|:---|:---|:---|
| TC-ADM-01 | Crear Instalación | `/facilities/create` | Formulario completo (Nombre, Deporte, Superficie) -> Save |
| TC-ADM-02 | Listar Instalaciones | `/facilities` | Ver la instalación creada en la lista |

#### Módulo: Calendario Maestro & Recurrencia
| ID | Caso de Prueba | Rutas | Flujo Detallado |
|:---|:---|:---|:---|
| TC-ADM-03 | Calendario Global | `/bookings/calendar` | Vista mensual/semanal de TODAS las canchas |
| TC-ADM-04 | **Bloqueo Administrativo** | `/bookings/calendar` | 1. Click slot vacío<br>2. Tipo "Bloqueo/Mantenimiento"<br>3. Guardar |
| TC-ADM-05 | Configurar Recurrencia | `/admin/recurring-rules` | Crear regla: "Cancha 1 - Lun/Mie/Vie - 19:00 - Todo el año" |
| TC-ADM-06 | Ver Reservas Recurrentes | `/admin/recurring-bookings` | Verificar que la regla anterior generó las reservas individuales |

#### Módulo: Usuarios y Accesos
| ID | Caso de Prueba | Rutas | Acción |
|:---|:---|:---|:---|
| TC-ADM-07 | Gestión Familia | `/users/family` | Vincular usuario A (Padre) con usuario B (Hijo) |
| TC-ADM-08 | Control de Acceso | `/access-control` | Simular input de QR/DNI -> Ver "Acceso Permitido/Denegado" |
| TC-ADM-09 | Ajustes Club | `/settings` | Cambiar nombre club, logo, colores |

---

### 2.5. Rol: SUPER_ADMIN (Plataforma)

#### Módulo: Multi-Tenant
| ID | Caso de Prueba | Rutas | Resultado Esperado |
|:---|:---|:---|:---|
| TC-SUP-01 | Dashboard Plataforma | `/admin/platform` | KPIs globales (Total MRR, Total Clubes) |
| TC-SUP-02 | Crear Nuevo Club | `/admin/clubs` | Wizard alta de tenant -> Provisioning DB/Storage |
| TC-SUP-03 | Config Global | `/admin/settings` | Feature flags globales, tiers de precios |

---

## 3. Pruebas de API (Módulos sin UI / Headless)

**Nota:** Estos módulos se detectaron en Backend (`migrations`) pero no tienen UI en Frontend. Se deben probar vía Postman/Curl o tests de integración.

### 3.1. Módulo: Ads & Sponsors (`/api/v1/club/ads`)
| ID | Caso de Prueba | Endpoint | Payload / Acción |
|:---|:---|:---|:---|
| API-ADS-01 | Crear Sponsor | `POST /sponsors` | `{ "name": "Nike", "logo_url": "..." }` |
| API-ADS-02 | Crear Placement | `POST /placements` | `{ "sponsor_id": "...", "location": "SIDEBAR", "active": true }` |
| API-ADS-03 | Consultar Ads (Member) | `GET /club/ads` | Verificar que retorna ads activos para el club del usuario |

### 3.2. Módulo: Gamification Phase 2 (`/api/v1/gamification`)
| ID | Caso de Prueba | Endpoint | Payload / Acción |
|:---|:---|:---|:---|
| API-GAM-01 | Listar Badges | `GET /badges` | Retorna catálogo de insignias (DB `badges`) |
| API-GAM-02 | Asignar Badge (Admin) | `POST /users/{id}/badges` | Asignar badge manual a usuario |
| API-GAM-03 | Consultar Misiones | `GET /missions` | Retorna misiones diarias/semanales activas |
| API-GAM-04 | Claim Misión | `POST /missions/{id}/claim` | Marcar misión como reclamada -> Sumar XP |
| API-GAM-05 | Ver Leaderboard Global | `GET /gamification/leaderboard` | Retorna ranking global paginado (Daily/Weekly/Monthly) |
| API-GAM-06 | Ver Contexto Leaderboard | `GET /gamification/leaderboard/context` | Retorna posición del usuario +/- vecinos |

### 3.3. Módulo: Wallet (`/api/v1/users/{id}/wallet`)
| ID | Caso de Prueba | Endpoint | Payload / Acción |
|:---|:---|:---|:---|
| API-WAL-01 | Consultar Saldo Propio | `GET /users/me/wallet` | Retorna objeto { balance, points } |
| API-WAL-02 | Consultar Saldo Usuario (Admin) | `GET /users/{id}/wallet` | Admin consulta billetera de un tercero |
| API-WAL-03 | Transacción (Simulada) | `POST /users/me/wallet/transaction` | Simular débito/crédito (si existe endpoint dev) |

### 3.4. Módulo: Incidentes (`/api/v1/incidents`)
| ID | Caso de Prueba | Endpoint | Payload / Acción |
|:---|:---|:---|:---|
| API-INC-01 | Reportar Incidente | `POST /users/incidents` | Payload: `{ injured_user_id, description, witnesses, action_taken }` |
| API-INC-02 | Listar Incidentes (Admin) | `GET /admin/incidents` | Ver log de incidentes reportados (si existe endpoint de listado) |

---

## 4. Estrategia de Ejecución y Pre-condiciones

### 4.1. Setup de Datos (Seeder)
Antes de ejecutar, correr el seeder para asegurar usuarios y roles:
```bash
cd backend && go run cmd/seeder/main.go
```
Esto crea:
- `admin@clubpulse.com` / `Admin123!`
- `coach@clubpulse.com`
- `member@clubpulse.com`

### 4.2. Ejecución Automatizada
El proyecto cuenta con Playwright para E2E.
```bash
cd frontend
npx playwright test
```
**Cobertura Actual:** `auth.spec.ts`, `booking-flow.spec.ts`
**Cobertura Faltante:** `coach.spec.ts`, `admin-settings.spec.ts`

### 4.3. Verificación Manual
Para los módulos **sin UI** (Ads, Gamification Ph2), usar carpeta de colección Postman `docs/postman/ClubPulse.postman_collection.json` (si existe) o `curl`.
