# Guía del Desarrollador (Contributing)

¡Bienvenido al equipo de desarrollo de Club Pulse! Esta guía te ayudará a configurar tu entorno y entender nuestros estándares.

## Configuración del Entorno (Setup)

### Prerrequisitos
- **Go** 1.21+
- **Docker** y Docker Compose
- **Node.js** 18+ (para el Frontend)

### Pasos Iniciales
1.  **Clonar el repositorio**:
    ```bash
    git clone <repo-url>
    cd club-pulse-system-api
    ```

2.  **Iniciar Infraestructura (Base de Datos)**:
    ```bash
    docker-compose up -d postgres redis
    ```

3.  **Ejecutar Backend**:
    ```bash
    cd backend
    go mod download
    # Iniciar servidor
    DB_PASSWORD=pulse_secret PORT=8081 go run cmd/api/main.go
    ```

4.  **Ejecutar Frontend**:
    ```bash
    cd frontend
    npm install
    npm run dev
    ```

### Seed Data (Datos de Prueba)
Para poblar la base de datos con usuarios admin y facilidades de prueba:
```bash
cd backend
DB_PASSWORD=pulse_secret go run cmd/seeder/main.go
```

## Estándares de Código

### Backend (Go)
- **SOLID**: Respetar el Principio de Responsabilidad Única. Cada archivo/función debe hacer una sola cosa.
- **Nombres**: Variables en `camelCase`, Exportadas en `PascalCase`.
- **Errores**: Manejar errores explícitamente. No ignorar errores (`_`).
- **Logs**: Usar el logger estructurado, no `fmt.Println` en producción.

### Frontend (React/Next.js)
- **Componentes**: Pequeños y reutilizables.
- **Hooks**: Extraer lógica compleja a custom hooks (`useAuth`, `useBookings`).
- **Tipado**: Usar TypeScript estricto. Evitar `any`.

## Flujo de Trabajo (Git)
1.  Crear rama desde `main`: `feature/nombre-de-la-tarea`.
2.  Commit cambios: `git commit -m "feat: descripción"`.
3.  Abrir Pull Request.
4.  Mergear a `main` tras revisión.
