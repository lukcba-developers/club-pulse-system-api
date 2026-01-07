# Protocolos de Seguridad

La seguridad es una prioridad en **Club Pulse**. Implementamos múltiples capas de protección para garantizar la integridad de los datos y la seguridad de los usuarios.

## 1. Autenticación y Autorización
- **Cookies HttpOnly**: Los tokens de acceso y refresco se almacenan en cookies `HttpOnly` y `Secure`, mitigando el riesgo de robo de tokens vía XSS.
- **RBAC (Control de Acceso Basado en Roles)**:
    - Middleware de autorización verifica roles (`SUPER_ADMIN`, `ADMIN`, `MEMBER`).
- **Protección BOLA (Broken Object Level Authorization)**: Validación estricta del `club_id` en cada solicitud para asegurar el aislamiento de tenants. Un middleware inyecta el `club_id` del usuario autenticado en el contexto de la petición, y la capa de repositorio lo utiliza en cada consulta a la base de datos, haciendo imposible por diseño que un usuario de un club acceda a datos de otro.

## 2. Protección de API
- **Headers OWASP**: Implementados vía middleware de seguridad (`SecurityHeadersMiddleware`).
    - `X-Content-Type-Options: nosniff`
    - `X-Frame-Options: DENY`
    - `Strict-Transport-Security` (HSTS).
- **CORS**: Configurado dinámicamente para orígenes permitidos con soporte de credenciales (Cookies).

## 3. Validación de Entradas
- DTOs (Data Transfer Objects) estrictos en la capa de entrada.
- Prevención de **SQL Injection** mediante consultas parametrizadas (GORM).

## 4. Limitación de Tasa (Rate Limiting)
- **Implementado**: Middleware de Rate Limiting respaldado por Redis para prevenir ataques de fuerza bruta y DDoS.
    - Límite global por defecto: 100 req/min.

## 5. Gestión de Secretos
- Las credenciales (DB Password, JWT Secret) se inyectan exclusivamente vía **Variables de Entorno**. NO se guardan en el código fuente (hardcoded).
N