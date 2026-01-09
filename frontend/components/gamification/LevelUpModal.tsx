"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
    Dialog,
    DialogContent,
    DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface LevelUpModalProps {
    isOpen: boolean;
    onClose: () => void;
    newLevel: number;
    xpGained?: number;
}

export function LevelUpModal({
    isOpen,
    onClose,
    newLevel,
    xpGained,
}: LevelUpModalProps) {
    const [showConfetti, setShowConfetti] = useState(false);
    const [particles, setParticles] = useState<Array<{
        left: number;
        colorIndex: number;
        rotate: number;
        duration: number;
        delay: number;
    }>>([]);

    useEffect(() => {
        if (isOpen) {
            // Generate deterministic random values for one animation cycle
            // eslint-disable-next-line react-hooks/set-state-in-effect
            setParticles([...Array(50)].map(() => ({
                left: Math.random() * 100,
                colorIndex: Math.floor(Math.random() * 5),
                rotate: Math.random() * 360,
                duration: 2 + Math.random() * 2,
                delay: Math.random() * 0.5,
            })));

            setShowConfetti(true);
            // Auto-hide confetti after animation
            const timer = setTimeout(() => setShowConfetti(false), 3000);
            return () => clearTimeout(timer);
        }
    }, [isOpen]);

    return (
        <AnimatePresence>
            {isOpen && (
                <Dialog open={isOpen} onOpenChange={onClose}>
                    <DialogContent className="sm:max-w-md overflow-hidden bg-gradient-to-b from-purple-900 via-indigo-900 to-slate-900 border-purple-500/50 text-white">
                        <DialogTitle className="sr-only">¡Subiste de nivel!</DialogTitle>

                        {/* Confetti Effect */}
                        {showConfetti && (
                            <div className="absolute inset-0 pointer-events-none overflow-hidden">
                                {particles.map((p, i) => (
                                    <motion.div
                                        key={i}
                                        className="absolute w-2 h-2 rounded-full"
                                        style={{
                                            left: `${p.left}%`,
                                            backgroundColor: [
                                                "#FFD700",
                                                "#FF6B6B",
                                                "#4ECDC4",
                                                "#A855F7",
                                                "#3B82F6",
                                            ][p.colorIndex],
                                        }}
                                        initial={{ y: -20, opacity: 1 }}
                                        animate={{
                                            y: 400,
                                            opacity: 0,
                                            rotate: p.rotate,
                                        }}
                                        transition={{
                                            duration: p.duration,
                                            delay: p.delay,
                                            ease: "easeOut",
                                        }}
                                    />
                                ))}
                            </div>
                        )}

                        {/* Content */}
                        <div className="text-center py-8 relative z-10">
                            {/* Level Ring Animation */}
                            <motion.div
                                className="relative mx-auto w-32 h-32 mb-6"
                                initial={{ scale: 0, rotate: -180 }}
                                animate={{ scale: 1, rotate: 0 }}
                                transition={{
                                    type: "spring",
                                    stiffness: 200,
                                    damping: 15,
                                    delay: 0.2,
                                }}
                            >
                                <div className="absolute inset-0 rounded-full bg-gradient-to-r from-yellow-400 via-orange-500 to-red-500 animate-pulse" />
                                <div className="absolute inset-2 rounded-full bg-slate-900 flex items-center justify-center">
                                    <span className="text-5xl font-bold bg-gradient-to-r from-yellow-400 to-orange-500 bg-clip-text text-transparent">
                                        {newLevel}
                                    </span>
                                </div>
                            </motion.div>

                            {/* Title */}
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: 0.4 }}
                            >
                                <h2 className="text-3xl font-bold mb-2">
                                    🎉 ¡Nivel {newLevel}!
                                </h2>
                                <p className="text-purple-200 text-lg">
                                    ¡Felicitaciones! Has subido de nivel
                                </p>
                            </motion.div>

                            {/* XP Gained */}
                            {xpGained && (
                                <motion.div
                                    className="mt-4 inline-block px-4 py-2 rounded-full bg-purple-500/30 border border-purple-400/50"
                                    initial={{ opacity: 0, scale: 0.8 }}
                                    animate={{ opacity: 1, scale: 1 }}
                                    transition={{ delay: 0.6 }}
                                >
                                    <span className="text-purple-200">+{xpGained} XP</span>
                                </motion.div>
                            )}

                            {/* Rewards Preview (Future Enhancement) */}
                            <motion.div
                                className="mt-6 p-4 rounded-lg bg-white/5 border border-white/10"
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: 0.8 }}
                            >
                                <p className="text-sm text-purple-300 mb-2">
                                    Beneficios desbloqueados:
                                </p>
                                <div className="flex justify-center gap-4 text-sm">
                                    <span className="flex items-center gap-1">
                                        <span className="text-yellow-400">⚡</span> +5% XP Bonus
                                    </span>
                                    <span className="flex items-center gap-1">
                                        <span className="text-green-400">🏆</span> Nuevo Rango
                                    </span>
                                </div>
                            </motion.div>

                            {/* Close Button */}
                            <motion.div
                                className="mt-8"
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                transition={{ delay: 1 }}
                            >
                                <Button
                                    onClick={onClose}
                                    className="bg-gradient-to-r from-purple-500 to-pink-500 hover:from-purple-600 hover:to-pink-600 text-white px-8 py-2 rounded-full font-semibold"
                                >
                                    ¡Genial!
                                </Button>
                            </motion.div>
                        </div>
                    </DialogContent>
                </Dialog>
            )}
        </AnimatePresence>
    );
}
