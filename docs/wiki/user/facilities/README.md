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

### 3. Consultar Disponibilidad

Antes de hacer una reserva, un socio puede verificar cuándo una instalación específica está libre.

-   **Flujo**:
    1.  El usuario selecciona una instalación y un rango de fechas.
    2.  El sistema devuelve un calendario o una lista de los horarios ocupados y disponibles para esa instalación en el período seleccionado.
-   **Endpoint relacionado**: `GET /facilities/:id/availability`
-   **Nota**: Este es un endpoint clave que es consumido por el módulo de Reservas.

## Casos de Uso para Administradores (Admins)

### 1. Gestión de Instalaciones (CRUD)

Los administradores tienen control total sobre el catálogo de instalaciones del club.

-   **Crear**: Añadir una nueva instalación al sistema, definiendo todas sus propiedades.
    -   `POST /facilities`
-   **Modificar**: Actualizar los detalles de una instalación existente.
    -   `PUT /facilities/:id`
-   **Eliminar**: Quitar una instalación del sistema.
    -   `DELETE /facilities/:id`

### 2. Preparar Datos para Búsqueda Semántica

Para que la búsqueda en lenguaje natural funcione, un administrador debe ejecutar un proceso que "lee" y "entiende" las descripciones de todas las instalaciones, convirtiéndolas en vectores de búsqueda (embeddings).

-   **Flujo**:
    1.  Después de añadir o modificar varias instalaciones, el administrador ejecuta esta acción.
    2.  El sistema procesa los datos para actualizar el índice de búsqueda.
-   **Endpoint relacionado**: `POST /facilities/embeddings/generate`
-   **Nota**: Esta es una operación de mantenimiento técnico.

---
## Funcionalidad Interna (Sin API)

### Gestión de Equipamiento y Préstamos

El sistema está diseñado para soportar la gestión de equipamiento asociado a las instalaciones y el préstamo de dicho equipamiento a los socios.

-   **Concepto de Equipamiento**: Se puede definir equipamiento (ej: "Raqueta de Tenis", "Balón de Fútbol") y asociarlo a una instalación específica. Cada pieza de equipamiento tiene un estado (`disponible`, `en_uso`) y una condición (`bueno`, `dañado`).

-   **Concepto de Préstamo (`Loan`)**: El sistema puede registrar un préstamo, que vincula una pieza de equipamiento con un socio. El préstamo registra la fecha del préstamo, la fecha de devolución esperada y el estado (`ACTIVO`, `DEVUELTO`, `VENCIDO`). El sistema incluye lógica para detectar automáticamente si un préstamo está vencido.

-   **Gestión**: Al igual que con las Becas, actualmente **no existen endpoints en la API** para que los usuarios o administradores gestionen el equipamiento o los préstamos. La creación de equipamiento, el registro de un préstamo y su devolución son capacidades internas del sistema que requerirían manipulación directa de la base de datos o la implementación de nuevos endpoints para ser funcionales.
