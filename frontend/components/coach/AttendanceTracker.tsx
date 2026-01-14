"use client"

import { useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Checkbox } from "@/components/ui/checkbox"
import { useToast } from "@/hooks/use-toast"
import { Save, AlertTriangle } from "lucide-react"

interface Player {
    id: string
    first_name: string
    last_name: string
    avatar_url?: string
    status_flags: {
        medical_status: string
        financial_status: string
    }
}

interface AttendanceTrackerProps {
    teamId: string
    teamName: string
}

export function AttendanceTracker({ teamId, teamName }: AttendanceTrackerProps) {
    const [players, setPlayers] = useState<Player[]>([])
    const [selectedPlayers, setSelectedPlayers] = useState<Set<string>>(new Set())
    const [loading, setLoading] = useState(false)
    const [date] = useState(new Date().toISOString().split('T')[0])
    const { toast } = useToast()

    // Simulación de carga de jugadores (usaríamos el mismo endpoint que PlayerStatusTable)
    const loadPlayers = async () => {
        setLoading(true)
        try {
            const response = await fetch(`/api/teams/${teamId}/players`, {
                headers: { "Authorization": `Bearer ${localStorage.getItem("token")}` }
            })
            if (response.ok) {
                const data = await response.json()
                setPlayers(data.map((d: { user: Player; status_flags: Player['status_flags'] }) => ({ ...d.user, status_flags: d.status_flags })))
            }
        } catch (e) {
            console.error(e)
        } finally {
            setLoading(false)
        }
    }

    const handleToggle = (playerId: string) => {
        const newSelected = new Set(selectedPlayers)
        if (newSelected.has(playerId)) {
            newSelected.delete(playerId)
        } else {
            newSelected.add(playerId)
        }
        setSelectedPlayers(newSelected)
    }

    const handleSave = async () => {
        setLoading(true)
        try {
            const response = await fetch("/api/attendance/batch", {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    team_id: teamId,
                    date: date,
                    player_ids: Array.from(selectedPlayers),
                    type: "TRAINING"
                })
            })

            if (response.ok) {
                toast({
                    title: "Asistencia Guardada",
                    description: `Se registró la asistencia de ${selectedPlayers.size} jugadores.`
                })
                setSelectedPlayers(new Set())
            } else {
                throw new Error("Error al guardar")
            }
        } catch {
            toast({
                title: "Error",
                description: "No se pudo actualizar la asistencia",
                variant: "destructive",
            });
        } finally {
            setLoading(false)
        }
    }

    // Alerta si intentamos marcar presente a alguien con problemas médicos
    const getAlert = (playerId: string) => {
        const player = players.find(p => p.id === playerId)
        if (player?.status_flags.medical_status !== "VALID" && selectedPlayers.has(playerId)) {
            return (
                <div className="flex items-center text-amber-500 text-xs mt-1">
                    <AlertTriangle className="h-3 w-3 mr-1" />
                    Sin Apto Médico
                </div>
            )
        }
        return null
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Tomar Asistencia - {teamName}</CardTitle>
                <CardDescription>Registro del entrenamiento del día</CardDescription>
                <div className="flex justify-end">
                    <Button variant="outline" size="sm" onClick={loadPlayers}>Cargar Jugadores</Button>
                </div>
            </CardHeader>
            <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {players.map(player => (
                        <div
                            key={player.id}
                            className={`
                flex items-start space-x-3 p-3 rounded-lg border 
                ${selectedPlayers.has(player.id) ? "bg-primary/5 border-primary" : "bg-card"}
              `}
                        >
                            <Checkbox
                                id={player.id}
                                checked={selectedPlayers.has(player.id)}
                                onCheckedChange={() => handleToggle(player.id)}
                            />
                            <div className="grid gap-1.5 leading-none">
                                <label
                                    htmlFor={player.id}
                                    className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 flex items-center gap-2"
                                >
                                    <Avatar className="h-6 w-6">
                                        <AvatarImage src={player.avatar_url} />
                                        <AvatarFallback>{player.first_name[0]}</AvatarFallback>
                                    </Avatar>
                                    {player.first_name} {player.last_name}
                                </label>
                                {getAlert(player.id)}
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
            <CardFooter className="flex justify-between">
                <div className="text-sm text-muted-foreground">
                    {selectedPlayers.size} presentes seleccionados
                </div>
                <Button onClick={handleSave} disabled={loading || players.length === 0}>
                    <Save className="h-4 w-4 mr-2" />
                    Guardar Asistencia
                </Button>
            </CardFooter>
        </Card>
    )
}
