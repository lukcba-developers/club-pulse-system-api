# Módulo de Identidad y Seguridad (Auth & User)

El guardián de la plataforma, asegurando que cada acceso sea legítimo y cada dato permanezca privado.

## 🌟 Funcionalidades Principales

### 1. Autenticación Robusta
-   **Credenciales Clásicas**: Registro con Email y Password (almacenamiento seguro con **Bcrypt**).
-   **Login Social**: Integración OAuth 2.0 con **Google** para acceso friction-less.

### 2. Seguridad de Sesiones (Best-in-Class)
Enfoque moderno alejado del inseguro `localStorage`.
-   **HttpOnly Cookies**: Tokens inaccesibles para scripts del navegador (inmunidad XSS).
-   **Secure + SameSite**: Configuración estricta de cookies para prevenir CSRF.
-   **Redis Session Store**: Almacenamiento centralizado de sesiones activas. Permite revocación instantánea desde el backend.

### 3. Control de Acceso (RBAC)
Sistema de roles jerárquico.
-   **Roles**:
    -   `SUPER_ADMIN`: Acceso total multi-tenant.
    -   `ADMIN`: Gestión completa de SU club.
    -   `MEMBER`: Acceso a reservas y perfil propio.
-   **Middleware**: Verificación en cada request protegido.

### 4. Perfil de Usuario
Gestión centrada en el socio.
-   **Datos Personales**: Teléfono, Dirección, Preferencias.
-   **Grupo Familiar**: Capacidad de gestionar cuentas de hijos/dependientes (Roadmap).

### 5. Protección de Datos (Tenant Isolation)
-   **BOLA Protection**: Arquitectura diseñada para prevenir *Broken Object Level Authorization*. Cada consulta a la base de datos inyecta obligatoriamente el `club_id` del usuario autenticado, haciendo matemáticamente imposible acceder a datos de otro club.
