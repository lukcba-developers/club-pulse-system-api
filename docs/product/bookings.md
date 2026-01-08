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

-   **Reservas Recurrentes:** Capacidad de crear patrones de reserva (ej: todos los lunes de 18:00 a 19:00) mediante reglas de recurrencia que materializan reservas automáticamente.
-   **Validación de Salud:** Bloqueo preventivo de reservas si el socio no cuenta con un **Certificado Médico** vigente.

-   **Ciclo de Vida de la Reserva:**
    -   El modelo de dominio soporta los estados: `CONFIRMED` y `CANCELLED`.

## 3. Funcionalidades en Desarrollo (No Implementadas)

-   **Seguimiento de Asistencia (No-Show):** Detección automática de ausencias para penalizaciones o liberación de slots.
-   **Matchmaking:** Sistema para encontrar compañeros de juego basado en nivel y proximidad.

*Este documento refleja el estado actual del código y se actualizará a medida que las funcionalidades en desarrollo se completen.*
