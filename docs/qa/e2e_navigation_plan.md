# Plan de Pruebas de Navegabilidad E2E (End-to-End)

Este documento define la suite de pruebas para verificar la navegabilidad y funcionalidad completa del proyecto `club-pulse-system-api` utilizando un agente de navegación (simulado o automatizado).

Las pruebas están organizadas por **Roles de Usuario** y **Módulos Funcionales**, cubriendo los casos de uso documentados.

## 1. Matriz de Roles y Permisos

| Rol | Descripción | Alcance |
| :--- | :--- | :--- |
| **`MEMBER`** (Socio) | Usuario final del club. | Acceso a gestión personal, reservas, visualización de instalaciones, inscripción a equipos, tienda. |
| **`COACH`** (Entrenador) | Profesor o instructor. | Gestión de asistencia, visualización de grupos asignados. |
| **`ADMIN`** (Administrador) | Gestor de un club específico. | Gestión de instalaciones, usuarios del club, reportes, configuración del club. |
| **`SUPER_ADMIN`** | Administrador del sistema. | Gestión multi-tenant, creación de clubes, acceso global. |

---

## 2. Suite de Pruebas por Rol

### 2.1. Rol: MEMBER (Socio)

**Objetivo:** Verificar que un socio puede realizar sus actividades diarias sin errores.

#### Módulo: Autenticación (Auth)
- [ ] **Caso TC-MEM-01: Login Exitoso**
    - **Acción:** Navegar a `/login` -> Ingresar credenciales válidas -> Submit.
    - **Esperado:** Redirección al Dashboard (`/`). Bienvenida visible.
- [ ] **Caso TC-MEM-02: Logout**
    - **Acción:** Click en Avatar -> "Cerrar Sesión".
    - **Esperado:** Redirección a `/login`. Cookies limpias.

#### Módulo: Perfil (Profile)
- [ ] **Caso TC-MEM-03: Ver Perfil**
    - **Acción:** Navegar a `/profile`.
    - **Esperado:** Datos personales visibles (Nombre, Email, Categoría).
- [ ] **Caso TC-MEM-04: Ver Grupo Familiar**
    - **Acción:** En Dashboard o Perfil, buscar sección "Familia".
    - **Esperado:** Lista de miembros familiares visible (si aplica).

#### Módulo: Membresía (Membership) y Pagos
- [ ] **Caso TC-MEM-05: Estado de Membresía**
    - **Acción:** Navegar a `/membership`.
    - **Esperado:** Tarjeta con estado "Activo" o "Vencido", próxima fecha de cobro.
- [ ] **Caso TC-MEM-06: Ver Planes**
    - **Acción:** Navegar a `/membership/plans`.
    - **Esperado:** Lista de Tiers (Gold, Platinum, etc.) con precios.
- [ ] **Caso TC-MEM-07: Iniciar Pago (Simulado)**
    - **Acción:** Seleccionar un plan o deuda -> Click "Pagar".
    - **Esperado:** Redirección a URL externa (Pasarela) o simulación de éxito si es entorno de test (`POST /payments/checkout` exitoso).

#### Módulo: Disciplinas y Grupos (Disciplines) - **[NUEVO]**
- [ ] **Caso TC-MEM-08: Explorar Disciplinas**
    - **Acción:** Navegar a `/disciplines`.
    - **Esperado:** Lista de deportes con imágenes (Tenis, Fútbol, etc.).
- [ ] **Caso TC-MEM-09: Ver Grupos de Entrenamiento**
    - **Acción:** Click en una disciplina -> Ver grupos disponibles.
    - **Esperado:** Lista de grupos filtrables por categoría/edad.

#### Módulo: Instalaciones y Reservas (Facilities & Booking)
- [ ] **Caso TC-MEM-10: Listar Instalaciones**
    - **Acción:** Navegar a `/facilities`.
    - **Esperado:** Lista de canchas/espacios. Buscador funcional.
- [ ] **Caso TC-MEM-11: Crear Reserva Exitosa**
    - **Acción:** Seleccionar slot libre en calendario -> Confirmar reserva.
    - **Esperado:** Redirección a pago o confirmación directa. Notificación de éxito.
- [ ] **Caso TC-MEM-12: Cancelar Reserva**
    - **Acción:** Ir a `/bookings` (Mis Reservas) -> Click "Cancelar".
    - **Esperado:** Estado cambia a "Cancelada". Slot liberado.

