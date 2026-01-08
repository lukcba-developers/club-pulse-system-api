import Link from "next/link"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Calendar, Trophy, Users } from "lucide-react"
import { Button } from "@/components/ui/button"

interface Tournament {
    id: string
    name: string
    sport: string
    category: string
    status: string
    start_date: string
    logo_url?: string
}

async function getTournaments(slug: string) {
    try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/public/clubs/${slug}/championships`, {
            next: { revalidate: 60 }
        })
        if (!res.ok) return []
        return await res.json()
    } catch (error) {
        console.error("Error fetching tournaments", error)
        return []
    }
}

export default async function ChampionshipsPage({ params }: { params: { clubSlug: string } }) {
    const tournaments: Tournament[] = await getTournaments(params.clubSlug)

    return (
        <div className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold mb-8 flex items-center gap-2">
                <Trophy className="h-8 w-8 text-primary" />
                Campeonatos y Torneos
            </h1>

            {tournaments.length === 0 ? (
                <div className="text-center py-12 bg-muted rounded-lg border-2 border-dashed">
                    <Trophy className="mx-auto h-12 w-12 text-muted-foreground opacity-50 mb-4" />
                    <h3 className="text-lg font-semibold">No hay torneos activos</h3>
                    <p className="text-muted-foreground">Pronto se anunciarán nuevas competencias.</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {tournaments.map((t) => (
                        <Card key={t.id} className="flex flex-col hover:shadow-lg transition-shadow">
                            <CardHeader>
                                <div className="flex justify-between items-start">
                                    <div className="space-y-1">
                                        <CardTitle className="text-xl">{t.name}</CardTitle>
                                        <CardDescription>{t.sport} • {t.category}</CardDescription>
                                    </div>
                                    <Badge variant={t.status === 'ACTIVE' ? "default" : "secondary"}>
                                        {t.status === 'ACTIVE' ? 'En Curso' : t.status === 'DRAFT' ? 'Próximamente' : 'Finalizado'}
                                    </Badge>
                                </div>
                            </CardHeader>
                            <CardContent className="flex-grow space-y-4">
                                <div className="flex items-center text-sm text-muted-foreground gap-2">
                                    <Calendar className="h-4 w-4" />
                                    <span>Inicio: {new Date(t.start_date).toLocaleDateString()}</span>
                                </div>
                                {/* Placeholder for enrolled teams count if available */}
                            </CardContent>
                            <CardFooter>
                                <Link href={`/${params.clubSlug}/championships/${t.id}`} className="w-full">
                                    <Button className="w-full" variant="outline">Ver Detalles</Button>
                                </Link>
                            </CardFooter>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    )
}
