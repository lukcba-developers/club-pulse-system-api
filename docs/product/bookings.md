# Módulo de Reservas (Booking Engine)

## 1. Propósito

El motor de **Reservas** es una de las funcionalidades centrales para el socio. Permite a los usuarios consultar la disponibilidad y reservar instalaciones (como canchas de tenis, pádel) o su lugar en clases grupales, garantizando la integridad de los datos y previniendo la sobreventa (overbooking).

## 2. Funcionalidades Implementadas

-   **Calendario de Disponibilidad en Tiempo Real:**
    -   Un endpoint (`GET /bookings/availability`) calcula y devuelve los horarios disponibles para una instalación en una fecha específica.
    -   Utiliza **cacheo en Redis** para ofrecer respuestas de alta velocidad en consultas repetidas.

-   **Creación de Reservas:**
    -   Los usuarios pueden crear una reserva para un slot disponible.
    -   **Gestión de Invitados:** El sistema permite añadir los datos de un invitado al crear la reserva. El backend valida cualquier tarifa asociada.

-   **Cancelación de Reservas:**
    -   Los usuarios pueden cancelar sus propias reservas a través de la API.

-   **Listas de Espera Automatizadas:**
    -   Si un horario está ocupado, un usuario puede unirse a una lista de espera (`POST /bookings/waitlist`).
    -   Si la reserva original se cancela, el sistema automáticamente promueve al primer usuario de la lista de espera, creándole una reserva y notificándole.

-   **Ciclo de Vida de la Reserva:**
    -   El modelo de dominio actual soporta los estados: `CONFIRMED` y `CANCELLED`.

## 3. Funcionalidades en Desarrollo (No Implementadas)

-   **Reservas Recurrentes:** La capacidad de crear patrones de reserva (ej. "todos los lunes") está definida en el código pero no es funcional.
-   **Estados Adicionales:** Los estados como `Pending` (esperando pago), `Completed` (completada) o `No-Show` (ausente) son parte de la visión del producto pero no están implementados en el modelo de dominio actual.

*Este documento refleja el estado actual del código y se actualizará a medida que las funcionalidades en desarrollo se completen.*
