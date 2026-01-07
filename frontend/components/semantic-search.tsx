'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import { Search, X, Loader2, MapPin, Users, Sparkles } from 'lucide-react';
import { facilityService, SearchResult } from '@/services/facility-service';

interface SemanticSearchProps {
    onResultSelect?: (facilityId: string) => void;
    placeholder?: string;
}

export function SemanticSearch({
    onResultSelect,
    placeholder = "Buscar: canchas techadas, piscina nocturna..."
}: SemanticSearchProps) {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState<SearchResult[]>([]);
    const [loading, setLoading] = useState(false);
    const [isOpen, setIsOpen] = useState(false);
    const [error, setError] = useState('');
    const debounceRef = useRef<NodeJS.Timeout | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);

    // Debounced search
    const handleSearch = useCallback(async (searchQuery: string) => {
        if (searchQuery.length < 2) {
            setResults([]);
            setIsOpen(false);
            return;
        }

        setLoading(true);
        setError('');

        try {
            const response = await facilityService.search(searchQuery, 5);
            setResults(response.results);
            setIsOpen(true);
        } catch (err) {
            console.error('Search failed:', err);
            setError('Error buscando instalaciones');
            setResults([]);
        } finally {
            setLoading(false);
        }
    }, []);

    // Debounce input changes
    useEffect(() => {
        if (debounceRef.current) {
            clearTimeout(debounceRef.current);
        }

        debounceRef.current = setTimeout(() => {
            handleSearch(query);
        }, 300);

        return () => {
            if (debounceRef.current) {
                clearTimeout(debounceRef.current);
            }
        };
    }, [query, handleSearch]);

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
                setIsOpen(false);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleResultClick = (facilityId: string) => {
        setIsOpen(false);
        setQuery('');
        onResultSelect?.(facilityId);
    };

    const clearSearch = () => {
        setQuery('');
        setResults([]);
        setIsOpen(false);
    };

    return (
        <div ref={containerRef} className="relative w-full max-w-2xl">
            {/* Search Input */}
            <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                    {loading ? (
                        <Loader2 className="h-5 w-5 text-brand-500 animate-spin" />
                    ) : (
                        <Search className="h-5 w-5 text-gray-400" />
                    )}
                </div>
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onFocus={() => results.length > 0 && setIsOpen(true)}
                    placeholder={placeholder}
                    className="w-full pl-12 pr-12 py-3 bg-white dark:bg-zinc-900 border border-gray-200 dark:border-zinc-700 rounded-xl shadow-sm focus:ring-2 focus:ring-brand-500 focus:border-transparent transition-all duration-200 text-sm placeholder:text-gray-400"
                />
                {query && (
                    <button
                        onClick={clearSearch}
                        className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                    >
                        <X className="h-5 w-5" />
                    </button>
                )}
            </div>

            {/* AI Badge */}
            <div className="absolute -top-2 -right-2">
                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-gradient-to-r from-purple-500 to-pink-500 text-white shadow-lg">
                    <Sparkles className="h-3 w-3" />
                    IA
                </span>
            </div>

            {/* Results Dropdown */}
            {isOpen && (
                <div className="absolute z-50 w-full mt-2 bg-white dark:bg-zinc-900 rounded-xl shadow-xl border border-gray-200 dark:border-zinc-700 overflow-hidden">
                    {error ? (
                        <div className="p-4 text-sm text-red-500">{error}</div>
                    ) : results.length === 0 ? (
                        <div className="p-4 text-sm text-gray-500 text-center">
                            No se encontraron resultados para &quot;{query}&quot;
                        </div>
                    ) : (
                        <ul className="divide-y divide-gray-100 dark:divide-zinc-800">
                            {results.map((result) => (
                                <li
                                    key={result.facility.id}
                                    onClick={() => handleResultClick(result.facility.id)}
                                    className="p-4 hover:bg-gray-50 dark:hover:bg-zinc-800 cursor-pointer transition-colors"
                                >
                                    <div className="flex justify-between items-start">
                                        <div className="flex-1">
                                            <h4 className="font-medium text-gray-900 dark:text-gray-100">
                                                {result.facility.name}
                                            </h4>
                                            <div className="flex items-center gap-3 mt-1 text-xs text-gray-500">
                                                <span className="inline-flex items-center gap-1">
                                                    <MapPin className="h-3 w-3" />
                                                    {result.facility.location.name}
                                                </span>
                                                <span className="inline-flex items-center gap-1">
                                                    <Users className="h-3 w-3" />
                                                    {result.facility.capacity}
                                                </span>
                                                <span className="capitalize px-2 py-0.5 rounded-full bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400">
                                                    {result.facility.type}
                                                </span>
                                            </div>
                                            {/* Specifications badges */}
                                            <div className="flex flex-wrap gap-1 mt-2">
                                                {result.facility.specifications.covered && (
                                                    <span className="px-2 py-0.5 text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 rounded-full">
                                                        Techada
                                                    </span>
                                                )}
                                                {result.facility.specifications.lighting && (
                                                    <span className="px-2 py-0.5 text-xs bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400 rounded-full">
                                                        Iluminación
                                                    </span>
                                                )}
                                                {result.facility.specifications.surface_type && (
                                                    <span className="px-2 py-0.5 text-xs bg-gray-100 dark:bg-zinc-700 text-gray-600 dark:text-gray-300 rounded-full">
                                                        {result.facility.specifications.surface_type}
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                        <div className="text-right ml-4">
                                            <span className="text-lg font-bold text-brand-600">
                                                ${result.facility.hourly_rate}
                                            </span>
                                            <span className="text-xs text-gray-500 block">/hora</span>
                                        </div>
                                    </div>
                                </li>
                            ))}
                        </ul>
                    )}

                    {/* Footer */}
                    <div className="px-4 py-2 bg-gray-50 dark:bg-zinc-800 text-xs text-gray-500 text-center">
                        Búsqueda semántica potenciada por IA
                    </div>
                </div>
            )}
        </div>
    );
}
