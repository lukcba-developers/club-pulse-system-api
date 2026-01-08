# Módulo de Membresías (Membership)

Este módulo es el sistema central para la gestión de la relación con los socios, sus planes de suscripción y los ciclos de ingresos recurrentes del club.

## 🌟 Funcionalidades Implementadas

### 1. Gestión de Planes de Membresía (Tiers)
-   **Planes Configurables:** Permite a los administradores crear y gestionar diferentes tipos de membresía (ej: "Plan Individual", "Plan Familiar").
-   **Atributos del Plan:** Cada plan define su precio, ciclo de facturación (mensual, anual, etc.) y los beneficios o restricciones asociados.

### 2. Ciclo de Vida de la Membresía
-   **Gestión de Estado:** El sistema maneja el estado de la membresía de cada socio, que puede ser `Activa`, `Pendiente de Pago`, `Vencida` o `Cancelada`.
-   **Control de Acceso:** El estado de la membresía se utiliza para determinar el acceso a los servicios del club, como la creación de nuevas reservas.

### 3. Integración con Facturación y Becas
-   **Generación de Deuda:** Integrado con el módulo de Pagos para facturación automática.
-   **Becas (Scholarships):** Soporte para descuentos porcentuales sobre la cuota mensual con motivos y fechas de validez.
-   **Facturación Robusta:** Manejo coordinado de fechas de cierre (ej: 31 de enero -> 28 de febrero).

### 4. Automatización
-   **Proceso de Cobro Mensual:** Motor para procesar masivamente la deuda de socios activos (`ProcessMonthlyBilling`) aplicando descuentos de becas automáticamente.

## 5. Funcionalidades en Desarrollo
-   **Historial de Cambios:** Log detallado de cambios de plan (upgrades/downgrades).
