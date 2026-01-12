# Plan de Pruebas de Navegabilidad E2E (End-to-End)

Este documento define la suite completa de pruebas para verificar la navegabilidad y funcionalidad de la aplicación `club-pulse-system-api`. Las pruebas están organizadas por **Roles de Usuario** y **Módulos Funcionales**.

---

## 1. Matriz de Roles y Permisos

| Rol | Descripción | Alcance Principal |
| :--- | :--- | :--- |
| **`SUPER_ADMIN`** | Administrador de plataforma | Gestión multi-tenant, clubes, configuración global |
| **`ADMIN`** | Administrador de club | Instalaciones, usuarios, reservas, pagos, campeonatos |
| **`COACH`** | Entrenador | Gestión de equipos, asistencia, viajes, jugadores |
| **`MEMBER`** | Socio/Jugador | Reservas, perfil, membresía, tienda, campeonatos |
| **`MEDICAL_STAFF`** | Personal médico | Acceso a datos sensibles (GDPR Art. 9) |
| **`GUEST`** | Invitado/Público | Registro, páginas legales, landing de club |

---

## 2. Suite de Pruebas por Rol

### 2.1. Navegación Pública (Sin Autenticación)

**Objetivo:** Verificar accesibilidad de páginas públicas y flujos de entrada.

#### Módulo: Legal
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-PUB-01 | Ver Términos y Condiciones | Navegar a `/legal/terms` | Texto legal visible, título "Términos y Condiciones" | Alta |
| TC-PUB-02 | Ver Política de Privacidad | Navegar a `/legal/privacy` | Texto de privacidad visible, secciones GDPR | Alta |

#### Módulo: Registro de Jugador
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-PUB-03 | Ver formulario de registro | Navegar a `/register-player` | Formulario visible con campos obligatorios | Alta |
| TC-PUB-04 | Registro con datos válidos | Completar formulario → Submit | Redirección a login o confirmación | Alta |
| TC-PUB-05 | Validación de campos requeridos | Submit formulario vacío | Mensajes de error en campos obligatorios | Alta |
| TC-PUB-06 | Aceptar términos y privacidad | Marcar checkboxes de consentimiento | Checkboxes marcados, botón habilitado | Alta |
| TC-PUB-07 | Validación de email duplicado | Registrar con email existente | Mensaje de error "Email ya registrado" | Media |
| TC-PUB-08 | Consentimiento parental (menores) | Registrar usuario < 18 años | Campo de consentimiento parental visible | Alta |

#### Módulo: Landing Pública por Club
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-PUB-09 | Ver landing de club | Navegar a `/{clubSlug}` | Hero section, cards de torneos/eventos/tienda | Media |
| TC-PUB-10 | Ver noticias de club | Navegar a `/{clubSlug}` | Sección "Últimas Novedades" visible | Baja |
| TC-PUB-11 | Acceder a tienda pública | Click en "Ir a la Tienda" | Redirección a `/{clubSlug}/store` | Media |
| TC-PUB-12 | Acceder a registro desde landing | Click en "Hacerme Socio" | Redirección a `/register` | Media |
| TC-PUB-13 | Ver campeonatos públicos | Navegar a `/{clubSlug}/championships` | Lista de campeonatos públicos visible | Baja |

---

### 2.2. Rol: MEMBER (Socio/Jugador)

**Objetivo:** Verificar flujos completos del socio para gestión personal y reservas.

#### Módulo: Autenticación
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-01 | Login exitoso | Navegar a `/login` → Ingresar credenciales → Submit | Redirección a `/dashboard`, saludo visible | Crítica |
| TC-MEM-02 | Login con credenciales inválidas | Ingresar credenciales incorrectas | Mensaje de error "Credenciales inválidas" | Alta |
| TC-MEM-03 | Login con Google OAuth | Click en "Iniciar con Google" | Flujo OAuth, redirección a dashboard | Media |
| TC-MEM-04 | Logout | Click en Avatar → "Cerrar Sesión" | Redirección a `/login`, cookies limpiadas | Alta |
| TC-MEM-05 | Sesión expirada | Esperar expiración de JWT | Redirección automática a `/login` | Media |
| TC-MEM-06 | Recordar sesión | Login con "Recordarme" marcado | Sesión persistente tras cerrar navegador | Baja |

