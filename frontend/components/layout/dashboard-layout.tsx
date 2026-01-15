'use client';

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, Calendar, Building2, Settings, LogOut, Menu, User as UserIcon, CreditCard, Globe, Users, ShoppingBag, Trophy } from 'lucide-react';
import { cn } from "@/lib/utils";
import { useAuth } from "@/hooks/use-auth";
import { useBrand } from "@/components/providers/BrandProvider";
import { NotificationToast } from "@/components/notification-toast";
import Image from "next/image";

interface DashboardLayoutProps {
    children: React.ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
    const [isSidebarOpen, setSidebarOpen] = useState(false);
    const pathname = usePathname();
    const { user, logout } = useAuth();
    const { club, isLoading } = useBrand();

    // Define menus per role
    const menus = {
        SUPER_ADMIN: [
            { name: "Dashboard Global", href: "/admin/platform", icon: Globe },
            { name: "Gestión Clubes", href: "/admin/clubs", icon: Building2 },
            { name: "Configuración", href: "/admin/settings", icon: Settings },
        ],
        ADMIN: [
            { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
            { name: "Calendario", href: "/bookings/calendar", icon: Calendar },
            { name: "Instalaciones", href: "/facilities", icon: Building2 },
            { name: "Usuarios", href: "/users", icon: Users },
            { name: "Configuración", href: "/settings", icon: Settings },
        ],
        MEMBER: [
            { name: "Inicio", href: "/dashboard", icon: LayoutDashboard },
            { name: "Reservar", href: "/bookings/new", icon: Calendar },
            { name: "Mis Reservas", href: "/bookings", icon: Users }, // Reuse Users icon or Calendar
            { name: "Mi Perfil", href: "/profile", icon: UserIcon },
            { name: "Membresía", href: "/membership", icon: CreditCard },
            { name: "Tienda", href: "/store", icon: ShoppingBag },
            { name: "Campeonatos", href: "/championships", icon: Trophy },
        ],
        // Default fallback
        GUEST: [
            { name: "Inicio", href: "/", icon: LayoutDashboard },
        ]
    };

    // Select menu based on role
    // user.role is string "SUPER_ADMIN" | "ADMIN" | "MEMBER"
    const role = user?.role as keyof typeof menus || "GUEST";
    const navigation = menus[role] || menus["GUEST"];

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-zinc-900 flex text-sm">
            {/* Mobile Sidebar Overlay */}
            {isSidebarOpen && (
                <div
                    className="fixed inset-0 z-40 bg-black/50 lg:hidden glass-backdrop"
                    onClick={() => setSidebarOpen(false)}
                />
            )}

            {/* Sidebar */}
            <aside
                className={cn(
                    "fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-zinc-950 border-r border-gray-200 dark:border-zinc-800 transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0 shadow-xl",
                    isSidebarOpen ? "translate-x-0" : "-translate-x-full"
                )}
            >
                <div className="flex flex-col h-full">
                    {/* Logo */}
                    <div className="h-16 flex items-center px-6 border-b border-gray-100 dark:border-zinc-800">
                        {isLoading ? (
                            <div className="flex items-center gap-2 animate-pulse">
                                <div className="w-8 h-8 rounded-lg bg-gray-200 dark:bg-zinc-800"></div>
                                <div className="h-4 w-24 bg-gray-200 dark:bg-zinc-800 rounded"></div>
                            </div>
                        ) : (
                            <div className="flex items-center gap-2 font-bold text-xl text-brand-600 dark:text-brand-400">
                                {club?.logo_url ? (
                                    <div className="relative w-8 h-8">
                                        <Image
                                            src={club.logo_url}
                                            alt={club.name}
                                            fill
                                            className="object-contain rounded-lg"
                                            unoptimized
                                        />
                                    </div>
                                ) : (
                                    <div className="w-8 h-8 rounded-lg bg-brand-600 flex items-center justify-center text-white">
                                        {club?.name ? club.name.substring(0, 2).toUpperCase() : 'CP'}
                                    </div>
                                )}
                                <span className="truncate max-w-[150px]">{club?.name || 'Club Pulse'}</span>
                            </div>
                        )}
                    </div>

                    {/* Navigation */}
                    <nav className="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
                        <div className="mb-4 px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider">
                            {user?.role ? user.role.replace('_', ' ') : 'Menu'}
                        </div>
                        {navigation.map((item) => {
                            const isActive = pathname === item.href;
                            return (
                                <Link
                                    key={item.name}
                                    href={item.href}
                                    className={cn(
                                        "flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all duration-200 group font-medium",
                                        isActive
                                            ? "bg-brand-50 text-brand-700 dark:bg-brand-900/20 dark:text-brand-300 shadow-sm"
                                            : "text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-zinc-900 hover:text-gray-900 dark:hover:text-gray-200"
                                    )}
                                >
                                    <item.icon className={cn("w-5 h-5", isActive ? "text-brand-600 dark:text-brand-400" : "text-gray-400 group-hover:text-gray-600")} />
                                    {item.name}
                                </Link>
                            );
                        })}
                    </nav>

                    {/* User Profile / Logout */}
                    <div className="p-4 border-t border-gray-100 dark:border-zinc-800">
                        <div className="flex items-center gap-3 px-3 py-3 rounded-xl bg-gray-50 dark:bg-zinc-900/50 mb-2">
                            <div className="w-10 h-10 rounded-full bg-brand-100 dark:bg-brand-900/50 flex items-center justify-center text-brand-600 dark:text-brand-400 font-semibold">
                                {user?.name ? user.name.substring(0, 2).toUpperCase() : 'Guest'}
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-sm font-semibold text-gray-900 dark:text-white truncate">
                                    {user?.name || 'Guest User'}
                                </p>
                                <p className="text-xs text-gray-500 dark:text-gray-400 truncate">
                                    {user?.email || 'Please Login'}
                                </p>
                            </div>
                        </div>
                        <button
                            onClick={logout}
                            className="w-full flex items-center gap-3 px-3 py-2 text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/10 rounded-lg transition-colors text-sm font-medium"
                        >
                            <LogOut className="w-4 h-4" />
                            Cerrar Sesión
                        </button>
                    </div>
                </div>
            </aside>

            {/* Main Content */}
            <div className="flex-1 flex flex-col min-w-0 bg-gray-50/50 dark:bg-zinc-900">
                {/* Top Header (Mobile Only for Menu) + Page Title area if needed */}
                <header className="h-16 lg:hidden flex items-center justify-between px-4 bg-white dark:bg-zinc-950 border-b border-gray-200 dark:border-zinc-800 sticky top-0 z-30">
                    <div className="flex items-center gap-2 font-bold text-lg text-brand-600 dark:text-brand-400">
                        {club?.logo_url ? (
                            <div className="relative w-7 h-7">
                                <Image
                                    src={club.logo_url}
                                    alt={club.name}
                                    fill
                                    className="object-contain rounded-md"
                                    unoptimized
                                />
                            </div>
                        ) : (
                            <div className="w-7 h-7 rounded-md bg-brand-600 flex items-center justify-center text-white text-xs">
                                {club?.name ? club.name.substring(0, 2).toUpperCase() : 'CP'}
                            </div>
                        )}
                        <span>{club?.name || 'Club Pulse'}</span>
                    </div>
                    <div className="flex items-center gap-2 relative">
                        <NotificationToast />
                        <button
                            onClick={() => setSidebarOpen(true)}
                            className="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-zinc-800 rounded-lg"
                        >
                            <Menu className="w-6 h-6" />
                        </button>
                    </div>
                </header>

                {/* Content Scroll Area */}
                <main className="flex-1 overflow-auto p-4 lg:p-8">
                    <div className="max-w-7xl mx-auto space-y-6">
                        {children}
                    </div>
                </main>
            </div >
        </div >
    );
}
