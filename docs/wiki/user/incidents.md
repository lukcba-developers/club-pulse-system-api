# Módulo de Registro de Incidentes

## 1. Propósito

El módulo de Registro de Incidentes proporciona un mecanismo formal para que los socios y el personal del club reporten cualquier evento adverso, como lesiones, accidentes o disputas, que ocurra dentro de las instalaciones. Esto es crucial para la seguridad, el cumplimiento normativo y la gestión de riesgos.

## 2. Funcionalidades

-   **Reporte Formal:** Permite a cualquier usuario autenticado registrar un incidente de manera estructurada.
-   **Trazabilidad:** Almacena un registro de qué sucedió, quién estuvo involucrado, qué acciones se tomaron y quién lo reportó.
-   **Análisis de Seguridad:** Proporciona datos valiosos para que la administración del club identifique áreas de riesgo y tome medidas preventivas.

## 3. Modelo de Datos (`IncidentLog`)

| Campo             | Tipo      | Descripción                                                              |
| ----------------- | --------- | ------------------------------------------------------------------------ |
| `ID`              | `UUID`    | Identificador único del registro del incidente.                          |
| `ClubID`          | `string`  | ID del club donde ocurrió el incidente.                                  |
| `InjuredUserID`   | `string`  | ID del usuario que resultó herido (si aplica).                           |
| `Description`     | `string`  | Descripción detallada de lo que sucedió.                                 |
| `Witnesses`       | `string`  | Nombres o descripción de los testigos del incidente.                     |
| `ActionTaken`     | `string`  | Descripción de las acciones inmediatas que se tomaron (ej: "Se aplicó hielo", "Se llamó a emergencias"). |
| `ReportedAt`      | `datetime`| Fecha y hora en que se reportó el incidente.                             |
| `CreatedBy`       | `string`  | ID del usuario que creó el registro del incidente.                       |

## 4. Endpoint de la API

### `POST /users/incidents`

-   **Acción:** Crea un nuevo registro de incidente.
-   **Permisos:** Cualquier usuario autenticado (`MEMBER`, `COACH`, `ADMIN`, etc.) puede reportar un incidente.
-   **Request Body (JSON):**
    ```json
    {
      "injured_user_id": "user-uuid-of-injured-person",
      "description": "El socio se resbaló en la cancha 3 debido a una mancha de agua y se torció el tobillo.",
      "witnesses": "Juan Pérez, María Gómez",
      "action_taken": "Se le proporcionó una bolsa de hielo y se le ayudó a salir de la cancha."
    }
    ```
-   **Respuesta Exitosa (201 Created):** El objeto del incidente creado.
