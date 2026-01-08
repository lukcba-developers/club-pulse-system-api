# Manual de Usuario: Módulo de Disciplinas

## 1. Propósito

Este módulo organiza toda la oferta deportiva y de actividades del club. Aquí puedes explorar qué deportes se practican, ver los diferentes grupos de entrenamiento y encontrar el que mejor se adapte a tu nivel y edad.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede ver las disciplinas y los grupos.
-   **Administrador (`ADMIN`):** Puede gestionar las disciplinas, los grupos y asignar entrenadores.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Explorar las Disciplinas y Grupos

**Paso a paso:**
1.  **Navega a la sección "Disciplinas" o "Deportes"** en la aplicación.
2.  Verás una lista de todas las actividades que ofrece el club (ej: Tenis, Fútbol, Natación, Yoga).
3.  **Haz clic en una disciplina** que te interese.
4.  Se mostrará una página con información sobre esa disciplina y una lista de los **grupos de entrenamiento** disponibles.
5.  Cada grupo tendrá detalles como:
    -   **Nivel o categoría** (ej: "Infantil", "Adulto Principiante").
    -   **Entrenador** a cargo.
    -   **Horarios** de las clases.
    -   **Instalación** donde se realiza.
6.  Desde aquí, podrás solicitar tu inscripción a un grupo.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Crear o Editar una Disciplina

**Paso a paso:**
1.  **Accede al Panel de Administración** y ve a la sección de **"Disciplinas"**.
2.  Para crear una nueva, haz clic en **"Nueva Disciplina"**. Para editar, búscala en la lista y haz clic en "Editar".
3.  **Completa el formulario** con el nombre del deporte o actividad.
4.  **Guarda los cambios.**

### 🔸 Cómo Gestionar los Grupos de un Deporte

**Paso a paso:**
1.  En la lista de disciplinas, haz clic en la que deseas gestionar.
2.  Verás una opción para **"Añadir Grupo de Entrenamiento"**.
3.  **Completa el formulario del grupo:**
    -   Nombre del grupo (ej: "Competición Sub-16").
    -   Asigna un **entrenador** de la lista de usuarios con rol `COACH`.
    -   Define los **horarios y días** de entrenamiento.
    -   Selecciona la **instalación** que utilizará el grupo. El sistema puede bloquear automáticamente esos horarios en el calendario de reservas.
4.  **Gestiona los miembros:** Desde la página del grupo, podrás ver la lista de socios inscritos, aceptar nuevas solicitudes o añadir miembros manualmente.

---

## 5. Diagrama de Flujo: Organización de Disciplinas (Admin)

```mermaid
graph TD
    A[Inicio: Panel de Admin] --> B[Ir a "Disciplinas"];
    B --> C{¿Disciplina ya existe?};
    C -- No --> D[Crear Nueva Disciplina];
    D --> E[Guardar Disciplina];
    E --> F[Seleccionar Disciplina];
    C -- Sí --> F;
    F --> G[Clic en "Añadir Grupo"];
    G --> H[Rellenar Detalles del Grupo (horario, coach, etc.)];
    H --> I[Guardar Grupo ✅];
    I --> J[Gestionar Miembros del Grupo];
```
