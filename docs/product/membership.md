# Módulo de Membresías (Membership)

Este módulo es el sistema central para la gestión de la relación con los socios, sus planes de suscripción y los ciclos de ingresos recurrentes del club.

## 🌟 Funcionalidades Implementadas

### 1. Gestión de Planes de Membresía (Tiers)
-   **Planes Configurables:** Permite a los administradores crear y gestionar diferentes tipos de membresía (ej: "Plan Individual", "Plan Familiar").
-   **Atributos del Plan:** Cada plan define su precio, ciclo de facturación (mensual, anual, etc.) y los beneficios o restricciones asociados.

### 2. Ciclo de Vida de la Membresía
-   **Gestión de Estado:** El sistema maneja el estado de la membresía de cada socio, que puede ser `Activa`, `Pendiente de Pago`, `Vencida` o `Cancelada`.
-   **Control de Acceso:** El estado de la membresía se utiliza para determinar el acceso a los servicios del club, como la creación de nuevas reservas.

### 3. Integración con Facturación
-   **Generación de Deuda:** Se integra con el módulo de Pagos para la generación automática de las facturas recurrentes de las cuotas de membresía.

## 4. Funcionalidades en Desarrollo

-   **Automatización de Cobros:** Aunque se genera la deuda, la automatización completa del proceso de cobro (ej: `ProcessMonthlyBilling`) es una funcionalidad del roadmap.
-   **Historial de Cambios:** Un log detallado de los cambios de plan (upgrades/downgrades) por socio es una mejora futura.
