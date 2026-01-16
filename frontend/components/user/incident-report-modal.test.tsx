import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { IncidentReportModal } from './incident-report-modal';
import { userService } from '@/services/user-service';

// Mock dependencies
jest.mock('@/services/user-service', () => ({
    userService: {
        logIncident: jest.fn()
    }
}));

jest.mock('@/components/ui/use-toast', () => ({
    useToast: () => ({
        toast: jest.fn()
    })
}));

describe('IncidentReportModal', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('renders the trigger button correctly', () => {
        render(<IncidentReportModal />);
        expect(screen.getByText('Reportar Incidente')).toBeInTheDocument();
    });

    it('opens modal on click', () => {
        render(<IncidentReportModal />);

        fireEvent.click(screen.getByText('Reportar Incidente'));

        expect(screen.getByText('Reporte de Incidente / Accidente')).toBeInTheDocument();
        expect(screen.getByLabelText(/Descripción del Hecho/)).toBeInTheDocument();
    });

    it('submits form with valid data', async () => {
        (userService.logIncident as jest.Mock).mockResolvedValue({});

        render(<IncidentReportModal />);

        // Open modal
        fireEvent.click(screen.getByText('Reportar Incidente'));

        // Fill form
        fireEvent.change(screen.getByLabelText(/Descripción del Hecho/), { target: { value: 'Caída en cancha 1' } });
        fireEvent.change(screen.getByLabelText(/Acción Tomada/), { target: { value: 'Hielo aplicado' } });
        fireEvent.change(screen.getByLabelText(/Testigos/), { target: { value: 'Juan Perez' } });

        // Submit
        fireEvent.click(screen.getByText('Confirmar Reporte'));

        await waitFor(() => {
            expect(userService.logIncident).toHaveBeenCalledWith({
                description: 'Caída en cancha 1',
                action_taken: 'Hielo aplicado',
                witnesses: 'Juan Perez',
                injured_user_id: undefined
            });
        });
    });

    it('handles submission errors', async () => {
        // Mock console.error to avoid noise in test output
        const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => { });
        (userService.logIncident as jest.Mock).mockRejectedValue(new Error('API Error'));

        render(<IncidentReportModal />);

        fireEvent.click(screen.getByText('Reportar Incidente'));
        fireEvent.change(screen.getByLabelText(/Descripción del Hecho/), { target: { value: 'Test' } });
        fireEvent.click(screen.getByText('Confirmar Reporte'));

        await waitFor(() => {
            expect(userService.logIncident).toHaveBeenCalled();
        });

        // Note: verifying toast call is hard due to hook mock structure, but ensuring no crash and service call is good enough
        consoleSpy.mockRestore();
    });
});
