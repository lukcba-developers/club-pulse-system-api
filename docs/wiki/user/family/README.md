# Módulo de Grupos Familiares (Family Groups)

## 1. Propósito

El módulo de Grupos Familiares permite a un usuario (`Head of Family`) agrupar y gestionar las cuentas de otros miembros del club (generalmente, sus hijos o dependientes) bajo una única entidad. Esto centraliza la administración, las membresías y la facturación.

## 2. Funcionalidades Principales

-   **Creación de un Grupo Familiar:** Un usuario puede crear un grupo y convertirse en el "Cabeza de Familia".
-   **Gestión de Miembros:** El cabeza de familia puede agregar o eliminar miembros de su grupo.
-   **Visibilidad Centralizada:** Permite al cabeza de familia ver información relevante de los miembros del grupo, como sus próximas clases o el estado de su documentación.
-   **Facturación Unificada (Futuro):** La intención es que los pagos de membresías y actividades de todos los miembros se puedan consolidar en una única factura para el cabeza de familia.

## 3. Modelo de Datos

La información se gestiona a través de la entidad `FamilyGroup`.

| Campo         | Tipo      | Descripción                                                              |
| ------------- | --------- | ------------------------------------------------------------------------ |
| `ID`          | `UUID`    | Identificador único del grupo familiar.                                  |
| `ClubID`      | `string`  | ID del club donde existe el grupo.                                       |
| `Name`        | `string`  | Nombre del grupo familiar (ej. "Familia Pérez").                         |
| `HeadUserID`  | `string`  | ID del usuario que es el cabeza de familia y administrador del grupo.    |
| `Members`     | `[]User`  | Lista de los objetos de usuario que son miembros del grupo.              |

## 4. Flujo de Uso

### 🔹 Crear un Grupo Familiar

1.  Un usuario navega a la sección "Mi Familia" en su perfil.
2.  Hace clic en "Crear Grupo Familiar".
3.  Asigna un nombre al grupo (ej. "Familia García").
4.  El sistema crea el grupo y asigna al usuario actual como `HeadUserID`.

### 🔹 Agregar un Miembro

1.  El cabeza de familia busca a un socio existente en el club (que no pertenezca ya a otro grupo).
2.  Envía una invitación o lo agrega directamente (dependiendo de la configuración del club).
3.  Una vez aceptado/agregado, el `FamilyGroupID` del miembro se actualiza para vincularlo al grupo.

### 🔹 Eliminar un Miembro

1.  El cabeza de familia selecciona a un miembro de su lista de grupo.
2.  Hace clic en "Eliminar del grupo".
3.  El `FamilyGroupID` del miembro se establece en `null`, desvinculándolo del grupo.

## 5. Endpoints de la API

La gestión de Grupos Familiares se realiza a través de los siguientes endpoints:

-   `POST /users/family-groups`: Crea un nuevo grupo familiar. El usuario que lo crea se convierte en el `HeadUserID`.
-   `GET /users/family-groups/me`: Obtiene los detalles del grupo familiar al que pertenece o que administra el usuario autenticado.
-   `POST /users/family-groups/:id/members`: Agrega un nuevo miembro al grupo familiar especificado por `:id`. Solo el `HeadUserID` puede realizar esta acción.
