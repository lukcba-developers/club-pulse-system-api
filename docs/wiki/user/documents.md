# Módulo de Gestión de Documentos de Usuario

El módulo de `Documentos de Usuario` es una parte crítica del sistema para garantizar que todos los miembros del club tengan su documentación personal y deportiva al día. Esto es fundamental para la elegibilidad en competiciones, la seguridad y el cumplimiento normativo.

## 1. Propósito y Funcionalidad

Este módulo permite a los usuarios subir documentos requeridos por el club y a los administradores validarlos. Las funcionalidades clave son:

-   **Subida de Documentos:** Los usuarios pueden subir archivos (imágenes, PDF) para diferentes tipos de documentos.
-   **Validación por Administradores:** El personal autorizado puede revisar los documentos, aprobarlos o rechazarlos.
-   **Control de Vencimientos:** El sistema maneja fechas de expiración para documentos como certificados médicos o seguros.
-   **Notificaciones Automáticas:** Un proceso en segundo plano notifica a los usuarios sobre documentos próximos a vencer o ya vencidos.
-   **Consulta de Elegibilidad:** Permite verificar si un usuario cumple con todos los requisitos de documentación para participar en actividades.

## 2. Modelo de Datos

La información se almacena en la tabla `user_documents`.

### Campos Principales

| Campo            | Tipo        | Descripción                                                                  |
| ---------------- | ----------- | ---------------------------------------------------------------------------- |
| `id`             | `UUID`      | Identificador único del documento.                                           |
| `club_id`        | `VARCHAR`   | ID del club al que pertenece el contexto del documento.                      |
| `user_id`        | `VARCHAR`   | ID del usuario que subió el documento.                                       |
| `type`           | `ENUM`      | Tipo de documento. Ver `DocumentType`.                                       |
| `file_url`       | `VARCHAR`   | URL donde está almacenado el archivo (ej. en un bucket S3/MinIO).            |
| `status`         | `ENUM`      | Estado actual del documento. Ver `DocumentStatus`.                           |
| `expiration_date`| `TIMESTAMPTZ`| Fecha y hora en que el documento expira (opcional).                          |
| `rejection_notes`| `TEXT`      | Notas del administrador en caso de que el documento sea rechazado (opcional).|
| `uploaded_at`    | `TIMESTAMPTZ`| Fecha y hora de la subida.                                                   |
| `validated_at`   | `TIMESTAMPTZ`| Fecha y hora de la validación (opcional).                                    |
| `validated_by`   | `VARCHAR`   | ID del administrador que realizó la validación (opcional).                   |

### `DocumentType` (Tipos de Documento)

-   `DNI_FRONT`: Foto del anverso del Documento Nacional de Identidad.
-   `DNI_BACK`: Foto del reverso del Documento Nacional de Identidad.
-   `EMMAC_MEDICAL`: Certificado Médico de Aptitud (EMMAC).
-   `LEAGUE_FORM`: Formulario de inscripción a la liga/federación.
-   `INSURANCE`: Póliza de seguro de accidentes.

### `DocumentStatus` (Estados del Documento)

-   `PENDING`: El documento ha sido subido y está pendiente de revisión.
-   `VALID`: El documento ha sido aprobado por un administrador y está vigente.
-   `REJECTED`: El documento ha sido rechazado. El usuario debe subir uno nuevo.
-   `EXPIRED`: La fecha de vencimiento del documento ha pasado.

## 3. Endpoints de la API

Las siguientes rutas están disponibles para interactuar con el módulo. Todas requieren autenticación y operan dentro del contexto de un `club_id`.

---

### `POST /users/:userId/documents`

-   **Acción:** Sube un nuevo documento para un usuario.
-   **Permisos:**
    -   Un usuario puede subir sus propios documentos.
    -   Un `ADMIN` o `SUPER_ADMIN` puede subir documentos para cualquier usuario.
-   **Request Body (multipart/form-data):**
    -   `file`: El archivo a subir.
    -   `type`: (string) Uno de los `DocumentType` válidos.
    -   `expiration_date`: (string, opcional) Fecha en formato `YYYY-MM-DD`.
-   **Respuesta Exitosa (201 Created):** El objeto del documento creado.

---

### `GET /users/:userId/documents`

-   **Acción:** Lista todos los documentos de un usuario específico.
-   **Permisos:** Abierto a usuarios autenticados que consultan sus propios datos o a administradores.
-   **Respuesta Exitosa (200 OK):** Un array de objetos de documentos.

---

### `GET /users/:userId/documents/summary`

-   **Acción:** Devuelve un resumen del estado de la documentación de un usuario, indicando qué documentos faltan o están inválidos.
-   **Permisos:** Similar a la lista de documentos.
-   **Respuesta Exitosa (200 OK):** Un objeto que detalla el estado de cada tipo de documento requerido.

---

### `GET /users/:userId/documents/:docId`

-   **Acción:** Obtiene los detalles de un documento específico por su ID.
-   **Permisos:** Similar a la lista de documentos.
-   **Respuesta Exitosa (200 OK):** El objeto del documento solicitado.

---

### `PUT /users/:userId/documents/:docId/validate`

-   **Acción:** Aprueba o rechaza un documento pendiente.
-   **Permisos:** `ADMIN` o `SUPER_ADMIN`.
-   **Request Body (JSON):**
    ```json
    {
      "approve": true, // o false para rechazar
      "notes": "El número de DNI no es legible." // Opcional, requerido si se rechaza
    }
    ```
-   **Respuesta Exitosa (200 OK):** Un mensaje de confirmación y el objeto del documento actualizado.

---

### `DELETE /users/:userId/documents/:docId`

-   **Acción:** Elimina un documento.
-   **Permisos:**
    -   Un usuario puede eliminar sus propios documentos.
    -   Un `ADMIN` o `SUPER_ADMIN` puede eliminar cualquier documento.
-   **Respuesta Exitosa (200 OK):** Un mensaje de confirmación.

---

### `GET /users/:userId/eligibility`

-   **Acción:** Realiza una verificación completa de la elegibilidad de un usuario, que incluye el estado de sus documentos.
-   **Permisos:** `ADMIN` o `SUPER_ADMIN`.
-   **Respuesta Exitosa (200 OK):** Un objeto con el resultado de la verificación.

## 4. Proceso Automatizado (Job de Vencimiento)

Un proceso automatizado (`DocumentExpirationJob`) se ejecuta periódicamente (ej. una vez al día) para gestionar los vencimientos.

-   **Verificación:** El job busca en la base de datos:
    1.  Documentos que están a punto de vencer (ej. en 30 y 7 días).
    2.  Documentos que ya han vencido.
-   **Acciones:**
    -   **Notificación:** Envía notificaciones por email (u otro canal) a los usuarios cuyos documentos están próximos a vencer.
    -   **Actualización de Estado:** Cambia el estado de los documentos vencidos a `EXPIRED`.
    -   **Notificación de Vencimiento:** Informa a los usuarios que su documento ha vencido y que deben subir uno nuevo.
