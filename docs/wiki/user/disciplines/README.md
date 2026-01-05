# Módulo de Disciplinas

Este módulo define los deportes y actividades que ofrece el club (ej: Tenis, Fútbol, Natación) y organiza a los socios en grupos de entrenamiento.

La gestión de competencias y torneos se realiza en el [Módulo de Campeonatos (Championship)](../championship/README.md).

---

## Casos de Uso

### 1. Listar Disciplinas

Permite a los usuarios ver todos los deportes y actividades que se practican en el club.

-   **Flujo**: Un usuario explora la sección de deportes y ve una lista con todas las disciplinas disponibles.
-   **Endpoint relacionado**: `GET /disciplines`

### 2. Listar Grupos de Entrenamiento

Dentro de cada disciplina, los socios suelen estar organizados en grupos, por ejemplo, por nivel o edad (ej: "Tenis Infantil", "Fútbol Adultos Avanzado").

-   **Flujo**: Un usuario o administrador puede filtrar y ver los distintos grupos que existen para una disciplina y categoría específicas.
-   **Endpoint relacionado**: `GET /disciplines/:id/groups`

### 3. Ver Alumnos de un Grupo (Admin/Entrenador)

Un administrador o el entrenador a cargo de un grupo puede consultar la lista de todos los socios inscritos en él.

-   **Endpoint relacionado**: `GET /groups/:id/students`
