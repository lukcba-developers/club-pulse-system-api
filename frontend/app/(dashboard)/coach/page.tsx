"use client"

import { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { PlayerStatusTable } from "@/components/coach/PlayerStatusTable"
import { TravelEvents } from "@/components/coach/TravelEvents"
import { TravelCalendar } from "@/components/coach/TravelCalendar"
import { AttendanceTracker } from "@/components/coach/AttendanceTracker"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import {
    Users,
    Bus,
    Calendar,
    BarChart3,
    Shield,
    ClipboardList
} from "lucide-react"

// TODO: Obtener estos datos del backend o contexto de auth
const MOCK_TEAMS = [
    { id: "team-1", name: "Sub-15 Masculino" },
    { id: "team-2", name: "Sub-17 Femenino" },
]

export default function CoachDashboardPage() {
    const [selectedTeam, setSelectedTeam] = useState(MOCK_TEAMS[0])

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Dashboard del Entrenador</h1>
                    <p className="text-muted-foreground">
                        Gestiona tu equipo, verifica elegibilidad, organiza viajes y toma asistencia.
                    </p>
                </div>
            </div>

            {/* Stats Overview */}
            <div className="grid gap-4 md:grid-cols-4">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Jugadores</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">24</div>
                        <p className="text-xs text-muted-foreground">
                            En tu plantilla
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Habilitados</CardTitle>
                        <Shield className="h-4 w-4 text-green-500" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-green-600">20</div>
                        <p className="text-xs text-muted-foreground">
                            Listos para jugar
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Inhabilitados</CardTitle>
                        <Shield className="h-4 w-4 text-destructive" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold text-destructive">4</div>
                        <p className="text-xs text-muted-foreground">
                            Requieren atención
                        </p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Próximos Viajes</CardTitle>
                        <Bus className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">2</div>
                        <p className="text-xs text-muted-foreground">
                            Esta semana
                        </p>
                    </CardContent>
                </Card>
            </div>

            {/* Team Selector */}
            {MOCK_TEAMS.length > 1 && (
                <div className="flex gap-2">
                    {MOCK_TEAMS.map((team) => (
                        <button
                            key={team.id}
                            onClick={() => setSelectedTeam(team)}
                            className={`px-4 py-2 rounded-lg transition-colors ${selectedTeam.id === team.id
                                ? "bg-primary text-primary-foreground"
                                : "bg-muted hover:bg-muted/80"
                                }`}
                        >
                            {team.name}
                        </button>
                    ))}
                </div>
            )}

            {/* Main Tabs */}
            <Tabs defaultValue="players" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="players" className="gap-2">
                        <Users className="h-4 w-4" />
                        Jugadores
                    </TabsTrigger>
                    <TabsTrigger value="attendance" className="gap-2">
                        <ClipboardList className="h-4 w-4" />
                        Asistencia
                    </TabsTrigger>
                    <TabsTrigger value="travels" className="gap-2">
                        <Bus className="h-4 w-4" />
                        Viajes
                    </TabsTrigger>
                    <TabsTrigger value="calendar" className="gap-2">
                        <Calendar className="h-4 w-4" />
                        Calendario
                    </TabsTrigger>
                    <TabsTrigger value="stats" className="gap-2">
                        <BarChart3 className="h-4 w-4" />
                        Estadísticas
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="players">
                    <PlayerStatusTable
                        teamId={selectedTeam.id}
                        teamName={selectedTeam.name}
                    />
                </TabsContent>

                <TabsContent value="attendance">
                    <AttendanceTracker
                        teamId={selectedTeam.id}
                        teamName={selectedTeam.name}
                    />
                </TabsContent>

                <TabsContent value="travels">
                    <TravelEvents teamId={selectedTeam.id} />
                </TabsContent>

                <TabsContent value="calendar">
                    <TravelCalendar teamId={selectedTeam.id} />
                </TabsContent>

                <TabsContent value="stats">
                    <Card>
                        <CardHeader>
                            <CardTitle>Estadísticas</CardTitle>
                            <CardDescription>
                                Métricas y rendimiento del equipo
                            </CardDescription>
                        </CardHeader>
                        <CardContent className="flex items-center justify-center py-12">
                            <div className="text-center text-muted-foreground">
                                <BarChart3 className="h-12 w-12 mx-auto mb-4" />
                                <p>Próximamente: Estadísticas del equipo</p>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    )
}
