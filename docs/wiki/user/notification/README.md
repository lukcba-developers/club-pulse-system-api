# Módulo de Notificaciones (Notification)

**Nota Importante:** Este es un módulo de servicio que opera exclusivamente en el backend. No expone una API pública para ser consumida directamente por el frontend o por usuarios. Su función es ser utilizado internamente por otros módulos del sistema para enviar comunicaciones a los socios.

## Propósito

El propósito del Módulo de Notificaciones es centralizar y gestionar el envío de todas las comunicaciones salientes del sistema a los usuarios. Al tener un servicio dedicado, se puede cambiar fácilmente de proveedor de envío (ej: cambiar de SendGrid a otro servicio de email) sin tener que modificar los demás módulos.

## Canales de Notificación Soportados

El servicio está diseñado para enviar notificaciones a través de diferentes canales:

-   **`EMAIL`**: Para comunicaciones formales como confirmaciones de pago, recibos, o reseteo de contraseñas. Se integra con proveedores como **SendGrid**.
-   **`SMS`**: Para alertas urgentes o recordatorios, como un cambio de último minuto en una reserva. Se integra con proveedores como **Twilio**.
-   **`PUSH`**: Para notificaciones en tiempo real a la aplicación móvil (funcionalidad futura).

## Ejemplos de Casos de Uso (Internos)

Otros módulos del sistema utilizan este servicio para notificar a los usuarios sobre eventos importantes. Por ejemplo:

-   El **Módulo de Reservas (Booking)** podría llamar a este servicio para:
    -   Enviar un email de confirmación cuando una reserva es creada.
    -   Enviar un SMS de recordatorio 24 horas antes de una reserva.

-   El **Módulo de Pagos (Payment)** podría usarlo para:
    -   Enviar un email con el recibo después de un pago exitoso.
    -   Notificar al usuario si un pago ha fallado.

-   El **Módulo de Autenticación (Auth)** lo usaría para:
    -   Enviar un email con el enlace para restablecer la contraseña.

En resumen, aunque los usuarios no interactúan directamente con este módulo, reciben el resultado de su trabajo cada vez que el sistema se comunica con ellos.
