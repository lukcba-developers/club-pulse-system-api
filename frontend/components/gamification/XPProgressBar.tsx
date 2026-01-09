"use client";

import { Progress } from "@/components/ui/progress";
import { cn } from "@/lib/utils";
import { useStreakStatus } from "@/hooks/useGamification";

interface XPProgressBarProps {
    level: number;
    currentXP: number;
    requiredXP: number;
    totalXP?: number;
    currentStreak?: number;
    className?: string;
    showDetails?: boolean;
}

export function XPProgressBar({
    level,
    currentXP,
    requiredXP,
    totalXP,
    currentStreak = 0,
    className,
    showDetails = true,
}: XPProgressBarProps) {
    const progress = Math.min((currentXP / requiredXP) * 100, 100);
    const streakInfo = useStreakStatus(currentStreak);

    return (
        <div className={cn("space-y-2", className)}>
            {/* Level & XP Header */}
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                    {/* Level Badge */}
                    <div className="relative">
                        <div className="w-12 h-12 rounded-full bg-gradient-to-r from-purple-500 to-pink-500 flex items-center justify-center shadow-lg">
                            <span className="text-white font-bold text-lg">{level}</span>
                        </div>
                        {currentStreak >= 3 && (
                            <span className="absolute -top-1 -right-1 text-sm">
                                {streakInfo.emoji}
                            </span>
                        )}
                    </div>

                    <div>
                        <p className="font-semibold text-lg">Nivel {level}</p>
                        {currentStreak >= 3 && (
                            <p className={cn("text-xs font-medium", streakInfo.color)}>
                                {streakInfo.label} • x{streakInfo.multiplier} XP
                            </p>
                        )}
                    </div>
                </div>

                {/* XP Counter */}
                {showDetails && (
                    <div className="text-right">
                        <p className="text-sm font-medium">
                            {currentXP.toLocaleString()} / {requiredXP.toLocaleString()} XP
                        </p>
                        {totalXP !== undefined && (
                            <p className="text-xs text-muted-foreground">
                                Total: {totalXP.toLocaleString()} XP
                            </p>
                        )}
                    </div>
                )}
            </div>

            {/* Progress Bar */}
            <div className="relative">
                <Progress value={progress} className="h-3 bg-slate-200 dark:bg-slate-700" />

                {/* Milestone markers */}
                <div className="absolute inset-0 flex">
                    {[25, 50, 75].map((milestone) => (
                        <div
                            key={milestone}
                            className="absolute top-0 bottom-0 w-0.5 bg-white/30"
                            style={{ left: `${milestone}%` }}
                        />
                    ))}
                </div>
            </div>

            {/* Streak Info */}
            {currentStreak > 0 && showDetails && (
                <div className="flex items-center justify-between text-xs">
                    <span className="text-muted-foreground">
                        Racha actual: <strong>{currentStreak}</strong> días
                    </span>
                    <span className={cn("font-medium", streakInfo.color)}>
                        {streakInfo.multiplier > 1 && `+${((streakInfo.multiplier - 1) * 100).toFixed(0)}% XP Bonus`}
                    </span>
                </div>
            )}
        </div>
    );
}

// Compact version for navbar/sidebar
interface XPProgressCompactProps {
    level: number;
    progress: number;
    className?: string;
}

export function XPProgressCompact({ level, progress, className }: XPProgressCompactProps) {
    return (
        <div className={cn("flex items-center gap-2", className)}>
            <div className="w-8 h-8 rounded-full bg-gradient-to-r from-purple-500 to-pink-500 flex items-center justify-center">
                <span className="text-white font-bold text-xs">{level}</span>
            </div>
            <div className="flex-1 min-w-[60px]">
                <Progress value={progress} className="h-2" />
            </div>
        </div>
    );
}
