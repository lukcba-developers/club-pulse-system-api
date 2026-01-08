# Informe Final de Auditoría de Seguridad - club-pulse-system-api

## Resumen Ejecutivo

El sistema presenta una base de seguridad sólida, especialmente en el backend. El uso correcto de GORM previene la inyección de SQL, y la implementación de un bloqueo distribuido con Redis para la concurrencia es excelente. Sin embargo, se han identificado varias vulnerabilidades de severidad media a alta que requieren atención, principalmente relacionadas con la gestión de tokens, la configuración de cabeceras de seguridad y el almacenamiento de tokens en el frontend. **Una auditoría posterior reveló fallos críticos en el Control de Acceso Basado en Roles (RBAC).**

---

## Vulnerabilidades Identificadas

### 1. Almacenamiento Inseguro de Tokens de Autenticación en el Frontend

*   **Severidad:** Alta
*   **Archivo/Línea:** `frontend/context/auth-context.tsx` (Implícito)
*   **Descripción:** El archivo `auth-context.tsx` no gestiona directamente los tokens, sino que depende de que las llamadas a la API (`/users/me`) funcionen. Esto implica que el token de autenticación (JWT) se almacena en un lugar accesible para el código JavaScript del lado del cliente, muy probablemente en `localStorage`. Almacenar tokens en `localStorage` es una práctica insegura, ya que cualquier script de terceros malicioso (introducido a través de un ataque XSS o una dependencia comprometida) puede robar el token y suplantar la identidad del usuario.
*   **Código Corregido (Conceptual):** La solución es migrar a un sistema basado en cookies `HttpOnly`.

    **Paso 1: Backend - Establecer la cookie en el login (`auth/service/auth_service.go`)**
    ```go
    // En el servicio de autenticación, después de generar el token
    func (s *AuthService) Login(...) (*domain.Token, error) {
        // ... lógica de login ...
        token, err := s.tokenService.GenerateToken(user)
        if err != nil {
            return nil, err
        }

        // --> INICIO DE LA CORRECCIÓN
        // En lugar de (o además de) devolver el token en el JSON,
        // establécelo en una cookie HttpOnly y segura.
        http.SetCookie(w, &http.Cookie{
            Name:     "access_token",
            Value:    token.AccessToken,
            Expires:  time.Now().Add(15 * time.Minute), // Coincide con la vida del token
            HttpOnly: true,                             // Previene acceso desde JS
            Secure:   true,                             // Solo enviar sobre HTTPS
            Path:     "/",
            SameSite: http.SameSiteLaxMode,
        })
        // <-- FIN DE LA CORRECCIÓN

        return token, nil // Puedes seguir devolviendo el token si la app móvil lo necesita
    }
    ```

    **Paso 2: Frontend - Eliminar la gestión manual de tokens.**
    El frontend ya parece estar configurado para esto. `axios` (`lib/axios.ts`) debe configurarse con `withCredentials: true` para que envíe automáticamente la cookie en cada petición. El `auth-context.tsx` ya depende de una llamada a `/users/me` que funcionará si la cookie es válida.

### 2. Tiempo de Expiración Excesivamente Largo para Access Token

*   **Severidad:** Media
*   **Archivo/Línea:** `backend/internal/modules/auth/infrastructure/token/jwt.go` (Línea 27 aprox.)
*   **Descripción:** El token de acceso (Access Token) tiene una duración de 24 horas. Si un token es robado, el atacante tiene un día completo para realizar acciones en nombre del usuario. Los access tokens deben tener una vida corta para minimizar la ventana de oportunidad de un atacante.
*   **Código Corregido:**
    ```go
    // backend/internal/modules/auth/infrastructure/token/jwt.go

    func (s *JWTService) GenerateToken(user *domain.User) (*domain.Token, error) {
        // expiration := time.Now().Add(24 * time.Hour) // ANTERIOR
        expiration := time.Now().Add(15 * time.Minute) // CORREGIDO: 15 minutos

        claims := jwt.MapClaims{
            // ...
            "exp": expiration.Unix(),
        }
        // ...
        return &domain.Token{
            AccessToken:  signedToken,
            RefreshToken: refreshToken,
            ExpiresIn:    900, // CORREGIDO: 15 minutos en segundos
        }, nil
    }
    ```

### 3. Política de Seguridad de Contenido (CSP) Débil

*   **Severidad:** Media
*   **Archivo/Línea:** `backend/internal/platform/http/middlewares/security.go` (Línea 26 aprox.)
*   **Descripción:** La cabecera `Content-Security-Policy` contiene `script-src 'self' 'unsafe-inline'`. El valor `'unsafe-inline'` permite la ejecución de scripts en línea (código JS dentro de etiquetas `<script>` o en atributos como `onclick`), lo que anula en gran medida la protección contra ataques de Cross-Site Scripting (XSS).
*   **Código Corregido:** Idealmente, se debe eliminar `'unsafe-inline'`. Esto puede requerir refactorizar el frontend para evitar scripts en línea. Si se usa Next.js con `nonce`, la configuración sería más compleja. Una solución pragmática es ser más restrictivo.

    ```go
    // backend/internal/platform/http/middlewares/security.go

    // ...
    // c.Header("Content-Security-Policy", "default-src 'self'; ... script-src 'self' 'unsafe-inline'; ...") // ANTERIOR

    // CORREGIDO (si no se pueden eliminar los inline scripts)
    // Añadir hashes o nonces es la mejor opción. Si no es posible, al menos restringir otros orígenes.
    // La mejor corrección es eliminar 'unsafe-inline' y refactorizar el frontend.
    c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; font-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'self' ws://localhost:3000")
    // ...
    ```
    **Nota:** La corrección ideal implica un esfuerzo en el frontend para eliminar todo el código "inline".

