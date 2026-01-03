# Plan de Experiencia de Usuario por Rol (UX Role Matrix)

Este documento define el alcance, la navegación y los elementos visuales (Dashboard) específicos para cada uno de los tres perfiles de usuario del sistema.

---

## 1. 🌍 Super Admin (Plataforma Global)
**Perfil**: Dueño o Gerente de la Franquicia / Plataforma SaaS.
**Objetivo**: Supervisar la salud del negocio global y gestionar los inquilinos (clubes).
**Scope**: Multi-Tenant (Cross-Club).

### 🖥️ Dashboard (Landing Page)
1.  **Métricas Globales (KPIs)**:
    -   *Total MRR (Monthly Recurring Revenue)* de todos los clubes.
    -   *Clubes Activos* vs *Inactivos*.
    -   *Total Usuarios Registrados* en la plataforma.
2.  **Mapa de Estado**:
    -   Lista de clubes con indicadores de salud (Semaforo: Verde/Rojo según errores o pagos).

### 📍 Menú de Navegación (Sidebar)
-   **Dashboard Global**: Vista de águila.
-   **Gestión de Clubes (Tenants)**:
    -   *Alta de Club*: Formulario para crear nuevo Tenant (ID, Logo, Configuración Regional).
    -   *Ajustes de Club*: Suspender servicio, resetear password de admin local.
-   **Configuración del Sistema**:
    -   Feature Flags globales (activar/desactivar módulos por club).
    -   Ver Logs de Auditoría del Sistema.

---

## 2. 🏢 Administrador de Club (Role: ADMIN)
**Perfil**: Gerente o Recepcionista de una sede específica.
**Objetivo**: Maximizar la ocupación, gestionar el día a día y resolver problemas de socios.
**Scope**: Single-Tenant (Solo datos de su `club_id`).

### 🖥️ Dashboard (Landing Page)
1.  **Operación del Día**:
    -   *Ocupación Hoy*: % de canchas reservadas.
    -   *Próximos Ingresos*: Lista de reservas confirmadas para la próxima hora (Check-in rápido).
2.  **Alertas**:
    -   Reservas pendientes de pago (si aplica).
    -   Instalaciones en Mantenimiento.
3.  **Financial Snapshot**:
    -   Facturación del mes en curso.

### 📍 Menú de Navegación (Sidebar)
-   **Dashboard**: Vista operativa.
-   **Calendario Maestro**: Grid visual de todas las canchas (Drag & Drop para mover reservas - *Roadmap*).
-   **Instalaciones**:
    -   ABM (Alta/Baja/Modificación) de Canchas.
    -   Gestión de Bloqueos (Mantenimiento).
-   **Usuarios**:
    -   Lista de Socios.
    -   Gestión de Membresías (Asignar Tier, Ajustar Saldo).
-   **Configuración Sede**:
    -   Horarios de Apertura/Cierre.
    -   Reglas de Cancelación.

---

## 3. 👤 Usuario Miembro (Role: MEMBER)
**Perfil**: Jugador o Socio del club.
**Objetivo**: Reservar rápido, ver sus jugadas y pagar sin fricción.
**Scope**: Personal (Solo sus datos y disponibilidad pública).

### 🖥️ Dashboard (Landing Page)
1.  **Mi Próximo Partido**:
    -   Card destacada con fecha, hora, cancha y clima pronosticado.
    -   Botón "Cancelar" o "Invitar Amigos".
2.  **Reserva Rápida**:
    -   Accesos directos a sus deportes favoritos ("Reservar Padel", "Reservar Tenis").
3.  **Estado de Cuenta**:
    -   Aviso si hay cuota vencida o saldo en billetera.

### 📍 Menú de Navegación (Sidebar / Bottom Bar en Mobile)
-   **Inicio**: Dashboard personal.
-   **Reservar**: Buscador de disponibilidad (Filtros por deporte/fecha).
-   **Mis Reservas**: Historial y futuros turnos.
-   **Mi Perfil**:
    -   Datos Personales.
    -   Métodos de Pago (Tarjetas guardadas).
    -   Membresía (Ver Tier actual y beneficios).

---

## 🎨 Mejoras de Experiencia Propuestas (Action Items)

Para lograr esta segmentación efectiva:
1.  **Frontend (Role Guard)**: Implementar un `RoleBasedLayout` que renderice un Sidebar distinto según `user.role` (JWT).
2.  **Super Admin Dashboard**: Crear una ruta `/admin/platform` exclusiva para el Super Admin.
3.  **Onboarding**:
    -   *Admin*: Tour guiado para configurar la primera cancha.
    -   *Member*: Tutorial rápido "Cómo reservar en 3 clics".
