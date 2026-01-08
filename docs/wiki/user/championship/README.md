# Manual de Usuario: Módulo de Campeonatos

## 1. Propósito

Este módulo te permite participar y seguir los torneos y competiciones organizados por el club. Fomenta la competencia amistosa y la interacción entre socios.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede ver torneos, inscribir equipos y seguir los resultados.
-   **Administrador (`ADMIN`):** Puede crear y gestionar todos los aspectos de un torneo.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Ver los Campeonatos Disponibles

**Paso a paso:**
1.  **Navega a la sección "Campeonatos"** en la aplicación.
2.  Verás una lista de los torneos actuales y futuros.
3.  Haz clic en un torneo para ver sus detalles:
    -   **Reglamento:** Las reglas específicas del torneo.
    -   **Fechas:** Fecha de inicio, fin e inscripción.
    -   **Formato:** (ej: Liga, Eliminación Directa).
    -   **Equipos Inscritos.**

### 🔹 Cómo Inscribirte a un Campeonato

**Paso a paso:**
1.  Dentro de la página de detalles de un torneo abierto, busca el botón **"Inscribirme"** o **"Inscribir Equipo"**.
2.  Si el torneo es por equipos, se te pedirá que selecciones un equipo que hayas creado previamente en el **Módulo de Equipos** o que crees uno nuevo.
3.  Confirma la inscripción. Puede que se te redirija al **Módulo de Pagos** si la inscripción tiene un costo.
4.  Una vez inscrito, tu equipo aparecerá en la lista de participantes.

### 🔹 Cómo Seguir un Torneo

**Paso a paso:**
1.  Entra a la página de detalles del torneo que deseas seguir.
2.  Navega por las diferentes pestañas para ver:
    -   **Fixture:** El calendario de todos los partidos.
    -   **Tabla de Posiciones:** La clasificación de los equipos actualizada en tiempo real.
    -   **Resultados:** Los marcadores de los partidos que ya se han jugado.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Crear un Nuevo Campeonato

**Paso a paso:**
1.  **Accede al Panel de Administración** y ve a la sección de **"Campeonatos"**.
2.  Haz clic en **"Nuevo Campeonato"**.
3.  **Completa el formulario** con toda la información: nombre, disciplina, fechas, formato, reglas, costo de inscripción, etc.
4.  **Guarda los cambios.** El torneo se publicará y los socios podrán empezar a inscribirse.

### 🔸 Cómo Gestionar un Torneo en Curso

**Paso a paso:**
1.  Desde el panel de "Campeonatos", selecciona el torneo que deseas gestionar.
2.  Desde aquí podrás:
    -   **Aprobar o rechazar inscripciones** de equipos.
    -   **Generar el fixture** (calendario de partidos) una vez cerradas las inscripciones.
    -   **Cargar los resultados** de los partidos a medida que se van jugando. La tabla de posiciones se actualizará automáticamente.

---

## 5. Diagrama de Flujo: Inscripción a un Torneo (Socio)

```mermaid
graph TD
    A[Inicio] --> B[Explorar Campeonatos];
    B --> C{Elige un Torneo Abierto};
    C --> D[Clic en "Inscribir Equipo"];
    D --> E{¿Equipo ya creado?};
    E -- Sí --> F[Selecciona tu Equipo];
    E -- No --> G[Ir a Módulo de Equipos y Crear Equipo];
    G --> F;
    F --> H{¿Inscripción tiene costo?};
    H -- Sí --> I[Ir a Módulo de Pagos];
    I --> J[Confirmación de Inscripción ✅];
    H -- No --> J;
```

---

## 6. Gestión de Voluntarios

Para el buen desarrollo de los eventos y partidos, el sistema permite la gestión de voluntarios que pueden cubrir diferentes roles durante un encuentro.

### 🔹 Propósito

La gestión de voluntarios permite a los administradores asignar socios a roles específicos para un partido, asegurando que haya suficiente personal para tareas como seguridad, primeros auxilios, etc. Los socios también pueden ver dónde se necesita ayuda y posiblemente ofrecerse como voluntarios.

### 🔹 Modelo de Datos

La información se almacena en la tabla `volunteer_assignments`.

| Campo         | Tipo           | Descripción                                                              |
| ------------- | -------------- | ------------------------------------------------------------------------ |
| `id`          | `UUID`         | Identificador único de la asignación.                                    |
| `club_id`     | `VARCHAR`      | ID del club.                                                             |
| `match_id`    | `UUID`         | ID del partido al que se asigna el voluntario.                           |
| `user_id`     | `VARCHAR`      | ID del socio que actuará como voluntario.                                |
| `role`        | `ENUM`         | Rol que desempeñará el voluntario. Ver `VolunteerRole`.                  |
| `notes`       | `TEXT`         | Notas adicionales sobre la asignación.                                   |
| `assigned_by` | `VARCHAR`      | ID del administrador que realizó la asignación.                          |
| `assigned_at` | `TIMESTAMPTZ`  | Fecha y hora de la asignación.                                           |

#### `VolunteerRole` (Roles de Voluntario)

-   `BUFFET`: Atención en el buffet o cantina.
-   `SECURITY`: Tareas de seguridad y control de acceso.
-   `TRANSPORT`: Ayuda con el transporte de equipos o materiales.
-   `FIRST_AID`: Asistencia de primeros auxilios.
-   `COACH`: Asistente técnico o de campo.

### 🔹 Endpoints de la API

---

#### `POST /championships/matches/:id/volunteers`

-   **Acción:** Asigna un socio como voluntario a un partido específico.
-   **Permisos:** `ADMIN` o `SUPER_ADMIN`.
-   **`:id`:** Corresponde al `match_id`.
-   **Request Body (JSON):**
    ```json
    {
      "user_id": "ID_DEL_SOCIO",
      "role": "BUFFET", // Uno de los VolunteerRole
      "notes": "Encargado de la caja." // Opcional
    }
    ```
-   **Respuesta Exitosa (201 Created):** Un mensaje de confirmación.

---

#### `GET /championships/matches/:id/volunteers`

-   **Acción:** Obtiene la lista y un resumen de los voluntarios asignados a un partido.
-   **Permisos:** Abierto a usuarios autenticados.
-   **`:id`:** Corresponde al `match_id`.
-   **Respuesta Exitosa (200 OK):** Un objeto `VolunteerSummary` que contiene la lista de voluntarios, el total de cupos y cuántos están cubiertos.

---

#### `DELETE /championships/volunteers/:id`

-   **Acción:** Elimina una asignación de voluntario.
-   **Permisos:** `ADMIN` o `SUPER_ADMIN`.
-   **`:id`:** Corresponde al `id` de la asignación (`volunteer_assignments.id`).
-   **Respuesta Exitosa (200 OK):** Un mensaje de confirmación.
