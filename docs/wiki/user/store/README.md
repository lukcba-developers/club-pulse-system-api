# Módulo de Tienda (Store)

Este módulo proporciona una funcionalidad de tienda o punto de venta dentro del club. Permite a los socios comprar productos como merchandising, artículos del buffet/cantina o equipamiento deportivo.

## Conceptos Clave

-   **Producto (`Product`):** Representa un artículo individual que se puede vender. Cada producto tiene un nombre, precio, categoría (ej: "Merch", "Buffet") y una cantidad de stock.
-   **Orden (`Order`):** Representa una compra realizada por un socio. Contiene la información del comprador, los productos y cantidades, y el monto total. El pago de la orden probablemente se gestiona a través del [Módulo de Pagos](../payment/README.md) o se debita del saldo en la [Billetera](../user/README.md) del usuario.

---

## Casos de Uso

### 1. Ver el Catálogo de Productos

Los socios pueden explorar los productos que están a la venta.

-   **Flujo**:
    1.  El usuario accede a la sección "Tienda" del sistema.
    2.  Se muestra una lista de todos los productos disponibles.
    3.  Opcionalmente, el usuario puede filtrar los productos por categoría (ej: ver solo la comida del "Buffet").
-   **Endpoint relacionado**: `GET /store/products`

### 2. Realizar una Compra

Un socio puede seleccionar uno o más productos y realizar una compra.

-   **Flujo**:
    1.  El usuario añade productos a su carrito.
    2.  Al confirmar la compra, el sistema crea una `Orden`.
    3.  El backend calcula el monto total, verifica el stock y lo descuenta.
    4.  Se procesa el pago de la orden.
-   **Endpoint relacionado**: `POST /store/purchase`
