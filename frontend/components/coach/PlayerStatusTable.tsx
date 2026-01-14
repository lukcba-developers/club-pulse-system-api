import { useState, useEffect, useCallback } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from "@/components/ui/tooltip"
import {
    Download,
    AlertTriangle,
    CheckCircle,
    XCircle,
    DollarSign,
    Heart,
    TrendingUp,
    Filter,
    RefreshCw
} from "lucide-react"
import { useToast } from "@/hooks/use-toast"

interface PlayerStatusFlags {
    financial_status: "ACTIVE" | "DEBTOR" | "NO_MEMBERSHIP" | "UNKNOWN"
    medical_status: "VALID" | "EXPIRED" | "MISSING"
    attendance_rate: number
    is_inhabilitado: boolean
}

interface PlayerWithStatus {
    user: {
        id: string
        first_name: string
        last_name: string
        email: string
        avatar_url?: string
        date_of_birth?: string
    }
    status_flags: PlayerStatusFlags
}

interface PlayerStatusTableProps {
    teamId: string
    teamName: string
}

export function PlayerStatusTable({ teamId, teamName }: PlayerStatusTableProps) {
    const [players, setPlayers] = useState<PlayerWithStatus[]>([])
    const [loading, setLoading] = useState(true)
    const [downloading, setDownloading] = useState(false)
    const [filter, setFilter] = useState<"all" | "inhabilitados">("all")
    const { toast } = useToast()


    const fetchPlayers = useCallback(async () => {
        setLoading(true)
        try {
            const response = await fetch(`/api/teams/${teamId}/players`, {
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            })
            if (response.ok) {
                const data = await response.json()
                setPlayers(data)
            }
        } catch (error) {
            console.error("Error fetching players:", error)
            toast({
                title: "Error",
                description: "No se pudieron cargar los jugadores",
                variant: "destructive"
            })
        } finally {
            setLoading(false)
        }
    }, [teamId, toast])

    useEffect(() => {
        if (teamId) {
            fetchPlayers()
        }
    }, [fetchPlayers, teamId])

    const handleExportLeagueFolder = async () => {
        setDownloading(true)
        try {
            const response = await fetch(`/api/teams/${teamId}/league-export?only_eligible=true`, {
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            })

            if (response.ok) {
                const blob = await response.blob()
                const url = window.URL.createObjectURL(blob)
                const a = document.createElement("a")
                a.href = url
                a.download = `carpeta_liga_${teamName.replace(/\s+/g, "_")}.pdf`
                document.body.appendChild(a)
                a.click()
                window.URL.revokeObjectURL(url)
                document.body.removeChild(a)

                toast({
                    title: "Éxito",
                    description: "Carpeta de Liga descargada correctamente"
                })
            } else {
                throw new Error("Error al descargar")
            }
        } catch {
            toast({
                title: "Error",
                description: "No se pudo descargar la Carpeta de Liga",
                variant: "destructive"
            })
        } finally {
            setDownloading(false)
        }
    }

    const getFinancialBadge = (status: string) => {
        switch (status) {
            case "ACTIVE":
                return (
                    <Badge variant="success" className="gap-1">
                        <DollarSign className="h-3 w-3" />
                        Al día
                    </Badge>
                )
            case "DEBTOR":
                return (
                    <Badge variant="destructive" className="gap-1">
                        <DollarSign className="h-3 w-3" />
                        Deudor
                    </Badge>
                )
            case "NO_MEMBERSHIP":
                return (
                    <Badge variant="outline" className="gap-1">
                        <DollarSign className="h-3 w-3" />
                        Sin membresía
                    </Badge>
                )
            default:
                return (
                    <Badge variant="secondary" className="gap-1">
                        <DollarSign className="h-3 w-3" />
                        Desconocido
                    </Badge>
                )
        }
    }

    const getMedicalBadge = (status: string) => {
        switch (status) {
            case "VALID":
                return (
                    <Badge variant="success" className="gap-1">
                        <Heart className="h-3 w-3" />
                        Vigente
                    </Badge>
                )
            case "EXPIRED":
                return (
                    <Badge variant="destructive" className="gap-1">
                        <Heart className="h-3 w-3" />
                        Vencido
                    </Badge>
                )
            case "MISSING":
                return (
                    <Badge variant="warning" className="gap-1">
                        <Heart className="h-3 w-3" />
                        Faltante
                    </Badge>
                )
            default:
                return (
                    <Badge variant="secondary" className="gap-1">
                        <Heart className="h-3 w-3" />
                        Desconocido
                    </Badge>
                )
        }
    }

    const getAttendanceBadge = (rate: number) => {
        const percentage = Math.round(rate * 100)
        if (rate >= 0.8) {
            return (
                <Badge variant="success" className="gap-1">
                    <TrendingUp className="h-3 w-3" />
                    {percentage}%
                </Badge>
            )
        } else if (rate >= 0.5) {
            return (
                <Badge variant="warning" className="gap-1">
                    <TrendingUp className="h-3 w-3" />
                    {percentage}%
                </Badge>
            )
        } else {
            return (
                <Badge variant="destructive" className="gap-1">
                    <TrendingUp className="h-3 w-3" />
                    {percentage}%
                </Badge>
            )
        }
    }

    const filteredPlayers = filter === "inhabilitados"
        ? players.filter(p => p.status_flags.is_inhabilitado)
        : players

    const inhabilitadosCount = players.filter(p => p.status_flags.is_inhabilitado).length

    return (
        <Card>
            <CardHeader>
                <div className="flex justify-between items-start">
                    <div>
                        <CardTitle className="flex items-center gap-2">
                            {teamName}
                            {inhabilitadosCount > 0 && (
                                <Badge variant="destructive">
                                    <AlertTriangle className="h-3 w-3 mr-1" />
                                    {inhabilitadosCount} Inhabilitados
                                </Badge>
                            )}
                        </CardTitle>
                        <CardDescription>
                            Semáforo del jugador - Estado documental, financiero y de asistencia
                        </CardDescription>
                    </div>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setFilter(filter === "all" ? "inhabilitados" : "all")}
                        >
                            <Filter className="h-4 w-4 mr-2" />
                            {filter === "all" ? "Ver Inhabilitados" : "Ver Todos"}
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={fetchPlayers}
                            disabled={loading}
                        >
                            <RefreshCw className={`h-4 w-4 mr-2 ${loading ? "animate-spin" : ""}`} />
                            Actualizar
                        </Button>
                        <Button
                            onClick={handleExportLeagueFolder}
                            disabled={downloading}
                        >
                            <Download className="h-4 w-4 mr-2" />
                            {downloading ? "Descargando..." : "Exportar Carpeta Liga"}
                        </Button>
                    </div>
                </div>
            </CardHeader>
            <CardContent>
                <TooltipProvider>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Jugador</TableHead>
                                <TableHead className="text-center">💰 Cuota</TableHead>
                                <TableHead className="text-center">🏥 Apto Médico</TableHead>
                                <TableHead className="text-center">📊 Asistencia</TableHead>
                                <TableHead className="text-center">Estado</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {loading ? (
                                <TableRow>
                                    <TableCell colSpan={5} className="text-center py-8">
                                        <RefreshCw className="h-8 w-8 animate-spin mx-auto text-muted-foreground" />
                                    </TableCell>
                                </TableRow>
                            ) : filteredPlayers.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                                        {filter === "inhabilitados"
                                            ? "¡Todos los jugadores están habilitados! 🎉"
                                            : "No hay jugadores en este equipo"
                                        }
                                    </TableCell>
                                </TableRow>
                            ) : (
                                filteredPlayers.map((player) => (
                                    <TableRow
                                        key={player.user.id}
                                        className={player.status_flags.is_inhabilitado ? "bg-destructive/5" : ""}
                                    >
                                        <TableCell>
                                            <div className="flex items-center gap-3">
                                                <Avatar>
                                                    <AvatarImage src={player.user.avatar_url} />
                                                    <AvatarFallback>
                                                        {player.user.first_name[0]}{player.user.last_name[0]}
                                                    </AvatarFallback>
                                                </Avatar>
                                                <div>
                                                    <p className="font-medium">
                                                        {player.user.first_name} {player.user.last_name}
                                                    </p>
                                                    <p className="text-sm text-muted-foreground">
                                                        {player.user.email}
                                                    </p>
                                                </div>
                                            </div>
                                        </TableCell>
                                        <TableCell className="text-center">
                                            {getFinancialBadge(player.status_flags.financial_status)}
                                        </TableCell>
                                        <TableCell className="text-center">
                                            {getMedicalBadge(player.status_flags.medical_status)}
                                        </TableCell>
                                        <TableCell className="text-center">
                                            {getAttendanceBadge(player.status_flags.attendance_rate)}
                                        </TableCell>
                                        <TableCell className="text-center">
                                            <Tooltip>
                                                <TooltipTrigger>
                                                    {player.status_flags.is_inhabilitado ? (
                                                        <XCircle className="h-6 w-6 text-destructive mx-auto" />
                                                    ) : (
                                                        <CheckCircle className="h-6 w-6 text-green-500 mx-auto" />
                                                    )}
                                                </TooltipTrigger>
                                                <TooltipContent>
                                                    {player.status_flags.is_inhabilitado
                                                        ? "Inhabilitado para jugar"
                                                        : "Habilitado para jugar"
                                                    }
                                                </TooltipContent>
                                            </Tooltip>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </TooltipProvider>
            </CardContent>
        </Card>
    )
}
