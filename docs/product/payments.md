# Módulo de Pagos (Payments)

Pasarela financiera centralizada para el procesamiento seguro de transacciones.

## 🌟 Funcionalidades Principales

### 1. Integración MercadoPago
Implementación nativa del proveedor líder en Latam.
-   **Checkout Pro**: Redirección segura a la pasarela de MP para completar el pago (Tarjeta, Saldo MP, etc).
-   **Preferencias**: Creación dinámica de órdenes de pago con items detallados.

### 2. Procesamiento de Webhooks (IPN)
Sistema de conciliación automática en tiempo real.
-   **Funcionamiento**: MercadoPago notifica a nuestro endpoint `/webhook` cuando un pago cambia de estado.
-   **Lógica**: El sistema busca el pago interno por ID externo y actualiza su estado (`Approved`, `Rejected`), disparando acciones (ej. confirmar la reserva asociada).

### 3. Muelle de Transacciones (Transaction Log)
Libro mayor inmutable de intentos de pago.
-   **Estados**: `Pending` -> `Approved` / `Rejected` / `Refunded`.
-   **Moneda**: Soporte inicial ARS, extensible multi-divisa.
-   **Métodos**: Registro del método utilizado (CC, Debit, Cash, Transfer).

### 4. Contexto de Referencia
Cada pago nace con un propósito.
-   **Vinculación**: Campos `ReferenceID` y `ReferenceType`.
-   **Tipos**:
    -   `BOOKING`: Pago de un turno de cancha.
    -   `MEMBERSHIP`: Pago de cuota social.
    -   `WALLET_TOPUP`: Carga de saldo en billetera virtual (Roadmap).
