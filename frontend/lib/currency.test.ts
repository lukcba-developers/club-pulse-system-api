import { formatARS, parseARS, compareARS } from './currency';

describe('Currency Utilities', () => {
    describe('formatARS', () => {
        it('formats numbers correctly as ARS', () => {
            const formatted = formatARS(1500.5);
            expect(formatted.replace(/\s/g, ' ')).toContain('$ 1.500,50'); // Normalize spaces
            expect(formatARS(100).replace(/\s/g, ' ')).toContain('$ 100,00');
        });

        it('formats strings correctly', () => {
            expect(formatARS('2500.00').replace(/\s/g, ' ')).toContain('$ 2.500,00');
        });

        it('handles null/undefined/empty', () => {
            expect(formatARS(null).replace(/\s/g, ' ')).toContain('$ 0,00');
        });
    });

    describe('parseARS', () => {
        it('parses formatted strings', () => {
            // ARS uses comma for decimal separator
            expect(parseARS('$ 1.500,50')).toBe(1500.5);
        });

        it('parses plain strings (European format assumption)', () => {
            // The utility assumes strings are ARS formatted (dot=thousands, comma=decimal)
            // So '1500,50' is 1500.5
            expect(parseARS('1500,50')).toBe(1500.5);
        });

        it('handles numbers directly', () => {
            expect(parseARS(123.45)).toBe(123.45);
        });

        it('handles invalid inputs', () => {
            expect(parseARS(null)).toBe(0);
            expect(parseARS('abc')).toBe(0);
        });
    });

    describe('compareARS', () => {
        it('compares correctly', () => {
            expect(compareARS('100.00', '200.00')).toBe(-1);
            expect(compareARS('200.00', '100.00')).toBe(1);
            expect(compareARS('$ 1.000,00', '$ 1.000,00')).toBe(0);
        });
    });
});
