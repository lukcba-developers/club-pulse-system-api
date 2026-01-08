# Manual de Usuario: Módulo de Campeonatos

## 1. Propósito

Este módulo te permite participar y seguir los torneos y competiciones organizados por el club. Fomenta la competencia amistosa y la interacción entre socios.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede ver torneos, inscribir equipos y seguir los resultados.
-   **Administrador (`ADMIN`):** Puede crear y gestionar todos los aspectos de un torneo.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Ver los Campeonatos Disponibles

**Paso a paso:**
1.  **Navega a la sección "Campeonatos"** en la aplicación.
2.  Verás una lista de los torneos actuales y futuros.
3.  Haz clic en un torneo para ver sus detalles:
    -   **Reglamento:** Las reglas específicas del torneo.
    -   **Fechas:** Fecha de inicio, fin e inscripción.
    -   **Formato:** (ej: Liga, Eliminación Directa).
    -   **Equipos Inscritos.**

### 🔹 Cómo Inscribirte a un Campeonato

**Paso a paso:**
1.  Dentro de la página de detalles de un torneo abierto, busca el botón **"Inscribirme"** o **"Inscribir Equipo"**.
2.  Si el torneo es por equipos, se te pedirá que selecciones un equipo que hayas creado previamente en el **Módulo de Equipos** o que crees uno nuevo.
3.  Confirma la inscripción. Puede que se te redirija al **Módulo de Pagos** si la inscripción tiene un costo.
4.  Una vez inscrito, tu equipo aparecerá en la lista de participantes.

### 🔹 Cómo Seguir un Torneo

**Paso a paso:**
1.  Entra a la página de detalles del torneo que deseas seguir.
2.  Navega por las diferentes pestañas para ver:
    -   **Fixture:** El calendario de todos los partidos.
    -   **Tabla de Posiciones:** La clasificación de los equipos actualizada en tiempo real.
    -   **Resultados:** Los marcadores de los partidos que ya se han jugado.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Crear un Nuevo Campeonato

**Paso a paso:**
1.  **Accede al Panel de Administración** y ve a la sección de **"Campeonatos"**.
2.  Haz clic en **"Nuevo Campeonato"**.
3.  **Completa el formulario** con toda la información: nombre, disciplina, fechas, formato, reglas, costo de inscripción, etc.
4.  **Guarda los cambios.** El torneo se publicará y los socios podrán empezar a inscribirse.

### 🔸 Cómo Gestionar un Torneo en Curso

**Paso a paso:**
1.  Desde el panel de "Campeonatos", selecciona el torneo que deseas gestionar.
2.  Desde aquí podrás:
    -   **Aprobar o rechazar inscripciones** de equipos.
    -   **Generar el fixture** (calendario de partidos) una vez cerradas las inscripciones.
    -   **Cargar los resultados** de los partidos a medida que se van jugando. La tabla de posiciones se actualizará automáticamente.

---

## 5. Diagrama de Flujo: Inscripción a un Torneo (Socio)

```mermaid
graph TD
    A[Inicio] --> B[Explorar Campeonatos];
    B --> C{Elige un Torneo Abierto};
    C --> D[Clic en "Inscribir Equipo"];
    D --> E{¿Equipo ya creado?};
    E -- Sí --> F[Selecciona tu Equipo];
    E -- No --> G[Ir a Módulo de Equipos y Crear Equipo];
    G --> F;
    F --> H{¿Inscripción tiene costo?};
    H -- Sí --> I[Ir a Módulo de Pagos];
    I --> J[Confirmación de Inscripción ✅];
    H -- No --> J;
```
