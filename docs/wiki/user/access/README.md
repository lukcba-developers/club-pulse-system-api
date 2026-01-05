# Módulo de Acceso (Access)

Este módulo gestiona el control de acceso físico a las instalaciones del club. Su principal responsabilidad es determinar si un socio o usuario tiene permiso para entrar en un momento dado.

No es un módulo con el que los usuarios interactúen directamente a través de una interfaz web, sino que es un servicio consumido por dispositivos físicos como lectores de códigos QR en los tornos de entrada.

## Casos de Uso

### 1. Validación de Entrada

Este es el caso de uso central del módulo.

-   **Flujo**:
    1.  Un socio presenta su identificador en un punto de acceso (ej: muestra un código QR de la app en un lector).
    2.  El dispositivo lector envía el identificador del socio al sistema Club Pulse.
    3.  El sistema verifica en tiempo real el estado del socio:
        -   ¿Tiene una [membresía activa](../membership/README.md)?
        -   ¿Tiene una [reserva para una instalación](../booking/README.md) en este momento?
        -   ¿Tiene alguna [deuda pendiente](../payment/README.md) que le impida el acceso?
    4.  Basado en estas reglas, el sistema responde al dispositivo con una de dos respuestas:
        -   `GRANTED` (Acceso Concedido): El torno se abre.
        -   `DENIED` (Acceso Denegado): El torno permanece cerrado y se puede mostrar un mensaje al socio (ej: "Membresía vencida").
-   **Endpoint relacionado**: `POST /access/entry`

### 2. Simulación de Entrada (Admin)

El frontend proporciona una función para que un administrador pueda simular la entrada de un socio. Esto es útil para realizar pruebas o para verificar por qué a un socio se le podría estar denegando el acceso, sin necesidad de que esté físicamente en el torno.

-   **Flujo**:
    1.  Un administrador busca a un socio en el panel de administración.
    2.  Hace clic en "Simular Entrada".
    3.  El sistema realiza la misma validación que haría en un torno real y muestra el resultado (`GRANTED` o `DENIED`) y el motivo.
-   **Función del Frontend**: `accessService.simulateEntry`
