# Lista de Tareas y Deuda Técnica

Este documento contiene una lista de tareas de mantenimiento, refactorización y finalización de funcionalidades que han sido identificadas durante el análisis del código para la creación de la documentación.

---

## 1. Refactorización del Módulo de Usuario

### Tarea: Renombrar `gamification.go` y separar conceptos

-   **Análisis:** El archivo `backend/internal/modules/user/domain/gamification.go` está mal nombrado. Actualmente define los modelos de `Wallet` (billetera) y `Transaction` (transacciones), que son conceptos financieros. El verdadero modelo de gamificación (`UserStats`) se encuentra en `stats.go`.
-   **Acción Sugerida:**
    1.  Renombrar `gamification.go` a `wallet.go`.
    2.  Considerar si los conceptos de Billetera y Transacciones deberían pertenecer a un módulo más financiero (como `payment`) en lugar de `user`, para mejorar la cohesión de los módulos.

## 2. Finalización de Funcionalidades Incompletas

### Tarea: Implementar la lógica de actualización de Gamificación

-   **Análisis:** El modelo `UserStats` y la lógica para subir de nivel (`AddExperience`) están definidos. Sin embargo, no existe código que llame a esta lógica. Por ejemplo, cuando un partido finaliza, no se actualizan las estadísticas del jugador (partidos ganados, XP, etc.).
-   **Acción Sugerida:**
    1.  En el módulo `disciplines`, al momento de actualizar el resultado de un partido (`UpdateMatchResult`), se debe obtener a los usuarios ganadores y perdedores.
    2.  Llamar a una función del `UserUseCases` para actualizar las `UserStats` correspondientes (incrementar `MatchesPlayed`, `MatchesWon`, y llamar a `AddExperience`).

### Tarea: Exponer la gestión de Becas (Scholarships) en la API

-   **Análisis:** El sistema tiene la capacidad de crear y aplicar becas (descuentos porcentuales) a los usuarios. La lógica de dominio y el repositorio están implementados. Sin embargo, no existen endpoints en la API para que un administrador pueda crear, asignar o revocar estas becas.
-   **Acción Sugerida:**
    1.  Crear un nuevo `handler` en el módulo `membership`.
    2.  Añadir endpoints protegidos para administradores (`ADMIN` role) que permitan realizar operaciones CRUD sobre las becas de un usuario (ej: `POST /users/{id}/scholarship`, `DELETE /scholarships/{id}`).

### Tarea: Exponer la funcionalidad de Lista de Espera (Waitlist) en la API

-   **Análisis:** El frontend (`booking-service.ts`) intenta llamar a un endpoint `POST /bookings/waitlist`, y la lógica de negocio y de base de datos para una lista de espera existe. Sin embargo, la ruta no está registrada en el `handler` principal de `booking`.
-   **Acción Sugerida:**
    1.  Verificar por qué la ruta no está expuesta.
    2.  Añadir el endpoint `POST /bookings/waitlist` al `handler` de `booking`, protegiéndolo con el middleware de autenticación.

## 3. Corrección de Inconsistencias

### Tarea: Sincronizar el `membership-service.ts` del Frontend con el Backend

-   **Análisis:** El servicio de membresías del frontend tiene llamadas a endpoints que no coinciden con los del backend.
    -   Llama a `/memberships/plans` en lugar de `/memberships/tiers`.
    -   Llama a `/memberships/subscribe` en lugar de usar el endpoint `POST /memberships`.
    -   Intenta usar `DELETE /memberships` que no existe en el backend.
-   **Acción Sugerida:**
    1.  Refactorizar `membership-service.ts` para que utilice los endpoints correctos definidos en el `handler` del backend.

### Tarea: Consolidar el endpoint de Disponibilidad

-   **Análisis:** Existen dos endpoints para consultar la disponibilidad: uno en el módulo `facilities` y otro en `booking`. El de `booking` es más completo ya que incluye una capa de caché.
-   **Acción Sugerida:**
    1.  Eliminar el endpoint de disponibilidad del `handler` de `facilities`.
    2.  Mantener y usar exclusivamente `GET /bookings/availability` para esta funcionalidad, asegurando que todos los componentes del frontend lo utilicen.
