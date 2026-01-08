"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
    ChevronLeft,
    ChevronRight,
    Bus,
    MapPin,
    Clock
} from "lucide-react"
import {
    format,
    startOfMonth,
    endOfMonth,
    eachDayOfInterval,
    isSameMonth,
    isSameDay,
    addMonths,
    subMonths,
    isToday
} from "date-fns"
import { es } from "date-fns/locale"

interface TravelEvent {
    id: string
    title: string
    destination: string
    departure_date: string
    meeting_time: string
    type: "TRAVEL" | "MATCH" | "TOURNAMENT" | "TRAINING"
}

interface TravelCalendarProps {
    teamId: string
}

export function TravelCalendar({ teamId }: TravelCalendarProps) {
    const [currentMonth, setCurrentMonth] = useState(new Date())
    const [events, setEvents] = useState<TravelEvent[]>([])
    const [selectedDate, setSelectedDate] = useState<Date | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        fetchEvents()
    }, [teamId, currentMonth])

    const fetchEvents = async () => {
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
    }

    const monthStart = startOfMonth(currentMonth)
    const monthEnd = endOfMonth(currentMonth)
    const days = eachDayOfInterval({ start: monthStart, end: monthEnd })

    // Obtener el día de la semana del primer día (0 = Domingo)
    const startDay = monthStart.getDay()
    const emptyDays = Array(startDay).fill(null)

    const getEventsForDay = (date: Date) => {
        return events.filter(event =>
            isSameDay(new Date(event.departure_date), date)
        )
    }

    const selectedEvents = selectedDate ? getEventsForDay(selectedDate) : []

    const getEventColor = (type: string) => {
        switch (type) {
            case "TRAVEL": return "bg-blue-500"
            case "MATCH": return "bg-green-500"
            case "TOURNAMENT": return "bg-purple-500"
            case "TRAINING": return "bg-orange-500"
            default: return "bg-gray-500"
        }
    }

    return (
        <div className="grid gap-6 lg:grid-cols-3">
            {/* Calendario */}
            <Card className="lg:col-span-2">
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <CardTitle>
                            {format(currentMonth, "MMMM yyyy", { locale: es })}
                        </CardTitle>
                        <div className="flex gap-2">
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => setCurrentMonth(subMonths(currentMonth, 1))}
                            >
                                <ChevronLeft className="h-4 w-4" />
                            </Button>
                            <Button
                                variant="outline"
                                size="icon"
                                onClick={() => setCurrentMonth(addMonths(currentMonth, 1))}
                            >
                                <ChevronRight className="h-4 w-4" />
                            </Button>
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    {/* Días de la semana */}
                    <div className="grid grid-cols-7 gap-1 mb-2">
                        {["Dom", "Lun", "Mar", "Mié", "Jue", "Vie", "Sáb"].map(day => (
                            <div key={day} className="text-center text-sm font-medium text-muted-foreground py-2">
                                {day}
                            </div>
                        ))}
                    </div>

                    {/* Días del mes */}
                    <div className="grid grid-cols-7 gap-1">
                        {/* Días vacíos al inicio */}
                        {emptyDays.map((_, index) => (
                            <div key={`empty-${index}`} className="h-20" />
                        ))}

                        {/* Días del mes */}
                        {days.map(day => {
                            const dayEvents = getEventsForDay(day)
                            const isSelected = selectedDate && isSameDay(day, selectedDate)

                            return (
                                <button
                                    key={day.toISOString()}
                                    onClick={() => setSelectedDate(day)}
                                    className={`
                    h-20 p-1 border rounded-lg text-left transition-colors
                    ${!isSameMonth(day, currentMonth) ? "text-muted-foreground" : ""}
                    ${isToday(day) ? "border-primary" : "border-border"}
                    ${isSelected ? "bg-primary/10 border-primary" : "hover:bg-muted"}
                  `}
                                >
                                    <span className={`
                    text-sm font-medium
                    ${isToday(day) ? "text-primary" : ""}
                  `}>
                                        {format(day, "d")}
                                    </span>

                                    {/* Indicadores de eventos */}
                                    <div className="mt-1 space-y-0.5">
                                        {dayEvents.slice(0, 2).map(event => (
                                            <div
                                                key={event.id}
                                                className={`
                          text-xs truncate px-1 py-0.5 rounded text-white
                          ${getEventColor(event.type)}
                        `}
                                            >
                                                {event.title}
                                            </div>
                                        ))}
                                        {dayEvents.length > 2 && (
                                            <div className="text-xs text-muted-foreground px-1">
                                                +{dayEvents.length - 2} más
                                            </div>
                                        )}
                                    </div>
                                </button>
                            )
                        })}
                    </div>

                    {/* Leyenda */}
                    <div className="flex gap-4 mt-4 pt-4 border-t">
                        <div className="flex items-center gap-1.5">
                            <div className="w-3 h-3 rounded bg-blue-500" />
                            <span className="text-xs">Viaje</span>
                        </div>
                        <div className="flex items-center gap-1.5">
                            <div className="w-3 h-3 rounded bg-green-500" />
                            <span className="text-xs">Partido</span>
                        </div>
                        <div className="flex items-center gap-1.5">
                            <div className="w-3 h-3 rounded bg-purple-500" />
                            <span className="text-xs">Torneo</span>
                        </div>
                        <div className="flex items-center gap-1.5">
                            <div className="w-3 h-3 rounded bg-orange-500" />
                            <span className="text-xs">Entrenamiento</span>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* Panel lateral - Eventos del día seleccionado */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-lg">
                        {selectedDate
                            ? format(selectedDate, "d 'de' MMMM", { locale: es })
                            : "Selecciona un día"
                        }
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {!selectedDate ? (
                        <p className="text-muted-foreground text-sm">
                            Haz clic en un día para ver los eventos
                        </p>
                    ) : selectedEvents.length === 0 ? (
                        <p className="text-muted-foreground text-sm">
                            No hay eventos programados
                        </p>
                    ) : (
                        <div className="space-y-4">
                            {selectedEvents.map(event => (
                                <div key={event.id} className="p-3 border rounded-lg space-y-2">
                                    <div className="flex items-start justify-between">
                                        <h4 className="font-medium">{event.title}</h4>
                                        <Badge variant="outline" className="text-xs">
                                            {event.type}
                                        </Badge>
                                    </div>

                                    <div className="space-y-1 text-sm text-muted-foreground">
                                        <div className="flex items-center gap-2">
                                            <MapPin className="h-3 w-3" />
                                            {event.destination}
                                        </div>
                                        <div className="flex items-center gap-2">
                                            <Clock className="h-3 w-3" />
                                            {format(new Date(event.meeting_time), "HH:mm")}hs
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}
