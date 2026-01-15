"use client"

import { useState } from "react"
import { championshipService, Team } from "@/services/championship-service"
import { useToast } from "@/components/ui/use-toast"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Plus, Users } from "lucide-react"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"

export function TeamManagement() {
    const { toast } = useToast()
    const [teams, setTeams] = useState<Team[]>([])
    const [isCreateOpen, setIsCreateOpen] = useState(false)
    const [newTeamName, setNewTeamName] = useState("")
    const [selectedTeamIdForMember, setSelectedTeamIdForMember] = useState<string | null>(null)

    // NOTE: Currently backend doesn't have ListTeams endpoint in Championship module explicitly exposed for Admin listing 
    // without context. We might only be able to create for now or need to add List functionality.
    // For this gap analysis fix, we focus on CREATION which was the blocker.

    const handleCreateTeam = async () => {
        try {
            const team = await championshipService.createTeam({
                name: newTeamName,
                // club_id is handled by context/auth in backend usually, or mapped.
            })
            setTeams([...teams, team])
            toast({ title: "Equipo Creado", description: `Equipo ${team.name} listo para inscribirse.` })
            setIsCreateOpen(false)
            setNewTeamName("")
        } catch (error) {
            console.error(error)
            toast({ title: "Error", description: "No se pudo crear el equipo.", variant: "destructive" })
        }
    }

    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                    <Users className="w-5 h-5" /> Gestión de Equipos
                </CardTitle>
                <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
                    <DialogTrigger asChild>
                        <Button size="sm" className="gap-2">
                            <Plus className="w-4 h-4" /> Nuevo Equipo
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Registrar Nuevo Equipo</DialogTitle>
                            <DialogDescription>
                                Crea un equipo para luego inscribirlo en los torneos.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4 py-4">
                            <div className="space-y-2">
                                <Label>Nombre del Equipo</Label>
                                <Input
                                    placeholder="Ej: Los Rayos FC"
                                    value={newTeamName}
                                    onChange={(e) => setNewTeamName(e.target.value)}
                                />
                            </div>
                        </div>
                        <DialogFooter>
                            <Button onClick={handleCreateTeam} disabled={!newTeamName}>Crear Equipo</Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </CardHeader>
            <CardContent>
                <div className="text-sm text-muted-foreground mb-4">
                    Equipos creados recientemente (Sesión actual):
                </div>
                {teams.length === 0 ? (
                    <div className="text-center py-6 text-gray-500 bg-gray-50 rounded-lg border border-dashed">
                        No has creado equipos en esta sesión.
                    </div>
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                        {teams.map(team => (
                            <div key={team.id} className="p-3 bg-white border rounded shadow-sm flex justify-between items-center group">
                                <div className="flex flex-col">
                                    <span className="font-medium">{team.name}</span>
                                    <span className="text-xs text-gray-400 font-mono select-all" title="Click para copiar ID" onClick={() => {
                                        navigator.clipboard.writeText(team.id)
                                        toast({ title: "ID Copiado" })
                                    }}>
                                        {team.id.substring(0, 8)}...
                                    </span>
                                </div>
                                <Button size="sm" variant="ghost" className="opacity-0 group-hover:opacity-100 transition-opacity" onClick={() => setSelectedTeamIdForMember(team.id)}>
                                    <Plus className="w-3 h-3 mr-1" /> Miembro
                                </Button>
                            </div>
                        ))}
                    </div>
                )}
            </CardContent>

            {/* Add Member Dialog */}
            <AddMemberDialog
                isOpen={!!selectedTeamIdForMember}
                onClose={() => setSelectedTeamIdForMember(null)}
                teamId={selectedTeamIdForMember || ""}
                teamName={teams.find(t => t.id === selectedTeamIdForMember)?.name || ""}
            />
        </Card>
    )
}

function AddMemberDialog({ isOpen, onClose, teamId, teamName }: { isOpen: boolean, onClose: () => void, teamId: string, teamName: string }) {
    const { toast } = useToast()
    const [userId, setUserId] = useState("")
    const [loading, setLoading] = useState(false)

    const handleAddMember = async () => {
        try {
            setLoading(true)
            await championshipService.addMember(teamId, userId)
            toast({ title: "Miembro Agregado", description: `Usuario agregado a ${teamName}.` })
            onClose()
            setUserId("")
        } catch (error) {
            console.error(error)
            toast({ title: "Error", description: "No se pudo agregar el miembro.", variant: "destructive" })
        } finally {
            setLoading(false)
        }
    }

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Agregar Miembro a {teamName}</DialogTitle>
                    <DialogDescription>
                        Ingresa el ID del usuario para agregarlo al equipo.
                        (Nota: En una versión completa, usaríamos un buscador por nombre/email)
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                    <div className="space-y-2">
                        <Label>ID de Usuario</Label>
                        <Input
                            placeholder="UUID del Usuario"
                            value={userId}
                            onChange={(e) => setUserId(e.target.value)}
                        />
                    </div>
                </div>
                <DialogFooter>
                    <Button onClick={handleAddMember} disabled={!userId || loading}>
                        {loading ? "Agregando..." : "Agregar Miembro"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    )
}