#### Módulo: Dashboard y Menú
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-07 | Ver Dashboard de Member | Login como MEMBER | Vista `MemberDashboardView` con métricas | Alta |
| TC-MEM-08 | Menú de navegación | Verificar sidebar | Items: Inicio, Reservar, Mis Reservas, Mi Perfil, Membresía, Tienda, Campeonatos | Alta |
| TC-MEM-09 | Navegación a cada item del menú | Click en cada opción del menú | Redirección correcta sin errores | Alta |
| TC-MEM-10 | Ver carnet digital | En dashboard | Carnet de socio con QR visible | Media |
| TC-MEM-11 | Ver próximas reservas | En dashboard | Lista de próximas reservas del usuario | Alta |
| TC-MEM-12 | Ver alerta de balance | En dashboard con saldo bajo | Alerta de balance visible | Media |
| TC-MEM-13 | Ver gamificación (nivel/XP) | En dashboard | Barra de progreso de nivel y XP | Media |

#### Módulo: Perfil (Profile)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-14 | Ver perfil | Navegar a `/profile` | Datos personales visibles (nombre, email, teléfono) | Alta |
| TC-MEM-15 | Editar perfil | Modificar campos → Guardar | Cambios guardados, toast de confirmación | Alta |
| TC-MEM-16 | Cambiar foto de perfil | Subir nueva imagen | Avatar actualizado | Media |
| TC-MEM-17 | Ver datos de GDPR | En perfil | Fechas de consentimiento visibles | Media |

#### Módulo: Membresía (Membership)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-18 | Ver estado de membresía | Navegar a `/membership` | Plan actual, estado, fechas de facturación | Alta |
| TC-MEM-19 | Ver planes disponibles (sin membresía) | Navegar sin membresía activa | Cards de precios con opciones de plan | Alta |
| TC-MEM-20 | Suscribirse a un plan | Click en plan → Confirmar | Membresía creada, toast de éxito | Alta |
| TC-MEM-21 | Ver historial de facturas | Click en "Ver Facturas" | Modal o página con historial | Baja |
| TC-MEM-22 | Cancelar membresía | Click en "Cancelar" → Confirmar | Estado cambiado a CANCELLED | Media |

#### Módulo: Reservas (Bookings)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-23 | Ver lista de reservas | Navegar a `/bookings` | Lista de reservas futuras y pasadas | Alta |
| TC-MEM-24 | Crear nueva reserva | Navegar a `/bookings/new` | Selector de fecha, instalación y horario | Crítica |
| TC-MEM-25 | Seleccionar instalación | En `/bookings/new` → Elegir cancha | Horarios disponibles mostrados | Alta |
| TC-MEM-26 | Seleccionar fecha y hora | Elegir slot disponible | Slot marcado como seleccionado | Alta |
| TC-MEM-27 | Confirmar reserva | Click en "Confirmar Reserva" | Reserva creada, redirección a `/bookings` | Crítica |
| TC-MEM-28 | Validar horario en el pasado | Intentar reservar hora pasada | Mensaje de error, slot no seleccionable | Alta |
| TC-MEM-29 | Validar conflicto de reserva | Intentar reservar slot ocupado | Mensaje de error "Slot no disponible" | Alta |
| TC-MEM-30 | Cancelar reserva | En lista → Click "Cancelar" → Confirmar | Reserva cancelada, toast de confirmación | Alta |
| TC-MEM-31 | Ver detalles de reserva | Click en reserva de la lista | Modal con detalles completos | Media |
| TC-MEM-32 | Filtrar reservas por estado | Usar filtros (Activas/Pasadas/Canceladas) | Lista filtrada correctamente | Baja |

