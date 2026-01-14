# Manual de Usuario: Módulo de Perfil de Usuario (User)

## 1. Propósito

Este módulo te permite gestionar toda tu información personal, de contacto y deportiva. Mantener tu perfil actualizado es importante para que el club pueda comunicarse contigo, gestionar tu elegibilidad para competiciones y ofrecerte una experiencia personalizada.

## 2. Modelo de Datos del Perfil

El perfil de un usuario contiene una variedad de campos para gestionar su información de manera integral.

| Campo                     | Tipo        | Descripción                                                              |
| ------------------------- | ----------- | ------------------------------------------------------------------------ |
| `ID`                      | `string`    | Identificador único del usuario.                                         |
| `Name`                    | `string`    | Nombre completo del usuario.                                             |
| `Email`                   | `string`    | Correo electrónico (usado para inicio de sesión y comunicaciones).       |
| `Role`                    | `string`    | Rol principal del usuario en el sistema. Ver sección de Roles.           |
| `DateOfBirth`             | `date`      | Fecha de nacimiento, usada para calcular la categoría deportiva.         |
| `ClubID`                  | `string`    | Identificador del club al que pertenece el usuario.                      |
| `FamilyGroupID`           | `uuid`      | ID del grupo familiar al que pertenece (si aplica).                      |
| `MedicalCertStatus`       | `string`    | Estado del certificado médico (`VALID`, `EXPIRED`, `PENDING`).           |
| `MedicalCertExpiry`       | `date`      | Fecha de vencimiento del certificado médico.                             |
| `EmergencyContactName`    | `string`    | Nombre de un contacto de emergencia.                                     |
| `EmergencyContactPhone`   | `string`    | Teléfono del contacto de emergencia.                                     |
| `InsuranceProvider`       | `string`    | Proveedor del seguro médico o de accidentes.                             |
| `InsuranceNumber`         | `string`    | Número de póliza del seguro.                                             |
| `TermsAcceptedAt`         | `datetime`  | Fecha y hora en que el usuario aceptó los términos y condiciones.        |
| `PrivacyPolicyVersion`    | `string`    | Versión de la política de privacidad que el usuario aceptó.              |

## 3. Roles del Sistema

El sistema define varios roles con diferentes niveles de acceso y permisos.

| Rol              | Descripción                                                                                             |
| ---------------- | ------------------------------------------------------------------------------------------------------- |
| `MEMBER`         | Socio general del club. Puede ver y editar su propio perfil, hacer reservas y gestionar su membresía.     |
| `COACH`          | Entrenador. Puede gestionar equipos, registrar asistencia y ver perfiles de los miembros de su equipo.    |
| `ADMIN`          | Administrador del club. Puede gestionar usuarios, membresías, instalaciones y otros aspectos del club.  |
| `SUPER_ADMIN`    | Super-administrador con acceso total a todos los clubes y configuraciones del sistema.                  |
| `MEDICAL_STAFF`  | Personal médico. Tiene acceso restringido a información de salud (certificados médicos) por razones de cumplimiento (GDPR). |

---

## 4. Guía para Socios (Rol: `MEMBER`)

### 🔹 Cómo Ver y Editar tu Perfil

**Paso a paso:**
1.  **Inicia sesión** en la plataforma.
2.  **Navega a la sección "Mi Perfil"**. Generalmente, puedes acceder a ella haciendo clic en tu nombre o avatar en la esquina superior derecha.
3.  **Visualiza tu información:** Verás todos tus datos personales, de contacto y deportivos.
4.  **Haz clic en el botón "Editar Perfil"**.
5.  **Modifica los campos** que desees actualizar.
6.  **Guarda los cambios.** Haz clic en "Guardar" para aplicar las modificaciones.

### 🔹 Cómo Gestionar tu Información de Emergencia

Puedes registrar información de contacto que el club utilizará en caso de una emergencia.

**Paso a paso:**
1.  Ve a tu perfil.
2.  Busca la sección "Contacto de Emergencia".
3.  Rellena los campos: Nombre del Contacto, Teléfono del Contacto, Proveedor de Seguro y Número de Póliza.
4.  Guarda los cambios.

**Endpoint de la API:** `PUT /users/me/emergency`

### 🔹 Cómo Gestionar tus Hijos (Dependientes)

El sistema te permite registrar y gestionar las cuentas de tus hijos o dependientes directamente desde tu perfil.

**Paso a paso para registrar un hijo:**
1.  Ve a la sección "Mis Hijos" en tu perfil.
2.  Haz clic en "Registrar Nuevo Hijo".
3.  Completa el formulario con el nombre y la fecha de nacimiento de tu hijo.
4.  El sistema creará una nueva cuenta de usuario para tu hijo, vinculada a la tuya como padre/madre.

**Endpoints de la API:**
-   `GET /users/me/children`: Lista los hijos asociados a tu cuenta.
-   `POST /users/me/children`: Registra un nuevo hijo.

### 🔹 Estadísticas y Billetera

Para entender tu progreso, rendimiento y saldo en el club, consulta los siguientes módulos:
-   [**Estadísticas de Usuario (Stats)**](./stats.md): Revisa tu rendimiento en partidos, ranking y nivel.
-   [**Billetera (Wallet)**](./wallet.md): Consulta tu saldo monetario y tus puntos de recompensa.

---

## 5. Guía para Administradores (Rol: `ADMIN` / `SUPER_ADMIN`)

### 🔸 Cómo Buscar y Ver el Perfil de un Socio

**Paso a paso:**
1.  **Accede al Panel de Administración.**
2.  Navega a la sección de **"Usuarios"** o **"Socios"**.
3.  Utiliza la **barra de búsqueda** para encontrar a un socio por su nombre, apellido o correo electrónico.
4.  **Haz clic en el socio** en los resultados de búsqueda.
5.  Serás dirigido a una vista de solo lectura de su perfil, donde podrás consultar toda su información.

---

## 6. Diagrama de Flujo: Actualización de Perfil (Socio)

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
