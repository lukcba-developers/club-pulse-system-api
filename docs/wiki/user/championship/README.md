# Módulo de Campeonatos (Championship)

Este módulo se encarga de todo lo relacionado con la creación y gestión de torneos y competencias dentro del club.

La gestión de los [Equipos (Team)](../team/README.md) que participan en estos campeonatos se detalla en su propio módulo.

---

## Casos de Uso

### 1. Listar Campeonatos

Los socios pueden ver todos los torneos que están abiertos para inscripción, en curso o que ya han finalizado.

-   **Endpoint relacionado**: `GET /championships`

### 2. Crear un Campeonato (Admin)

Los administradores pueden crear nuevos torneos, definiendo el nombre, la disciplina, las fechas y el formato (ej: Liga, Eliminatoria).

-   **Endpoint relacionado**: `POST /championships`

### 3. Programar un Partido (Admin)

Los administradores pueden definir los enfrentamientos del torneo, programando qué equipos juegan entre sí, cuándo y dónde.

-   **Endpoint relacionado**: `POST /championships/:id/matches`

### 4. Ver Partidos de un Campeonato

Todos los participantes y socios pueden ver la lista de partidos programados para un torneo, tanto los futuros como los que ya se han jugado.

-   **Endpoint relacionado**: `GET /championships/:id/matches`

### 5. Actualizar Resultado de un Partido (Admin)

Después de que un partido ha finalizado, un administrador registra el resultado final (marcador).

-   **Endpoint relacionado**: `PUT /championships/matches/:id/result`

### 6. Ver la Tabla de Posiciones (Clasificación)

En cualquier momento durante un torneo, los socios pueden consultar la tabla de posiciones (leaderboard) para ver el rendimiento de cada equipo: puntos, partidos ganados, perdidos, etc.

-   **Endpoint relacionado**: `GET /championships/:id/standings`
