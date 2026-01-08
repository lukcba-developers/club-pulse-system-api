# Matriz de Casos de Uso y Roles del Sistema Club Pulse

## 1. Resumen de Roles

-   **SUPER_ADMIN**: Propietario de la plataforma SaaS. Gestiona los clubes (tenants) y la configuración global del sistema.
-   **ADMIN**: Gerente o administrador de un club específico. Gestiona las operaciones diarias de su sede.
-   **MEMBER**: Socio/cliente del club. Utiliza el sistema para reservar, gestionar su perfil y ver información.

---

## 2. Módulos y Casos de Uso

A continuación se detallan las funcionalidades por módulo y los roles que tienen permiso para realizar cada acción.

### Módulo: Autenticación y Usuarios (Auth & User)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Registrarse en la plataforma | ✅ | ✅ | ✅ | Creación de una nueva cuenta de usuario. |
| Iniciar Sesión (Login) | ✅ | ✅ | ✅ | Autenticación con email y contraseña. |
| Ver y modificar su propio perfil | ✅ | ✅ | ✅ | Cada usuario puede gestionar sus datos personales. |
| Ver lista de usuarios del club | | ✅ | | El ADMIN puede ver todos los miembros de su club. |
| Gestionar datos de un usuario | | ✅ | | El ADMIN puede modificar datos de un miembro de su club. |

### Módulo: Clubes (Club)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Crear un nuevo Club (Tenant) | ✅ | | | Dar de alta una nueva sede en la plataforma. |
| Configurar datos de un Club | ✅ | ✅ | | `SUPER_ADMIN` ajusta configuraciones globales, `ADMIN` ajusta las de su sede (horarios, etc.). |
| Suspender/Activar un Club | ✅ | | | Control del estado del servicio para un tenant. |

### Módulo: Instalaciones (Facilities)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Crear/Modificar una instalación | | ✅ | | Dar de alta canchas, piscinas, etc., con sus atributos. |
| Cambiar estado de una instalación | | ✅ | | Poner una instalación como "En Mantenimiento" o "Cerrada". |
| Consultar instalaciones y su estado | | ✅ | ✅ | Ver la lista de instalaciones disponibles para reservar. |
| Definir horarios de funcionamiento | | ✅ | | Establecer las horas en que una instalación está operativa. |

### Módulo: Reservas (Bookings)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Consultar disponibilidad en tiempo real | | ✅ | ✅ | Ver el calendario con los slots horarios disponibles. |
| Crear una reserva para un miembro | | ✅ | ✅ | El `MEMBER` crea su propia reserva, el `ADMIN` puede crearla para cualquier miembro. |
| Cancelar una reserva | | ✅ | ✅ | El `MEMBER` cancela su propia reserva, el `ADMIN` puede cancelar cualquiera. |
| Crear reserva recurrente | | ✅ | | Programar una reserva que se repite en el tiempo (ej. todos los lunes). |
| Unirse a lista de espera | | | ✅ | Anotarse si un horario deseado está ocupado. |
| Gestionar bloqueos de calendario | | ✅ | | Bloquear horarios por eventos especiales o mantenimiento. |

### Módulo: Membresías (Membership)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Crear/Modificar planes de membresía | | ✅ | | Definir los Tiers (Individual, Familiar) con sus precios y beneficios. |
| Asignar o cambiar la membresía de un socio | | ✅ | | Mover a un usuario de un plan a otro. |
| Consultar estado de su membresía | | | ✅ | Ver el estado actual de su plan (Activa, Vencida). |
| Procesar facturación mensual | | ✅ | | Ejecutar el proceso automático que genera la deuda a los socios. |
| Aplicar becas o descuentos | | ✅ | | Asignar un descuento porcentual a la cuota de un socio. |