#### Módulo: Tienda (Store)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-33 | Ver catálogo de productos | Navegar a `/store` | Grid de productos con precio y stock | Alta |
| TC-MEM-34 | Agregar al carrito | Click "Agregar" en producto | Toast "Agregado al carrito", contador actualizado | Alta |
| TC-MEM-35 | Ver carrito | Click en ícono de carrito | Lista de productos en carrito | Alta |
| TC-MEM-36 | Modificar cantidad en carrito | Cambiar cantidad | Total actualizado | Media |
| TC-MEM-37 | Eliminar del carrito | Click en eliminar | Producto removido | Media |
| TC-MEM-38 | Realizar compra | Click "Pagar Carrito" | Compra procesada, toast de éxito | Alta |
| TC-MEM-39 | Producto sin stock | Intentar agregar producto agotado | Botón deshabilitado o mensaje de error | Media |

#### Módulo: Campeonatos (Championships)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MEM-40 | Ver lista de torneos | Navegar a `/championships` | Lista de torneos activos del club | Alta |
| TC-MEM-41 | Ver tabla de posiciones | Seleccionar torneo | Tabla con standings ordenados | Alta |
| TC-MEM-42 | Ver fixture y resultados | En torneo | Lista de partidos con resultados | Alta |
| TC-MEM-43 | Inscribirse en equipo | Click "Inscribirse" | Confirmación de inscripción | Media |

---

### 2.3. Rol: COACH (Entrenador)

**Objetivo:** Verificar gestión de equipos, asistencia y viajes.

#### Módulo: Dashboard del Coach
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-COA-01 | Ver Dashboard de Coach | Navegar a `/coach` | Dashboard con métricas (jugadores, habilitados, viajes) | Crítica |
| TC-COA-02 | Ver estadísticas generales | En dashboard | Cards con totales: Jugadores, Habilitados, Inhabilitados, Viajes | Alta |
| TC-COA-03 | Selector de equipo | Si tiene múltiples equipos | Botones para cambiar entre equipos | Alta |
| TC-COA-04 | Ver tabs de navegación | En dashboard | Tabs: Jugadores, Asistencia, Viajes, Calendario, Estadísticas | Alta |

#### Módulo: Gestión de Jugadores
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-COA-05 | Ver lista de jugadores | Tab "Jugadores" | Tabla con jugadores del equipo | Crítica |
| TC-COA-06 | Ver estado de habilitación | En tabla de jugadores | Badge de estado (Habilitado/Inhabilitado) | Alta |
| TC-COA-07 | Ver detalle de jugador | Click en jugador | Modal con información detallada | Alta |
| TC-COA-08 | Buscar jugador | Usar buscador | Filtrado de tabla por nombre | Media |
| TC-COA-09 | Filtrar por estado | Usar filtro de estado | Solo jugadores con estado seleccionado | Media |
| TC-COA-10 | Ver razón de inhabilitación | Jugador inhabilitado | Razón visible (cuota, documentación, etc.) | Alta |
| TC-COA-11 | Exportar lista de jugadores | Click "Exportar" | Descarga de archivo CSV/Excel | Baja |

#### Módulo: Control de Asistencia
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-COA-12 | Ver tracker de asistencia | Tab "Asistencia" | Componente `AttendanceTracker` visible | Crítica |
| TC-COA-13 | Tomar asistencia | Marcar presente/ausente | Estado actualizado, guardado automático | Crítica |
| TC-COA-14 | Ver historial de asistencia | Seleccionar fecha pasada | Asistencia del día visible | Alta |
| TC-COA-15 | Agregar observación | En registro de asistencia | Nota guardada | Media |
| TC-COA-16 | Ver porcentaje de asistencia | Por jugador | Estadística calculada visible | Media |
| TC-COA-17 | Navegar a `/coach/attendance` | Acceso directo | Página de asistencia dedicada | Alta |

