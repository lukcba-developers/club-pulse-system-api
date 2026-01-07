# Módulo de Instalaciones (Facilities)

Este módulo gestiona toda la información y las operaciones relacionadas con las instalaciones físicas del club, como canchas, piscinas, salones, etc. Es la base sobre la cual operan otros módulos, especialmente el de Reservas (Booking).

## Casos de Uso para Socios (Members)

### 1. Listar y Ver Instalaciones

Los socios pueden explorar todas las instalaciones que el club ofrece.

-   **Flujo**:
    1.  El usuario accede a la sección "Instalaciones" del sistema.
    2.  Se muestra una lista completa de las instalaciones disponibles, con su nombre y tipo.
    3.  El usuario puede hacer clic en una instalación para ver sus detalles completos: capacidad, precio por hora, especificaciones (ej: tipo de superficie, si tiene iluminación) y su ubicación dentro del club.
-   **Endpoints relacionados**:
    -   `GET /facilities` (para la lista)
    -   `GET /facilities/:id` (para los detalles)

### 2. Búsqueda Avanzada de Instalaciones (Búsqueda Semántica)

En lugar de usar filtros rígidos, el sistema permite a los usuarios buscar instalaciones usando lenguaje natural.

-   **Flujo**:
    1.  El usuario escribe en una barra de búsqueda lo que necesita, por ejemplo: "cancha de tenis techada con luz" o "piscina para niños".
    2.  El sistema interpreta la petición y devuelve una lista de las instalaciones que mejor coinciden con la descripción, ordenadas por relevancia.
-   **Endpoint relacionado**: `GET /facilities/search`

## Casos de Uso para Administradores (Admins)

### 1. Gestión de Instalaciones (CRUD)

Los administradores tienen control total sobre el catálogo de instalaciones del club.

-   **Crear**: Añadir una nueva instalación al sistema, definiendo todas sus propiedades.
    -   `POST /facilities`
-   **Modificar**: Actualizar los detalles de una instalación existente.
    -   `PUT /facilities/:id`
-   **Eliminar**: Quitar una instalación del sistema.
    -   `DELETE /facilities/:id` (No implementado en el handler actual, pero es parte de un CRUD estándar).

### 2. Gestión de Equipamiento y Préstamos

Los administradores pueden gestionar el inventario de equipamiento asociado a las instalaciones y sus préstamos a los socios.

-   **Añadir Equipamiento a una Instalación**: Permite registrar nuevas piezas de equipamiento (ej: raquetas, balones) y asociarlas a una instalación.
    -   `POST /facilities/:id/equipment`
-   **Listar Equipamiento de una Instalación**: Muestra todo el equipamiento asociado a una instalación específica.
    -   `GET /facilities/:id/equipment`
-   **Prestar Equipamiento a un Socio**: Registra un nuevo préstamo, vinculando una pieza de equipamiento a un socio y estableciendo una fecha de devolución.
    -   `POST /equipment/:id/loan`
-   **Registrar Devolución de un Préstamo**: Marca un préstamo como "devuelto" y permite registrar la condición del equipamiento al momento de la devolución.
    -   `POST /loans/:id/return`

### 3. Preparar Datos para Búsqueda Semántica

Para que la búsqueda en lenguaje natural funcione, un administrador debe ejecutar un proceso que "lee" y "entiende" las descripciones de todas las instalaciones, convirtiéndolas en vectores de búsqueda (embeddings).

-   **Flujo**:
    1.  Después de añadir o modificar varias instalaciones, el administrador ejecuta esta acción.
    2.  El sistema procesa los datos para actualizar el índice de búsqueda.
-   **Endpoint relacionado**: `POST /facilities/embeddings/generate`
-   **Nota**: Esta es una operación de mantenimiento técnico.
