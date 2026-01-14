# Módulo de Derechos del Usuario (GDPR)

## 1. Propósito

Este módulo proporciona a los usuarios las herramientas para ejercer sus derechos de protección de datos, en cumplimiento con el Reglamento General de Protección de Datos (GDPR) de la Unión Europea. Específicamente, implementa el derecho a la portabilidad de los datos (Artículo 20) y el derecho al olvido (Artículo 17).

## 2. Derecho a la Portabilidad de Datos (Artículo 20)

Esta funcionalidad permite a los usuarios solicitar y descargar una copia de todos sus datos personales que el club ha almacenado en el sistema.

### Endpoint de la API

#### `GET /users/me/data-export`

-   **Acción:** Inicia una exportación de todos los datos asociados con el usuario autenticado.
-   **Permisos:** Solo el propio usuario puede solicitar su exportación de datos.
-   **Respuesta Exitosa (200 OK):**
    -   El sistema devuelve un archivo `my_data_export.json`.
    -   La respuesta incluye las cabeceras `Content-Disposition: attachment; filename=my_data_export.json` para forzar la descarga en el navegador.
    -   El JSON contiene una recopilación completa de la información del usuario, como su perfil, membresías, reservas, historial de asistencia, etc.

---

## 3. Derecho al Olvido (Artículo 17)

Esta funcionalidad permite a un usuario solicitar la eliminación de su cuenta y de sus datos personales identificables. Para mantener la integridad de los registros históricos (ej: resultados de partidos pasados), el sistema no borra las filas de la base de datos, sino que las **anonimiza**.

### Proceso de Anonimización

-   Los datos personales como nombre, email, dirección, teléfono, etc., son reemplazados por valores genéricos o aleatorios (ej: "Usuario Anónimo", "deleted@user.com").
-   El identificador único del usuario (`ID`) se mantiene para preservar las relaciones en la base de datos.
-   La cuenta del usuario queda desactivada y ya no puede iniciar sesión.

### Endpoint de la API

#### `DELETE /users/me/gdpr-erasure`

-   **Acción:** Inicia el proceso de anonimización para la cuenta del usuario autenticado.
-   **Permisos:** Solo el propio usuario puede solicitar la eliminación de su cuenta.
-   **Respuesta Exitosa (200 OK):** Un mensaje de confirmación indicando que la cuenta ha sido desactivada y los datos anonimizados.

```json
{
  "message": "Your data has been anonymized. Your account is now deactivated.",
  "details": "Personal identifying information has been removed in compliance with GDPR Article 17."
}
```
