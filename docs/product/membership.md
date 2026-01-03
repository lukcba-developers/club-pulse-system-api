# Módulo de Membresías (Membership)

Sistema integral para la gestión de la fidelización de socios y flujos de ingresos recurrentes.

## 🌟 Funcionalidades Principales

### 1. Niveles de Membresía (Tiers)
Configuración de productos de suscripción diferenciados.
-   **Personalización**:
    -   *Nombre*: "Socio Pleno", "Pase Fin de Semana", "Estudiante".
    -   *Fee*: Costo base del plan.
    -   *Beneficios*: Descuentos en reservas, prioridad de booking, acceso a áreas exclusivas.

### 2. Ciclos de Facturación (Billing Cycles)
Flexibilidad total en la periodicidad de cobro.
-   **Mensual**: El estándar de la industria.
-   **Trimestral / Semestral**: Para promociones estacionales.
-   **Anual**: Para socios vitalicios o largo plazo.

### 3. Automatización de Cobros
Job automatizado (`ProcessMonthlyBilling`) que corre periódicamente.
-   **Detección**: Identifica socios activos cuya fecha de `NextBilling` ha llegado.
-   **Generación de Deuda**: Calcula el monto a pagar y actualiza el saldo deudor (`OutstandingBalance`) del socio.
-   **Auditoría**: Registra cada evento de ciclo de facturación.

### 4. Historial y Estado
-   **Trazabilidad**: Registro histórico de cambios de plan (ej. Upgrade de Gold a Platinum).
-   **Estados**:
    -   `Active`: Socio al día.
    -   `Inactive`: Baja voluntaria.
    -   `Suspended`: Por falta de pago (bloquea acceso a reservas).
