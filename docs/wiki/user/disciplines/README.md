# Módulo de Disciplinas y Torneos

Este módulo es el corazón de la gestión deportiva del club. Se encarga de dos áreas principales:
1.  **Disciplinas y Grupos**: Define los deportes y actividades que ofrece el club (ej: Tenis, Fútbol, Natación) y organiza a los socios en grupos de entrenamiento.
2.  **Torneos**: Proporciona un sistema completo para crear y gestionar torneos, desde la inscripción de equipos hasta el seguimiento de los resultados.

---

## 1. Disciplinas y Grupos de Entrenamiento

### Listar Disciplinas

Permite a los usuarios ver todos los deportes y actividades que se practican en el club.

-   **Flujo**: Un usuario explora la sección de deportes y ve una lista con todas las disciplinas disponibles.
-   **Endpoint relacionado**: `GET /disciplines`

### Listar Grupos de Entrenamiento

Dentro de cada disciplina, los socios suelen estar organizados en grupos, por ejemplo, por nivel o edad (ej: "Tenis Infantil", "Fútbol Adultos Avanzado").

-   **Flujo**: Un usuario o administrador puede filtrar y ver los distintos grupos que existen para una disciplina y categoría específicas.
-   **Endpoint relacionado**: `GET /groups`

### Ver Alumnos de un Grupo (Admin/Entrenador)

Un administrador o el entrenador a cargo de un grupo puede consultar la lista de todos los socios inscritos en él.

-   **Endpoint relacionado**: `GET /groups/:id/students`

---

## 2. Gestión de Torneos

### Listar Torneos

Los socios pueden ver todos los torneos que están abiertos para inscripción, en curso o que ya han finalizado.

-   **Endpoint relacionado**: `GET /tournaments`

### Crear un Torneo (Admin)

Los administradores pueden crear nuevos torneos, definiendo el nombre, la disciplina, las fechas y el formato (ej: Liga, Eliminatoria).

-   **Endpoint relacionado**: `POST /tournaments`

### Inscribir un Equipo a un Torneo

Un socio (actuando como capitán) puede inscribir a su equipo en un torneo abierto.

-   **Flujo**: El capitán le da un nombre al equipo y selecciona a los socios que formarán parte de él.
-   **Endpoint relacionado**: `POST /tournaments/:id/teams`

### Programar un Partido (Admin)

Los administradores pueden definir los enfrentamientos del torneo, programando qué equipos juegan entre sí, cuándo y dónde.

-   **Endpoint relacionado**: `POST /tournaments/:id/matches`

### Ver Partidos de un Torneo

Todos los participantes y socios pueden ver la lista de partidos programados para un torneo, tanto los futuros como los que ya se han jugado.

-   **Endpoint relacionado**: `GET /tournaments/:id/matches`

### Actualizar Resultado de un Partido (Admin)

Después de que un partido ha finalizado, un administrador registra el resultado final (marcador).

-   **Endpoint relacionado**: `PUT /matches/:id/result`

### Ver la Tabla de Posiciones (Clasificación)

En cualquier momento durante un torneo, los socios pueden consultar la tabla de posiciones (leaderboard) para ver el rendimiento de cada equipo: puntos, partidos ganados, perdidos, etc.

-   **Endpoint relacionado**: `GET /tournaments/:id/standings`
