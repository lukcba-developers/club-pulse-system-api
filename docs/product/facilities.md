# Módulo de Instalaciones (Facilities)

El núcleo del sistema para administrar espacios físicos y recursos, permitiendo una representación digital fiel de la infraestructura del club.

## 🌟 Funcionalidades Principales

### 1. Inventario de Espacios
Soporte flexible para múltiples tipologías de instalaciones. Cada instalación es una entidad independiente con reglas propias.
-   **Tipos Soportados**: Canchas de Tenis (Polvo/Rápida), Padel (Cristal/Muro), Fútbol (5/7/11), Gimnasios, Piscinas, Salones de Usos Múltiples.
-   **Metadata**: Capacidad de etiquetar instalaciones (ej. "Outdoor", "Climatizada", "Iluminación LED").

### 2. Gestión de Estados Operativos
Control total sobre la disponibilidad de los activos.
-   **✅ Activo**: La instalación está operativa y listada en el motor de reservas.
-   **🛠️ Mantenimiento**: Bloqueo temporal.
    -   *Efecto*: Impide nuevas reservas durante el periodo designado.
    -   *Automatización*: Puede disparar alertas o cancelaciones si se solapa con reservas existentes (configurable).
-   **⛔ Clausurado**: Fuera de servicio indefinidamente (ej. reformas mayores).

### 3. Tarificación Flexible (Pricing)
-   **Hourly Rate**: Configuración de tarifa base por hora.
-   **Override**: Capacidad de ajustar precios para slots específicos (ej. "Hora Pico" vs "Hora Valle" - *Roadmap*).

### 4. Búsqueda Semántica (Vector Search)
Implementación avanzada utilizando **PostgreSQL + pgvector**.
-   **Caso de Uso**: Un usuario busca *"cancha techada para jugar de noche barata"*.
-   **Funcionamiento**: El sistema interpreta la intención ("techada", "noche" -> iluminación, "barata" -> precio bajo) y devuelve las mejores coincidencias ordenadas por relevancia, no solo por coincidencia de texto exacto.

### 5. Gestión de Equipamiento (Equipment)
Inventario de ítems físicos asociados a las instalaciones.
-   **Relación**: Trazabilidad de qué equipamiento pertenece a qué instalación.
-   **Estados**: `Nuevo`, `Usado`, `Dañado`, `En Reparación`.
-   **Uso**: Permite bloquear equipamiento si está dañado, afectando la disponibilidad de la instalación asociada si es crítico (ej. Red de tenis rota).
