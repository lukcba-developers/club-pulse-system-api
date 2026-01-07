# Módulo de Membresías (Membership)

## 1. Propósito

El módulo de **Membresías** es el núcleo de la relación entre el club y sus socios. Se encarga de gestionar los diferentes tipos de planes, el estado de los socios y los ciclos de facturación.

## 2. Funcionalidades Principales

-   **Gestión de Planes de Membresía:** Permite a los administradores crear y configurar diferentes planes o tipos de membresía (ej: "Plan Individual", "Plan Familiar", "Plan Fin de Semana"). Para cada plan se puede definir:
    -   Precio y ciclo de facturación (mensual, anual).
    -   Beneficios y restricciones (ej: acceso a ciertas instalaciones).
    -   Matrícula de inscripción.
-   **Gestión de Socios:** Mantiene un registro de todos los socios del club, asociando a cada uno con su plan de membresía correspondiente.
-   **Ciclo de Vida de la Membresía:** Gestiona el estado de la membresía de un socio:
    -   `Activa`: El socio está al día con sus pagos.
    -   `Pendiente de Pago`: Se ha generado una factura pero aún no ha sido pagada.
    -   `Vencida`: El socio no ha pagado y se le restringen los beneficios.
    -   `Cancelada`: La membresía ha sido dada de baja.
-   **Integración con Pagos:** Se integra estrechamente con el **Módulo de Pagos** para la generación automática de facturas recurrentes y el procesamiento de los pagos.