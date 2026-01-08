# 🔍 SemanticSearch

El componente `SemanticSearch` eleva la experiencia de descubrimiento en la plataforma, permitiendo a los usuarios realizar búsquedas en lenguaje natural potenciadas por IA.

## 🚀 Propósito
Permitir que los usuarios encuentren instalaciones no solo por nombre, sino por características semánticas como "canchas con techo para lluvia" o "lugares con equipo para 12 personas".

## ⚙️ Props

| Prop | Tipo | Descripción | Obligatorio |
| :--- | :--- | :--- | :--- |
| `onResultSelect` | `(facilityId: string) => void` | Callback que se dispara al seleccionar una instalación de la lista de resultados. | No |
| `placeholder` | `string` | Texto de ayuda en el input. | No |

## 💡 Características Premium
- **AI-Powered:** Identificado visualmente con un badge de gradiente para indicar que la búsqueda es semántica.
- **Auto-Debounce:** Retrasa las llamadas a la API 300ms tras el último teclado para optimizar el tráfico de red.
- **Smart Cards:** Los resultados muestran badges dinámicos con especificaciones técnicas (`Techada`, `Iluminación`, `Superficie`).
- **Click Outside:** El menú de resultados se cierra automáticamente al perder el foco, mejorando la usabilidad.

## 🛠️ Ejemplo de Implementación

```tsx
import { SemanticSearch } from '@/components/semantic-search';
import { useRouter } from 'next/navigation';

export default function Home() {
  const router = useRouter();

  const handleSelect = (id: string) => {
    router.push(`/facilities/${id}`);
  };

  return (
    <div className="py-20 flex flex-col items-center">
      <h1 className="text-4xl font-bold mb-8">¿Qué quieres jugar hoy?</h1>
      <SemanticSearch onResultSelect={handleSelect} />
    </div>
  );
}
```

## 🧩 Integraciones
- **FacilityService:** Utiliza el método `search` del servicio de instalaciones.
- **Lucide Icons:** Iconografía moderna para tipos de datos (Mapa, Usuarios, Destellos).
- **Tailwind CSS:** Diseño responsivo y modo oscuro integrado.

## ⚠️ Notas Técnicas
- **Longitud Mínima:** Solo dispara búsquedas si el texto tiene al menos 2 caracteres.
- **Ref Handling:** Usa `useRef` para manejar el temporizador del debounce y el contenedor para la detección de clics externos.
- **Accesibilidad:** Soporta navegación básica y cierre con el botón `X`.

⚠️ **Deuda Técnica:** Los resultados de búsqueda están limitados a 5 por defecto en el componente. Se recomienda pasar el límite como una prop si se planea usar en diferentes áreas con necesidades de espacio distintas.
