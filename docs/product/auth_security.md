# Módulo de Identidad y Seguridad (Auth & User)

Este módulo es el responsable de la seguridad de la plataforma, la gestión de la identidad de los usuarios y el control de acceso.

## 🌟 Funcionalidades Implementadas

### 1. Autenticación
-   **Registro y Login con Credenciales**: Soporte para registro y autenticación mediante email y contraseña. Las contraseñas se almacenan de forma segura utilizando el algoritmo **Bcrypt**.
-   **Gestión de Sesiones Segura**: El sistema utiliza **tokens JWT** para gestionar las sesiones de usuario, que se transmiten de forma segura.

### 2. Control de Acceso (RBAC)
-   **Sistema de Roles**: Se ha implementado un sistema de Control de Acceso Basado en Roles. Los roles definidos en el código incluyen `MEMBER`, `ADMIN`, y `SUPER_ADMIN`.
-   **Autorización por Middleware**: Un middleware en el backend se encarga de verificar el rol del usuario en cada petición a un endpoint protegido, garantizando que solo los usuarios con los permisos adecuados puedan acceder a los recursos.

### 3. Perfil de Usuario
-   **Gestión de Datos Personales**: Los usuarios pueden gestionar su información de perfil básica.
-   **Aislamiento de Datos (Tenant Isolation)**: La arquitectura asegura que un usuario de un club no pueda acceder a la información de otro. Cada consulta a la base de datos está estrictamente segmentada por `club_id`.

## 4. Funcionalidades en Desarrollo
-   **Login Social (OAuth)**: La capacidad de iniciar sesión con proveedores como Google está planificada pero aún no implementada.
-   **Grupos Familiares**: La funcionalidad para que un usuario gestione las cuentas de sus dependientes está en el roadmap.
