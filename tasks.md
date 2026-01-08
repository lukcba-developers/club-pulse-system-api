# Lista de Tareas y Deuda Técnica

Este archivo centraliza las tareas pendientes, mejoras y deuda técnica identificadas en el proyecto.

##  приоритет: Crítico 🟥

-   [ ] **Corregir Fuga de Datos en Módulo de Autenticación (Auth)**
    -   **Problema:** Las funciones `FindUserByEmail` y `FindUserByID` en `backend/internal/modules/auth/infrastructure/repository/postgres.go` aceptan el `club_id`, pero no filtran las querys por este campo (`r.db.Where("email = ?", email)`), permitiendo potencialmente que un usuario de otro club sea accedido si se conoce su email o ID.
    -   **Solución:**
        1.  Actualizar la firma de `FindUserByEmail` y `FindUserByID` para aceptar `club_id` (o extraerlo del contexto).
        2.  Añadir `.Where("club_id = ?", clubID)` a las consultas GORM.
        3.  Verificar que el login solo permita acceso si el usuario pertenece al club del dominio/contexto actual.

-   [ ] **Automatizar Proceso de Facturación (Membership)**
    -   **Problema:** El proceso de facturación mensual se dispara con una llamada manual a la API.
    -   **Solución:** Implementar un CRON Job que ejecute el proceso automáticamente todos los días.

-   [ ] **Validar Firma de Webhooks de Pago (Payment)**
    -   **Problema:** El endpoint `HandleWebhook` en `backend/internal/modules/payment/infrastructure/http/handler.go` procesa notificaciones confiando ciegamente en los parámetros `type` y `data.id` sin validar que la petición provenga realmente de Mercado Pago.
    -   **Solución:** Implementar validación de firma (`x-signature` o `x-request-id`) comparando con el secreto del webhook configurado en el Dashboard de MP. Retornar `403` si la firma es inválida.

-   [ ] **Refactorizar Lógica de Precios de Invitados (Booking)**
    -   **Problema:** La tarifa por invitado se calcula en el frontend.
    -   **Solución:** Mover la lógica de precios al backend.

## приоритет: Medio 🟧

-   [ ] **Forzar Filtro Multi-Tenant en Repositorios (Varios)**
    -   **Problema:** Varios repositorios (`Championship`, `Facilities`) tienen funciones que buscan registros solo por `id`, sin filtrar por `club_id`, delegando la seguridad a la capa de servicio.
    -   **Solución:** Refactorizar todas las funciones `Get...ByID` para que siempre requieran y apliquen el filtro `club_id`, añadiendo una capa de defensa en profundidad.

-   [ ] **Refactorizar Lógica de Webhook (Payment)**
    -   **Problema:** La lógica de negocio del webhook está en el `handler` en lugar del `service`.
    -   **Solución:** Mover la lógica a la capa de aplicación para seguir el patrón de Clean Architecture.

-   [ ] **Mejorar Manejo de Errores en Webhook (Payment)**
    -   **Problema:** El endpoint responde con `200 OK` (`c.Status(http.StatusOK)`) al final del `func` sin importar si el procesamiento (`h.gateway.ProcessWebhook` o `h.repo.Update`) falló. Esto impide que Mercado Pago reintente la notificación.
    -   **Solución:**
        1.  Si `ProcessWebhook` falla, retornar `500` o `502`.
        2.  Si `repo.Update` falla, retornar `500`.
        3.  Solo retornar `200` si la actualización fue exitosa.

-   [ ] **Mejorar Reporte de Errores en Facturación (Membership)**
    -   **Problema:** No hay un reporte consolidado de errores durante el proceso de facturación en lote.
    -   **Solución:** Generar un resumen de los socios que no pudieron ser procesados.

-   [ ] **Centralizar Configuración de Horarios (Facilities/Booking)**
    -   **Problema:** Las horas de operación de las instalaciones están hardcodeadas.
    -   **Solución:** Añadir campos de configuración de horarios al modelo de `Facility`.

-   [ ] **Implementar Creación de Pagos Manuales (Payment)**
    -   **Problema:** La funcionalidad para registrar pagos offline no está implementada.
    -   **Solución:** Desarrollar el caso de uso correspondiente.

-   [ ] **Implementar Reservas Recurrentes (Booking)**
    -   **Problema:** La funcionalidad de reservas recurrentes no está implementada.
    -   **Solución:** Desarrollar los casos de uso correspondientes.

## приоритет: Bajo 🟦

-   [ ] **Robustecer Cálculo de Fechas de Facturación (Membership)**
    -   **Problema:** El cálculo de la siguiente fecha de facturación puede ser impreciso con los fines de mes.
    -   **Solución:** Usar una librería de manejo de fechas más robusta.

-   [ ] **Implementar Funcionalidad de Becas (Membership)**
    -   **Problema:** La funcionalidad `AssignScholarship` no está implementada.
    -   **Solución:** Desarrollar el caso de uso correspondiente.

-   [ ] **Mejorar la Gestión de Secretos de API**
    -   **Problema:** Verificar que las claves de API externas se gestionan de forma segura.
    -   **Solución:** Asegurar que todas las claves se inyectan a través de variables de entorno.

-   [ ] **Implementar Grupos Familiares (User)**
    -   **Problema:** La funcionalidad de grupos familiares no está implementada.
    -   **Solución:** Desarrollar la lógica correspondiente.

## Auditoría y Pruebas Generales

-   [x] **Auditar Implementación Multi-Tenant**
    -   **Objetivo:** Verificar que todas las consultas a la base de datos en todos los módulos incluyan un filtro por `club_id` para garantizar el aislamiento de datos entre clubes.
    -   **Riesgo si no se hace:** Potencial fuga de datos entre clientes (vulnerabilidad crítica).
    -   **Resultado:** Auditoría completada. Se encontraron varias fugas potenciales. Tareas de corrección creadas.

-   [ ] **Crear Suite de Tests End-to-End (Multi-Tenant)**
    -   **Objetivo:** Desarrollar pruebas de integración que simulen flujos de usuario reales a través de múltiples módulos, validando la lógica de negocio y el aislamiento multi-tenant.
    -   **Ejemplo de Flujo:** Crear 2 clubes y 2 usuarios, y verificar que el usuario de un club no puede ver ni operar con datos del otro.
