import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export function getAuthHeader(): Record<string, string> {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    return token ? { Authorization: `Bearer ${token}` } : {};
}

/**
 * Formats a monetary amount for display.
 * Backend sends decimal.Decimal as string (e.g., "150.00").
 * This helper normalizes and formats it with currency symbol.
 * 
 * @param amount - Amount as string or number from API
 * @param currency - Currency code (default: ARS)
 * @returns Formatted string like "$150.00"
 */
export function formatMoney(amount: string | number, currency: string = 'ARS'): string {
    const numericAmount = typeof amount === 'string' ? parseFloat(amount) : amount;

    if (isNaN(numericAmount)) {
        return '$0.00';
    }

    // Use Intl for proper locale formatting
    return new Intl.NumberFormat('es-AR', {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    }).format(numericAmount);
}
