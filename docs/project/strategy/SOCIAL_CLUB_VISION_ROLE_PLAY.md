# Definición de Visión: Plataforma de Gestión de Club Social

**Participantes:**
- **Product Owner (PO):** Define la visión de negocio (Club Social, Fútbol Infantil, Organización).
- **Tech Lead / Arquitecto (TL):** Alinea la visión con la arquitectura del sistema legado.
- **Frontend Lead (FE):** Define la experiencia de usuario para padres y administradores.
- **Backend Lead (BE):** Valida la viabilidad técnica y los modelos de datos.

---

## Acta de Reunión de Definición de Producto

### 1. El Problema: "Caos en Fútbol Infantil"
**PO:** "El problema principal es la desorganización. Los coordinadores de fútbol infantil manejan planillas de excel. No saben quién pagó la cuota, quién debe, o a qué categoría pertenece el chico (2012, 2013). Necesitamos que el sistema organice automáticamente por año de nacimiento."

**TL:** "Entiendo. El sistema actual (`club-management-system-api`) ya tiene `DateOfBirth` en el usuario. Podemos automatizar la asignación de categoría."

### 2. Solución Propuesta: Gestión por Categorías y Cuota Social

**BE:** "Analicé el modelo legado.
- Tenemos `User.DateOfBirth`.
- Tenemos `Membership` con lógica de ciclos de facturación (`Monthly`, `Quarterly`).
- Lo que FALTA es el concepto de 'Equipo' o 'División' automática. En el legado, `Championship` maneja equipos, pero es para torneos. Necesitamos un concepto de **'Grupo de Entrenamiento'** basado en el año."

**Propuesta Técnica (Backend):**
1.  **Categorías Dinámicas:** Crear una lógica que calcule la categoría del usuario al vuelo basada en `DateOfBirth` (ej. si nació en 2012 -> Categoría 2012).
2.  **Relación Padre-Hijo:** Necesitamos modelar que un usuario (Padre) paga por otro (Hijo). Actualmente el `User` es individual.
    - *Decisión:* Para la Fase 1, el usuario registrado es el socio (el chico o el adulto). Si es menor, los datos de contacto son del padre. En Fase 2 implementamos "Grupos Familiares".
3.  **Control de Morosos:** El sistema de `Membership` ya tiene `CalculateLateFee` y estado `Suspended`. Usaremos eso para bloquear asistencia si no hay pago.

### 3. Experiencia de Usuario (Frontend)

**FE:** "Para los padres, la interfaz debe ser simple:
- **Dashboard:** 'Estado de Cuota' (Al día / Vencida) bien grande.
- **Perfil:** Ver la categoría asignada (ej. 'Pre-Novena 2013').
- **Pagos:** Botón 'Pagar Cuota' integrado con MercadoPago."

**FE:** "Para el Coordinador/Profe:
- **Lista de Asistencia:** Filtrada por Categoría (Año).
- **Indicador de Deuda:** Al lado del nombre del chico, un ícono rojo si debe la cuota."

---

## Plan de Acción Ajustado (Social Club Pivot)

1.  **Refinamiento de Entidad Usuario (`User`):**
    - Asegurar que `DateOfBirth` sea obligatorio para socios deportivos.
    - Añadir campo virtual `Category` (calculado).

2.  **Módulo de Disciplinas (`Disciplines` - Nuevo/Refactor):**
    - En lugar de solo 'Facilities' (Canchas), necesitamos 'Actividades'.
    - *Migración:* Adaptar `Championship` o crear `TrainingGroups` para manejar las listas de jugadores por año.

3.  **Cobro de Cuotas (Prioridad Alta):**
    - Activar la lógica de `Membership` legacy inmediatamente (Fase 2 del Plan General).
    - Implementar el recargo por mora (`CalculateLateFee`).

**Conclusión:** El sistema dejará de ser solo un "booking de canchas" para convertirse en un **ERP de Club Social**.
