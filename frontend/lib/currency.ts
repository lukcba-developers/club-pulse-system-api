/**
 * Currency utilities for handling Argentine Peso (ARS) amounts.
 * Uses string representation to preserve decimal precision from backend.
 */

/**
 * Formats a decimal string value as Argentine Peso currency.
 * @param value - The decimal value as a string (e.g., "1500.50")
 * @returns Formatted currency string (e.g., "$ 1.500,50")
 */
export function formatARS(value: string | number | null | undefined): string {
    if (value === null || value === undefined || value === '') {
        return '$ 0,00';
    }

    const numValue = typeof value === 'string' ? parseFloat(value) : value;

    if (isNaN(numValue)) {
        return '$ 0,00';
    }

    return new Intl.NumberFormat('es-AR', {
        style: 'currency',
        currency: 'ARS',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    }).format(numValue);
}

/**
 * Parses a currency string or decimal string to a number for calculations.
 * Use with caution - prefer string operations for financial data.
 * @param value - The value to parse
 * @returns The numeric value
 */
export function parseARS(value: string | number | null | undefined): number {
    if (value === null || value === undefined || value === '') {
        return 0;
    }

    if (typeof value === 'number') {
        return value;
    }

    // Remove currency symbols and formatting
    const cleaned = value
        .replace(/[$.]/g, '')
        .replace(',', '.')
        .trim();

    const result = parseFloat(cleaned);
    return isNaN(result) ? 0 : result;
}

/**
 * Compares two decimal string values.
 * @returns -1 if a < b, 0 if equal, 1 if a > b
 */
export function compareARS(a: string, b: string): number {
    const numA = parseARS(a);
    const numB = parseARS(b);
    if (numA < numB) return -1;
    if (numA > numB) return 1;
    return 0;
}
