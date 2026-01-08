# Manual de Usuario: Módulo de Equipos (Team)

## 1. Propósito

Este módulo te permite crear y gestionar tus propios equipos para competir en los torneos del club. Puedes juntarte con tus amigos, elegir un nombre y un escudo, y prepararse para la competición.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede crear equipos, ser capitán, unirse a equipos e invitar a otros socios.

---

## 3. Guía de Usuario (Rol: `MEMBER`)

### 🔹 Cómo Crear un Nuevo Equipo

Si quieres ser el capitán, puedes crear tu propio equipo.

**Paso a paso:**
1.  **Navega a la sección "Mis Equipos"** en tu perfil o en el menú de "Campeonatos".
2.  Haz clic en el botón **"Crear Nuevo Equipo"**.
3.  **Completa el formulario:**
    -   **Nombre del Equipo:** ¡Elige un nombre original!
    -   **Logo/Escudo:** Sube una imagen para representar a tu equipo.
4.  **Guarda los cambios.** ¡Tu equipo ha sido creado y tú eres el capitán!

### 🔹 Cómo Invitar Jugadores a tu Equipo

Como capitán, eres el encargado de reclutar a tus compañeros.

**Paso a paso:**
1.  Ve a la página de gestión de tu equipo.
2.  Busca la opción **"Invitar Jugadores"**.
3.  Se abrirá un buscador donde podrás **encontrar a otros socios** del club por su nombre.
4.  Selecciona a los socios que quieres invitar y haz clic en **"Enviar Invitación"**.
5.  Los socios recibirán una notificación para unirse a tu equipo.

### 🔹 Cómo Aceptar o Rechazar una Invitación

Si un capitán te invita a su equipo, recibirás una notificación.

**Paso a paso:**
1.  Ve a tu panel de notificaciones o a la sección "Mis Equipos".
2.  Verás la invitación pendiente con el nombre del equipo.
3.  Tendrás los botones **"Aceptar"** y **"Rechazar"**. Haz clic en la opción que prefieras.
4.  Si aceptas, pasarás a formar parte del equipo.

### 🔹 Cómo Salir de un Equipo

**Paso a paso:**
1.  Ve a la página del equipo del que formas parte.
2.  Busca la opción **"Abandonar Equipo"**.
3.  Confirma tu decisión. Dejarás de ser miembro de ese equipo.

---

## 4. Diagrama de Flujo: Creación y Formación de un Equipo

```mermaid
graph TD
    A[Capitán: Clic en "Crear Equipo"] --> B[Rellena Nombre y Logo];
    B --> C[Equipo Creado ✅];
    C --> D[Capitán: Invita a Jugadores];
    D --> E[Jugador Invitado: Recibe Notificación];
    E --> F{¿Acepta la Invitación?};
    F -- Sí --> G[Jugador se une al Equipo];
    F -- No --> H[Invitación Rechazada];
    G --> I[Equipo listo para competir];
```

---

## 5. Eventos de Viaje y Convocatorias (RSVP)

Este sub-módulo permite a los capitanes y administradores organizar eventos para un equipo, como viajes a torneos, partidos amistosos o entrenamientos especiales. Los miembros del equipo pueden confirmar o declinar su asistencia.

### 🔹 Propósito

-   **Organización Centralizada:** Planificar la logística de un viaje o partido, incluyendo destinos, fechas y puntos de encuentro.
-   **Gestión de Asistencia (RSVP):** Saber con antelación qué jugadores asistirán a un evento.
-   **Cálculo de Costos:** Estimar y calcular los costos del evento y dividirlos entre los participantes confirmados.

### 🔹 Modelo de Datos

La funcionalidad se basa en dos tablas principales: `travel_events` y `event_rsvps`.

#### `travel_events`

| Campo             | Tipo            | Descripción                                                              |
| ----------------- | --------------- | ------------------------------------------------------------------------ |
| `id`              | `UUID`          | ID único del evento.                                                     |
| `team_id`         | `UUID`          | ID del equipo para el que se organiza el evento.                         |
| `type`            | `ENUM`          | Tipo de evento: `TRAVEL`, `MATCH`, `TOURNAMENT`, `TRAINING`.             |
| `title`           | `VARCHAR`       | Título del evento.                                                       |
| `destination`     | `VARCHAR`       | Lugar de destino del evento.                                             |
| `departure_date`  | `TIMESTAMPTZ`   | Fecha y hora de salida/inicio.                                           |
| `estimated_cost`  | `DECIMAL`       | Costo total estimado del evento.                                         |
| `actual_cost`     | `DECIMAL`       | Costo real final del evento (se actualiza a posteriori).                 |
| `cost_per_person` | `DECIMAL`       | Costo por persona, calculado automáticamente.                            |
| `max_participants`| `INT`           | Número máximo de asistentes (opcional).                                  |
| `created_by`      | `VARCHAR`       | ID del usuario que creó el evento.                                       |

#### `event_rsvps`

| Campo         | Tipo          | Descripción                                                              |
| ------------- | ------------- | ------------------------------------------------------------------------ |
| `id`          | `UUID`        | ID único de la respuesta.                                                |
| `event_id`    | `UUID`        | ID del evento al que se responde.                                        |
| `user_id`     | `VARCHAR`     | ID del usuario que responde.                                             |
| `status`      | `ENUM`        | Estado de la respuesta: `PENDING`, `CONFIRMED`, `DECLINED`.              |
| `notes`       | `TEXT`        | Notas adicionales del usuario (ej. "Llego 15 minutos tarde").            |
| `responded_at`| `TIMESTAMPTZ` | Fecha y hora de la respuesta.                                            |

### 🔹 Endpoints de la API

---

#### `POST /events`

-   **Acción:** Crea un nuevo evento de viaje/partido para un equipo.
-   **Permisos:** `ADMIN`, `SUPER_ADMIN` o Capitán del equipo.
-   **Request Body (JSON):** Un objeto `TravelEvent` con los detalles del evento.
-   **Respuesta Exitosa (201 Created):** El objeto del evento creado.

---

#### `GET /teams/:teamId/events`

-   **Acción:** Lista todos los eventos asociados a un equipo específico.
-   **Permisos:** Miembros del equipo, `ADMIN`, `SUPER_ADMIN`.
-   **Respuesta Exitosa (200 OK):** Un array de objetos de eventos.

---

#### `GET /events/:eventId`

-   **Acción:** Obtiene los detalles de un evento específico.
-   **Permisos:** Miembros del equipo, `ADMIN`, `SUPER_ADMIN`.
-   **Respuesta Exitosa (200 OK):** El objeto del evento.

---

#### `POST /events/:eventId/rsvp`

-   **Acción:** Permite a un usuario responder a la convocatoria de un evento (confirmar o declinar asistencia).
-   **Permisos:** Miembro del equipo invitado.
-   **Request Body (JSON):**
    ```json
    {
      "status": "CONFIRMED", // o "DECLINED"
      "notes": "Puedo llevar a 2 personas en mi auto." // Opcional
    }
    ```
-   **Respuesta Exitosa (200 OK):** Un mensaje de confirmación.

---

#### `GET /events/:eventId/summary`

-   **Acción:** Devuelve un resumen completo del evento, incluyendo estadísticas de asistencia y el costo calculado por persona.
-   **Permisos:** Miembros del equipo, `ADMIN`, `SUPER_ADMIN`.
-   **Respuesta Exitosa (200 OK):** Un objeto `EventSummary` con todos los detalles.
