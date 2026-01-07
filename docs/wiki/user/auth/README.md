# Módulo de Autenticación (Auth)

## 1. Propósito

El módulo de **Autenticación** es la piedra angular de la seguridad del sistema. Se encarga de verificar la identidad de los usuarios y de gestionar los permisos que determinan qué acciones pueden realizar dentro de la plataforma.

## 2. Funcionalidades Principales

-   **Inicio de Sesión (Login):**
    -   **Credenciales:** Permite a los usuarios iniciar sesión con su correo electrónico y una contraseña.
    -   **OAuth 2.0 (Próximamente):** Permitirá el inicio de sesión a través de proveedores externos como Google, Facebook, etc.
-   **Gestión de Contraseñas:** Incluye funcionalidades para restablecer contraseñas olvidadas de forma segura.
-   **Control de Acceso Basado en Roles (RBAC):** El sistema define varios roles, cada uno con un conjunto específico de permisos. Los roles principales son:
    -   `MEMBER`: Rol estándar para todos los socios del club. Tienen acceso a sus perfiles, reservas, etc.
    -   `STAFF`: Personal del club (entrenadores, recepción) con permisos para gestionar asistencias, reservas de otros, etc.
    -   `ADMIN`: Administradores del club con acceso a la configuración de módulos, gestión de usuarios y finanzas.
    -   `SUPER_ADMIN`: Rol de más alto nivel para la administración total del sistema.
-   **Gestión de Sesiones:** Maneja de forma segura los tokens de sesión (JWT) que mantienen al usuario conectado y autorizan sus peticiones a la API.