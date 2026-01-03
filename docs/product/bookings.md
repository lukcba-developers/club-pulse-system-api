# Módulo de Reservas (Booking Engine)

Un motor transaccional de alta concurrencia diseñado para gestionar la agenda del club con precisión de milisegundos.

## 🌟 Funcionalidades Principales

### 1. Reservas en Tiempo Real
El corazón del sistema. Permite a los usuarios asegurar un espacio en segundos.
-   **Integridad de Datos**: Garantía absoluta de no sobreventa (overbooking).
-   **Mecanismo**: Utiliza bloqueos a nivel de base de datos (`CheckAvailability` con rangos de tiempo) antes de confirmar cualquier transacción.

### 2. Motor de Disponibilidad Inteligente
Endpoint core (`/availability`) que responde a la pregunta *"¿Qué hay libre?"*.
-   **Cálculo Dinámico**: Intersecta el horario operativo del club con las reservas existentes.
-   **Lógica de Slots**: Fragmenta el tiempo en bloques jugables (ej. 60 min, 90 min) según la configuración del deporte.

### 3. Reglas de Recurrencia
Para usuarios habituales y escuelas deportivas.
-   **Patrones Flexibles**:
    -   "Todos los Lunes a las 19:00".
    -   "Martes y Jueves por 2 meses".
-   **Materialización**: El sistema genera las instancias (bookings individuales) automáticamente, validando disponibilidad para cada una.

### 4. Ciclo de Vida de la Reserva
Estados claros para gestión operativa:
-   `Confirmed`: Pago o seña realizada (o usuario de confianza).
-   `Pending`: Reservado temporalmente esperando pago (TTL de 15 min).
-   `Cancelled`: Liberada por usuario o admin.
-   `Completed`: El turno ya ocurrió.
-   `No-Show`: El usuario no asistió (para estadísticas y penalizaciones).

### 5. Validación de Políticas
-   **Ventana de Reserva**: "¿Con cuánta anticipación puedo reservar?" (ej. máx 14 días).
-   **Límite de Reservas**: Restricción de reservas simultáneas por usuario para equidad.
-   **Conflictos**: Verificación automática contra mantenimientos programados.
