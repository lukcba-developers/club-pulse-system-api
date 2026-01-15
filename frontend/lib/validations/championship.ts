import { z } from 'zod';

/**
 * Tournament Status Enum
 * Matches backend domain.TournamentStatus exactly
 */
export const TOURNAMENT_STATUS = {
    DRAFT: 'DRAFT',
    ACTIVE: 'ACTIVE',
    COMPLETED: 'COMPLETED',
    CANCELLED: 'CANCELLED'
} as const;

export type TournamentStatus = typeof TOURNAMENT_STATUS[keyof typeof TOURNAMENT_STATUS];

/**
 * Zod Schema for Tournament Status Validation
 * Use this in forms and API validation
 */
export const TournamentStatusSchema = z.enum(['DRAFT', 'ACTIVE', 'COMPLETED', 'CANCELLED']);

/**
 * Zod Schema for Match Status Validation
 */
export const MatchStatusSchema = z.enum(['SCHEDULED', 'COMPLETED', 'CANCELLED']);

/**
 * Zod Schema for Tournament Creation/Update
 */
export const TournamentFormSchema = z.object({
    name: z.string().min(1, 'El nombre es requerido').max(255),
    description: z.string().optional(),
    sport: z.string().min(1, 'El deporte es requerido'),
    category: z.string().optional(),
    status: TournamentStatusSchema.default('DRAFT'),
    start_date: z.string().datetime(), // ISO 8601 format
    end_date: z.string().datetime().optional(),
    settings: z.record(z.string(), z.unknown()).optional(), // JSONB flexible
});

export type TournamentFormData = z.infer<typeof TournamentFormSchema>;

/**
 * Zod Schema for Match Result Update
 */
export const MatchResultSchema = z.object({
    match_id: z.string().uuid(),
    home_score: z.number().min(0, 'El puntaje no puede ser negativo'),
    away_score: z.number().min(0, 'El puntaje no puede ser negativo'),
});

export type MatchResultData = z.infer<typeof MatchResultSchema>;

/**
 * Helper: Validate tournament status before API call
 */
export function validateTournamentStatus(status: unknown): TournamentStatus {
    return TournamentStatusSchema.parse(status);
}

/**
 * Helper: Safe parse for form validation (returns errors instead of throwing)
 */
export function safeParseTournamentForm(data: unknown) {
    return TournamentFormSchema.safeParse(data);
}