#### Módulo: Gestión de Viajes
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-COA-18 | Ver eventos de viaje | Tab "Viajes" | Componente `TravelEvents` con lista | Alta |
| TC-COA-19 | Crear nuevo viaje | Click "Nuevo Viaje" | Formulario de creación | Alta |
| TC-COA-20 | Definir destino y fecha | En formulario | Campos de destino, fecha, hora | Alta |
| TC-COA-21 | Agregar jugadores al viaje | Seleccionar jugadores | Lista de convocados | Alta |
| TC-COA-22 | Confirmar asistencia a viaje | Marcar jugadores confirmados | Estado de confirmación actualizado | Alta |
| TC-COA-23 | Ver calendario de viajes | Tab "Calendario" | Componente `TravelCalendar` | Media |
| TC-COA-24 | Editar viaje existente | Click en viaje → Editar | Formulario con datos cargados | Media |
| TC-COA-25 | Cancelar viaje | Click "Cancelar" → Confirmar | Viaje cancelado, notificaciones enviadas | Media |
| TC-COA-26 | Ver detalles de viaje | Click en viaje | Modal con información completa | Alta |

#### Módulo: Estadísticas del Equipo
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-COA-27 | Ver estadísticas | Tab "Estadísticas" | Card con mensaje "Próximamente" o gráficos | Baja |

---

### 2.4. Rol: ADMIN (Administrador de Club)

**Objetivo:** Verificar gestión completa del club y sus recursos.

#### Módulo: Dashboard Administrativo
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-01 | Ver Dashboard de Admin | Login como ADMIN → `/dashboard` | Vista `AdminDashboardView` con métricas | Crítica |
| TC-ADM-02 | Menú de navegación | Verificar sidebar | Items: Dashboard, Calendario, Instalaciones, Usuarios, Configuración | Alta |
| TC-ADM-03 | Ver métricas principales | En dashboard | Cards con ingresos, reservas, usuarios activos | Alta |
| TC-ADM-04 | Ver calendario resumido | En dashboard | Vista rápida de ocupación | Media |

#### Módulo: Gestión de Instalaciones
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-05 | Ver lista de instalaciones | Navegar a `/facilities` | Grid de cards con instalaciones del club | Alta |
| TC-ADM-06 | Crear nueva instalación | Click "Nueva Instalación" → Completar form | Instalación creada, visible en lista | Crítica |
| TC-ADM-07 | Ver detalles de instalación | Click en card de instalación | Información completa (capacidad, precio, horario) | Alta |
| TC-ADM-08 | Editar instalación | Click "Editar" → Modificar → Guardar | Cambios persistidos | Alta |
| TC-ADM-09 | Cambiar estado de instalación | Toggle activo/inactivo | Estado actualizado, badge cambiado | Alta |
| TC-ADM-10 | Eliminar instalación | Click "Eliminar" → Confirmar | Instalación removida de lista | Media |
| TC-ADM-11 | Reservar desde admin | Click "Reservar" en instalación | Modal `BookingModal` visible | Alta |

#### Módulo: Calendario Global de Reservas
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-12 | Ver calendario global | Navegar a `/bookings/calendar` | Vista de todas las reservas del club | Crítica |
| TC-ADM-13 | Filtrar por instalación | Usar selector de instalación | Solo reservas de instalación seleccionada | Alta |
| TC-ADM-14 | Ver detalle de reserva | Click en reserva del calendario | Modal con información del socio/hora | Alta |
| TC-ADM-15 | Crear reserva manual | Click en slot vacío | Formulario de reserva para cualquier socio | Alta |
| TC-ADM-16 | Cancelar reserva de otro usuario | Click "Cancelar" en reserva ajena | Confirmación, notificación al socio | Alta |
| TC-ADM-17 | Navegar entre fechas | Usar controles de fecha | Calendario actualizado | Alta |
| TC-ADM-18 | Vista diaria/semanal/mensual | Cambiar tipo de vista | Layout actualizado correctamente | Media |

#### Módulo: Reservas Recurrentes
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-19 | Ver reservas recurrentes | Navegar a `/admin/recurring-bookings` | Lista de reglas de reservas recurrentes | Alta |
| TC-ADM-20 | Crear regla recurrente | Click "Nueva Regla" → Completar | Regla creada, reservas generadas | Alta |
| TC-ADM-21 | Definir patrón de recurrencia | Seleccionar días/frecuencia | Patrón configurado (diario, semanal, mensual) | Alta |
| TC-ADM-22 | Editar regla recurrente | Click "Editar" → Modificar | Cambios aplicados a futuras reservas | Media |
| TC-ADM-23 | Eliminar regla recurrente | Click "Eliminar" → Confirmar | Regla eliminada, opción de cancelar futuras | Media |

