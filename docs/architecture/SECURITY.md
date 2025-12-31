# Protocolos de Seguridad

La seguridad es una prioridad en **Club Pulse**. Implementamos múltiples capas de protección para garantizar la integridad de los datos y la seguridad de los usuarios.

## 1. Autenticación y Autorización
- **JWT (JSON Web Tokens)**: Utilizamos tokens JWT firmados (HS256) para la autenticación sin estado.
    - **Access Token**: Vida corta (24 horas para desarrollo, usualmente 15-60 min en prod).
- **RBAC (Role-Based Access Control)**:
    - Middleware de autorización que verifica roles (`ADMIN`, `MEMBER`) antes de permitir el acceso a endpoints sensibles.
- **Passwords**: Hasheadas usando **Bcrypt** con un costo computacional adecuado.

## 2. Protección de API
- **Headers OWASP**:
    - `X-Content-Type-Options: nosniff`
    - `X-Frame-Options: DENY`
    - `X-XSS-Protection: 1; mode=block`
    - `Strict-Transport-Security` (HSTS) habilitado.
- **CORS**: Configurado restrictivamente para permitir solo orígenes confiables (Frontend).

## 3. Validación de Entradas
- DTOs (Data Transfer Objects) estrictos en la capa de entrada.
- Validación de tipos de datos y formatos (email, fechas, UUIDs) antes de procesar lógica de negocio.
- Prevención de **SQL Injection** mediante el uso de ORM (GORM) y consultas parametrizadas.

## 4. Limitación de Tasa (Rate Limiting)
- *Planificado*: Implementar middleware de Rate Limiting por IP para prevenir ataques de fuerza bruta en endpoints de login.

## 5. Gestión de Secretos
- Las credenciales (DB Password, JWT Secret) se inyectan exclusivamente vía **Variables de Entorno**. NO se guardan en el código fuente (hardcoded).
N