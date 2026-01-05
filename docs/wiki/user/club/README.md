# Módulo de Club (Gestión Multi-Tenant)

**Nota Importante:** Este módulo es de carácter puramente administrativo y está restringido para uso exclusivo del rol `SUPER_ADMIN`. No gestiona los detalles del día a día de un club (como noticias o reglas), sino que se encarga de la gestión de los diferentes "tenants" o instancias de clubes dentro de la plataforma.

En una arquitectura multi-tenant, cada club opera en su propio entorno aislado para garantizar la privacidad y seguridad de sus datos. Este módulo proporciona las herramientas para administrar dichos entornos.

## Casos de Uso (Exclusivo para Super Admin)

### 1. Crear un Nuevo Club (Tenant)

Un `SUPER_ADMIN` puede dar de alta un nuevo club en la plataforma, creando un nuevo entorno aislado para él.

-   **Flujo**:
    1.  El `SUPER_ADMIN` completa un formulario con los datos del nuevo club (nombre, dominio, etc.).
    2.  El sistema crea una nueva entrada para el club en la base de datos, asignándole un identificador único y preparando su entorno.
-   **Endpoint relacionado**: `POST /clubs`

### 2. Listar Todos los Clubs

El `SUPER_ADMIN` puede obtener una lista de todos los clubes (tenants) que existen en la plataforma.

-   **Flujo**:
    1.  El `SUPER_ADMIN` accede a un panel de control global.
    2.  El sistema muestra una tabla con todos los clubes, su estado (`ACTIVE`, `INACTIVE`), y otros detalles.
-   **Endpoint relacionado**: `GET /clubs`

### 3. Actualizar un Club

El `SUPER_ADMIN` puede modificar la configuración de un club existente.

-   **Flujo**:
    1.  El `SUPER_ADMIN` selecciona un club de la lista y edita sus propiedades.
    2.  El sistema actualiza la información del club correspondiente.
-   **Endpoint relacionado**: `PUT /clubs/:id`