#### Módulo: Reglas Recurrentes
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-24 | Ver reglas recurrentes | Navegar a `/admin/recurring-rules` | Tabla con reglas configuradas | Alta |
| TC-ADM-25 | Crear nueva regla | Click "Nueva Regla" | Formulario de configuración | Alta |
| TC-ADM-26 | Editar regla existente | Click en regla | Modal de edición | Media |
| TC-ADM-27 | Activar/Desactivar regla | Toggle de estado | Estado actualizado | Media |

#### Módulo: Gestión de Usuarios
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-28 | Ver lista de usuarios | Navegar a `/users` | Tabla con usuarios del club | Crítica |
| TC-ADM-29 | Buscar usuario | Usar barra de búsqueda | Tabla filtrada por nombre/email | Alta |
| TC-ADM-30 | Filtrar por rol | Usar selector de rol | Solo usuarios con rol seleccionado | Alta |
| TC-ADM-31 | Ver detalle de usuario | Click en usuario | Modal/página con información completa | Alta |
| TC-ADM-32 | Editar rol de usuario | Cambiar rol → Guardar | Rol actualizado | Alta |
| TC-ADM-33 | Desactivar usuario | Click "Desactivar" | Usuario inactivo, no puede loguearse | Media |
| TC-ADM-34 | Ver historial de usuario | En detalle de usuario | Reservas, pagos, asistencia | Media |
| TC-ADM-35 | Exportar usuarios | Click "Exportar" | Descarga CSV/Excel | Baja |

#### Módulo: Usuarios Familia
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-36 | Ver grupos familiares | Navegar a `/users/family` | Lista de grupos familiares | Media |
| TC-ADM-37 | Vincular miembros | Click "Agregar a familia" | Usuarios vinculados como grupo | Media |
| TC-ADM-38 | Ver miembros de familia | Click en grupo familiar | Lista de miembros relacionados | Media |

#### Módulo: Pagos (Payment)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-39 | Ver dashboard de pagos | Navegar a `/payment` | Resumen de ingresos, pendientes, historial | Crítica |
| TC-ADM-40 | Registrar pago offline | Click "Registrar Pago" → Completar | Pago registrado, balance actualizado | Crítica |
| TC-ADM-41 | Buscar socio para pago | Usar buscador en modal | Resultados de búsqueda visibles | Alta |
| TC-ADM-42 | Seleccionar método de pago | Elegir Efectivo/Transferencia/Canje | Método seleccionado | Alta |
| TC-ADM-43 | Ver historial de transacciones | En tabla de pagos | Lista con fecha, monto, método, estado | Alta |
| TC-ADM-44 | Procesar reembolso | Click "Reembolsar" → Confirmar | Pago reembolsado, estado REFUNDED | Alta |
| TC-ADM-45 | Filtrar pagos por estado | Usar filtros | Tabla filtrada (PENDING, COMPLETED, REFUNDED) | Media |
| TC-ADM-46 | Ver resultado de pago | Navegar a `/payment/result` | Página de confirmación de pago | Media |

#### Módulo: Configuración del Club
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-47 | Ver configuración | Navegar a `/settings` | Página de ajustes del club | Alta |
| TC-ADM-48 | Editar datos del club | Modificar nombre/dirección | Cambios guardados | Alta |
| TC-ADM-49 | Configurar horarios | Definir horario de apertura/cierre | Horarios actualizados | Alta |
| TC-ADM-50 | Configurar notificaciones | Toggle de tipos de notificación | Preferencias guardadas | Media |

