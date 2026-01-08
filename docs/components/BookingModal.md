# 🎫 BookingModal

El componente `BookingModal` es el orquestador principal del proceso de reserva desde la perspectiva del usuario (frontend). Utiliza un diálogo modal para recopilar toda la información necesaria y enviarla al servidor.

## 🚀 Propósito
Ser el punto de entrada unificado para realizar reservas, integrando la selección de horarios, la gestión de invitados y la validación de reglas de negocio en una sola interfaz cohesiva.

## ⚙️ Props

| Prop | Tipo | Descripción | Obligatorio |
| :--- | :--- | :--- | :--- |
| `isOpen` | `boolean` | Controla la visibilidad del diálogo. | Sí |
| `onClose` | `() => void` | Función para cerrar el modal y limpiar estados internos. | Sí |
| `facilityId` | `string` | ID de la instalación a reservar. | Sí |
| `facilityName` | `string` | Nombre comercial de la instalación (ej. "Cancha de Tenis 1"). | Sí |

## 💡 Flujo de Trabajo
1. **Selección de Horario:** El usuario navega por el calendario interno (`AvailabilityCalendar`) y selecciona un slot.
2. **Invitados (Opcional):** Se habilita una sección para cargar los datos de un invitado si el usuario marca la casilla.
3. **Validación:** Verifica si el usuario está autenticado y si se ha elegido un horario.
4. **Reserva:** Envía un `POST` a `/bookings`.
5. **Confirmación:** Muestra una vista de éxito temporal antes de cerrarse automáticamente.

## 🛠️ Ejemplo de Implementación

```tsx
import { useState } from 'react';
import { BookingModal } from '@/components/booking-modal';
import { Button } from '@/components/ui/button';

export default function FacilityCard() {
  const [modalOpen, setModalOpen] = useState(false);

  return (
    <>
      <Button onClick={() => setModalOpen(true)}>Reservar Ahora</Button>
      
      <BookingModal 
        isOpen={modalOpen} 
        onClose={() => setModalOpen(false)} 
        facilityId="uuid-123" 
        facilityName="Pádel Pro 4" 
      />
    </>
  );
}
```

## 🧩 Integraciones
- **Auth:** Consume el hook `useAuth` para obtener el ID del usuario actual.
- **API:** Se comunica con el backend mediante `lib/axios`.
- **UI:** Basado en componentes de `Radix UI` via `shadcn/ui` (Dialog).

## ⚠️ Notas Técnicas
- **Duración Fija:** Actualmente, el componente asume turnos de 1 hora exacta.
- **Manejo de Conflictos:** Detecta errores `409 Conflict` del servidor para informar al usuario si alguien más reservó el turno mientras completaba el formulario.
- **Reset de Estado:** El formulario se limpia automáticamente después de un éxito para evitar duplicados.

⚠️ **Deuda Técnica:** El costo del invitado está hardcodeado en `$1500`. Debería obtenerse dinámicamente desde la configuración de la instalación si se desea flexibilidad.
