# Manual de Usuario: Módulo de Notificaciones

## 1. Propósito

Este módulo es el servicio de mensajería del club. Funciona de manera automática para mantenerte informado sobre todas las actividades importantes relacionadas con tu cuenta, enviando comunicaciones a tu correo electrónico.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Recibe las notificaciones.
-   **Administrador (`ADMIN`):** Puede gestionar las plantillas de los correos.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 ¿Qué Notificaciones Recibirás?

No necesitas hacer nada para que este módulo funcione, simplemente te mantendrá al día. Recibirás correos electrónicos automáticos para eventos como:

-   **Confirmación de Reserva:** Inmediatamente después de que reserves una instalación.
-   **Recordatorio de Reserva:** Un tiempo antes de tu reserva (ej: 24 horas antes).
-   **Promoción de Lista de Espera:** Cuando un lugar se libera y se te asigna automáticamente.
-   **Nueva Factura:** Cuando se genera la cuota mensual de tu membresía.
-   **Confirmación de Pago:** Tan pronto como tu pago sea procesado con éxito.
-   **Restablecimiento de Contraseña:** Cuando lo solicites desde la página de inicio de sesión.
-   **Anuncios del Club:** Noticias importantes publicadas por los administradores.

### 🔹 Cómo Gestionar tus Preferencias

En la sección **"Mi Perfil"**, podrás encontrar (o se añadirá en el futuro) una sección de **"Preferencias de Comunicación"** donde podrás elegir qué notificaciones deseas recibir.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Personalizar las Plantillas de Correo

Para mantener una imagen de marca consistente, puedes editar el contenido de los correos que reciben los socios.

**Paso a paso:**
1.  **Accede al Panel de Administración.**
2.  Navega a "Configuración" -> **"Plantillas de Notificaciones"**.
3.  Verás una lista de todas las plantillas de correo electrónico (ej: "Confirmación de Reserva", "Recordatorio de Pago").
4.  **Selecciona la plantilla** que deseas editar.
5.  Se abrirá un editor donde podrás **modificar el texto y el asunto** del correo, utilizando variables como `{{nombre_usuario}}` o `{{fecha_reserva}}` que el sistema reemplazará con los datos correctos.
6.  **Guarda los cambios.** A partir de ese momento, todos los correos de ese tipo usarán la nueva plantilla.

---

## 5. Diagrama de Flujo: Notificación de Reserva

```mermaid
graph TD
    A[Socio crea una reserva] --> B[Módulo de Reservas confirma la reserva en la BD];
    B --> C[Módulo de Reservas le pide al Módulo de Notificaciones que envíe un email];
    C --> D[Módulo de Notificaciones busca la plantilla "Confirmación de Reserva"];
    D --> E[Rellena la plantilla con los datos de la reserva];
    E --> F[Envía el email al correo del socio];
    F --> G[Socio recibe la confirmación ✅];
```
