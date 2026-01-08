# Módulo de Pagos (Payments)

Este módulo actúa como la pasarela financiera del sistema, gestionando todas las transacciones, la facturación y la comunicación con proveedores de pago externos.

## 🌟 Funcionalidades Implementadas

### 1. Integración con Pasarelas de Pago
-   **Checkout Externo:** Integración con Mercado Pago para pagos digitales.
-   **Pagos Offline:** Soporte para registro manual de pagos en **Efectivo (Cash)**, **Transferencia** e **Intercambio de Mano de Obra (Labor Exchange)**.
-   **Notas de Auditoría:** Cada pago offline permite adjuntar una nota justificativa del cobro o del trabajo realizado.

### 2. Conciliación Automática (Webhooks)
-   **Recepción de Notificaciones:** Webhook para Mercado Pago (`POST /payments/webhook`).
-   **Actualización de Estado:** Sincronización automática de estados: `PENDING`, `COMPLETED`, `FAILED`, `REFUNDED`.

### 3. Trazabilidad de Pagos
-   **Contexto Multi-tenant:** Cada transacción está aislada por `club_id` y vinculada a un `payer_id`.
-   **Referencias:** Vínculo con el ID de la reserva o membresía correspondiente.

## 4. Funcionalidades en Desarrollo

-   **Billetera Virtual (Wallet):** Saldo prepago para socios.
-   **Suscripciones Recurrentes Automáticas:** Gestión de débitos automáticos.
