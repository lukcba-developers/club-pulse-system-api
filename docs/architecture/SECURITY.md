# Protocolos de Seguridad

La seguridad es una prioridad en **Club Pulse**. Implementamos múltiples capas de protección para garantizar la integridad de los datos y la seguridad de los usuarios.

## 1. Autenticación y Autorización
- **Cookies HttpOnly**: Los tokens de acceso y refresco se almacenan en cookies `HttpOnly` y `Secure` con protección `SameSite=Strict`. No son accesibles vía JavaScript, lo que mitiga drásticamente ataques XSS de robo de sesión.
- **RBAC (Control de Acceso Basado en Roles)**: El sistema implementa roles como `SUPER_ADMIN`, `ADMIN`, `COACH` y `MEMBER`.
- **Protección BOLA (Broken Object Level Authorization)**: El aislamiento de datos se garantiza por diseño forzando que cada método de la capa de `Repository` requiera un `club_id`. El `TenantMiddleware` extrae este ID de la sesión o del host y lo propaga a través de todas las capas, asegurando que un usuario nunca pueda consultar datos que no pertenezcan a su club.

## 2. Gestión de Sesiones
- **Revocación Activa**: Los usuarios pueden listar sus sesiones activas y revocarlas individualmente (`GET /auth/sessions`).
- **Google OAuth**: Integración segura vía backend que intercambia códigos de autorización por tokens, evitando la exposición de secretos en el cliente.

## 3. Validación de Entradas
- DTOs (Data Transfer Objects) estrictos en la capa de entrada.
- Prevención de **SQL Injection** mediante consultas parametrizadas (GORM).

## 4. Limitación de Tasa (Rate Limiting)
- **Implementado**: Middleware de Rate Limiting respaldado por Redis para prevenir ataques de fuerza bruta y DDoS.
    - Límite global por defecto: 100 req/min.

## 5. Gestión de Secretos
- Las credenciales (DB Password, JWT Secret) se inyectan exclusivamente vía **Variables de Entorno**. NO se guardan en el código fuente (hardcoded).
N