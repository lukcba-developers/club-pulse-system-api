# Módulo de Billetera (Wallet)

## 1. Propósito

La Billetera es un componente personal de cada usuario que gestiona sus fondos y puntos dentro del club. Permite a los usuarios mantener un saldo para pagar servicios y acumular puntos que pueden ser canjeados por recompensas.

## 2. Funcionalidades

-   **Saldo Monetario:** Muestra el balance de dinero real (ej: en EUR, USD) que el usuario ha cargado en su cuenta del club. Este saldo puede ser utilizado para pagar reservas, inscripciones a torneos o productos de la tienda.
-   **Puntos de Recompensa:** Muestra la cantidad de puntos que el usuario ha ganado a través de actividades de gamificación, como participar en partidos o alcanzar logros. Estos puntos suelen ser canjeables por descuentos u otros beneficios.

## 3. Modelo de Datos

| Campo     | Tipo      | Descripción                               |
| --------- | --------- | ----------------------------------------- |
| `ID`      | `string`  | Identificador único de la billetera.      |
| `Balance` | `number`  | Saldo monetario disponible.               |
| `Points`  | `number`  | Cantidad de puntos de recompensa acumulados. |

## 4. Endpoint de la API

### `GET /users/:id/wallet`

-   **Acción:** Obtiene el estado de la billetera de un usuario, incluyendo su saldo y puntos.
-   **Permisos:**
    -   Un usuario puede consultar su propia billetera usando el alias `/users/me/wallet`.
    -   Un `ADMIN` o `SUPER_ADMIN` puede consultar la billetera de cualquier usuario.
-   **Respuesta Exitosa (200 OK):** Un objeto `Wallet`.

```json
{
  "data": {
    "id": "wallet-uuid-string",
    "balance": 15.50,
    "points": 1200
  }
}
```
