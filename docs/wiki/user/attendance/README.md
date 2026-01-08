# Manual de Usuario: Módulo de Asistencia (Attendance)

## 1. Propósito

Este módulo es la herramienta digital para los entrenadores. Permite pasar lista de forma rápida y sencilla, manteniendo un registro histórico de la asistencia de los socios a las clases y entrenamientos.

## 2. Roles Implicados

-   **Entrenador (`COACH`):** Es el usuario principal. Pasa lista en sus clases.
-   **Administrador (`ADMIN`):** Puede supervisar los registros de asistencia de todo el club.
-   **Socio (`MEMBER`):** Puede consultar su propio historial de asistencia.

---

## 3. Guía para Entrenadores (Rol: `COACH`)

### 🔹 Cómo Tomar Asistencia para una Clase

**Paso a paso:**
1.  **Inicia sesión** con tu cuenta de entrenador.
2.  **Navega a la sección "Mis Grupos" o "Asistencia"**.
3.  Verás una lista de los grupos de entrenamiento que tienes asignados.
4.  **Selecciona el grupo** para el cual deseas tomar asistencia.
5.  El sistema mostrará la **lista de socios inscritos** en ese grupo para la fecha actual.
6.  Para cada socio, **selecciona su estado**:
    -   `Presente`
    -   `Ausente`
    -   `Tarde`
    -   `Justificado`
7.  Una vez que hayas marcado a todos los socios, haz clic en **"Guardar Asistencia"**. El registro quedará guardado.

---

## 4. Guía para Socios (Rol: `MEMBER`)

### 🔸 Cómo Ver tu Historial de Asistencia

**Paso a paso:**
1.  **Inicia sesión** en tu cuenta.
2.  Ve a **"Mi Perfil"** y busca la pestaña de **"Asistencia"** o "Mi Progreso".
3.  Verás un resumen de tu historial de asistencia a las clases en las que estás inscrito.

---

## 5. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Ver Reportes de Asistencia

**Paso a paso:**
1.  **Accede al Panel de Administración.**
2.  Navega a la sección de **"Reportes"** -> **"Asistencia"**.
3.  Podrás **filtrar los registros de asistencia** por grupo, entrenador o rango de fechas para analizar la concurrencia a las clases.

---

## 6. Diagrama de Flujo: Toma de Asistencia (Entrenador)

```mermaid
graph TD
    A[Inicio: Entrenador inicia sesión] --> B[Ir a "Mis Grupos"];
    B --> C[Seleccionar Grupo y Fecha];
    C --> D[Sistema muestra la lista de alumnos];
    D --> E[Entrenador marca el estado de cada alumno];
    E --> F[Clic en "Guardar Asistencia"];
    F --> G[Registro Guardado ✅];
```
