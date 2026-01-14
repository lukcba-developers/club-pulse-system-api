# Módulo de Estadísticas de Usuario (User Stats)

## 1. Propósito

El módulo de Estadísticas de Usuario proporciona una visión general del rendimiento y la actividad de un socio en las competiciones y actividades del club. Estas estadísticas son un componente clave de la gamificación y el seguimiento del progreso.

## 2. Funcionalidades

-   **Seguimiento de Rendimiento:** Registra métricas clave sobre la participación y el éxito del usuario en los eventos del club.
-   **Perfil Deportivo:** Ofrece una instantánea del nivel y la clasificación de un jugador, visible para otros socios y entrenadores.

## 3. Modelo de Datos

| Campo              | Tipo      | Descripción                                                                 |
| ------------------ | --------- | --------------------------------------------------------------------------- |
| `matches_played`   | `number`  | Número total de partidos jugados.                                           |
| `matches_won`      | `number`  | Número total de partidos ganados.                                           |
| `ranking_points`   | `number`  | Puntos de ranking acumulados (diferente de los XP de la tabla de líderes).    |
| `level`            | `number`  | Nivel actual del jugador, calculado a partir de su actividad y rendimiento. |
| `current_streak`   | `number`  | Racha actual de victorias o participación.                                  |

## 4. Endpoint de la API

### `GET /users/:id/stats`

-   **Acción:** Obtiene el objeto de estadísticas de un usuario.
-   **Permisos:**
    -   Un usuario puede consultar sus propias estadísticas usando el alias `/users/me/stats`.
    -   Un `ADMIN` o `SUPER_ADMIN` puede consultar las estadísticas de cualquier usuario.
-   **Respuesta Exitosa (200 OK):** Un objeto `UserStats`.

```json
{
  "data": {
    "matches_played": 58,
    "matches_won": 32,
    "ranking_points": 850,
    "level": 12,
    "current_streak": 4
  }
}
```
