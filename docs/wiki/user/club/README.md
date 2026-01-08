# Manual de Usuario: Módulo de Club (Tenant)

## 1. Propósito y Rol Central

El módulo de **Club** es la entidad fundamental del sistema. En la arquitectura multi-tenant de la plataforma, cada `Club` funciona como un **inquilino (tenant)** independiente y aislado.

**Funcionalidades clave:**
-   **Aislamiento de Datos:** Cada club tiene su propio conjunto de usuarios, productos, campeonatos, etc. La información no se comparte entre clubes.
-   **Punto de Acceso:** Cada club tiene un `slug` único (un identificador para la URL, ej: `mi-club-favorito`), que define el punto de acceso a toda su información (ej: `plataforma.com/mi-club-favorito`).
-   **Configuración General:** Actúa como el centro de control donde los administradores del club gestionan la información general, publican noticias y configuran parámetros para otros módulos.

## 2. Roles Implicados

-   **Super Administrador (`SUPER_ADMIN`):** Gestiona la creación y el ciclo de vida de todos los clubes en la plataforma.
-   **Administrador de Club (`ADMIN`):** Gestiona la información y configuración de su propio club.
-   **Socio (`MEMBER`):** Ve la información y las noticias publicadas para su club.

---

## 3. Guía para Administradores de Club (Rol: `ADMIN`)

### 🔸 Cómo Editar la Información General del Club

**Paso a paso:**
1.  **Accede al Panel de Administración** de tu club.
2.  Navega a la sección **"Configuración del Club"** o "Información General".
3.  Desde aquí, podrás **editar los campos** principales:
    -   Nombre del Club.
    -   Dirección y teléfono de contacto.
    -   Horarios generales de apertura y cierre.
    -   Subir o cambiar el logotipo del club.
4.  **Guarda los cambios.** La información se actualizará en toda la plataforma para los usuarios de tu club.

### 🔸 Cómo Publicar una Noticia o Anuncio

**Paso a paso:**
1.  En el Panel de Administración, ve a la sección de **"Noticias"** o **"Anuncios"**.
2.  Haz clic en **"Crear Noticia"**.
3.  Escribe un **título y el contenido** del anuncio.
4.  Haz clic en **"Publicar"**. La noticia aparecerá en el panel principal para todos los socios de tu club.

---

## 4. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Ver la Información y Noticias del Club

**Paso a paso:**
1.  **Inicia sesión** en la plataforma en el contexto de tu club.
2.  En el **panel principal o dashboard**, verás la sección de **"Últimas Noticias"** con los anuncios más recientes.
3.  Para ver la información de contacto o los horarios, generalmente encontrarás un enlace en el pie de página o en una sección llamada **"El Club"** o **"Contacto"**.

---

## 5. Diagrama de Flujo: Publicar una Noticia (Admin)

```mermaid
graph TD
    A[Inicio: Panel de Admin del Club] --> B[Ir a "Gestión de Noticias"];
    B --> C[Clic en "Crear Noticia"];
    C --> D[Escribir Título y Contenido];
    D --> E[Clic en "Publicar"];
    E --> F[Noticia visible para los socios del Club ✅];
```
