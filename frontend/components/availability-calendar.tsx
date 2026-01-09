'use client';

import { useState, useCallback, useMemo, useEffect } from 'react';
import { format, addDays, startOfDay, isBefore, isAfter } from 'date-fns';
import { Loader2, ChevronLeft, ChevronRight, Bell, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import api from '@/lib/axios';
import { bookingService } from '@/services/booking-service';

// Limit booking to 14 days in advance (matches backend validation)
const MAX_BOOKING_DAYS = 14;

// --- Types ---

export interface TimeSlot {
    start_time: string;
    end_time: string;
    status: 'available' | 'booked' | 'maintenance' | 'closed';
}

interface AvailabilityCalendarProps {
    facilityId: string;
    onSlotSelect: (date: string, time: string) => void;
}

// --- Custom Hook (Logic Separation) ---

function useAvailability(facilityId: string, selectedDate: Date) {
    const [availability, setAvailability] = useState<TimeSlot[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<unknown>(null);

    // Static operational hours config (in real app, fetch from Facility config)
    const hours = useMemo(() => Array.from({ length: 15 }, (_, i) => i + 8), []);

    const fetchAvailability = useCallback(async () => {
        if (!facilityId) return;

        setLoading(true);
        setError(null);
        try {
            const dateStr = format(selectedDate, 'yyyy-MM-dd');
            // Using the High Performance Booking endpoint we refactored
            const res = await api.get<{ data: unknown[] }>(`/bookings/availability?facility_id=${facilityId}&date=${dateStr}`);

            const serverSlots = res.data.data;

            // Safer mapping:
            if (Array.isArray(serverSlots)) {
                // Check if it looks like the new format (has 'start_time' and 'status')
                const first = serverSlots[0] as object;
                const isNewFormat = serverSlots.length > 0 && first && 'status' in first;

                if (isNewFormat) {
                    setAvailability(serverSlots as TimeSlot[]);
                    return;
                }
            }

            // Fallback to legacy block calculation if API mismatch
            const blocked = Array.isArray(serverSlots) ? serverSlots : [];
            const computedSlots: TimeSlot[] = hours.map(h => {
                const timeStr = `${h.toString().padStart(2, '0')}:00`;
                // Check range with safe casting inside
                const isBlocked = blocked.some((b) => {
                    const blockItem = b as { start_time: string; end_time: string };
                    const start = new Date(blockItem.start_time).getHours();
                    const end = new Date(blockItem.end_time).getHours();
                    return h >= start && h < end;
                });
                return {
                    start_time: timeStr,
                    end_time: `${(h + 1).toString().padStart(2, '0')}:00`,
                    status: isBlocked ? 'booked' : 'available'
                };
            });
            setAvailability(computedSlots);

        } catch (err: unknown) {
            console.error("Failed to fetch availability:", err);
            setError(err);
        } finally {
            setLoading(false);
        }
    }, [facilityId, selectedDate, hours]);

    useEffect(() => {
        fetchAvailability();
    }, [fetchAvailability]);

    return { availability, loading, error, refresh: fetchAvailability };
}

// --- Component (View) ---

export function AvailabilityCalendar({ facilityId, onSlotSelect }: AvailabilityCalendarProps) {
    const [selectedDate, setSelectedDate] = useState<Date>(new Date());
    const [selectedSlot, setSelectedSlot] = useState<string | null>(null);
    const [joiningWaitlist, setJoiningWaitlist] = useState<string | null>(null);
    const [waitlistSuccess, setWaitlistSuccess] = useState<string | null>(null);

    // Use Custom Hook
    const { availability, loading } = useAvailability(facilityId, selectedDate);

    // Computed date constraints
    const today = useMemo(() => startOfDay(new Date()), []);
    const maxDate = useMemo(() => addDays(today, MAX_BOOKING_DAYS), [today]);
    const canGoBack = !isBefore(startOfDay(addDays(selectedDate, -1)), today);
    const canGoForward = !isAfter(startOfDay(addDays(selectedDate, 1)), maxDate);

    const handleDateChange = (days: number) => {
        const newDate = addDays(selectedDate, days);
        const newDateStart = startOfDay(newDate);

        // Prevent past dates
        if (isBefore(newDateStart, today)) {
            return;
        }

        // Prevent dates beyond 14 days
        if (isAfter(newDateStart, maxDate)) {
            return;
        }

        setSelectedDate(newDate);
        setSelectedSlot(null);
    };

    const handleSlotClick = (time: string, status: string) => {
        if (status !== 'available') return;
        setSelectedSlot(time);
        onSlotSelect(format(selectedDate, 'yyyy-MM-dd'), time);
    };

    const handleJoinWaitlist = async (slot: TimeSlot) => {
        setJoiningWaitlist(slot.start_time);
        try {
            await bookingService.addToWaitlist({
                resource_id: facilityId,
                target_date: format(selectedDate, 'yyyy-MM-dd') + 'T' + slot.start_time + ':00'
            });
            setWaitlistSuccess(slot.start_time);
            setTimeout(() => setWaitlistSuccess(null), 3000);
        } catch (err) {
            console.error('Failed to join waitlist', err);
        } finally {
            setJoiningWaitlist(null);
        }
    };

    return (
        <div className="space-y-4">
            <DateNavigator
                date={selectedDate}
                onPrev={() => handleDateChange(-1)}
                onNext={() => handleDateChange(1)}
                canGoBack={canGoBack}
                canGoForward={canGoForward}
                daysRemaining={Math.ceil((maxDate.getTime() - startOfDay(selectedDate).getTime()) / (1000 * 60 * 60 * 24))}
            />

            <div className="grid grid-cols-3 gap-2 max-h-60 overflow-y-auto p-1">
                {loading ? (
                    <div className="col-span-3 flex justify-center py-8">
                        <Loader2 className="h-6 w-6 animate-spin text-brand-600" />
                    </div>
                ) : (
                    availability.map((slot) => (
                        <TimeSlotButton
                            key={slot.start_time}
                            slot={slot}
                            isSelected={selectedSlot === slot.start_time}
                            onClick={() => handleSlotClick(slot.start_time, slot.status)}
                            onWaitlist={() => handleJoinWaitlist(slot)}
                            joiningWaitlist={joiningWaitlist === slot.start_time}
                            waitlistSuccess={waitlistSuccess === slot.start_time}
                        />
                    ))
                )}
                {!loading && availability.length === 0 && (
                    <p className="col-span-3 text-center text-xs text-gray-400 py-4">No slots available</p>
                )}
            </div>

            <Legend />
        </div>
    );
}

// --- Sub-components (Component Splitting) ---

function DateNavigator({
    date,
    onPrev,
    onNext,
    canGoBack,
    canGoForward,
    daysRemaining
}: {
    date: Date;
    onPrev: () => void;
    onNext: () => void;
    canGoBack: boolean;
    canGoForward: boolean;
    daysRemaining: number;
}) {
    return (
        <div className="space-y-2">
            <div className="flex items-center justify-between bg-gray-50 dark:bg-zinc-800 p-2 rounded-lg">
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={onPrev}
                    disabled={!canGoBack}
                    className={cn(!canGoBack && "opacity-50 cursor-not-allowed")}
                >
                    <ChevronLeft className="h-4 w-4" />
                </Button>
                <span className="font-semibold text-sm">
                    {format(date, 'EEEE, MMMM d')}
                </span>
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={onNext}
                    disabled={!canGoForward}
                    className={cn(!canGoForward && "opacity-50 cursor-not-allowed")}
                >
                    <ChevronRight className="h-4 w-4" />
                </Button>
            </div>
            {daysRemaining <= 3 && canGoForward && (
                <div className="flex items-center gap-1 text-xs text-amber-600 dark:text-amber-400 px-2">
                    <AlertCircle className="h-3 w-3" />
                    <span>Podés reservar hasta {daysRemaining} {daysRemaining === 1 ? 'día' : 'días'} más de anticipación.</span>
                </div>
            )}
            {!canGoForward && (
                <div className="flex items-center gap-1 text-xs text-red-600 dark:text-red-400 px-2">
                    <AlertCircle className="h-3 w-3" />
                    <span>Alcanzaste el límite máximo de reserva anticipada (14 días).</span>
                </div>
            )}
        </div>
    );
}

function TimeSlotButton({ slot, isSelected, onClick, onWaitlist, joiningWaitlist, waitlistSuccess }: {
    slot: TimeSlot;
    isSelected: boolean;
    onClick: () => void;
    onWaitlist?: () => void;
    joiningWaitlist?: boolean;
    waitlistSuccess?: boolean;
}) {
    if (slot.status === 'booked' && onWaitlist) {
        return (
            <div className="relative">
                <button
                    disabled
                    className="w-full px-2 py-2 text-xs font-medium rounded-md border bg-gray-100 dark:bg-zinc-800 text-gray-400 border-transparent cursor-not-allowed"
                >
                    {slot.start_time}
                </button>
                {waitlistSuccess ? (
                    <span className="absolute -top-1 -right-1 text-[8px] bg-green-500 text-white px-1 rounded">✓</span>
                ) : (
                    <button
                        onClick={(e) => { e.stopPropagation(); onWaitlist(); }}
                        disabled={joiningWaitlist}
                        className="absolute -top-1 -right-1 p-0.5 bg-amber-500 hover:bg-amber-600 text-white rounded-full transition-colors"
                        title="Avisarme si se libera"
                    >
                        {joiningWaitlist ? (
                            <Loader2 className="h-3 w-3 animate-spin" />
                        ) : (
                            <Bell className="h-3 w-3" />
                        )}
                    </button>
                )}
            </div>
        );
    }

    return (
        <button
            disabled={slot.status !== 'available'}
            onClick={onClick}
            className={cn(
                "px-2 py-2 text-xs font-medium rounded-md border transition-all relative",
                slot.status === 'available'
                    ? "bg-white dark:bg-zinc-900 border-gray-200 dark:border-zinc-700 hover:border-brand-500 hover:text-brand-600 dark:hover:text-brand-400 text-gray-700 dark:text-gray-300"
                    : "bg-gray-100 dark:bg-zinc-800 text-gray-400 border-transparent cursor-not-allowed decoration-slice",
                isSelected && "ring-2 ring-brand-600 border-brand-600 bg-brand-50 dark:bg-brand-900/20 text-brand-700 dark:text-brand-300"
            )}
        >
            {slot.start_time}
        </button>
    );
}

function Legend() {
    return (
        <div className="flex items-center gap-4 text-[10px] text-gray-500 justify-center pt-2">
            <div className="flex items-center gap-1">
                <div className="w-2 h-2 rounded-full bg-white border border-gray-300"></div> Available
            </div>
            <div className="flex items-center gap-1">
                <div className="w-2 h-2 rounded-full bg-gray-100 border border-gray-200"></div> Occupied
            </div>
            <div className="flex items-center gap-1">
                <div className="w-2 h-2 rounded-full bg-brand-50 border border-brand-600"></div> Selected
            </div>
        </div>
    );
}
