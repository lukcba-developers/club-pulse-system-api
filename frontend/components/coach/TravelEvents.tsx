import { useState, useEffect, useCallback } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogFooter,
} from "@/components/ui/dialog"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Calendar } from "@/components/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import {
    Bus,
    MapPin,
    Clock,
    Users,
    DollarSign,
    CheckCircle,
    XCircle,
    HelpCircle,
    Plus,
    CalendarIcon,
    RefreshCw
} from "lucide-react"
import { format } from "date-fns"
import { es } from "date-fns/locale"
import { useToast } from "@/hooks/use-toast"
import { cn } from "@/lib/utils"

interface TravelEvent {
    id: string
    team_id: string
    title: string
    description?: string
    destination: string
    departure_date: string
    return_date?: string
    meeting_point: string
    meeting_time: string
    estimated_cost: number
    actual_cost: number
    cost_per_person: number
    max_participants?: number
    type: "TRAVEL" | "MATCH" | "TOURNAMENT" | "TRAINING"
}

interface EventSummary {
    event: TravelEvent
    total_invited: number
    total_confirmed: number
    total_declined: number
    total_pending: number
    cost_per_person: number
}

interface TravelEventsProps {
    teamId: string
}

export function TravelEvents({ teamId }: TravelEventsProps) {
    const [events, setEvents] = useState<TravelEvent[]>([])
    const [loading, setLoading] = useState(true)
    const [selectedEvent, setSelectedEvent] = useState<EventSummary | null>(null)
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
    const { toast } = useToast()


    const fetchEvents = useCallback(async () => {
        setLoading(true)
        try {
            const response = await fetch(`/api/teams/${teamId}/events`, {
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            })
            if (response.ok) {
                const data = await response.json()
                setEvents(data || [])
            }
        } catch (error) {
            console.error("Error fetching events:", error)
        } finally {
            setLoading(false)
        }
    }, [teamId])

    useEffect(() => {
        fetchEvents()
    }, [fetchEvents])

    const fetchEventSummary = useCallback(async (eventId: string) => {
        try {
            const response = await fetch(`/api/events/${eventId}/summary`, {
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`
                }
            })
            if (response.ok) {
                const data = await response.json()
                setSelectedEvent(data)
            }
        } catch (error) {
            console.error("Error fetching summary:", error)
        }
    }, [])

    const handleRSVP = async (eventId: string, status: "CONFIRMED" | "DECLINED") => {
        try {
            const response = await fetch(`/api/events/${eventId}/rsvp`, {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ status })
            })

            if (response.ok) {
                toast({
                    title: "Éxito",
                    description: status === "CONFIRMED"
                        ? "¡Confirmaste tu asistencia!"
                        : "Registraste que no asistirás"
                })
                fetchEvents()
                if (selectedEvent?.event.id === eventId) {
                    fetchEventSummary(eventId)
                }
            }
        } catch {
            toast({
                title: "Error",
                description: "No se pudo registrar tu respuesta",
                variant: "destructive"
            })
        }
    }


    const getEventTypeBadge = (type: string) => {
        switch (type) {
            case "TRAVEL":
                return <Badge variant="secondary"><Bus className="h-3 w-3 mr-1" />Viaje</Badge>
            case "MATCH":
                return <Badge variant="default">⚽ Partido</Badge>
            case "TOURNAMENT":
                return <Badge variant="outline">🏆 Torneo</Badge>
            case "TRAINING":
                return <Badge variant="secondary">🏃 Entrenamiento</Badge>
            default:
                return <Badge>{type}</Badge>
        }
    }

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat("es-AR", {
            style: "currency",
            currency: "ARS"
        }).format(amount)
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight">Viajes y Eventos</h2>
                    <p className="text-muted-foreground">
                        Gestiona los viajes del equipo y confirma tu asistencia
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={fetchEvents} disabled={loading}>
                        <RefreshCw className={`h-4 w-4 mr-2 ${loading ? "animate-spin" : ""}`} />
                        Actualizar
                    </Button>
                    <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                        <DialogTrigger asChild>
                            <Button>
                                <Plus className="h-4 w-4 mr-2" />
                                Nuevo Evento
                            </Button>
                        </DialogTrigger>
                        <DialogContent>
                            <CreateEventForm
                                teamId={teamId}
                                onSuccess={() => {
                                    setIsCreateDialogOpen(false)
                                    fetchEvents()
                                }}
                            />
                        </DialogContent>
                    </Dialog>
                </div>
            </div>

            {/* Events Grid */}
            {loading ? (
                <div className="flex justify-center py-12">
                    <RefreshCw className="h-8 w-8 animate-spin text-muted-foreground" />
                </div>
            ) : events.length === 0 ? (
                <Card>
                    <CardContent className="flex flex-col items-center justify-center py-12">
                        <Bus className="h-12 w-12 text-muted-foreground mb-4" />
                        <p className="text-lg font-medium">No hay eventos programados</p>
                        <p className="text-muted-foreground">Crea un nuevo evento para comenzar</p>
                    </CardContent>
                </Card>
            ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                    {events.map((event) => (
                        <Card key={event.id} className="hover:shadow-lg transition-shadow">
                            <CardHeader>
                                <div className="flex justify-between items-start">
                                    <div>
                                        <CardTitle className="text-lg">{event.title}</CardTitle>
                                        <CardDescription className="flex items-center gap-1 mt-1">
                                            <MapPin className="h-3 w-3" />
                                            {event.destination}
                                        </CardDescription>
                                    </div>
                                    {getEventTypeBadge(event.type)}
                                </div>
                            </CardHeader>
                            <CardContent className="space-y-3">
                                <div className="flex items-center gap-2 text-sm">
                                    <CalendarIcon className="h-4 w-4 text-muted-foreground" />
                                    <span>
                                        {format(new Date(event.departure_date), "PPP", { locale: es })}
                                    </span>
                                </div>
                                <div className="flex items-center gap-2 text-sm">
                                    <Clock className="h-4 w-4 text-muted-foreground" />
                                    <span>
                                        Encuentro: {format(new Date(event.meeting_time), "HH:mm")}hs
                                    </span>
                                </div>
                                {event.meeting_point && (
                                    <div className="flex items-center gap-2 text-sm">
                                        <MapPin className="h-4 w-4 text-muted-foreground" />
                                        <span>{event.meeting_point}</span>
                                    </div>
                                )}
                                {event.cost_per_person > 0 && (
                                    <div className="flex items-center gap-2 text-sm font-medium">
                                        <DollarSign className="h-4 w-4 text-green-600" />
                                        <span>{formatCurrency(event.cost_per_person)} por persona</span>
                                    </div>
                                )}
                            </CardContent>
                            <CardFooter className="flex gap-2">
                                <Button
                                    variant="default"
                                    size="sm"
                                    className="flex-1"
                                    onClick={() => handleRSVP(event.id, "CONFIRMED")}
                                >
                                    <CheckCircle className="h-4 w-4 mr-1" />
                                    Voy
                                </Button>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    className="flex-1"
                                    onClick={() => handleRSVP(event.id, "DECLINED")}
                                >
                                    <XCircle className="h-4 w-4 mr-1" />
                                    No voy
                                </Button>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => fetchEventSummary(event.id)}
                                >
                                    <Users className="h-4 w-4" />
                                </Button>
                            </CardFooter>
                        </Card>
                    ))}
                </div>
            )}

            {/* Event Summary Dialog */}
            {selectedEvent && (
                <Dialog open={!!selectedEvent} onOpenChange={() => setSelectedEvent(null)}>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>{selectedEvent.event.title}</DialogTitle>
                            <DialogDescription>
                                Resumen de confirmaciones
                            </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4">
                            <div className="grid grid-cols-3 gap-4">
                                <Card>
                                    <CardContent className="pt-4 text-center">
                                        <CheckCircle className="h-8 w-8 text-green-500 mx-auto mb-2" />
                                        <p className="text-2xl font-bold">{selectedEvent.total_confirmed}</p>
                                        <p className="text-sm text-muted-foreground">Confirmados</p>
                                    </CardContent>
                                </Card>
                                <Card>
                                    <CardContent className="pt-4 text-center">
                                        <XCircle className="h-8 w-8 text-red-500 mx-auto mb-2" />
                                        <p className="text-2xl font-bold">{selectedEvent.total_declined}</p>
                                        <p className="text-sm text-muted-foreground">No asisten</p>
                                    </CardContent>
                                </Card>
                                <Card>
                                    <CardContent className="pt-4 text-center">
                                        <HelpCircle className="h-8 w-8 text-yellow-500 mx-auto mb-2" />
                                        <p className="text-2xl font-bold">{selectedEvent.total_pending}</p>
                                        <p className="text-sm text-muted-foreground">Pendientes</p>
                                    </CardContent>
                                </Card>
                            </div>

                            {selectedEvent.cost_per_person > 0 && (
                                <Card>
                                    <CardContent className="pt-4">
                                        <div className="flex justify-between items-center">
                                            <span className="font-medium">Costo por persona:</span>
                                            <span className="text-xl font-bold text-green-600">
                                                {formatCurrency(selectedEvent.cost_per_person)}
                                            </span>
                                        </div>
                                        <p className="text-sm text-muted-foreground mt-1">
                                            Calculado entre {selectedEvent.total_confirmed} confirmados
                                        </p>
                                    </CardContent>
                                </Card>
                            )}
                        </div>
                    </DialogContent>
                </Dialog>
            )}
        </div>
    )
}

// Formulario para crear eventos
function CreateEventForm({ teamId, onSuccess }: { teamId: string; onSuccess: () => void }) {
    const [loading, setLoading] = useState(false)
    const [date, setDate] = useState<Date>()
    const { toast } = useToast()

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        setLoading(true)

        const formData = new FormData(e.currentTarget)

        try {
            const response = await fetch("/api/events", {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${localStorage.getItem("token")}`,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    team_id: teamId,
                    title: formData.get("title"),
                    description: formData.get("description"),
                    destination: formData.get("destination"),
                    departure_date: date?.toISOString(),
                    meeting_point: formData.get("meeting_point"),
                    meeting_time: date?.toISOString(),
                    estimated_cost: parseFloat(formData.get("estimated_cost") as string) || 0,
                    type: formData.get("type")
                })
            })

            if (response.ok) {
                toast({
                    title: "Éxito",
                    description: "Evento creado correctamente"
                })
                onSuccess()
            } else {
                throw new Error("Error al crear evento")
            }
        } catch {
            toast({
                title: "Error",
                description: "No se pudo crear el evento",
                variant: "destructive"
            })
        } finally {
            setLoading(false)
        }
    }

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            <DialogHeader>
                <DialogTitle>Nuevo Evento</DialogTitle>
                <DialogDescription>
                    Crea un nuevo viaje o evento para el equipo
                </DialogDescription>
            </DialogHeader>

            <div className="space-y-4">
                <div className="space-y-2">
                    <Label htmlFor="title">Título</Label>
                    <Input id="title" name="title" placeholder="Partido vs Club Rival" required />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="type">Tipo de evento</Label>
                    <Select name="type" defaultValue="TRAVEL">
                        <SelectTrigger>
                            <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="TRAVEL">Viaje</SelectItem>
                            <SelectItem value="MATCH">Partido</SelectItem>
                            <SelectItem value="TOURNAMENT">Torneo</SelectItem>
                            <SelectItem value="TRAINING">Entrenamiento</SelectItem>
                        </SelectContent>
                    </Select>
                </div>

                <div className="space-y-2">
                    <Label htmlFor="destination">Destino</Label>
                    <Input id="destination" name="destination" placeholder="Villa Carlos Paz" required />
                </div>

                <div className="space-y-2">
                    <Label>Fecha de salida</Label>
                    <Popover>
                        <PopoverTrigger asChild>
                            <Button
                                variant="outline"
                                className={cn(
                                    "w-full justify-start text-left font-normal",
                                    !date && "text-muted-foreground"
                                )}
                            >
                                <CalendarIcon className="mr-2 h-4 w-4" />
                                {date ? format(date, "PPP", { locale: es }) : "Seleccionar fecha"}
                            </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0">
                            <Calendar
                                mode="single"
                                selected={date}
                                onSelect={setDate}
                                locale={es}
                            />
                        </PopoverContent>
                    </Popover>
                </div>

                <div className="space-y-2">
                    <Label htmlFor="meeting_point">Punto de encuentro</Label>
                    <Input id="meeting_point" name="meeting_point" placeholder="Sede del club" />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="estimated_cost">Costo estimado (ARS)</Label>
                    <Input
                        id="estimated_cost"
                        name="estimated_cost"
                        type="number"
                        placeholder="5000"
                    />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="description">Descripción (opcional)</Label>
                    <Textarea
                        id="description"
                        name="description"
                        placeholder="Detalles adicionales del viaje..."
                    />
                </div>
            </div>

            <DialogFooter>
                <Button type="submit" disabled={loading}>
                    {loading ? "Creando..." : "Crear Evento"}
                </Button>
            </DialogFooter>
        </form>
    )
}
