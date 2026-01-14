"use client";

import { cn } from "@/lib/utils";
import NextImage from "next/image";

interface BadgeRarity {
    COMMON: string;
    RARE: string;
    EPIC: string;
    LEGENDARY: string;
}

const rarityStyles: BadgeRarity = {
    COMMON: "border-gray-400 bg-gray-50 dark:bg-gray-800",
    RARE: "border-blue-500 bg-blue-50 dark:bg-blue-900/30 shadow-blue-500/20 shadow-lg",
    EPIC: "border-purple-500 bg-purple-50 dark:bg-purple-900/30 shadow-purple-500/30 shadow-lg animate-pulse",
    LEGENDARY: "border-yellow-500 bg-gradient-to-r from-yellow-50 to-orange-50 dark:from-yellow-900/30 dark:to-orange-900/30 shadow-yellow-500/40 shadow-xl",
};

const rarityGlow: BadgeRarity = {
    COMMON: "",
    RARE: "ring-2 ring-blue-400/30",
    EPIC: "ring-2 ring-purple-400/50",
    LEGENDARY: "ring-4 ring-yellow-400/60 animate-pulse",
};

export interface Badge {
    id: string;
    code: string;
    name: string;
    description: string;
    iconUrl?: string;
    rarity: keyof BadgeRarity;
    category: string;
    xpReward?: number;
}

interface BadgeDisplayProps {
    badge: Badge;
    size?: "sm" | "md" | "lg";
    showDetails?: boolean;
    featured?: boolean;
    onClick?: () => void;
}

export function BadgeDisplay({
    badge,
    size = "md",
    showDetails = true,
    featured = false,
    onClick,
}: BadgeDisplayProps) {
    const sizeClasses = {
        sm: "w-12 h-12",
        md: "w-16 h-16",
        lg: "w-24 h-24",
    };

    const iconSizeClasses = {
        sm: "text-xl",
        md: "text-2xl",
        lg: "text-4xl",
    };

    const defaultIcons: Record<string, string> = {
        PROGRESSION: "🏆",
        STREAK: "🔥",
        SOCIAL: "👥",
        TOURNAMENT: "⚔️",
        BOOKING: "📅",
        SPECIAL: "⭐",
    };

    const icon = badge.iconUrl || defaultIcons[badge.category] || "🏅";

    return (
        <div
            className={cn(
                "flex flex-col items-center gap-2 p-3 rounded-xl border-2 transition-all duration-300",
                "hover:scale-105 cursor-pointer",
                rarityStyles[badge.rarity],
                rarityGlow[badge.rarity],
                featured && "ring-4 ring-offset-2 ring-offset-background ring-primary"
            )}
            onClick={onClick}
        >
            {/* Badge Icon */}
            <div
                className={cn(
                    "rounded-full flex items-center justify-center bg-white dark:bg-slate-800",
                    sizeClasses[size],
                    iconSizeClasses[size]
                )}
            >
                {badge.iconUrl ? (
                    <div className="relative w-full h-full">
                        <NextImage
                            src={badge.iconUrl}
                            alt={badge.name}
                            fill
                            className="object-cover rounded-full"
                        />
                    </div>
                ) : (
                    <span>{icon}</span>
                )}
            </div>

            {/* Badge Details */}
            {showDetails && (
                <div className="text-center">
                    <p className="font-semibold text-sm dark:text-white">{badge.name}</p>
                    {size !== "sm" && (
                        <p className="text-xs text-muted-foreground line-clamp-2">
                            {badge.description}
                        </p>
                    )}
                    {badge.xpReward && badge.xpReward > 0 && (
                        <span className="mt-1 inline-block text-xs px-2 py-0.5 rounded-full bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300">
                            +{badge.xpReward} XP
                        </span>
                    )}
                </div>
            )}

            {/* Rarity Label */}
            <span
                className={cn(
                    "text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full",
                    {
                        "bg-gray-200 text-gray-600": badge.rarity === "COMMON",
                        "bg-blue-200 text-blue-700": badge.rarity === "RARE",
                        "bg-purple-200 text-purple-700": badge.rarity === "EPIC",
                        "bg-yellow-200 text-yellow-700": badge.rarity === "LEGENDARY",
                    }
                )}
            >
                {badge.rarity.toLowerCase()}
            </span>
        </div>
    );
}

interface BadgeGridProps {
    badges: Badge[];
    size?: "sm" | "md" | "lg";
    onBadgeClick?: (badge: Badge) => void;
}

export function BadgeGrid({ badges, size = "md", onBadgeClick }: BadgeGridProps) {
    // Group by category
    const grouped = badges.reduce((acc, badge) => {
        if (!acc[badge.category]) {
            acc[badge.category] = [];
        }
        acc[badge.category].push(badge);
        return acc;
    }, {} as Record<string, Badge[]>);

    const categoryNames: Record<string, string> = {
        PROGRESSION: "Progresión",
        STREAK: "Rachas",
        SOCIAL: "Social",
        TOURNAMENT: "Torneos",
        BOOKING: "Reservas",
        SPECIAL: "Especiales",
    };

    return (
        <div className="space-y-6">
            {Object.entries(grouped).map(([category, categoryBadges]) => (
                <div key={category}>
                    <h3 className="text-lg font-semibold mb-3 flex items-center gap-2">
                        <span>{categoryNames[category] || category}</span>
                        <span className="text-sm text-muted-foreground">
                            ({categoryBadges.length})
                        </span>
                    </h3>
                    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                        {categoryBadges.map((badge) => (
                            <BadgeDisplay
                                key={badge.id}
                                badge={badge}
                                size={size}
                                onClick={() => onBadgeClick?.(badge)}
                            />
                        ))}
                    </div>
                </div>
            ))}
        </div>
    );
}
