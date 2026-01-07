# Módulo de Instalaciones (Facilities)

Este módulo gestiona el inventario digital de todos los espacios físicos y recursos que el club ofrece a sus socios. Es la base sobre la cual operan otros módulos como el de Reservas.

## 🌟 Funcionalidades Implementadas

### 1. Catálogo de Instalaciones
-   **Gestión de Espacios:** Permite a los administradores dar de alta y configurar todas las instalaciones del club (ej: "Cancha de Pádel 1", "Piscina Olímpica").
-   **Atributos:** Cada instalación tiene propiedades como nombre, tipo, capacidad y ubicación.

### 2. Gestión de Estados Operativos
-   **Control de Disponibilidad:** Los administradores pueden definir el estado de una instalación para controlar si está disponible para reservas.
-   **Estados Soportados:**
    -   `Disponible`: Operativa y abierta para reservas.
    -   `En Mantenimiento`: Bloqueada temporalmente, no se puede reservar.
    -   `Cerrada`: Fuera de servicio por un periodo prolongado.

### 3. Configuración de Horarios
-   **Horarios de Funcionamiento:** Se puede definir un horario de apertura y cierre para cada instalación, que puede ser independiente del horario general del club. Esta información es crucial para el motor de disponibilidad.

## 4. Funcionalidades en Desarrollo

-   **Búsqueda Semántica (Vector Search):** La capacidad de buscar instalaciones usando lenguaje natural (ej: "cancha techada para jugar de noche") está prevista mediante el uso de `pgvector` pero no está completamente integrada.
-   **Tarifas Flexibles:** La configuración de precios dinámicos por franja horaria ("hora pico") es parte del roadmap.
-   **Gestión de Equipamiento:** Un inventario detallado del equipamiento asociado a cada instalación es una mejora futura.
