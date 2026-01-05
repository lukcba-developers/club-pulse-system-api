# Módulo de Usuario (User)

Este módulo se centra en la gestión de toda la información relacionada con los perfiles de los usuarios. Permite a los socios ver y actualizar sus datos, gestionar sus grupos familiares y consultar información relevante sobre su actividad en el club. También proporciona herramientas a los administradores para gestionar a los miembros del club.

## Casos de Uso para Socios (Members)

### 1. Gestión del Perfil Personal

Cada usuario tiene un perfil con su información personal, que puede consultar y modificar.

-   **Ver Perfil**: Un usuario puede ver su propia información de perfil, incluyendo nombre, email, y estado de su certificado médico.
    -   **Endpoint relacionado**: `GET /users/me`

-   **Actualizar Perfil**: Un usuario puede actualizar ciertos datos de su perfil, como su información de contacto.
    -   **Endpoint relacionado**: `PUT /users/me`

### 2. Gestión del Grupo Familiar

La plataforma permite que un titular de cuenta (`MEMBER`) gestione las cuentas de sus hijos o dependientes.

-   **Ver Hijos/Dependientes**: El titular puede ver una lista de todos los miembros asociados a su grupo familiar.
    -   **Endpoint relacionado**: `GET /users/me/children`

-   **Registrar un Hijo/Dependiente**: El titular puede crear una nueva cuenta para un hijo o dependiente, asociándola a su grupo familiar.
    -   **Endpoint relacionado**: `POST /users/me/children`

### 3. Consulta de Actividad y Finanzas

Los usuarios pueden hacer seguimiento de su progreso y estado de cuenta en el club.

-   **Ver Estadísticas (Gamification)**: Muestra las estadísticas de actividad del usuario. El sistema de gamificación está diseñado para que los socios ganen puntos de experiencia (XP) y suban de nivel al participar en actividades, como jugar y ganar partidos. Esto se refleja en sus estadísticas.
    -   **Endpoint relacionado**: `GET /users/:id/stats`

-   **Consultar Billetera (Wallet)**: Permite ver el saldo actual de la billetera virtual del usuario y sus puntos de fidelidad.
    -   **Endpoint relacionado**: `GET /users/:id/wallet`

## Casos de Uso para Administradores (Admins)

### 1. Gestión de Miembros del Club

Los administradores (`ADMIN` o `SUPER_ADMIN`) tienen herramientas para gestionar la base de usuarios de su club.

-   **Listar y Buscar Usuarios**: Un administrador puede obtener una lista de todos los usuarios de su club y realizar búsquedas por nombre o email para encontrarlos rápidamente.
    -   **Endpoint relacionado**: `GET /users`

-   **Eliminar un Usuario**: Un administrador puede eliminar la cuenta de un usuario del club. Esta es una acción destructiva y debe usarse con precaución.
    -   **Endpoint relacionado**: `DELETE /users/:id`

---

## Estructura de Datos (Interfaces)

A continuación se detallan las estructuras de datos clave que se manejan en este módulo.

### Interfaz `User`

Representa el perfil básico de un usuario.

-   `id`: Identificador único del usuario.
-   `name`: Nombre completo.
-   `email`: Correo electrónico.
-   `role`: Rol del usuario en el sistema (ej: `MEMBER`, `ADMIN`).
-   `medical_cert_status`: Estado del certificado médico (`VALID`, `EXPIRED`, `PENDING`).
-   `medical_cert_expiry`: Fecha de vencimiento del certificado médico.
-   `family_group_id`: Identificador del grupo familiar al que pertenece.

### Interfaz `UserStats`

Representa las estadísticas de "gamification" de un usuario. El sistema está diseñado para que estas estadísticas se actualicen automáticamente cuando ocurren ciertos eventos (ej: se registra el resultado de un partido).

-   `matches_played`: Número total de partidos jugados.
-   `matches_won`: Número total de partidos ganados.
-   `ranking_points`: Puntos de ranking acumulados.
-   `level`: Nivel actual del jugador, que sube al acumular experiencia.
-   `experience`: Puntos de experiencia acumulados para el nivel actual.
-   `current_streak`: Racha actual de victorias o participación.

**Nota sobre la implementación**: Aunque el modelo de datos y la lógica para subir de nivel (`AddExperience`) están definidos en el código, la funcionalidad que actualiza estas estadísticas (ej: otorgar XP después de un partido) no parece estar implementada o conectada en la versión actual del backend.

### Interfaz `Wallet`

Representa la billetera virtual de un usuario.

-   `id`: Identificador único de la billetera.
-   `balance`: Saldo monetario de la cuenta (ej: para pagar reservas).
-   `points`: Puntos de fidelidad acumulados.
