# Manual de Usuario: Módulo de Pagos (Payment)

## 1. Propósito

Este módulo te permite gestionar todas tus transacciones financieras con el club. Desde aquí puedes pagar tus cuotas de membresía, ver tu historial de pagos y asegurarte de que tu cuenta esté siempre al día.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Realiza pagos y consulta su historial.
-   **Administrador (`ADMIN`):** Supervisa todas las transacciones del club.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Pagar tu Cuota de Membresía

Cuando se genera una nueva cuota de membresía, recibirás una notificación y podrás pagarla a través de la plataforma.

**Paso a paso:**
1.  **Inicia sesión** en la plataforma.
2.  En tu panel principal o en la sección "Mi Membresía", verás un aviso si tienes una **factura pendiente**.
3.  Haz clic en el botón **"Pagar Ahora"**.
4.  Serás redirigido de forma segura a la **pasarela de pagos** (ej: Mercado Pago) para completar la transacción con tu tarjeta de crédito, débito u otros métodos disponibles.
5.  Una vez completado el pago, serás redirigido de vuelta al sitio del club.
6.  Tu estado de membresía se actualizará automáticamente a `Activa`.

### 🔹 Cómo Ver tu Historial de Pagos

Puedes consultar un registro de todas las transacciones que has realizado.

**Paso a paso:**
1.  Navega a la sección **"Mi Perfil"** o **"Pagos"**.
2.  Busca la pestaña o sección de **"Historial de Transacciones"**.
3.  Verás una lista de todos tus pagos, con la fecha, el concepto (ej: "Cuota Mensual") y el monto de cada uno.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Supervisar las Transacciones del Club

**Paso a paso:**
1.  **Accede al Panel de Administración.**
2.  Navega a la sección de **"Finanzas"** o **"Pagos"**.
3.  Verás un **panel con todas las transacciones** realizadas en el club.
4.  Puedes usar los **filtros** para buscar pagos por:
    -   Socio específico.
    -   Rango de fechas.
    -   Estado del pago (`Aprobado`, `Rechazado`, `Pendiente`).
5.  Esto te permite tener un control total sobre los ingresos y el estado financiero de cada miembro.

---

## 5. Diagrama de Flujo: Proceso de Pago (Socio)

```mermaid
graph TD
    A[Inicio: Usuario con cuota pendiente] --> B[Clic en "Pagar Ahora"];
    B --> C[Redirección a Pasarela de Pago Externa];
    C --> D[Usuario completa el pago en la pasarela];
    D --> E{¿Pago Exitoso?};
    E -- Sí --> F[Pasarela envía Webhook de confirmación];
    F --> G[El sistema actualiza el estado a "Pagado"];
    G --> H[Membresía del socio se marca como "Activa" ✅];
    E -- No --> I[Pasarela informa del fallo];
    I --> J[El sistema mantiene el estado como "Pendiente"];
```
