# Documentación del Proyecto: Club Pulse API

Este directorio contiene la documentación estratégica y técnica del proyecto. A continuación se describe el propósito de cada documento para facilitar su lectura.

## 📚 Índice de Documentación

### 1. Auditoría y Estrategia (Lectura Recomendada para Gerencia/Producto)
> 📂 Ubicación: `./strategy/`

- **[SYSTEM_AUDIT_AND_ROADMAP.md](./strategy/SYSTEM_AUDIT_AND_ROADMAP.md)** (⭐ Nuevo)
  - **Qué es:** Auditoría de alto nivel comparando el sistema actual vs. un Club de Alto Rendimiento.
  - **Qué contiene:** Matriz de Priorización (Impacto vs. Esfuerzo) y la Hoja de Ruta de 3 Meses para resultados económicos inmediatos.

- **[SOCIAL_CLUB_VISION_ROLE_PLAY.md](./strategy/SOCIAL_CLUB_VISION_ROLE_PLAY.md)**
  - **Qué es:** Simulación de discusión de producto.
  - **Qué contiene:** La visión pivotando hacia un "Club Social" (Fútbol Infantil, Cuotas), definiendo por qué y para qué de los cambios.

### 2. Planificación Técnica (Lectura para Desarrolladores)
> 📂 Ubicación: `./technical/`

- **[POST_MVP_MIGRATION_PLAN.md](./technical/POST_MVP_MIGRATION_PLAN.md)** (Plan Maestro)
  - **Qué es:** El plan de ejecución técnica detallado fase por fase.
  - **Qué contiene:** Tareas específicas, referencias al código legado, y desglose de módulos (Auth, User, Membership).

- **[GAP_ANALYSIS.md](./technical/GAP_ANALYSIS.md)** (Referencia)
  - **Qué es:** Comparativa técnica "feature vs feature".
  - **Qué contiene:** Tablas detalladas de qué falta migrar del sistema legado (`club-management-system-api`) al nuevo.

### 3. Archivo Histórico
> 📂 Ubicación: `./archive/`

- **[MVP_ANALYSIS_AND_PLAN.md](./archive/MVP_ANALYSIS_AND_PLAN.md)** (Histórico)
  - **Qué es:** El análisis inicial para construir el MVP.
  - **Estado:** ✅ Completado. Mantener como referencia histórica.

---

## 🚀 Resumen de Prioridades (Q1 2026)

Según la [Auditoría Estratégica](./strategy/SYSTEM_AUDIT_AND_ROADMAP.md), nuestro foco inmediato es:

1.  **Caja (Cashflow):** Integrar Pagos para cobro de deudas.
2.  **Organización:** Implementar lógica de Categorías (Año nacimiento) para deportes.
3.  **Fidelización:** App/Dashboard para que el socio vea su estado al día.
