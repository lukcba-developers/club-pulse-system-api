import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { LoginForm } from './login-form';
import { useAuth } from '@/hooks/use-auth';
import { useRouter } from 'next/navigation';
import api from '@/lib/axios';

// Mocks
jest.mock('@/hooks/use-auth');
jest.mock('next/navigation', () => ({
    useRouter: jest.fn(),
}));
jest.mock('@/lib/axios');

describe('LoginForm', () => {
    const mockLogin = jest.fn();
    const mockPush = jest.fn();

    beforeEach(() => {
        (useAuth as jest.Mock).mockReturnValue({ login: mockLogin });
        (useRouter as jest.Mock).mockReturnValue({ push: mockPush });
        jest.clearAllMocks();
    });

    it('renders login form correctly', () => {
        render(<LoginForm />);
        expect(screen.getByLabelText(/correo electrónico/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/contraseña/i)).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /iniciar sesión/i })).toBeInTheDocument();
    });

    it('submits valid credentials and calls login', async () => {
        // Mock successful API response
        (api.post as jest.Mock).mockResolvedValue({ data: { success: true } });

        render(<LoginForm />);

        fireEvent.change(screen.getByLabelText(/correo electrónico/i), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));

        await waitFor(() => {
            expect(api.post).toHaveBeenCalledWith('/auth/login', {
                email: 'test@example.com',
                password: 'password123',
            });
            expect(mockLogin).toHaveBeenCalled();
            expect(mockPush).toHaveBeenCalledWith('/');
        });
    });

    it('displays error message on failure', async () => {
        (api.post as jest.Mock).mockRejectedValue({
            response: { data: { error: 'Credenciales inválidas' } }
        });

        render(<LoginForm />);

        fireEvent.change(screen.getByLabelText(/correo electrónico/i), { target: { value: 'wrong@example.com' } });
        fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'wrongpass' } });
        fireEvent.click(screen.getByRole('button', { name: /iniciar sesión/i }));

        await waitFor(() => {
            expect(screen.getByText('Credenciales inválidas')).toBeInTheDocument();
        });
    });

    it('shows loading state while submitting', async () => {
        // Promise that resolves after a delay
        (api.post as jest.Mock).mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));

        render(<LoginForm />);

        fireEvent.change(screen.getByLabelText(/correo electrónico/i), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText(/contraseña/i), { target: { value: 'pass' } });
        const submitButton = screen.getByRole('button', { name: /iniciar sesión/i });
        fireEvent.click(submitButton);

        expect(submitButton).toBeDisabled();
    });
});
