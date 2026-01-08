# Manual de Usuario: Módulo de Tienda (Store)

## 1. Propósito

Este módulo funciona como la tienda del club, donde puedes comprar productos como merchandising, equipamiento deportivo o snacks del buffet.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede explorar el catálogo y comprar productos.
-   **Administrador (`ADMIN`):** Gestiona el catálogo de productos, el inventario y ve los reportes de ventas.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Comprar en la Tienda

**Paso a paso:**
1.  **Navega a la sección "Tienda"** en la aplicación.
2.  **Explora el catálogo de productos.** Puedes filtrarlos por categoría (ej: "Ropa", "Bebidas").
3.  **Haz clic en un producto** para ver sus detalles, como la descripción, el precio y las tallas disponibles.
4.  **Añade los productos** que deseas a tu carrito de compras.
5.  **Procede al pago.** Haz clic en el ícono del carrito y luego en "Finalizar Compra".
6.  Serás dirigido al **Módulo de Pagos** para completar la transacción. También podrías tener la opción de pagar con el saldo de tu billetera virtual.
7.  Una vez confirmado el pago, recibirás una confirmación y podrás retirar tu producto en el club.

---

## 4. Compra para Invitados (Público General)

El sistema también permite que personas que no son socias del club realicen compras a través de un catálogo público.

### 🔹 Cómo Comprar como Invitado

**Paso a paso:**
1.  **Accede al catálogo público** del club a través de su página web o un enlace directo.
2.  **Explora los productos** disponibles para la venta al público.
3.  **Añade los productos** a tu carrito.
4.  Al finalizar la compra, se te pedirá que **proporciones tus datos de contacto** (nombre y correo electrónico) para poder procesar el pedido.
5.  Serás dirigido a la pasarela de pagos (ej. Mercado Pago) para completar la transacción.
6.  Recibirás un correo electrónico de confirmación con los detalles de tu pedido y las instrucciones para retirarlo.

---

## 5. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Añadir un Nuevo Producto

**Paso a paso:**
1.  **Accede al Panel de Administración** y ve a la sección de **"Tienda"**.
2.  Haz clic en **"Añadir Producto"**.
3.  **Completa el formulario:**
    -   Nombre del producto.
    -   Descripción.
    -   Precio.
    -   Categoría.
    -   **Stock inicial:** La cantidad de unidades disponibles.
    -   Sube una o más fotos del producto.
4.  **Guarda los cambios.** El producto estará visible inmediatamente en la tienda para los socios.

### 🔸 Cómo Gestionar el Inventario

**Paso a paso:**
1.  En el panel de la "Tienda", busca el producto cuyo stock deseas ajustar.
2.  Haz clic en **"Editar"** o en una opción específica de "Gestionar Stock".
3.  **Actualiza el número de unidades disponibles.** El sistema también descontará el stock automáticamente con cada venta.
4.  Puedes configurar alertas para que el sistema te notifique cuando el stock de un producto esté bajo.

---

## 6. Endpoints de la API y Cambios

### Endpoints para Usuarios Autenticados

-   `POST /store/purchase`: Realiza una compra para el usuario autenticado.
-   `GET /store/products`: Obtiene el catálogo de productos para el club del usuario.

### Endpoints Públicos

-   `POST /public/clubs/:slug/store/purchase`: **(Nuevo)** Permite a un invitado realizar una compra. El `user_id` en la orden es nulo y se guardan `guest_name` y `guest_email`.
-   `GET /public/clubs/:slug/store/products`: Obtiene el catálogo público de productos de un club.

---

## 7. Diagrama de Flujo: Compra de un Producto

```mermaid
graph TD
    subgraph Flujo General
        A[Inicio] --> B[Accede a la Tienda];
        B --> C[Explorar Productos];
        C --> D[Añadir al Carrito];
        D --> E{¿Seguir Comprando?};
        E -- Sí --> C;
        E -- No --> F[Finalizar Compra];
    end

    subgraph Proceso de Pago
        F --> G{¿Usuario es Socio?};
        G -- Sí --> H[Pagar como Socio (Billetera/MP)];
        G -- No --> I[Ingresar Nombre y Email];
        I --> J[Pagar como Invitado (MP)];
    end

    subgraph Finalización
      H --> K[Compra de Socio Exitosa ✅];
      J --> L[Compra de Invitado Exitosa ✅];
    end
```
