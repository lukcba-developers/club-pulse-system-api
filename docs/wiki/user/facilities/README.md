# Módulo de Instalaciones (Facilities)

## 1. Propósito

El módulo de **Instalaciones** gestiona el catálogo de todos los recursos físicos que pueden ser reservados o utilizados por los socios. Es la base sobre la cual opera el **Módulo de Reservas (Booking)**.

## 2. Funcionalidades Principales

-   **Catálogo de Instalaciones:** Permite a los administradores crear y gestionar una lista de todas las instalaciones del club. Para cada una se puede definir:
    -   Nombre (ej: "Cancha de Pádel 1").
    -   Tipo (ej: "Cancha", "Piscina", "Sala de Musculación").
    -   Capacidad máxima.
    -   Ubicación dentro del club.
    -   Descripción y fotos.
-   **Gestión de Estado:** Permite cambiar el estado de una instalación. Los estados típicos son:
    -   `Disponible`: Abierta para reservas.
    -   `En Mantenimiento`: No disponible temporalmente.
    -   `Cerrada`: Fuera de servicio por un periodo prolongado.
-   **Horarios de Funcionamiento:** Se puede definir un horario de apertura específico para cada instalación, que puede ser diferente al horario general del club.
-   **Asociación a Disciplinas:** Permite asociar instalaciones con las disciplinas que las utilizan, facilitando la búsqueda y la programación.