#### Módulo: Campeonatos (Admin)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-51 | Crear nuevo torneo | En `/championships` → "Nuevo Torneo" | Wizard de creación visible | Alta |
| TC-ADM-52 | Configurar formato de torneo | En wizard | Opciones de liga/eliminación/grupos | Alta |
| TC-ADM-53 | Agregar equipos al torneo | Seleccionar equipos | Equipos agregados al torneo | Alta |
| TC-ADM-54 | Generar fixture | Click "Generar" | Partidos creados automáticamente | Alta |
| TC-ADM-55 | Cargar resultado de partido | Click en partido → Ingresar score | Resultado guardado, standings actualizados | Alta |
| TC-ADM-56 | Activar/Finalizar torneo | Cambiar estado | Estado del torneo actualizado | Alta |

#### Módulo: Control de Acceso (Kiosk)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ADM-57 | Ver punto de acceso | Navegar a `/access-control` | Pantalla de kiosk para escaneo | Alta |
| TC-ADM-58 | Escanear QR válido | Ingresar UUID de socio habilitado | "ACCESO PERMITIDO" en verde | Crítica |
| TC-ADM-59 | Escanear QR inválido | Ingresar UUID inexistente | "ACCESO DENEGADO" en rojo | Crítica |
| TC-ADM-60 | Escanear socio inhabilitado | Ingresar UUID de socio inactivo | "ACCESO DENEGADO" con razón | Alta |
| TC-ADM-61 | Auto-reset de pantalla | Después de 3 segundos | Pantalla vuelve a estado inicial | Media |

---

### 2.5. Rol: SUPER_ADMIN (Administrador de Plataforma)

**Objetivo:** Verificar gestión multi-tenant y configuración global.

#### Módulo: Dashboard Global
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-SUP-01 | Redirección automática | Login como SUPER_ADMIN | Redirección a `/admin/platform` | Crítica |
| TC-SUP-02 | Ver Dashboard Global | En `/admin/platform` | Métricas globales (MRR, clubes, usuarios, estado) | Crítica |
| TC-SUP-03 | Menú de navegación | Verificar sidebar | Items: Dashboard Global, Gestión Clubes, Configuración | Alta |
| TC-SUP-04 | Ver estado del sistema | Card "Estado Sistema" | Estado de health check (Operativo/Degradado) | Alta |
| TC-SUP-05 | Ver métricas financieras | Card "Total MRR" | Ingresos mensuales recurrentes | Alta |

#### Módulo: Gestión de Clubes (Tenants)
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-SUP-06 | Ver lista de clubes | Navegar a `/admin/clubs` | Tabla con todos los clubes de la plataforma | Crítica |
| TC-SUP-07 | Ver clubes en dashboard | En `/admin/platform` | Lista de tenants con status | Alta |
| TC-SUP-08 | Crear nuevo club | Click "Nuevo Club" → Completar | Club creado, visible en lista | Crítica |
| TC-SUP-09 | Definir dominio de club | En formulario de creación | Subdomain único asignado | Alta |
| TC-SUP-10 | Ver detalle de club | Click en "Gestionar" | Página con información del club | Alta |
| TC-SUP-11 | Editar información de club | Modificar datos → Guardar | Cambios persistidos | Alta |
| TC-SUP-12 | Cambiar estado de club | Activar/Desactivar | Estado actualizado (badge) | Alta |
| TC-SUP-13 | Ver usuarios de un club | En detalle de club | Lista de usuarios del tenant | Media |
| TC-SUP-14 | Asignar admin a club | Seleccionar usuario como admin | Rol actualizado | Alta |
| TC-SUP-15 | Eliminar club | Click "Eliminar" → Confirmar | Club y datos eliminados | Baja |

#### Módulo: Configuración Global
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-SUP-16 | Ver configuración global | Navegar a `/admin/settings` | Página de ajustes de sistema | Alta |
| TC-SUP-17 | Configurar planes de membresía | Definir tiers globales | Plans actualizados | Alta |
| TC-SUP-18 | Configurar límites de plataforma | Definir límites por plan | Límites aplicados | Media |
| TC-SUP-19 | Ver logs de sistema | Sección de auditoría | Logs de actividad recientes | Baja |
| TC-SUP-20 | Configurar integraciones | MercadoPago, Google OAuth | Credenciales guardadas | Media |

---

### 2.6. Rol: MEDICAL_STAFF (Personal Médico)