---

## Conclusiones y Recomendaciones Adicionales (Auditoría Exhaustiva)

### 4. Vulnerabilidad Crítica en RBAC - Eliminación de Usuarios sin Restricción de Rol

*   **Severidad:** Crítica
*   **Ruta:** `DELETE /users/:id`
*   **Archivo:** `backend/internal/modules/user/infrastructure/http/handler.go` (Línea 290 aprox., en `RegisterRoutes`)
*   **Descripción:** La ruta para eliminar un usuario (`DELETE /users/:id`) está protegida por autenticación, pero **carece de un middleware de autorización de roles**. Esto permite que cualquier usuario autenticado pueda intentar eliminar a otro usuario con solo conocer su ID. Un usuario con rol `USER` no debería tener acceso a esta funcionalidad bajo ninguna circunstancia.
*   **Código Corregido:** Se debe aplicar el middleware `RequireRole` a la ruta, permitiendo el acceso únicamente a los administradores.
    ```go
    // backend/internal/modules/user/infrastructure/http/handler.go

    // Dentro de la función RegisterRoutes
    func RegisterRoutes(r *gin.RouterGroup, handler *UserHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
        users := r.Group("/users")
        users.Use(authMiddleware, tenantMiddleware)
        {
            // ... otras rutas ...
            
            // ANTERIOR (Vulnerable)
            // users.DELETE("/:id", handler.DeleteUser)

            // CORREGIDO
            adminOnly := users.Group("")
            adminOnly.Use(middleware.RequireRole(domain.RoleAdmin, domain.RoleSuperAdmin))
            {
                adminOnly.DELETE("/:id", handler.DeleteUser)
                // También se debería mover ListUsers aquí
                adminOnly.GET("", handler.ListUsers)
            }
        }
    }
    ```

### 5. Vulnerabilidad Alta en RBAC - Fuga de Información de Usuarios

*   **Severidad:** Alta
*   **Ruta:** `GET /users`
*   **Archivo:** `backend/internal/modules/user/infrastructure/http/handler.go` (Línea 288 aprox., en `RegisterRoutes`)
*   **Descripción:** La ruta para listar todos los usuarios del club (`GET /users`) solo requiere autenticación, permitiendo a cualquier miembro del club obtener una lista completa de todos los demás miembros, incluyendo datos personales como nombre y correo electrónico.
*   **Recomendación:** Mover esta ruta bajo la protección del middleware `RequireRole`, como se muestra en la corrección anterior, para que solo los administradores puedan listar a todos los usuarios.

### 6. Vulnerabilidad Media en RBAC - Verificación de Rol en el Handler

*   **Severidad:** Media
*   **Ruta:** `POST /payments/offline`
*   **Archivo:** `backend/internal/modules/payment/infrastructure/http/handler.go` (Línea 130)
*   **Descripción:** La autorización para crear pagos offline se verifica dentro del código del manejador (`handler`) en lugar de usar un middleware a nivel de ruta. Esta práctica es propensa a errores y debilita la defensa en profundidad.
*   **Recomendación:** Refactorizar el registro de la ruta para que utilice el middleware `RequireRole`, garantizando que la política de seguridad se aplique de manera consistente y declarativa.
    ```go
    // backend/internal/modules/payment/infrastructure/http/handler.go

    // Dentro de la función RegisterRoutes
    func RegisterRoutes(r *gin.RouterGroup, handler *PaymentHandler, authMiddleware, tenantMiddleware gin.HandlerFunc) {
        payments := r.Group("/payments")
        {
            // ...
            // ANTERIOR
            // payments.POST("/offline", authMiddleware, tenantMiddleware, handler.CreateOfflinePayment)

            // CORREGIDO
            staffAndAdmin := payments.Group("")
            staffAndAdmin.Use(authMiddleware, tenantMiddleware, middleware.RequireRole(domain.RoleAdmin, domain.RoleStaff, domain.RoleSuperAdmin))
            {
                staffAndAdmin.POST("/offline", handler.CreateOfflinePayment)
                staffAndAdmin.GET("", handler.ListPayments) // Proteger también la lista de pagos
            }
            // ...
        }
    }
    ```

### 7. Seguridad de Dependencias

*   **Backend (Go):**
    *   **Estado:** No analizado.
    *   **Descripción:** La herramienta `govulncheck` no se encontró en el sistema, por lo que no se pudo realizar el escaneo de vulnerabilidades en las dependencias de Go.
    *   **Recomendación Crítica:** Instalar `govulncheck` con `go install golang.org/x/vuln/cmd/govulncheck@latest` y ejecutar `govulncheck ./...` dentro del directorio `backend`. Este es un paso indispensable para una auditoría completa.
*   **Frontend (Node.js):**
    *   **Estado:** Analizado y Limpio.
    *   **Descripción:** `npm audit` se ejecutó y no encontró ninguna vulnerabilidad en las dependencias del frontend.