# Módulo de Pagos (Payment)

## 1. Propósito

El módulo de **Pagos** gestiona todas las transacciones financieras dentro del club. Se encarga de la facturación, el procesamiento de pagos y el seguimiento del estado financiero de los socios.

## 2. Funcionalidades Principales

-   **Generación de Facturas:**
    -   **Recurrentes:** Genera automáticamente las facturas de las cuotas de membresía según el ciclo de cada socio (**Módulo de Membresías**).
    -   **Puntuales:** Permite crear facturas por conceptos específicos, como la inscripción a un torneo, la compra en la tienda o una reserva con coste.
-   **Procesamiento de Pagos:**
    -   **Pasarelas de Pago:** Se integra con proveedores de pago online (como Stripe, Mercado Pago) para que los socios puedan pagar sus facturas con tarjeta de crédito/débito.
    -   **Pagos en Recepción:** Permite al personal del club registrar pagos realizados en efectivo o con tarjeta en el punto de venta físico.
-   **Billetera Virtual (Wallet):** Cada socio dispone de un saldo o billetera virtual. Los socios pueden cargar saldo y utilizarlo para pagar servicios del club de forma rápida.
-   **Historial de Transacciones:** Tanto los socios como los administradores pueden consultar un historial detallado de todas las facturas y pagos realizados.
-   **Gestión de Deudas:** El sistema identifica automáticamente a los socios con pagos vencidos y puede aplicar las restricciones correspondientes (ej: bloquear nuevas reservas).