# Manual de Usuario: Módulo de Membresías (Membership)

## 1. Propósito

Este módulo gestiona tu relación con el club. Te permite ver qué tipo de membresía tienes, cuál es su estado (si está activa o vencida) y explorar otros planes disponibles.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede ver su propia membresía y los planes que ofrece el club.
-   **Administrador (`ADMIN`):** Puede gestionar todos los planes de membresía y el estado de los socios.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Consultar tu Membresía

**Paso a paso:**
1.  **Inicia sesión** en la plataforma.
2.  **Navega a la sección "Mi Membresía"** o "Mi Perfil".
3.  Encontrarás una tarjeta o sección que muestra:
    -   **El nombre de tu plan actual** (ej: "Plan Familiar").
    -   **El estado de tu membresía:**
        -   `Activa`: Estás al día con tus pagos. ¡Puedes usar todos los servicios!
        -   `Pendiente de Pago`: Tienes una cuota pendiente.
        -   `Vencida`: No has pagado tu cuota. Es posible que tengas el acceso restringido a algunos servicios como las reservas.
    -   **Próxima fecha de facturación.**

### 🔹 Cómo Ver Otros Planes de Membresía

Si estás pensando en cambiar de plan, puedes explorar las opciones que ofrece el club.

**Paso a paso:**
1.  Busca una sección o página llamada **"Planes de Membresía"** o "Ver Planes".
2.  Se mostrará una lista con todos los planes disponibles, su precio y los beneficios que incluye cada uno.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Crear o Editar un Plan de Membresía

**Paso a paso:**
1.  **Accede al Panel de Administración.**
2.  Navega a "Configuración" -> **"Planes de Membresía"**.
3.  Para crear un nuevo plan, haz clic en **"Nuevo Plan"**. Para editar uno existente, búscalo en la lista y haz clic en "Editar".
4.  Completa el formulario con los detalles del plan: nombre, precio, ciclo de facturación (mensual, anual) y una descripción de sus beneficios.
5.  **Guarda los cambios.** El plan aparecerá inmediatamente como una opción para los socios.

### 🔸 Cómo Ver o Modificar la Membresía de un Socio

**Paso a paso:**
1.  Desde el Panel de Administración, **busca al socio** a través de la sección "Usuarios".
2.  En el perfil del socio, encontrarás la información de su membresía.
3.  Desde aquí, podrás **cambiar su plan** a otro existente o **modificar su estado** manualmente si es necesario (ej: de `Vencida` a `Activa` tras recibir un pago en persona).

---

## 5. Diagrama de Flujo: Consulta de Membresía (Socio)

```mermaid
graph TD
    A[Inicio] --> B[Ir a "Mi Membresía"];
    B --> C{¿Usuario Autenticado?};
    C -- Sí --> D[Muestra Plan Actual y Estado];
    D --> E{Estado: Activa?};
    E -- Sí --> F[Acceso Completo ✅];
    E -- No --> G[Acceso Restringido ⚠️];
    C -- No --> H[Redirige a Login];
```
