# Especificaciones Técnicas y Operaciones

La maquinaria invisible que impulsa Club Pulse, diseñada para estabilidad, velocidad y escala.

## 🌟 Capacidades Técnicas

### 1. Observabilidad (OpenTelemetry)
Visibilidad de rayos-X sobre el sistema.
-   **Trazabilidad Distribuida**: Cada request recibe un `TraceID` único (W3C Standard) al ingresar al Load Balancer, permitiendo seguir su viaje por todos los microservicios y bases de datos.
-   **Logs Estructurados**: Salida JSON (`slog`) correlacionada con trazas para depuración instantánea.

### 2. Rendimiento y Caching (Redis)
Uso estratégico de memoria in-memory.
-   **Rate Limiting**: Protección capa de aplicación contra ataques de fuerza bruta (ej. login) y DDoS.
    -   *Policy*: Token Bucket algorithm (100 req/min default).
-   **Session Store**: Gestión de millones de sesiones activas con latencia sub-milisegundo.
-   **Cache (Roadmap)**: Caché de capa de aplicación para endpoints de alta lectura (`/availability`).

### 3. Base de Datos (PostgreSQL)
-   **Motor Relacional**: Integridad referencial fuerte para transacciones financieras y de reservas.
-   **Extensiones**: `pgvector` habilitado para búsqueda de similitud n-dimensional (IA/Semántica).

### 4. Arquitectura de Despliegue
-   **Docker Native**: Contenedores inmutables.
-   **Multi-Stage Builds**: Imágenes optimizadas (<50MB) para despliegue rápido.
-   **Health Checks**: Endpoints `/health` y `/healthz` para orquestación automática (Kubernetes/ECS).
