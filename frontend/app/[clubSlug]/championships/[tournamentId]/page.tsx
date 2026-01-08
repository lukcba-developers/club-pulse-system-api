import { TournamentClientView } from "@/components/championship/TournamentClientView"
import { Button } from "@/components/ui/button"
import { ChevronLeft } from "lucide-react"
import Link from "next/link"
import { notFound } from "next/navigation"

async function getTournament(slug: string, id: string) {
    try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/public/clubs/${slug}/championships/${id}`, {
            next: { revalidate: 60 }
        })
        if (!res.ok) {
            if (res.status === 404) return null
            throw new Error("Failed to fetch")
        }
        return await res.json()
    } catch (error) {
        console.error("Error fetching tournament detail", error)
        return null
    }
}

export default async function TournamentDetailPage({ params }: { params: { clubSlug: string, tournamentId: string } }) {
    const tournament = await getTournament(params.clubSlug, params.tournamentId)

    if (!tournament) {
        notFound()
    }

    return (
        <div className="container mx-auto px-4 py-8">
            <Link href={`/${params.clubSlug}/championships`} className="inline-block mb-6">
                <Button variant="ghost" className="gap-2 pl-0 hover:pl-0 hover:bg-transparent text-muted-foreground hover:text-foreground">
                    <ChevronLeft className="h-4 w-4" />
                    Volver a Campeonatos
                </Button>
            </Link>

            <TournamentClientView tournament={tournament} clubSlug={params.clubSlug} />
        </div>
    )
}
