# Manual de Usuario: Módulo de Club

## 1. Propósito

Este módulo es el centro de control para la información general de tu club. Aquí los administradores pueden configurar los datos de contacto, los horarios y publicar noticias para todos los socios.

## 2. Roles Implicados

-   **Administrador (`ADMIN`):** Gestiona toda la información y configuración del club.
-   **Socio (`MEMBER`):** Ve la información y las noticias publicadas.

---

## 3. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Editar la Información General del Club

**Paso a paso:**
1.  **Accede al Panel de Administración.**
2.  Navega a la sección **"Configuración del Club"** o "Información General".
3.  Desde aquí, podrás **editar los campos** principales:
    -   Nombre del Club.
    -   Dirección y teléfono de contacto.
    -   Horarios generales de apertura y cierre.
    -   Subir o cambiar el logotipo del club.
4.  **Guarda los cambios.** La información se actualizará en toda la plataforma (ej: en el pie de página y en la página de contacto).

### 🔸 Cómo Publicar una Noticia o Anuncio

**Paso a paso:**
1.  En el Panel de Administración, ve a la sección de **"Noticias"** o **"Anuncios"**.
2.  Haz clic en **"Crear Noticia"**.
3.  Escribe un **título y el contenido** del anuncio.
4.  Haz clic en **"Publicar"**. La noticia aparecerá en el panel principal para todos los socios cuando inicien sesión.

### 🔸 Cómo Configurar otros Módulos

Desde la configuración del club, también puedes ajustar parámetros de otros módulos.
-   **Ejemplo:** En la sección de "Reservas", puedes establecer la **política de cancelación** (ej: "cancelaciones permitidas hasta 24 horas antes").

---

## 4. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Ver la Información y Noticias del Club

**Paso a paso:**
1.  **Inicia sesión** en la plataforma.
2.  En el **panel principal o dashboard**, verás la sección de **"Últimas Noticias"** con los anuncios más recientes del club.
3.  Para ver la información de contacto o los horarios, generalmente encontrarás un enlace en el pie de página o en una sección llamada **"El Club"** o **"Contacto"**.

---

## 5. Diagrama de Flujo: Publicar una Noticia (Admin)

```mermaid
graph TD
    A[Inicio: Panel de Admin] --> B[Ir a "Gestión de Noticias"];
    B --> C[Clic en "Crear Noticia"];
    C --> D[Escribir Título y Contenido];
    D --> E[Clic en "Publicar"];
    E --> F[Noticia visible para todos los socios ✅];
```
