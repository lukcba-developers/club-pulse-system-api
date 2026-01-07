# Módulo de Reservas (Booking)

## 1. Propósito

El módulo de **Reservas** es una de las funcionalidades centrales para el socio. Permite a los usuarios consultar la disponibilidad y reservar instalaciones (como canchas de tenis, pádel) o su lugar en clases grupales.

## 2. Funcionalidades Principales

-   **Calendario de Disponibilidad:** Muestra una vista de calendario intuitiva donde los socios pueden ver qué instalaciones o clases están disponibles y en qué horarios. El sistema utiliza **cacheo en Redis** para ofrecer una respuesta de disponibilidad casi instantánea.

-   **Proceso de Reserva:** Un flujo sencillo para que el socio seleccione un horario y confirme la reserva.
    -   **Gestión de Invitados:** Al realizar una reserva, un socio puede añadir los datos de un invitado (nombre, DNI). El sistema puede aplicar una tarifa adicional por cada invitado, la cual es validada por el backend.

-   **Gestión de Reservas Propias:** Los socios pueden ver un listado de sus próximas reservas y cancelarlas si la política del club lo permite.

-   **Listas de Espera Automatizadas:**
    -   Si una franja horaria está ocupada, el socio puede apuntarse a una lista de espera.
    -   Si una reserva se cancela, el sistema **automáticamente** notifica y crea una nueva reserva para el primer usuario en la lista de espera, asegurando la máxima ocupación de las instalaciones.

-   **Reglas de Reserva:** El administrador del club puede configurar reglas como:
    -   Cuántos días de antelación se puede reservar.
    -   Número máximo de reservas activas por socio.
    -   Coste de la reserva y política de cancelación.

-   **Integración con Pagos:** Si una reserva tiene un coste (por la propia reserva o por invitados), el sistema se integra con el **Módulo de Pagos**.

## 3. Funcionalidades en Desarrollo

-   **Reservas Recurrentes:** La capacidad de crear reservas que se repiten semanalmente (ej: "todos los martes a las 10:00") está planificada pero **aún no se encuentra implementada**. Los endpoints y la lógica de negocio son placeholders.

## 4. Endpoints de la API

| Verbo  | Ruta                      | Descripción                                             | Rol Requerido |
| :----- | :------------------------ | :------------------------------------------------------ | :------------ |
| `POST` | `/bookings`               | Crear una nueva reserva.                                | `MEMBER`      |
| `GET`  | `/bookings`               | Listar las reservas del usuario autenticado.            | `MEMBER`      |
| `GET`  | `/bookings/all`           | Listar todas las reservas del club (con filtros).       | `ADMIN`       |
| `GET`  | `/bookings/availability`  | Obtener los horarios disponibles para una instalación.  | `MEMBER`      |
| `DELETE`| `/bookings/:id`           | Cancelar una reserva.                                   | `MEMBER`      |
| `POST` | `/bookings/waitlist`      | Apuntarse a la lista de espera para un horario.         | `MEMBER`      |

## 5. Puntos de Mejora Identificados

-   **Configuración de Horarios:** Las horas de operación y la duración de los slots están actualmente hardcodeadas. Deberían obtenerse de la configuración de cada instalación en el backend.
-   **Tarifas de Invitados:** La tarifa por invitado está hardcodeada en el frontend. Esta lógica de precios debe moverse completamente al backend para garantizar seguridad y consistencia.
