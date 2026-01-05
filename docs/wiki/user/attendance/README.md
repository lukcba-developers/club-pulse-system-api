# Módulo de Asistencia (Attendance)

Este módulo está diseñado para que los entrenadores y administradores puedan llevar un registro de la asistencia de los socios a las clases y [grupos de entrenamiento](../disciplines/README.md).

## Flujo de Trabajo para Tomar Asistencia

El proceso para registrar la asistencia se realiza en dos pasos principales:

### 1. Obtener la Lista de Asistencia

Antes de poder marcar quién vino y quién no, el entrenador debe obtener la lista de los alumnos inscritos en su grupo para una fecha específica.

-   **Flujo**:
    1.  El entrenador selecciona su [grupo de entrenamiento](../disciplines/README.md) (ej: "Tenis Infantil A") y la fecha de la clase.
    2.  El sistema genera una lista con los nombres de todos los socios que deberían asistir a esa clase. Si es la primera vez que se toma asistencia para ese día, el sistema crea la lista en el momento.
-   **Endpoint relacionado**: `GET /attendance/training-groups/:id`

### 2. Registrar la Asistencia

Una vez que el entrenador tiene la lista, puede marcar el estado de cada socio.

-   **Flujo**:
    1.  Para cada socio en la lista, el entrenador selecciona uno de los siguientes estados:
        -   `PRESENT` (Presente)
        -   `ABSENT` (Ausente)
        -   `LATE` (Tarde)
        -   `EXCUSED` (Justificado)
    2.  Una vez completado, envía el formulario. El sistema guarda el registro de asistencia para esa fecha.
-   **Endpoint relacionado**: `POST /attendance/:listID/records`

---

Este módulo es fundamental para el seguimiento del progreso de los socios y para la gestión de los grupos por parte de los entrenadores.
