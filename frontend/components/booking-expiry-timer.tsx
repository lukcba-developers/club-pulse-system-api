'use client';

import { useEffect, useState, useRef } from 'react';
import { Clock } from 'lucide-react';
import { cn } from '@/lib/utils';

interface BookingExpiryTimerProps {
    expiry: string; // ISO String
    onExpire?: () => void;
    className?: string;
}

export function BookingExpiryTimer({ expiry, onExpire, className }: BookingExpiryTimerProps) {
    const [timeLeft, setTimeLeft] = useState<string>('');
    const [isExpired, setIsExpired] = useState(false);
    const onExpireRef = useRef(onExpire);

    // Keep callback ref updated without re-running effect
    useEffect(() => {
        onExpireRef.current = onExpire;
    }, [onExpire]);

    useEffect(() => {
        const calculateTimeLeft = () => {
            const now = new Date().getTime();
            const expiryTime = new Date(expiry).getTime();
            const distance = expiryTime - now;

            if (distance < 0) {
                setIsExpired(true);
                setTimeLeft('00:00');
                onExpireRef.current?.();
                return true; // Signal expired
            }

            const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
            const seconds = Math.floor((distance % (1000 * 60)) / 1000);

            setTimeLeft(
                `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
            );
            return false;
        };

        // Initial call
        const expired = calculateTimeLeft();
        if (expired) return;

        const timer = setInterval(() => {
            const expired = calculateTimeLeft();
            if (expired) clearInterval(timer);
        }, 1000);

        return () => clearInterval(timer);
    }, [expiry]); // Only re-run when expiry changes

    if (isExpired) return null; // Or show "Expired"

    return (
        <div className={cn("flex items-center space-x-2 text-amber-600 font-medium", className)}>
            <Clock className="h-4 w-4" />
            <span>Pagar en: {timeLeft}</span>
        </div>
    );
}
