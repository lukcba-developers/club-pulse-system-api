import api from '../lib/axios';

export interface Discipline {
    id: string;
    name: string;
    description: string;
    is_active: boolean;
}

export interface TrainingGroup {
    id: string;
    name: string;
    discipline_id: string;
    discipline?: Discipline;
    category: string;
    coach_id: string;
    schedule: string;
}

export interface Tournament {
    id: string;
    club_id: string;
    name: string;
    discipline_id: string;
    start_date: string;
    end_date: string;
    status: string; // OPEN, IN_PROGRESS, COMPLETED
    format: string; // LEAGUE, KNOCKOUT
}

export interface Match {
    id: string;
    tournament_id: string;
    home_team_id: string;
    away_team_id: string;
    score_home: number;
    score_away: number;
    start_time: string;
    location: string;
    status: string; // SCHEDULED, PLAYED
}

export interface Standing {
    team_id: string;
    points: number;
    won: number;
    lost: number;
    drawn: number;
    goals_for: number;
    goals_against: number;
}

export const disciplineService = {
    listDisciplines: async () => {
        const response = await api.get<Discipline[]>('/disciplines');
        return response.data;
    },

    listGroups: async (disciplineId?: string, category?: string) => {
        const params = new URLSearchParams();
        if (disciplineId) params.append('discipline_id', disciplineId);
        if (category) params.append('category', category);

        const response = await api.get<TrainingGroup[]>(`/groups?${params.toString()}`);
        return response.data;
    },

    getGroupStudents: async (groupId: string) => {
        const response = await api.get(`/groups/${groupId}/students`);
        return response.data;
    },

    // --- Championships ---

    listTournaments: async () => {
        const response = await api.get<Tournament[]>('/tournaments');
        return response.data;
    },

    createTournament: async (data: Partial<Tournament>) => {
        const response = await api.post<Tournament>('/tournaments', data);
        return response.data;
    },

    registerTeam: async (tournamentId: string, name: string, captainId?: string, memberIds: string[] = []) => {
        const response = await api.post(`/tournaments/${tournamentId}/teams`, {
            name,
            captain_id: captainId,
            member_ids: memberIds
        });
        return response.data;
    },

    scheduleMatch: async (tournamentId: string, data: unknown) => {
        const response = await api.post(`/tournaments/${tournamentId}/matches`, data);
        return response.data;
    },

    listMatches: async (tournamentId: string) => {
        const response = await api.get<Match[]>(`/tournaments/${tournamentId}/matches`);
        return response.data;
    },

    updateMatchResult: async (matchId: string, home: number, away: number) => {
        const response = await api.put(`/matches/${matchId}/result`, {
            score_home: home,
            score_away: away
        });
        return response.data;
    },

    getStandings: async (tournamentId: string) => {
        const response = await api.get<Standing[]>(`/tournaments/${tournamentId}/standings`);
        return response.data;
    }
};
