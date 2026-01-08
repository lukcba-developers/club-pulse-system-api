# Manual de Usuario: Módulo de Instalaciones (Facilities)

## 1. Propósito

Este módulo contiene el catálogo de todos los espacios físicos que el club pone a tu disposición, como las canchas de tenis, pádel, piscinas, etc. Como socio, puedes explorar estas instalaciones, y como administrador, puedes gestionarlas.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede ver y buscar instalaciones.
-   **Administrador (`ADMIN`):** Puede crear, editar y gestionar el estado de todas las instalaciones.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Explorar las Instalaciones

**Paso a paso:**
1.  **Navega a la sección "Instalaciones"** en la aplicación.
2.  Verás una lista o una galería con todas las instalaciones disponibles en el club.
3.  **Haz clic en una instalación** para ver sus detalles, como:
    -   Fotos y descripción.
    -   Horarios de apertura.
    -   Reglas específicas de uso.
4.  Desde la vista de detalle, normalmente encontrarás un botón para **"Ver Disponibilidad"** o **"Reservar"**, que te llevará directamente al calendario del módulo de Reservas para esa instalación.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Crear o Editar una Instalación

**Paso a paso:**
1.  **Accede al Panel de Administración** y ve a la sección de **"Instalaciones"**.
2.  Para crear una nueva, haz clic en **"Nueva Instalación"**. Para editar, búscala en la lista y haz clic en "Editar".
3.  **Completa el formulario** con toda la información relevante:
    -   **Nombre:** (ej: "Cancha Central de Tenis").
    -   **Tipo:** (ej: "Cancha", "Piscina").
    -   **Configuración de Horarios:** Define las horas de apertura y cierre.
    -   **Duración de Slots:** Establece la duración estándar de una reserva (ej: 60 minutos).
4.  **Guarda los cambios.** La instalación aparecerá inmediatamente en la lista pública.

### 🔸 Cómo Poner una Instalación en Mantenimiento

Si una instalación necesita reparaciones o no está disponible temporalmente, puedes bloquearla.

**Paso a paso:**
1.  En el listado de instalaciones del Panel de Administración, busca la que deseas bloquear.
2.  **Cambia su estado** de `Disponible` a `En Mantenimiento`.
3.  **Define un rango de fechas** para el bloqueo si es necesario.
4.  Durante este período, la instalación **no aparecerá como disponible** en el calendario de reservas para los socios.

---

## 5. Diagrama de Flujo: Gestión de Instalaciones (Admin)

```mermaid
graph TD
    A[Inicio: Panel de Admin] --> B[Ir a "Gestión de Instalaciones"];
    B --> C{¿Qué deseas hacer?};
    C -- Crear Nueva --> D[Rellenar Formulario de Nueva Instalación];
    D --> E[Guardar];
    C -- Editar Existente --> F[Seleccionar Instalación de la Lista];
    F --> G[Modificar Datos];
    G --> E;
    C -- Bloquear por Mantenimiento --> H[Seleccionar Instalación];
    H --> I[Cambiar Estado a "En Mantenimiento"];
    I --> E;
    E --> J[Instalación Actualizada ✅];
    J --> B;
```
