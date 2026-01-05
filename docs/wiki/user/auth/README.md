# Módulo de Autenticación (Auth)

El módulo de Autenticación es el responsable de verificar la identidad de los usuarios y controlar su acceso a la plataforma. Garantiza que solo las personas autorizadas puedan ingresar al sistema y que sus datos estén seguros.

## Casos de Uso Principales

A continuación se describen las operaciones fundamentales que un usuario puede realizar a través de este módulo.

### 1. Registro de Usuario

Los nuevos usuarios pueden crear una cuenta proporcionando su correo electrónico y una contraseña.

-   **Flujo**:
    1.  El usuario completa el formulario de registro.
    2.  El sistema crea una nueva cuenta de usuario.
    3.  El usuario es automáticamente autenticado y redirigido al panel principal.
-   **Endpoint relacionado**: `POST /auth/register`

### 2. Inicio de Sesión (Login)

Los usuarios registrados pueden acceder a su cuenta.

-   **Flujo con Email/Contraseña**:
    1.  El usuario introduce su correo electrónico y contraseña.
    2.  El sistema valida las credenciales.
    3.  Si son correctas, el sistema establece cookies de sesión seguras (`HttpOnly`) en el navegador y redirige al usuario al panel principal.
-   **Endpoint relacionado**: `POST /auth/login`

-   **Flujo con Google (OAuth 2.0)**:
    1.  El usuario hace clic en "Iniciar sesión con Google".
    2.  Es redirigido a la página de autenticación de Google.
    3.  Después de autorizar, Google redirige de nuevo a la aplicación.
    4.  El sistema finaliza el proceso de login y establece las cookies de sesión.
-   **Endpoint relacionado**: `POST /auth/google`

### 3. Cierre de Sesión (Logout)

Un usuario puede terminar su sesión activa de forma segura.

-   **Flujo**:
    1.  El usuario hace clic en el botón de "Cerrar Sesión".
    2.  El sistema invalida el token de sesión actual.
    3.  Las cookies de autenticación son eliminadas del navegador.
    4.  El usuario es redirigido a la página de inicio de sesión.
-   **Endpoint relacionado**: `POST /auth/logout`

### 4. Gestión de Sesiones Activas

Para mayor seguridad, los usuarios pueden ver y gestionar todos los dispositivos donde han iniciado sesión.

-   **Listar Sesiones**: Un usuario puede ver una lista de todas sus sesiones activas (ej: "Laptop Chrome", "Móvil Safari").
    -   **Endpoint relacionado**: `GET /auth/sessions`
-   **Revocar una Sesión**: Un usuario puede cerrar la sesión de forma remota en un dispositivo específico si, por ejemplo, ha perdido el dispositivo o sospecha de un acceso no autorizado.
    -   **Endpoint relacionado**: `DELETE /auth/sessions/:id`

## Seguridad y Roles

### Cookies de Sesión Seguras
La sesión del usuario no se almacena en `localStorage` (vulnerable a ataques XSS). En su lugar, se utilizan cookies `HttpOnly`, `Secure` y con `SameSite=Strict`, que no son accesibles desde el código JavaScript del navegador, ofreciendo un nivel de seguridad muy superior.

### Control de Acceso Basado en Roles (RBAC)
Una vez autenticado, el sistema asigna un rol al usuario. Este rol determina a qué funcionalidades y datos puede acceder. El sistema verifica este rol en cada solicitud a una ruta protegida para asegurar que el usuario tiene los permisos necesarios.

A continuación se detallan las responsabilidades de cada rol:

#### `MEMBER` (Socio)
Es el rol estándar para un usuario del club. Sus permisos se centran en el uso y disfrute de las instalaciones y servicios.
-   Gestionar su propio perfil de usuario.
-   Ver y gestionar su [grupo familiar](../user/README.md).
-   Consultar su [membresía](../membership/README.md) y estado de cuenta.
-   Ver y [reservar instalaciones](../booking/README.md).
-   Inscribirse en [torneos y disciplinas](../disciplines/README.md).
-   Consultar su historial de [asistencia](../attendance/README.md).

#### `ADMIN` (Administrador de Club)
Este rol tiene permisos para gestionar todas las operaciones de **un club específico (tenant)**.
-   Tiene todos los permisos de un `MEMBER`.
-   **Gestión de Usuarios**: Listar, buscar y eliminar usuarios de su club.
-   **Gestión de Instalaciones**: Crear, modificar y eliminar instalaciones.
-   **Gestión de Reservas**: Ver y cancelar cualquier reserva dentro de su club. Gestionar reservas recurrentes.
-   **Gestión de Disciplinas y Torneos**: Crear y administrar torneos, programar partidos y registrar resultados.
-   **Gestión de Asistencia**: Supervisar los registros de asistencia.

#### `SUPER_ADMIN` (Super Administrador)
Es el rol más alto del sistema, con control total sobre toda la plataforma multi-tenant.
-   Tiene todos los permisos de un `ADMIN`.
-   **Gestión Multi-Tenant**: Puede operar a través de todos los clubes (tenants) del sistema.
-   **Gestión de Clubs**: Es el único rol que puede crear, listar y modificar las instancias de los [clubes (tenants)](../club/README.md) en la plataforma.
