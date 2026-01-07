# Lista de Tareas y Deuda Técnica

Este documento contiene una lista de tareas de mantenimiento, refactorización y finalización de funcionalidades que han sido identificadas durante el análisis del código para la creación de la documentación.

**La mayoría de las tareas sugeridas en la versión anterior de este documento han sido implementadas. Las siguientes son las tareas pendientes.**

---

## 1. Tareas Pendientes

### Tarea: Implementar la lógica de actualización de Gamificación

-   **Análisis:** El modelo `UserStats` y la lógica para subir de nivel (`AddExperience`) están definidos. Sin embargo, no existe código que llame a esta lógica. Por ejemplo, cuando un partido finaliza, no se actualizan las estadísticas del jugador (partidos ganados, XP, etc.).
-   **Acción Sugerida:**
    1.  En el módulo `championship`, al momento de actualizar el resultado de un partido (`UpdateMatchResult`), se debe obtener a los usuarios ganadores y perdedores.
    2.  Llamar a una función del `UserUseCases` para actualizar las `UserStats` correspondientes (incrementar `MatchesPlayed`, `MatchesWon`, y llamar a `AddExperience`).

## 2. Sugerencias de Refactorización

### Tarea: Mover la lógica de Billetera (Wallet) al módulo de Pagos

-   **Análisis:** El archivo `backend/internal/modules/user/domain/wallet.go` define los modelos de `Wallet` y `Transaction`. Aunque están relacionados con el usuario, conceptualmente pertenecen al dominio financiero, que está mayormente contenido en el módulo `payment`.
-   **Acción Sugerida:**
    1.  Mover los modelos `Wallet` y `Transaction` y su lógica asociada del módulo `user` al módulo `payment` para mejorar la cohesión y la separación de responsabilidades.
