import { Booking } from '@/services/booking-service';
import { Facility } from '@/services/facility-service';
import { cn } from '@/lib/utils';
import { format, parseISO, differenceInMinutes, getHours, getMinutes } from 'date-fns';

interface BookingCalendarProps {
    bookings: Booking[];
    facilities: Facility[];
    date: Date;
}

export function BookingCalendar({ bookings, facilities }: BookingCalendarProps) {
    const startHour = 8;
    const endHour = 23;
    const hours = Array.from({ length: endHour - startHour + 1 }, (_, i) => startHour + i);
    const hourHeight = 80;

    const getBookingStyle = (booking: Booking) => {
        const start = parseISO(booking.start_time); // "2024-01-01T10:00:00Z"
        const end = parseISO(booking.end_time);

        // Normalize simple hour calculation for local time display in grid
        // Assuming booking times are already compatible or UTC.
        // For MVP, simplistic parsing.

        const startH = getHours(start);
        const startM = getMinutes(start);
        const durationMin = differenceInMinutes(end, start);

        const top = ((startH - startHour) * hourHeight) + ((startM / 60) * hourHeight);
        const height = (durationMin / 60) * hourHeight;

        return {
            top: `${top}px`,
            height: `${height}px`,
        };
    };

    return (
        <div className="border rounded-lg bg-white dark:bg-zinc-900 shadow-sm overflow-hidden">
            {/* Header: Facilities */}
            <div className="flex border-b divide-x dark:border-zinc-800 dark:divide-zinc-800">
                <div className="w-16 flex-shrink-0 p-4 border-r bg-gray-50 dark:bg-zinc-950 dark:border-zinc-800"></div> {/* Time Col Header */}
                {facilities.map(fac => (
                    <div key={fac.id} className="flex-1 p-4 text-center font-medium text-sm text-gray-700 dark:text-gray-200 bg-gray-50 dark:bg-zinc-950 truncate">
                        {fac.name}
                    </div>
                ))}
            </div>

            {/* Grid */}
            <div className="relative overflow-y-auto" style={{ height: '600px' }}>
                <div className="flex" style={{ height: `${hours.length * hourHeight}px` }}>

                    {/* Time Column */}
                    <div className="w-16 flex-shrink-0 border-r bg-gray-50 dark:bg-zinc-950 dark:border-zinc-800 divide-y dark:divide-zinc-800">
                        {hours.map(h => (
                            <div key={h} className="text-xs text-gray-500 text-right pr-2 pt-1 border-b dark:border-zinc-800" style={{ height: `${hourHeight}px` }}>
                                {h}:00
                            </div>
                        ))}
                    </div>

                    {/* Columns per Facility */}
                    {facilities.map(fac => (
                        <div key={fac.id} className="flex-1 border-r relative dark:border-zinc-800">
                            {/* Background Lines */}
                            {hours.map(h => (
                                <div key={h} className="border-b border-dashed border-gray-100 dark:border-zinc-800 w-full absolute" style={{ top: `${(h - startHour) * hourHeight}px`, height: '1px' }}></div>
                            ))}

                            {/* Bookings */}
                            {bookings
                                .filter(b => b.facility_id === fac.id)
                                .map(booking => (
                                    <div
                                        key={booking.id}
                                        className={cn(
                                            "absolute left-1 right-1 rounded px-2 py-1 text-xs border overflow-hidden",
                                            "bg-indigo-100 border-indigo-200 text-indigo-700 dark:bg-indigo-900/50 dark:border-indigo-700 dark:text-indigo-200"
                                        )}
                                        style={getBookingStyle(booking)}
                                        title={`${format(parseISO(booking.start_time), 'HH:mm')} - ${format(parseISO(booking.end_time), 'HH:mm')}`}
                                    >
                                        <div className="font-semibold">{format(parseISO(booking.start_time), 'HH:mm')}</div>
                                        <div className="truncate opacity-75">Ocupado</div>
                                    </div>
                                ))
                            }
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
