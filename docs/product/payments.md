# Módulo de Pagos (Payments)

Este módulo actúa como la pasarela financiera del sistema, gestionando todas las transacciones, la facturación y la comunicación con proveedores de pago externos.

## 🌟 Funcionalidades Implementadas

### 1. Integración con Pasarelas de Pago
-   **Checkout Externo:** El sistema se integra con proveedores de pago (como Mercado Pago) para procesar los pagos de forma segura. Genera una orden de pago y redirige al usuario a la plataforma del proveedor para completar la transacción.
-   **Endpoint de Checkout:** Un endpoint `POST /payments/checkout` se encarga de crear la preferencia de pago con los detalles de la transacción (monto, descripción, etc.).

### 2. Conciliación Automática (Webhooks)
-   **Recepción de Notificaciones:** El sistema expone un endpoint de webhook (`POST /payments/webhook`) para recibir notificaciones en tiempo real desde la pasarela de pago cuando el estado de una transacción cambia.
-   **Actualización de Estado:** Al recibir una notificación, el sistema actualiza el estado del pago interno (ej: a `Approved` o `Rejected`) y dispara las acciones de negocio correspondientes (ej: confirmar una reserva o marcar una cuota como pagada).

### 3. Trazabilidad de Pagos
-   **Registro de Transacciones:** Cada intento de pago, ya sea exitoso o fallido, se registra en la base de datos, creando un historial financiero para cada socio.
-   **Contexto del Pago:** Cada transacción está vinculada a una referencia interna (como el ID de una reserva o una membresía), lo que permite saber qué se está pagando en cada momento.

## 4. Funcionalidades en Desarrollo

-   **Billetera Virtual (Wallet):** La capacidad de que los socios carguen saldo en una billetera virtual para realizar pagos más rápidos dentro de la plataforma es una funcionalidad del roadmap.
