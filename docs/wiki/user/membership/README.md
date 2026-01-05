# Módulo de Membresías (Membership)

El módulo de Membresías es el núcleo de la relación entre el socio y el club. Gestiona los diferentes planes de membresía, la suscripción de los socios a dichos planes y el estado de su afiliación.

## Casos de Uso para Socios (Members)

### 1. Ver Planes de Membresía Disponibles

Los usuarios, tanto existentes como potenciales, pueden ver los diferentes niveles o planes de membresía que ofrece el club. Esto les permite comparar precios, beneficios y elegir el que mejor se adapte a sus necesidades.

-   **Flujo**:
    1.  El usuario navega a la sección de "Planes" o "Hacerse Socio".
    2.  El sistema muestra una lista de todos los planes de membresía (`Tiers`) disponibles, con sus descripciones y precios.
-   **Endpoint relacionado**: `GET /memberships/tiers`

### 2. Adquirir o Cambiar de Membresía

Un usuario puede suscribirse a un plan de membresía. Esto formaliza su relación como socio activo del club.

-   **Flujo**:
    1.  El usuario selecciona un plan de la lista de `tiers`.
    2.  El sistema crea una nueva membresía que asocia al usuario con el plan elegido.
-   **Endpoint relacionado**: `POST /memberships`

### 3. Consultar Mis Membresías

Un socio puede revisar el estado actual de sus membresías activas.

-   **Flujo**:
    1.  El usuario accede a su perfil o a la sección "Mi Membresía".
    2.  El sistema muestra una lista de las membresías asociadas a su cuenta, incluyendo el nombre del plan y el estado actual (ej: "Activo", "Vencido").
-   **Endpoint relacionado**: `GET /memberships`

-   **Ver Detalle de una Membresía**: Para obtener información más detallada sobre una membresía específica.
-   **Endpoint relacionado**: `GET /memberships/:id`

## Casos de Uso para Administradores (Admins)

### 1. Procesar Facturación Mensual

Los administradores tienen la capacidad de ejecutar el ciclo de facturación para todos los miembros del club.

-   **Flujo**:
    1.  Un administrador (o un proceso automático) ejecuta la acción de procesar la facturación.
    2.  El sistema recorre todas las membresías activas, genera los cargos correspondientes en la billetera de cada socio y actualiza su estado de cuenta.
-   **Endpoint relacionado**: `POST /memberships/process-billing`
-   **Nota**: Esta es una operación administrativa crítica que afecta a todos los socios.

---
## Funcionalidad Interna (Sin API)

### Becas (Scholarships)

El sistema incluye una funcionalidad interna para gestionar becas que otorgan descuentos a los socios.

-   **Concepto**: Una beca (`Scholarship`) es un descuento (en porcentaje) que se le asigna a un usuario por un motivo específico (ej: "Beca por rendimiento deportivo"). Puede tener una fecha de vencimiento.
-   **Aplicación**: La lógica para aplicar este descuento (`ApplyDiscount`) existe en el backend. Se utiliza, por ejemplo, durante el proceso de facturación mensual (`ProcessBilling`) para calcular el monto final que se le debe cobrar a un socio becado.
-   **Gestión**: Actualmente, **no existen endpoints en la API** para que los usuarios o administradores creen, vean o gestionen las becas directamente. Esta es una capacidad del sistema que probablemente se gestiona de forma interna o directamente en la base de datos por un `SUPER_ADMIN`.
