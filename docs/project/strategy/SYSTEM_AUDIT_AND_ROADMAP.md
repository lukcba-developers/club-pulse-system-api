# Auditoría Estratégica y Hoja de Ruta: Club de Alto Rendimiento

**Documento:** `SYSTEM_AUDIT_AND_ROADMAP.md`
**Rol:** Arquitecto de Sistemas & Consultor Deportivo
**Fecha:** Enero 2026

---

## 1. Diagnóstico: Sistema Actual vs. Modelo de Alto Rendimiento

El siguiente cuadro compara la capacidad instalada actual (MVP Club Pulse) contra el estándar de la industria para clubes sociales modernos.

| Área | Sistema Actual (MVP) | Modelo de Alto Rendimiento (Objetivo) | Brecha (Gap) |
| :--- | :--- | :--- | :--- |
| **Ingresos** | Cobro manual o desconectado. Planes estáticos. | **Cobro Automatizado:** Débito automático (tarjetas), gestión de morosos, venta de "extras" (campus, indumentaria). | 🔴 Crítica |
| **Organización** | Listas de usuarios planas. Sin distinción de edad/categoría automática. | **Gestión por Categorías:** Asignación automática por año (ej. "2012"). Control de presentismo QR por disciplina. | 🔴 Crítica |
| **Administración** | Gestión de canchas (Booking). | **ERP Integral:** Liquidación de sueldos profes, control de inventario (pelotas/conos), alertas de certificados médicos. | 🟡 Media |
| **Socio (UX)** | Login y reserva de canchas. | **Autogestión Total (App):** Pagar cuota, ver deuda, carnet digital QR, ver asistencias. | 🟡 Media |
| **Infraestructura** | Monolito básico (costo bajo). | **Escalabilidad:** Notificaciones Push, Molinetes de acceso, Integración con contabilidad. | 🟢 Baja (Inicial) |

---

## 2. Matriz de Priorización (Impacto vs. Esfuerzo)

Clasificamos las funcionalidades faltantes para maximizar el retorno de inversión (ROI) a corto plazo.

| Cuadrante | Funcionalidad | Impacto Económico | Esfuerzo Dev | Acción |
| :--- | :--- | :---: | :---: | :--- |
| **💎 GANANCIAS RÁPIDAS** | **Link de Pago / Botón Pagar Cuota** | ⭐⭐⭐⭐⭐ | 📉 Bajo | **Hacer YA (Mes 1)** |
| **💎 GANANCIAS RÁPIDAS** | **Alertas de Deuda en Dashboard** | ⭐⭐⭐⭐ | 📉 Bajo | **Hacer YA (Mes 1)** |
| **🚀 ESTRATÉGICOS** | **Débito Automático (Suscripciones)** | ⭐⭐⭐⭐⭐ | 📈 Medio | **Planificar (Mes 2)** |
| **🚀 ESTRATÉGICOS** | **Gestión Categorías (Fútbol Infantil)** | ⭐⭐⭐⭐ | 📈 Medio | **Planificar (Mes 2)** |
| **🛠️ NECESARIOS** | Control de Acceso (QR Simple) | ⭐⭐⭐ | 📉 Bajo | Hacer (Mes 3) |
| 🐢 POSPONIBLES | Liquidación de Sueldos / Stock | ⭐⭐ | 📈 Alto | Backlog |

---

## 3. Hoja de Ruta de Ejecución (3 Meses)

Objetivo: Pasar de "Software de Reservas" a "Motor de Ingresos del Club".

### Mes 1: "Cashflow First" (Caja Inmediata)
**Foco:** Que el dinero entre al club lo más fácil posible.
14.  **Integración MercadoPago (Básica):** ✅ **[HECHO]**
    -   Botón "Pagar Deuda" en el perfil del socio (Backend listo).
    -   Webhook para impactar pago automáticamente (Backend listo).
2.  **Visualización de Deuda:** ✅ **[HECHO]**
    -   Cartel ROJO en el dashboard si debe cuota (Lógica `OutstandingBalance` lista).
    -   Bloqueo de Reserva de Canchas si hay deuda (Regla de negocio crítica).
3.  **Base de Datos Social:** ✅ **[HECHO]**
    -   Carga masiva de socios con `DateOfBirth` para preparar las categorías (Schema listo).

### Mes 2: "Orden Institucional" (Gestión Deportiva)
**Foco:** Organizar el caos de las disciplinas (ej. Fútbol Infantil).
1.  **Motor de Categorías:**
    -   Algoritmo: `Fecha Nacimiento` -> Asigna `Grupo de Entrenamiento` (ej. "Pre-Novena").
2.  **Listas de Asistencia Digitales:**
    -   Vista para "Profes": Lista de sus alumnos en el celular.
    -   Marcar presente/ausente con un tap.
3.  **Débito Automático (Suscripciones):**
    -   Migrar socios a suscripción recurrente (tarjeta guardada). Reduce la morosidad un 40%.

### Mes 3: "Experiencia y Control" (Fidelización)
**Foco:** Profesionalizar el acceso y la comunicación.
1.  **Carnet Digital (QR):**
    -   En la App del socio.
    -   Validación simple en portería (Escanear con celular del guardia/admin).
2.  **Notificaciones Push/Email:**
    -   "Tu cuota vence mañana".
    -   "Entrenamiento suspendido por lluvia".

---

## 4. Conclusión del Consultor

El sistema actual es una **base sólida tecnológica** (Go/Next.js), pero funcionalmente hoy es solo un "alquiler de canchas". Para convertirlo en un **Sistema de Gestión de Club**, la prioridad absoluta debe ser el **Módulo de Socio y Pagos**.

No construyas funciones complejas de torneos o tienda online todavía. **Cobra la cuota social de forma automática y organiza a los chicos por categoría.** Eso justifica el software ante la comisión directiva.
