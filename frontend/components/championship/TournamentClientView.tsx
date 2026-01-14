"use client"

import { useState, useEffect } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card"
import { PublicStandingsTable as StandingsTable } from "./PublicStandingsTable"
import { PublicFixtureList as FixtureList } from "./PublicFixtureList"
import { Calendar, Trophy, Users } from "lucide-react"

interface Tournament {
    id: string
    name: string
    sport: string
    category: string
    status: string
    start_date: string
    description?: string
    stages: Stage[]
}

interface Stage {
    id: string
    name: string
    type: string
    order: number
    groups: Group[]
}

interface Group {
    id: string
    name: string
}

export function TournamentClientView({ tournament, clubSlug }: { tournament: Tournament, clubSlug: string }) {
    // Determine initial defaults
    const initialStage = tournament.stages?.[0]
    const initialGroup = initialStage?.groups?.[0]

    const [selectedStageId, setSelectedStageId] = useState<string>(initialStage?.id || "")
    const [selectedGroupId, setSelectedGroupId] = useState<string>(initialGroup?.id || "")

    // Data State
    const [standings, setStandings] = useState([])
    const [matches, setMatches] = useState([])
    const [loadingData, setLoadingData] = useState(false)

    // Derived
    const currentStage = tournament.stages?.find(s => s.id === selectedStageId)
    const currentGroups = currentStage?.groups || []

    // Auto-select first group when stage changes
    useEffect(() => {
        if (currentStage && currentStage.groups?.length > 0) {
            // Only update if current selectedGroupId is not in this stage
            const groupExists = currentStage.groups.find(g => g.id === selectedGroupId)
            if (!groupExists) {
                setSelectedGroupId(currentStage.groups[0].id)
            }
        }
    }, [selectedStageId, currentStage, selectedGroupId])

    // Fetch Data when Group Changes or Tab mounts
    useEffect(() => {
        if (!selectedGroupId) return

        const fetchData = async () => {
            setLoadingData(true)
            try {
                const apiBase = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

                // Fetch Standings
                const resStandings = await fetch(`${apiBase}/public/clubs/${clubSlug}/championships/groups/${selectedGroupId}/standings`)
                if (resStandings.ok) {
                    const data = await resStandings.json()
                    setStandings(data)
                }

                // Fetch Fixture
                const resFixture = await fetch(`${apiBase}/public/clubs/${clubSlug}/championships/groups/${selectedGroupId}/fixture`)
                if (resFixture.ok) {
                    const data = await resFixture.json()
                    setMatches(data)
                }

            } catch (error) {
                console.error("Error fetching group data", error)
            } finally {
                setLoadingData(false)
            }
        }

        fetchData()
    }, [selectedGroupId, clubSlug])

    return (
        <div className="space-y-8">
            {/* Header Info */}
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 border-b pb-6">
                <div>
                    <h1 className="text-3xl font-bold flex items-center gap-3">
                        <Trophy className="h-8 w-8 text-primary" />
                        {tournament.name}
                    </h1>
                    <div className="flex items-center gap-4 mt-2 text-muted-foreground">
                        <span className="flex items-center gap-1"><Users className="h-4 w-4" /> {tournament.category}</span>
                        <span className="flex items-center gap-1"><Calendar className="h-4 w-4" /> {new Date(tournament.start_date).toLocaleDateString()}</span>
                        <span className="bg-primary/10 text-primary px-2 py-0.5 rounded text-xs font-semibold">{tournament.sport}</span>
                    </div>
                </div>
            </div>

            {/* Selectors */}
            {tournament.stages?.length > 0 && (
                <div className="flex flex-col sm:flex-row gap-4 bg-muted/30 p-4 rounded-lg">
                    <div className="flex-1">
                        <label className="text-sm font-medium mb-1 block text-muted-foreground">Fase / Etapa</label>
                        <Select value={selectedStageId} onValueChange={setSelectedStageId}>
                            <SelectTrigger>
                                <SelectValue placeholder="Selecciona fase" />
                            </SelectTrigger>
                            <SelectContent>
                                {tournament.stages.map(stage => (
                                    <SelectItem key={stage.id} value={stage.id}>{stage.name}</SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    {currentGroups.length > 0 && (
                        <div className="flex-1">
                            <label className="text-sm font-medium mb-1 block text-muted-foreground">Grupo</label>
                            <Select value={selectedGroupId} onValueChange={setSelectedGroupId}>
                                <SelectTrigger>
                                    <SelectValue placeholder="Selecciona grupo" />
                                </SelectTrigger>
                                <SelectContent>
                                    {currentGroups.map(group => (
                                        <SelectItem key={group.id} value={group.id}>{group.name}</SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>
                    )}
                </div>
            )}

            {/* Content Tabs */}
            <Tabs defaultValue="standings" className="w-full">
                <TabsList className="grid w-full grid-cols-3 md:w-[400px]">
                    <TabsTrigger value="standings">Tabla de Posiciones</TabsTrigger>
                    <TabsTrigger value="fixture">Fixture / Partidos</TabsTrigger>
                    <TabsTrigger value="info">Información</TabsTrigger>
                </TabsList>

                <TabsContent value="standings" className="mt-6">
                    <Card>
                        <CardHeader>
                            <CardTitle>Tabla de Posiciones</CardTitle>
                            <CardDescription>
                                {currentStage?.name} - {currentGroups.find(g => g.id === selectedGroupId)?.name}
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <StandingsTable standings={standings} loading={loadingData} />
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="fixture" className="mt-6">
                    <Card>
                        <CardHeader>
                            <CardTitle>Partidos</CardTitle>
                            <CardDescription>Calendario de encuentros</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <FixtureList matches={matches} loading={loadingData} />
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="info" className="mt-6">
                    <Card>
                        <CardHeader>
                            <CardTitle>Acerca del Torneo</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div>
                                <h4 className="font-semibold">Descripción</h4>
                                <p className="text-muted-foreground">{tournament.description || "Sin descripción disponible."}</p>
                            </div>
                            {/* More info placeholders */}
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    )
}
