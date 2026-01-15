import api from '@/lib/axios';
import {
    TOURNAMENT_STATUS,
    type TournamentStatus,
    TournamentStatusSchema,
    MatchStatusSchema,
    MatchResultSchema
} from '@/lib/validations/championship';

// Re-export for backward compatibility
export { TOURNAMENT_STATUS, TournamentStatusSchema, MatchStatusSchema };
export type { TournamentStatus };

export interface Tournament {
    id: string;
    name: string;
    description?: string;
    sport: string;
    category?: string;
    status: TournamentStatus;
    start_date: string;
    end_date?: string;
    stages?: TournamentStage[];
    club_id?: string;
    settings?: TournamentSettings;
}

export interface TournamentSettings {
    points_win: number;
    points_draw: number;
    points_loss: number;
    tiebreaker_criteria?: ('GOAL_DIFF' | 'GOALS_FOR' | 'HEAD_TO_HEAD')[];
    max_players_per_team?: number;
    [key: string]: unknown; // Allow flexibility but enforce core scoring
}

export interface TournamentStage {
    id: string;
    name: string;
    type: 'GROUP' | 'KNOCKOUT';
    order: number;
    groups?: Group[];
}

export interface Group {
    id: string;
    name: string;
    standings?: Standing[];
}

export interface Standing {
    id: string;
    team_id: string;
    points: number;
    played: number;
    won: number;
    drawn: number;
    lost: number;
    goals_for: number;
    goals_against: number;
    goal_difference: number;
    position: number;
    team_name?: string;
}

export interface Match {
    id: string;
    home_team_id: string;
    away_team_id: string;
    home_team_name?: string;
    away_team_name?: string;
    home_score: number | null;
    away_score: number | null;
    status: 'SCHEDULED' | 'COMPLETED' | 'CANCELLED';
    date: string;
    booking_id?: string;
}

export interface RequestTeam {
    name: string;
    logo_url?: string;
    contact?: string;
}

export interface Team {
    id: string;
    name: string;
    logo_url?: string;
    contact?: string;
    created_at: string;
}

export const championshipService = {
    async createTeam(data: RequestTeam): Promise<Team> {
        const response = await api.post('/championships/teams', data);
        return response.data;
    },

    async addMember(teamId: string, userId: string): Promise<void> {
        await api.post(`/championships/teams/${teamId}/members`, { user_id: userId });
    },

    async listTournaments(clubId?: string): Promise<Tournament[]> {
        const query = clubId ? `?club_id=${clubId}` : '';
        const response = await api.get(`/championships/${query}`);
        return response.data;
    },

    async createTournament(data: Omit<Tournament, "id">): Promise<Tournament> {
        const response = await api.post('/championships/', data);
        return response.data;
    },

    async addStage(tournamentId: string, data: { name: string; type: string; order: number }): Promise<TournamentStage> {
        const response = await api.post(`/championships/${tournamentId}/stages`, data);
        return response.data;
    },

    async addGroup(stageId: string, data: { name: string }): Promise<Group> {
        const response = await api.post(`/championships/stages/${stageId}/groups`, data);
        return response.data;
    },

    async registerTeam(groupId: string, teamId: string): Promise<Standing> {
        const response = await api.post(`/championships/groups/${groupId}/teams`, { team_id: teamId });
        return response.data;
    },

    async generateFixture(groupId: string): Promise<Match[]> {
        const response = await api.post(`/championships/groups/${groupId}/fixture`);
        return response.data;
    },

    async getMatches(groupId: string): Promise<Match[]> {
        const response = await api.get(`/championships/groups/${groupId}/matches`);
        return response.data;
    },

    async getStandings(groupId: string): Promise<Standing[]> {
        const response = await api.get(`/championships/groups/${groupId}/standings`);
        return response.data;
    },

    async updateMatchResult(matchId: string, homeScore: number, awayScore: number): Promise<void> {
        // Validate with Zod schema for type safety
        const validatedData = MatchResultSchema.parse({
            match_id: matchId,
            home_score: homeScore,
            away_score: awayScore
        });

        await api.post('/championships/matches/result', {
            match_id: validatedData.match_id,
            home_score: validatedData.home_score,
            away_score: validatedData.away_score
        });
    },

    async scheduleMatch(data: { club_id: string; match_id: string; court_id: string; start_time: string; end_time: string }): Promise<void> {
        await api.post('/championships/matches/schedule', data);
    },

    async getMyMatches(): Promise<Match[]> {
        const response = await api.get('/championships/my-matches');
        return response.data;
    },

    async getHeadToHead(groupId: string, teamAId: string, teamBId: string): Promise<HeadToHeadResult> {
        const response = await api.get(`/championships/groups/${groupId}/head-to-head`, {
            params: { team_a: teamAId, team_b: teamBId }
        });
        return response.data;
    }
};

export interface HeadToHeadResult {
    team_a_id: string;
    team_b_id: string;
    team_a_wins: number;
    team_b_wins: number;
    draws: number;
    team_a_goals: number;
    team_b_goals: number;
    matches: Match[];
}