### Módulo: Pagos (Payments)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Pagar una deuda (membresía, reserva) | | ✅ | ✅ | El `MEMBER` paga sus propias deudas, el `ADMIN` puede registrar un pago en su nombre. |
| Registrar un pago offline (efectivo, etc.) | | ✅ | | Marcar una deuda como pagada por medios no digitales. |
| Ver historial de pagos | | ✅ | ✅ | El `MEMBER` ve sus pagos, el `ADMIN` ve todos los del club. |
| Recibir y procesar webhooks de pago | ✅ | | | El sistema procesa notificaciones automáticas del proveedor de pago. |

### Módulo: Disciplinas y Grupos (Disciplines)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Crear/Eliminar una disciplina | | ✅ | | Define los deportes o actividades disponibles en el club (ej. Pádel, Yoga). |
| Listar disciplinas | | ✅ | ✅ | Muestra todas las disciplinas activas del club. |
| Crear un grupo de entrenamiento | | ✅ | | Crea un grupo para una disciplina, con categoría y horario (ej. "Pádel Infantil Lunes 18hs"). |
| Listar grupos de entrenamiento | | ✅ | ✅ | Muestra los grupos disponibles para unirse. |
| Listar alumnos de un grupo | | ✅ | | Muestra los socios inscritos en un grupo de entrenamiento. |

### Módulo: Equipos (Team)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Registrar un equipo en un torneo | | ✅ | ✅ | Inscribe un equipo con un nombre y sus miembros. Puede hacerlo un `ADMIN` o un `MEMBER` (capitán). |
| Ver los equipos de un torneo | | ✅ | ✅ | Lista los equipos que participan en un campeonato. |

### Módulo: Campeonatos (Championship)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Crear un campeonato | | ✅ | | Define un nuevo torneo para una disciplina, con fechas y formato. |
| Listar campeonatos | | ✅ | ✅ | Muestra todos los campeonatos (abiertos, en curso, finalizados). |
| Programar un partido | | ✅ | | Agenda un encuentro entre dos equipos con fecha, hora y lugar. |
| Listar partidos de un campeonato | | ✅ | ✅ | Muestra el fixture completo del torneo. |
| Cargar/Actualizar resultado de un partido | | ✅ | | Registra el marcador final de un encuentro. |
| Consultar la tabla de posiciones | | ✅ | ✅ | Muestra la clasificación de los equipos en el torneo. |

### Módulo: Tienda (Store)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Gestionar productos de la tienda | | ✅ | | **(Inferido)** Crear, editar, activar y controlar el stock de productos. |
| Consultar el catálogo de productos | | ✅ | ✅ | Ver los productos disponibles para la venta, con opción de filtrar por categoría. |
| Realizar una compra | | ✅ | ✅ | Comprar uno o más artículos del catálogo. La compra puede ser para un socio o un invitado. |

### Módulo: Control de Acceso (Access)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Registrar un intento de acceso | | ✅ | | El sistema registra (ej. con un QR o tarjeta) si un socio intenta ingresar. |
| Validar acceso según membresía/reserva | | ✅ | | El sistema comprueba si el socio tiene permiso para entrar. |

### Módulo: Asistencia (Attendance)

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Marcar asistencia a una clase/reserva | | ✅ | | Registrar que un socio asistió a un evento programado. |
| Consultar historial de asistencia | | ✅ | ✅ | Ver el registro de asistencias pasadas. |

### Módulo: Notificaciones (Notification)
*Este es un módulo transversal que probablemente no tiene casos de uso directos para el usuario, sino que es utilizado por otros módulos.*

| Caso de Uso | SUPER_ADMIN | ADMIN | MEMBER | Descripción |
| :--- | :---: | :---: | :---: | :--- |
| Enviar notificación de reserva confirmada | N/A | N/A | N/A | El sistema notifica al `MEMBER` automáticamente. |
| Enviar recordatorio de pago | N/A | N/A | N/A | El sistema notifica al `MEMBER` sobre deudas pendientes. |
