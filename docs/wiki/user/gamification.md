# Módulo de Gamificación

## 1. Propósito

El módulo de Gamificación introduce un sistema de competición y recompensas para aumentar la participación y el compromiso de los socios en el club. Se basa en la obtención de Puntos de Experiencia (XP), la consecución de insignias y la clasificación en tablas de líderes.

## 2. Sistema de Insignias (Badges)

Las insignias son logros que los usuarios pueden desbloquear al completar ciertas acciones o hitos dentro del club.

### Funcionalidades

-   **Obtención de Insignias:** Los usuarios reciben insignias automáticamente al cumplir los criterios definidos (ej: "Primer Partido Ganado", "Racha de 5 Días").
-   **Galería de Insignias:** Una sección donde los usuarios pueden ver todas las insignias disponibles en el sistema y cuáles han obtenido.
-   **Insignias Destacadas:** Los usuarios pueden seleccionar hasta 3 de sus insignias para que se muestren de forma prominente en su perfil público.

### Endpoints de la API

-   `GET /gamification/badges`: Lista todas las insignias disponibles en el club.
-   `GET /gamification/badges/my`: Devuelve todas las insignias obtenidas por el usuario autenticado.
-   `GET /gamification/badges/featured/:user_id`: Muestra las insignias destacadas de un usuario específico.
-   `PUT /gamification/badges/:badge_id/feature`: Permite al usuario marcar o desmarcar una de sus insignias como "destacada".

## 3. Tablas de Clasificación (Leaderboards)

Las tablas de clasificación (leaderboards) muestran el ranking de los socios basado en los Puntos de Experiencia (XP) acumulados.

### Funcionalidades

-   **Ranking Global:** Muestra una lista de los mejores jugadores del club.
-   **Filtros por Período:** El ranking se puede consultar en diferentes marcos de tiempo:
    -   `DAILY`: Ranking del día.
    -   `WEEKLY`: Ranking de la semana.
    -   `MONTHLY`: Ranking del mes (por defecto).
    -   `ALL_TIME`: Ranking histórico.
-   **Vista de Contexto:** Permite a un usuario ver su propia posición en el ranking, junto con los jugadores que están inmediatamente por encima y por debajo de él.

### Endpoints de la API

-   `GET /gamification/leaderboard`: Devuelve la tabla de clasificación global, con opciones de paginación y filtrado por período.
-   `GET /gamification/leaderboard/context`: Devuelve la vista de contexto del usuario en la tabla de clasificación.
-   `GET /gamification/leaderboard/rank`: Obtiene el número de ranking exacto del usuario autenticado.
