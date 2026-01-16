import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { FacilityCard } from './facility-card';
import { Facility } from '@/services/facility-service';

// Mock child components dynamically loaded
jest.mock('@/components/booking-modal', () => ({
    BookingModal: ({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) => (
        isOpen ? <div data-testid="booking-modal">
            <button onClick={onClose}>Close Booking</button>
        </div> : null
    )
}));

jest.mock('@/components/facility-schedule-modal', () => ({
    FacilityScheduleModal: ({ isOpen, onClose }: { isOpen: boolean; onClose: () => void }) => (
        isOpen ? <div data-testid="schedule-modal">
            <button onClick={onClose}>Close Schedule</button>
        </div> : null
    )
}));

// Mock Next.js Image
jest.mock('next/image', () => ({
    __esModule: true,
    default: (props: any) => <img {...props} alt={props.alt} />,
}));

describe('FacilityCard', () => {
    const mockFacility: Facility = {
        id: '123',
        name: 'Tennis Court 1',
        type: 'court',
        description: 'Best court',
        capacity: 4,
        hourly_rate: 50,
        status: 'active',
        location: {
            name: 'Central Club',
            description: 'Main Hall',
            address: '123 Main St',
            latitude: 0,
            longitude: 0
        },
        opening_time: '09:00',
        closing_time: '21:00',
        images: [],
        club_id: 'club-1',
        created_at: '',
        updated_at: ''
    };

    it('renders facility details correctly', () => {
        render(<FacilityCard facility={mockFacility} />);

        expect(screen.getByText('Tennis Court 1')).toBeInTheDocument();
        // Check capacity
        expect(screen.getByText(/4 Jugadores/)).toBeInTheDocument();
        // Check pricing
        expect(screen.getByText(/\$50\/h/)).toBeInTheDocument();
        // Check location
        expect(screen.getByText(/Central Club - Main Hall/)).toBeInTheDocument();
        // Check schedule
        expect(screen.getByText('09:00 - 21:00')).toBeInTheDocument();
        // Check status badge
        expect(screen.getByText('Disponible')).toBeInTheDocument();
    });

    it('displays "Ocupado" status when inactive', () => {
        const inactiveFacility = { ...mockFacility, status: 'inactive' };
        render(<FacilityCard facility={inactiveFacility} />);
        expect(screen.getByText('Ocupado')).toBeInTheDocument();
    });

    it('opens booking modal when clicking Reserve button', async () => {
        render(<FacilityCard facility={mockFacility} />);

        const reserveButton = screen.getAllByText('Reservar')[0]; // There are duplicates in the UI (mobile/desktop usually, or icon+text)
        // Actually code has: 
        // 1. Text button "Ver Disponibilidad"
        // 2. Button with CalendarPlus "Reservar" inside a flex-1 div
        // And inside BookingModal component usage.

        // Wait, the code has:
        // button > Ver Disponibilidad (onClick opens)
        // button > Reservar (onClick opens)

        fireEvent.click(reserveButton);

        await waitFor(() => {
            expect(screen.getByTestId('booking-modal')).toBeInTheDocument();
        });

        // Close it
        fireEvent.click(screen.getByText('Close Booking'));
        await waitFor(() => {
            expect(screen.queryByTestId('booking-modal')).not.toBeInTheDocument();
        });
    });

    it('opens schedule modal when clicking Edit button', async () => {
        render(<FacilityCard facility={mockFacility} />);

        const editButton = screen.getByText('Editar');
        fireEvent.click(editButton);

        await waitFor(() => {
            expect(screen.getByTestId('schedule-modal')).toBeInTheDocument();
        });
    });
});
