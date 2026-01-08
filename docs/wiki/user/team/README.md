# Manual de Usuario: Módulo de Equipos (Team)

## 1. Propósito

Este módulo te permite crear y gestionar tus propios equipos para competir en los torneos del club. Puedes juntarte con tus amigos, elegir un nombre y un escudo, y prepararse para la competición.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede crear equipos, ser capitán, unirse a equipos e invitar a otros socios.

---

## 3. Guía de Usuario (Rol: `MEMBER`)

### 🔹 Cómo Crear un Nuevo Equipo

Si quieres ser el capitán, puedes crear tu propio equipo.

**Paso a paso:**
1.  **Navega a la sección "Mis Equipos"** en tu perfil o en el menú de "Campeonatos".
2.  Haz clic en el botón **"Crear Nuevo Equipo"**.
3.  **Completa el formulario:**
    -   **Nombre del Equipo:** ¡Elige un nombre original!
    -   **Logo/Escudo:** Sube una imagen para representar a tu equipo.
4.  **Guarda los cambios.** ¡Tu equipo ha sido creado y tú eres el capitán!

### 🔹 Cómo Invitar Jugadores a tu Equipo

Como capitán, eres el encargado de reclutar a tus compañeros.

**Paso a paso:**
1.  Ve a la página de gestión de tu equipo.
2.  Busca la opción **"Invitar Jugadores"**.
3.  Se abrirá un buscador donde podrás **encontrar a otros socios** del club por su nombre.
4.  Selecciona a los socios que quieres invitar y haz clic en **"Enviar Invitación"**.
5.  Los socios recibirán una notificación para unirse a tu equipo.

### 🔹 Cómo Aceptar o Rechazar una Invitación

Si un capitán te invita a su equipo, recibirás una notificación.

**Paso a paso:**
1.  Ve a tu panel de notificaciones o a la sección "Mis Equipos".
2.  Verás la invitación pendiente con el nombre del equipo.
3.  Tendrás los botones **"Aceptar"** y **"Rechazar"**. Haz clic en la opción que prefieras.
4.  Si aceptas, pasarás a formar parte del equipo.

### 🔹 Cómo Salir de un Equipo

**Paso a paso:**
1.  Ve a la página del equipo del que formas parte.
2.  Busca la opción **"Abandonar Equipo"**.
3.  Confirma tu decisión. Dejarás de ser miembro de ese equipo.

---

## 4. Diagrama de Flujo: Creación y Formación de un Equipo

```mermaid
graph TD
    A[Capitán: Clic en "Crear Equipo"] --> B[Rellena Nombre y Logo];
    B --> C[Equipo Creado ✅];
    C --> D[Capitán: Invita a Jugadores];
    D --> E[Jugador Invitado: Recibe Notificación];
    E --> F{¿Acepta la Invitación?};
    F -- Sí --> G[Jugador se une al Equipo];
    F -- No --> H[Invitación Rechazada];
    G --> I[Equipo listo para competir];
```
