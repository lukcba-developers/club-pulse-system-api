# Manual de Usuario: Módulo de Perfil de Usuario (User)

## 1. Propósito

Este módulo te permite gestionar toda tu información personal y de contacto. Mantener tu perfil actualizado es importante para que el club pueda comunicarse contigo y ofrecerte una experiencia personalizada.

## 2. Roles Implicados

-   **Socio (`MEMBER`):** Puede ver y editar su propio perfil.
-   **Administrador (`ADMIN`):** Puede buscar y ver los perfiles de todos los socios del club.

---

## 3. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Ver y Editar tu Perfil

**Paso a paso:**
1.  **Inicia sesión** en la plataforma.
2.  **Navega a la sección "Mi Perfil"**. Generalmente, puedes acceder a ella haciendo clic en tu nombre o avatar en la esquina superior derecha.
3.  **Visualiza tu información:** Verás todos tus datos personales registrados, como nombre, email, teléfono, etc.
4.  **Haz clic en el botón "Editar Perfil"**.
5.  **Modifica los campos** que desees actualizar (por ejemplo, tu número de teléfono o dirección).
6.  **Guarda los cambios.** Haz clic en "Guardar" para aplicar las modificaciones.

### 🔹 Cómo Gestionar tu Grupo Familiar (Próximamente)

Esta funcionalidad te permitirá agrupar y gestionar las cuentas de tus familiares (ej: hijos) desde tu propio perfil. Podrás gestionar sus membresías y reservas de forma centralizada.

---

## 4. Guía para Administradores (Rol: `ADMIN`)

### 🔸 Cómo Buscar y Ver el Perfil de un Socio

**Paso a paso:**
1.  **Accede al Panel de Administración.**
2.  Navega a la sección de **"Usuarios"** o **"Socios"**.
3.  Utiliza la **barra de búsqueda** para encontrar a un socio por su nombre, apellido o correo electrónico.
4.  **Haz clic en el socio** en los resultados de búsqueda.
5.  Serás dirigido a una vista de solo lectura de su perfil, donde podrás consultar toda su información de contacto y estado de membresía.

---

## 5. Diagrama de Flujo: Actualización de Perfil (Socio)

```mermaid
graph TD
    A[Inicio] --> B[Ir a "Mi Perfil"];
    B --> C[Página de Perfil];
    C --> D[Clic en "Editar Perfil"];
    D --> E[Modificar Información en el Formulario];
    E --> F[Clic en "Guardar"];
    F --> G{¿Datos Válidos?};
    G -- Sí --> H[Perfil Actualizado ✅];
    G -- No --> I[Mostrar Error de Validación];
    I --> E;
    H --> C;
```
