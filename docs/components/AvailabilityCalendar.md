# 📅 AvailabilityCalendar

El componente `AvailabilityCalendar` proporciona una interfaz interactiva para visualizar y seleccionar turnos horarios disponibles para una instalación deportiva.

## 🚀 Propósito
Facilitar al usuario la elección de un horario específico para realizar una reserva, mostrando de forma clara qué slots están libres, cuáles ocupados y cuáles están bajo mantenimiento.

## ⚙️ Props

| Prop | Tipo | Descripción | Obligatorio |
| :--- | :--- | :--- | :--- |
| `facilityId` | `string` | Identificador único de la instalación (cancha, salón, etc.). | Sí |
| `onSlotSelect` | `(date: string, time: string) => void` | Callback que se ejecuta cuando el usuario selecciona un turno válido. Retorna la fecha (`YYYY-MM-DD`) y la hora (`HH:mm`). | Sí |

## 💡 Casos de Uso
- **Pantalla de Reserva Directa:** Cuando un socio desea alquilar una cancha específica.
- **Dashboard de Administrador:** Para visualizar rápidamente la ocupación diaria de una instalación.
- **Widget de Disponibilidad Rápida:** En listas de instalaciones para mostrar disponibilidad sin navegar fuera de la página.

## 🛠️ Ejemplo de Implementación

```tsx
import { AvailabilityCalendar } from '@/components/availability-calendar';

export default function BookingPage() {
  const handleSelect = (date: string, time: string) => {
    console.log(`Turno seleccionado: ${date} a las ${time}`);
    // Abrir modal de confirmación o enviar a la API
  };

  return (
    <div className="max-w-md mx-auto p-4 border rounded-xl shadow-lg">
      <h2 className="text-xl font-bold mb-4">Selecciona tu horario</h2>
      <AvailabilityCalendar 
        facilityId="uuid-cancha-1" 
        onSlotSelect={handleSelect} 
      />
    </div>
  );
}
```

## 🧩 Sub-componentes
- `DateNavigator`: Gestiona la navegación entre fechas.
- `TimeSlotButton`: Botón individual por cada slot horario con manejo de estados (`available`, `booked`, `maintenance`).
- `Legend`: Guía visual para que el usuario entienda el significado de los colores y marcas.

## ⚠️ Notas Técnicas
- **Consumo de API:** Realiza llamadas al endpoint `/bookings/availability`.
- **Performance:** Utiliza `useMemo` para la generación de la grilla horaria y `useCallback` para la función de refresco, minimizando re-renders innecesarios.
- **Responsive:** La grilla se ajusta automáticamente (layout de 3 columnas) y tiene scroll interno para no romper el layout del padre en pantallas pequeñas.
