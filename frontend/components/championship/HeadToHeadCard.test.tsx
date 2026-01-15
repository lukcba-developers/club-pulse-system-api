import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { HeadToHeadCard } from './HeadToHeadCard';
import { championshipService } from '@/services/championship-service';

// Mock the service
jest.mock('@/services/championship-service');

const mockData = {
    team_a_wins: 2,
    team_b_wins: 1,
    draws: 1,
    team_a_goals: 5,
    team_b_goals: 3,
    matches: [
        {
            id: 'm1',
            home_score: 2,
            away_score: 1,
            date: '2024-01-01T00:00:00Z',
            match_day: 1
        },
        {
            id: 'm2',
            home_score: 1,
            away_score: 1,
            date: '2024-01-08T00:00:00Z',
            match_day: 2
        }
    ]
};

describe('HeadToHeadCard', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders loading state initially', () => {
        (championshipService.getHeadToHead as jest.Mock).mockReturnValue(new Promise(() => { }));
        render(<HeadToHeadCard groupId="g1" teamAId="t1" teamBId="t2" />);
        // Look for the spinner or container, since Loader2 might not have aria-label by default in lucide
        // But the card header shouldn't be valid yet if loading covers it?
        // Actually the code returns early if loading.
        // We can check if `Enfrentamientos Directos` is NOT there yet or check for SVG
        const spinner = document.querySelector('.animate-spin');
        expect(spinner).toBeInTheDocument();
    });

    it('renders error state on failure', async () => {
        (championshipService.getHeadToHead as jest.Mock).mockRejectedValue(new Error('Failed'));
        render(<HeadToHeadCard groupId="g1" teamAId="t1" teamBId="t2" />);

        await waitFor(() => {
            expect(screen.getByText('No se pudo cargar el historial de enfrentamientos')).toBeInTheDocument();
        });
    });

    it('renders data correctly when loaded', async () => {
        (championshipService.getHeadToHead as jest.Mock).mockResolvedValue(mockData);

        render(
            <HeadToHeadCard
                groupId="g1"
                teamAId="t1"
                teamBId="t2"
                teamAName="Alpha"
                teamBName="Beta"
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Alpha vs Beta')).toBeInTheDocument();
        });

        // Check Stats
        // Check Stats - Using getAllByText since numbers might appear in match history too
        expect(screen.getAllByText('2').length).toBeGreaterThanOrEqual(1); // A Wins
        expect(screen.getAllByText('1').length).toBeGreaterThanOrEqual(1); // B Wins & Draws
        expect(screen.getByText('Empates')).toBeInTheDocument();

        // Check Goals
        expect(screen.getAllByText('5').length).toBeGreaterThanOrEqual(1);
        expect(screen.getAllByText('3').length).toBeGreaterThanOrEqual(1);

        // Check Match List
        expect(screen.getByText('Historial (4 partidos)')).toBeInTheDocument(); // 2+1+1
    });

    it('renders empty state when no matches', async () => {
        (championshipService.getHeadToHead as jest.Mock).mockResolvedValue({
            team_a_wins: 0,
            team_b_wins: 0,
            draws: 0,
            team_a_goals: 0,
            team_b_goals: 0,
            matches: []
        });

        render(<HeadToHeadCard groupId="g1" teamAId="t1" teamBId="t2" />);

        await waitFor(() => {
            expect(screen.getByText('No hay partidos jugados entre estos equipos')).toBeInTheDocument();
        });
    });
});
