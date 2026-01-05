# Módulo de Equipo (Team)

Este módulo gestiona la información de los equipos que participan en los [campeonatos](../championship/README.md) del club.

---

## Casos de Uso

### 1. Registrar un Equipo en un Campeonato

Un socio (actuando como capitán) puede inscribir a su equipo en un campeonato que esté abierto para inscripciones.

-   **Flujo**:
    1.  El capitán elige un campeonato.
    2.  Le da un nombre a su equipo.
    3.  Selecciona a los otros socios del club que formarán parte del equipo.
    4.  El sistema registra el equipo y lo asocia al campeonato.
-   **Endpoint relacionado**: `POST /championships/:id/teams`

### 2. Ver Detalles de un Equipo

Se puede consultar la información de un equipo específico, como sus miembros.

-   **Endpoint relacionado**: `GET /teams/:id`
