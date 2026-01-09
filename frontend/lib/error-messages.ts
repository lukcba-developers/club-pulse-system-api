/**
 * Diccionario de errores humanizados para UX Writing
 * Transforma códigos de error técnicos en mensajes amigables para el usuario
 */

export const ERROR_MESSAGES: Record<string, string> = {
    // ═══════════════════════════════════════════════════════════════
    // Errores de Autenticación
    // ═══════════════════════════════════════════════════════════════
    'invalid_credentials': 'El correo o la contraseña son incorrectos. Verificá los datos e intentá de nuevo.',
    'user_not_found': 'No encontramos una cuenta con ese correo. ¿Querés registrarte?',
    'account_locked': 'Tu cuenta está temporalmente bloqueada por seguridad. Intentá de nuevo en 15 minutos.',
    'session_expired': 'Tu sesión expiró. Por favor, volvé a iniciar sesión.',
    'unauthorized': 'No tenés permiso para realizar esta acción.',
    'token_invalid': 'Tu sesión ya no es válida. Por favor, volvé a iniciar sesión.',

    // ═══════════════════════════════════════════════════════════════
    // Errores de Reserva (Booking)
    // ═══════════════════════════════════════════════════════════════
    'medical_certificate_invalid': 'Tu certificado médico no está al día. Actualizalo desde tu perfil para poder reservar.',
    'medical_certificate_expired': 'Tu certificado médico expiró. Subí uno nuevo desde tu perfil.',
    'facility_inactive': 'Esta instalación no está disponible actualmente.',
    'cancel_unauthorized': 'Solo podés cancelar tus propias reservas.',
    'slot_unavailable': '¡Alguien reservó antes! Ese horario ya no está disponible.',
    'booking_conflict': 'Ese horario acaba de ser reservado. Probá con otro, ¡hay más opciones disponibles!',
    'insufficient_balance': 'Tu saldo no alcanza para esta reserva. Podés cargar saldo desde tu perfil.',
    'membership_required': 'Necesitás una membresía activa para reservar. Consultá los planes disponibles.',
    'facility_closed': 'La instalación está cerrada en ese horario. Revisá los horarios de operación.',
    'maintenance_scheduled': 'Hay mantenimiento programado para esa fecha. Elegí otro día.',
    'past_date': 'No podés reservar para fechas pasadas. Elegí una fecha futura.',
    'too_far_ahead': 'Las reservas se pueden hacer con hasta 14 días de anticipación.',

    // ═══════════════════════════════════════════════════════════════
    // Errores de Registro
    // ═══════════════════════════════════════════════════════════════
    'email_already_exists': 'Ese correo ya está registrado. ¿Querés iniciar sesión?',
    'invalid_email': 'El formato del correo electrónico no es válido.',
    'password_too_weak': 'La contraseña debe tener al menos 8 caracteres, una mayúscula y un número.',
    'consent_required': 'Debes aceptar los términos y condiciones para continuar.',
    'parental_consent_required': 'Se requiere consentimiento parental para registrar menores.',
    'invalid_club_id': 'El club especificado no existe. Verificá el enlace de registro.',

    // ═══════════════════════════════════════════════════════════════
    // Errores de Membresía
    // ═══════════════════════════════════════════════════════════════
    'membership_expired': 'Tu membresía expiró. Renovála para seguir disfrutando de los beneficios.',
    'payment_failed': 'El pago no se pudo procesar. Verificá los datos de tu tarjeta.',
    'subscription_cancelled': 'Tu suscripción fue cancelada. Contactá al club para más información.',

    // ═══════════════════════════════════════════════════════════════
    // Errores de Red y Servidor
    // ═══════════════════════════════════════════════════════════════
    'network_error': 'Parece que hay un problema de conexión. Verificá tu internet e intentá de nuevo.',
    'server_error': 'Algo falló de nuestro lado. Ya estamos trabajando para solucionarlo.',
    'timeout': 'La solicitud tardó demasiado. Intentá de nuevo.',
    'rate_limited': 'Estás haciendo demasiadas solicitudes. Esperá un momento e intentá de nuevo.',

    // ═══════════════════════════════════════════════════════════════
    // Errores de Validación
    // ═══════════════════════════════════════════════════════════════
    'required_field': 'Este campo es obligatorio.',
    'invalid_format': 'El formato no es válido. Revisá los datos ingresados.',
    'value_too_long': 'El texto es demasiado largo.',
    'invalid_date': 'La fecha ingresada no es válida.',

    // ═══════════════════════════════════════════════════════════════
    // Error por defecto
    // ═══════════════════════════════════════════════════════════════
    'default': 'Ocurrió un error inesperado. Por favor, intentá de nuevo.',
};