**Objetivo:** Verificar acceso a datos sensibles siguiendo GDPR Artículo 9.

#### Módulo: Acceso a Datos Médicos
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-MED-01 | Login como MEDICAL_STAFF | Navegar a `/login` → Credenciales | Acceso a dashboard con permisos especiales | Alta |
| TC-MED-02 | Ver datos de salud de jugadores | Acceder a ficha de jugador | Campos médicos visibles (alergias, condiciones) | Alta |
| TC-MED-03 | Registrar nota médica | Agregar observación | Nota guardada con timestamp | Alta |
| TC-MED-04 | Ver historial médico | En ficha de jugador | Historial de notas/eventos médicos | Alta |
| TC-MED-05 | Marcar apto/no apto | Cambiar estado de habilitación | Estado actualizado, notificación a coach | Alta |
| TC-MED-06 | Acceso denegado a datos financieros | Intentar acceder a `/payment` | Redirección o mensaje de acceso denegado | Alta |
| TC-MED-07 | Auditoría de acceso | Cada acceso a datos sensibles | Log registrado con fecha/usuario/acción | Alta |

---

## 3. Casos de Prueba Transversales

### 3.1. Manejo de Errores y Estados
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-ERR-01 | Página no encontrada (404) | Navegar a ruta inexistente | Página de error 404 con enlace a inicio | Alta |
| TC-ERR-02 | Error de servidor (500) | Simular falla de API | Página de error con mensaje amigable | Alta |
| TC-ERR-03 | Sin conexión a internet | Desconectar red | Mensaje "Sin conexión", retry disponible | Media |
| TC-ERR-04 | Timeout de API | Respuesta lenta (>30s) | Loader, luego mensaje de timeout | Media |
| TC-ERR-05 | Sesión inválida | Token expirado/inválido | Redirección a login con mensaje | Alta |

### 3.2. Responsividad y Accesibilidad
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-RES-01 | Vista móvil | Viewport 375px | Menú hamburguesa, layout adaptado | Alta |
| TC-RES-02 | Vista tablet | Viewport 768px | Layout responsive correcto | Media |
| TC-RES-03 | Vista desktop | Viewport 1440px | Sidebar visible, layout completo | Alta |
| TC-RES-04 | Navegación con teclado | Tab + Enter | Todos los elementos focuseables | Media |
| TC-RES-05 | Contraste de colores | Inspeccionar UI | Ratio mínimo 4.5:1 para texto | Media |

### 3.3. Notificaciones y Feedback
| ID | Caso de Prueba | Acción | Resultado Esperado | Prioridad |
|:---|:---|:---|:---|:---|
| TC-NOT-01 | Toast de éxito | Crear reserva exitosamente | Toast verde con mensaje | Alta |
| TC-NOT-02 | Toast de error | Error en operación | Toast rojo con mensaje descriptivo | Alta |
| TC-NOT-03 | Notificación en tiempo real | Recibir notificación push | Toast o badge en campana | Media |
| TC-NOT-04 | Ver historial de notificaciones | Click en campana | Lista de notificaciones recientes | Media |
| TC-NOT-05 | Marcar como leída | Click en notificación | Estado actualizado | Baja |

---

## 4. Matriz de Cobertura de Rutas

