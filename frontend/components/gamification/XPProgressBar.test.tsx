import { render, screen } from '@testing-library/react';
import { XPProgressBar, XPProgressCompact } from './XPProgressBar';

// Mock Progress component
jest.mock('@/components/ui/progress', () => ({
    Progress: ({ value, className }: { value: number; className?: string }) => (
        <div data-testid="progress-bar" role="progressbar" aria-valuenow={value} className={className} />
    )
}));

// Mock hooks
jest.mock('@/hooks/useGamification', () => ({
    useStreakStatus: (streak: number) => {
        if (streak >= 7) return { emoji: '🔥', label: 'Racha Épica', multiplier: 1.5, color: 'text-orange-500' };
        if (streak >= 3) return { emoji: '⚡', label: 'En Racha', multiplier: 1.2, color: 'text-yellow-500' };
        return { emoji: '', label: '', multiplier: 1, color: '' };
    }
}));

describe('XPProgressBar', () => {
    const defaultProps = {
        level: 5,
        currentXP: 500,
        requiredXP: 1000,
        totalXP: 4500,
        currentStreak: 0
    };

    it('renders level and xp info correctly', () => {
        render(<XPProgressBar {...defaultProps} />);

        expect(screen.getByText('Nivel 5')).toBeInTheDocument();
        expect(screen.getByText('500 / 1,000 XP')).toBeInTheDocument();
        expect(screen.getByText(/Total: 4,500 XP/)).toBeInTheDocument();

        // Progress should be 50%
        const progressBar = screen.getByTestId('progress-bar');
        expect(progressBar).toHaveAttribute('aria-valuenow', '50');
    });

    it('displays streak info when streak >= 3', () => {
        render(<XPProgressBar {...defaultProps} currentStreak={5} />);

        // Mock hook returns 'En Racha' for streak 5 (>=3 but <7)
        expect(screen.getByText(/En Racha/)).toBeInTheDocument();
        expect(screen.getByText(/x1.2 XP/)).toBeInTheDocument();
        expect(screen.getByText('⚡')).toBeInTheDocument();
    });

    it('hides details when showDetails is false', () => {
        render(<XPProgressBar {...defaultProps} showDetails={false} />);

        expect(screen.queryByText(/Total:/)).not.toBeInTheDocument();
    });

    it('calculates progress cap at 100%', () => {
        render(<XPProgressBar {...defaultProps} currentXP={1500} />);

        const progressBar = screen.getByTestId('progress-bar');
        expect(progressBar).toHaveAttribute('aria-valuenow', '100');
    });
});

describe('XPProgressCompact', () => {
    it('renders level and progress bar', () => {
        render(<XPProgressCompact level={10} progress={75} />);

        expect(screen.getByText('10')).toBeInTheDocument();
        const progressBar = screen.getByTestId('progress-bar');
        expect(progressBar).toHaveAttribute('aria-valuenow', '75');
    });
});
