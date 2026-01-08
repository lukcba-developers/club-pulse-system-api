# 📄 DocumentUpload

El componente `DocumentUpload` es fundamental para la gestión administrativa delegada al socio, permitiendo la carga y visualización del estado de sus documentos legales y médicos.

## 🚀 Propósito
Habilitar un flujo de autoservicio para que los socios cumplan con los requisitos documentales (ej. apto médico EMMAC) necesarios para estar habilitados en las actividades del club.

## ⚙️ Props

| Prop | Tipo | Descripción | Obligatorio |
| :--- | :--- | :--- | :--- |
| `userId` | `string` | ID del usuario al que pertenecen los documentos. | Sí |
| `documents` | `UserDocument[]` | Array con los documentos actuales del usuario para mostrar su estado e historial. | Sí |
| `onUploadSuccess` | `() => void` | Callback opcional que se ejecuta tras una subida exitosa (para refrescar la lista padre). | No |

## 🧩 Tipos de Documentos Soportados
- `DNI_FRONT` / `DNI_BACK`: Documento de identidad nacional.
- `EMMAC_MEDICAL`: Apto médico obligatorio para deportes competitivos.
- `LEAGUE_FORM`: Autorización para participar en ligas externas.
- `INSURANCE`: Póliza de seguro personal o de viaje.

## 🚦 Estados del Documento
- ⏳ `PENDING`: Subido, esperando revisión manual del administrador.
- ✅ `VALID`: Revisado y aprobado.
- ❌ `REJECTED`: Rechazado (muestra notas del motivo del rechazo).
- ⚠️ `EXPIRED`: Documento cuya fecha de vencimiento ha pasado.

## 🛠️ Ejemplo de Implementación

```tsx
import { DocumentUpload } from '@/components/DocumentUpload';

export default function ProfilePage({ user, userDocs }) {
  const handleRefresh = () => {
    // Lógica para recargar documentos desde la API
  };

  return (
    <DocumentUpload 
      userId={user.id} 
      documents={userDocs} 
      onUploadSuccess={handleRefresh}
    />
  );
}
```

## ⚠️ Notas Técnicas
- **Multipart Form Data:** Envía archivos binarios utilizando `FormData`.
- **Autorización:** Incluye el Bearer Token manualmente desde `localStorage` para la llamada `fetch`.
- **UI:** Utiliza `Lucide Icons` para feedback visual y componentes de `shadcn/ui` para la estructura visual (Cards, Selects).
- **Vencimiento:** La fecha de vencimiento solo se solicita para tipos de documentos que lo requieren (EMMAC y Seguros).

⚠️ **Deuda Técnica:** La carga de archivos utiliza `fetch` directo en lugar de la instancia de `api` (Axios) configurada globalmente, lo que duplica la lógica de headers y gestión de URLs base. Se recomienda refactorizar para usar `api.post`.