| Ruta | Roles | Casos de Prueba | Prioridad |
|:---|:---|:---|:---|
| `/login` | GUEST | TC-MEM-01, TC-MEM-02, TC-MEM-03 | Crítica |
| `/register-player` | GUEST | TC-PUB-03 a TC-PUB-08 | Alta |
| `/legal/terms` | GUEST | TC-PUB-01 | Alta |
| `/legal/privacy` | GUEST | TC-PUB-02 | Alta |
| `/{clubSlug}` | GUEST | TC-PUB-09 a TC-PUB-13 | Media |
| `/dashboard` | MEMBER, ADMIN | TC-MEM-07, TC-ADM-01 | Crítica |
| `/profile` | MEMBER | TC-MEM-14 a TC-MEM-17 | Alta |
| `/membership` | MEMBER | TC-MEM-18 a TC-MEM-22 | Alta |
| `/bookings` | MEMBER | TC-MEM-23 a TC-MEM-32 | Crítica |
| `/bookings/new` | MEMBER | TC-MEM-24 a TC-MEM-29 | Crítica |
| `/bookings/calendar` | ADMIN | TC-ADM-12 a TC-ADM-18 | Crítica |
| `/store` | MEMBER | TC-MEM-33 a TC-MEM-39 | Alta |
| `/championships` | MEMBER, ADMIN | TC-MEM-40 a TC-MEM-43, TC-ADM-51 a TC-ADM-56 | Alta |
| `/coach` | COACH | TC-COA-01 a TC-COA-27 | Crítica |
| `/coach/attendance` | COACH | TC-COA-12 a TC-COA-17 | Crítica |
| `/facilities` | ADMIN | TC-ADM-05 a TC-ADM-11 | Alta |
| `/users` | ADMIN | TC-ADM-28 a TC-ADM-35 | Alta |
| `/users/family` | ADMIN | TC-ADM-36 a TC-ADM-38 | Media |
| `/payment` | ADMIN | TC-ADM-39 a TC-ADM-46 | Crítica |
| `/settings` | ADMIN | TC-ADM-47 a TC-ADM-50 | Alta |
| `/access-control` | ADMIN | TC-ADM-57 a TC-ADM-61 | Alta |
| `/admin/platform` | SUPER_ADMIN | TC-SUP-01 a TC-SUP-05 | Crítica |
| `/admin/clubs` | SUPER_ADMIN | TC-SUP-06 a TC-SUP-15 | Crítica |
| `/admin/settings` | SUPER_ADMIN | TC-SUP-16 a TC-SUP-20 | Alta |
| `/admin/recurring-bookings` | ADMIN | TC-ADM-19 a TC-ADM-23 | Alta |
| `/admin/recurring-rules` | ADMIN | TC-ADM-24 a TC-ADM-27 | Alta |

---

## 5. Estrategia de Ejecución

### 5.1. Pruebas Automatizadas (Prioridad)

```bash
cd frontend
npx playwright test
```

**Archivos de test existentes:**
- `frontend/e2e/auth.spec.ts` - Flujos de autenticación
- `frontend/e2e/booking-flow.spec.ts` - Flujo de reservas

**Tests a crear (por prioridad):**
1. `coach-dashboard.spec.ts` - Flujos completos del rol COACH
2. `admin-payments.spec.ts` - Gestión de pagos
3. `super-admin.spec.ts` - Gestión multi-tenant
4. `access-control.spec.ts` - Kiosk de control de acceso
5. `championships.spec.ts` - Torneos y fixtures

### 5.2. Pruebas Manuales (Exploratorio)

1. **Validación Visual:** Verificar consistencia de menús según rol
2. **Responsividad:** Probar en dispositivos móviles reales
3. **Flujos Edge-Case:** Usuarios con múltiples roles, datos límite
4. **Accesibilidad:** Navegación con lector de pantalla

### 5.3. Datos de Prueba Requeridos

| Rol | Email de Prueba | Contraseña |
|:---|:---|:---|
| SUPER_ADMIN | `superadmin@clubpulse.com` | `SuperAdmin123!` |
| ADMIN | `admin@clubpulse.com` | `Admin123!` |
| COACH | `coach@clubpulse.com` | `Coach123!` |
| MEMBER | `member@clubpulse.com` | `Member123!` |
| MEDICAL_STAFF | `medical@clubpulse.com` | `Medical123!` |

---

## 6. Resumen de Métricas

| Categoría | Cantidad |
|:---|:---|
| **Total de Casos de Prueba** | 132 |
| **Prioridad Crítica** | 25 |
| **Prioridad Alta** | 67 |
| **Prioridad Media** | 31 |
| **Prioridad Baja** | 9 |
| **Roles Cubiertos** | 6 |
| **Rutas Únicas** | 27 |
| **Módulos Funcionales** | 18 |

---

*Última actualización: 2026-01-10*
*Versión: 2.0*