/**
 * Convierte un código o mensaje de error en un mensaje humanizado
 * @param errorCodeOrMessage - Código de error (ej: 'invalid_credentials') o mensaje técnico
 * @returns Mensaje amigable para el usuario
 */
export function humanizeError(errorCodeOrMessage: string): string {
    // Primero intentamos buscar por código exacto
    if (ERROR_MESSAGES[errorCodeOrMessage]) {
        return ERROR_MESSAGES[errorCodeOrMessage];
    }

    // Luego buscamos por código en lowercase/snake_case normalizado
    const normalizedCode = errorCodeOrMessage
        .toLowerCase()
        .replace(/\s+/g, '_')
        .replace(/-/g, '_');

    if (ERROR_MESSAGES[normalizedCode]) {
        return ERROR_MESSAGES[normalizedCode];
    }

    // Buscamos patrones conocidos en el mensaje
    const lowerMessage = errorCodeOrMessage.toLowerCase();

    if (lowerMessage.includes('conflict') || lowerMessage.includes('409')) {
        return ERROR_MESSAGES['booking_conflict'];
    }
    if (lowerMessage.includes('unauthorized') || lowerMessage.includes('401')) {
        return ERROR_MESSAGES['unauthorized'];
    }
    if (lowerMessage.includes('forbidden') || lowerMessage.includes('403')) {
        return ERROR_MESSAGES['unauthorized'];
    }
    if (lowerMessage.includes('not found') || lowerMessage.includes('404')) {
        return 'El recurso que buscás no existe o fue eliminado.';
    }
    if (lowerMessage.includes('network') || lowerMessage.includes('fetch')) {
        return ERROR_MESSAGES['network_error'];
    }
    if (lowerMessage.includes('timeout')) {
        return ERROR_MESSAGES['timeout'];
    }
    if (lowerMessage.includes('email') && lowerMessage.includes('exist')) {
        return ERROR_MESSAGES['email_already_exists'];
    }

    // Si el mensaje original está en español y es legible, lo devolvemos
    if (/^[A-ZÁÉÍÓÚÑ¡¿]/.test(errorCodeOrMessage) && errorCodeOrMessage.length > 10) {
        return errorCodeOrMessage;
    }

    // Fallback al mensaje por defecto
    return ERROR_MESSAGES['default'];
}

/**
 * Tipo de error para componentes de alerta
 */
export type ErrorSeverity = 'error' | 'warning' | 'info';

/**
 * Obtiene la severidad de un error basado en su tipo
 */
export function getErrorSeverity(errorCode: string): ErrorSeverity {
    const warnings = ['session_expired', 'membership_expired', 'too_far_ahead'];
    const info = ['slot_unavailable', 'facility_closed', 'maintenance_scheduled'];

    if (warnings.includes(errorCode)) return 'warning';
    if (info.includes(errorCode)) return 'info';
    return 'error';
}

/**
 * Extrae el mensaje humanizado de un error de Axios
 * Utiliza el campo 'errorType' inyectado por el interceptor de axios
 */
export function getErrorFromAxios(error: unknown): string {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const axiosError = error as any;

    // Primero intentamos usar el tipo de error extraído por el interceptor
    if (axiosError?.errorType) {
        return humanizeError(axiosError.errorType);
    }

    // Fallback a extraer de response.data
    if (axiosError?.response?.data?.type) {
        return humanizeError(axiosError.response.data.type);
    }

    // Fallback al mensaje de error directo
    if (axiosError?.response?.data?.error) {
        return humanizeError(axiosError.response.data.error);
    }

    // Último fallback
    if (axiosError?.message) {
        return humanizeError(axiosError.message);
    }

    return ERROR_MESSAGES['default'];
}