#### Módulo: Equipos y Campeonatos (Teams) - **[NUEVO]**
- [ ] **Caso TC-MEM-13: Inscribir Equipo**
    - **Acción:** Ir a `/championships` -> Seleccionar torneo abierto -> "Inscribir Equipo".
    - **Esperado:** Formulario para nombre de equipo y selección de miembros. Confirmación exitosa.

#### Módulo: Tienda (Store)
- [ ] **Caso TC-MEM-14: Ver Catálogo**
    - **Acción:** Navegar a `/store`.
    - **Esperado:** Productos listados con precio y stock.

---

### 2.2. Rol: COACH (Entrenador) - **[NUEVO]**

**Objetivo:** Verificar la gestión de clases y asistencia.

#### Módulo: Mis Grupos y Asistencia
- [ ] **Caso TC-COA-01: Ver Mis Grupos**
    - **Acción:** Login como Coach -> Navegar a `/coach/groups`.
    - **Esperado:** Lista de grupos asignados al profesor.
- [ ] **Caso TC-COA-02: Tomar Asistencia**
    - **Acción:** Entrar a un grupo -> Seleccionar fecha -> Cargar lista (`GET /attendance...`).
    - **Esperado:** Lista de alumnos con selectores de estado (Presente/Ausente).
- [ ] **Caso TC-COA-03: Guardar Asistencia**
    - **Acción:** Marcar estados -> Click "Guardar".
    - **Esperado:** Confirmación de guardado (`POST /attendance/...`).

---

### 2.3. Rol: ADMIN (Administrador de Club)

**Objetivo:** Verificar la capacidad de gestión y control del club.

#### Módulo: Dashboard Administrativo
- [ ] **Caso TC-ADM-01: Vista General**
    - **Acción:** Login como Admin -> Home.
    - **Esperado:** Métricas visibles. Menú lateral completo.

#### Módulo: Gestión de Instalaciones
- [ ] **Caso TC-ADM-02: ABM Instalaciones**
    - **Acción:** `/admin/facilities` -> Crear/Editar.
    - **Esperado:** Formulario de alta/edición funcional.
- [ ] **Caso TC-ADM-03: Bloqueo de Horario**
    - **Acción:** Marcar slot como "Mantenimiento".
    - **Esperado:** Slot no disponible para reservas de socios.

#### Módulo: Gestión de Usuarios
- [ ] **Caso TC-ADM-04: Buscar Usuario**
    - **Acción:** `/admin/users` -> Buscar por email.
    - **Esperado:** Detalle de usuario visible.

#### Módulo: Control de Acceso
- [ ] **Caso TC-ADM-05: Logs de Acceso**
    - **Acción:** `/access-control/logs`.
    - **Esperado:** Tabla histórica de accesos.

#### Módulo: Campeonatos
- [ ] **Caso TC-ADM-06: Gestionar Torneo**
    - **Acción:** Crear fixture o validar inscripciones de equipos.
    - **Esperado:** Cambios reflejados en la vista pública.

---

### 2.4. Rol: SUPER_ADMIN

**Objetivo:** Verificar la gestión multi-tenant.

#### Módulo: Gestión de Clubes
- [ ] **Caso TC-SUP-01: ABM de Clubes**
    - **Acción:** `/admin/clubs` -> Crear nuevo club.
    - **Esperado:** Nuevo tenant creado y accesible.
- [ ] **Caso TC-SUP-02: Impersonation**
    - **Acción:** "Loguearse como" admin de otro club.
    - **Esperado:** Cambio de contexto visual.

---

## 3. Flujos de Navegación Críticos (Cross-Module)

1.  **Onboarding Completo:**
    -   Registro -> Selección de Plan (Membership) -> Checkout (Payment) -> Reserva Inicial (Booking).

2.  **Ciclo de Clase (Coach):**
    -   Coach visualiza grupo (Disciplines) -> Toma asistencia (Attendance) -> Admin verifica ocupación.

3.  **Torneo (Team & Championship):**
    -   Socio crea equipo -> Paga inscripción -> Admin aprueba -> Equipo aparece en fixture.

## 4. Estrategia de Ejecución

1.  **Manual (Exploratorio):** Un agente humano o IA recorre estos pasos verificando visualmente la UI.
2.  **Automatizado (E2E):** Implementar estos casos usando herramientas como Cypress/Playwright, referenciando los tests de integración del backend (`backend/tests/e2e/payment_flow_test.go`, etc.) para asegurar consistencia de datos.
