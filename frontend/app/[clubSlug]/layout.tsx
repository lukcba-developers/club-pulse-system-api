import { Metadata } from "next"
import { notFound } from "next/navigation"
import Link from "next/link"
import Image from "next/image"
import { SponsorCarousel } from "@/components/club/SponsorCarousel"

// Función para obtener datos del club
async function getClub(slug: string) {
    try {
        const res = await fetch(`${process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/public/clubs/${slug}`, {
            next: { revalidate: 3600 } // Cache por 1 hora
        })

        if (!res.ok) return null
        return res.json()
    } catch (error) {
        console.error("Error fetching club:", error)
        return null
    }
}

export async function generateMetadata({ params }: { params: Promise<{ clubSlug: string }> }): Promise<Metadata> {
    const { clubSlug } = await params
    const club = await getClub(clubSlug)
    if (!club) return { title: "Club no encontrado" }

    return {
        title: club.name,
        description: `Bienvenido al sitio oficial de ${club.name}`,
        icons: {
            icon: club.logo_url || "/favicon.ico"
        }
    }
}

export default async function ClubLayout({
    children,
    params
}: {
    children: React.ReactNode
    params: Promise<{ clubSlug: string }>
}) {
    const { clubSlug } = await params
    const club = await getClub(clubSlug)

    if (!club) {
        notFound()
    }

    // Parsear configuración de tema
    let theme = {
        primaryColor: "#0f172a", // slate-900 default
        secondaryColor: "#475569", // slate-600
        accentColor: "#3b82f6", // blue-500
        fontFamily: "Inter"
    }

    if (club.theme_config) {
        try {
            // theme_config viene como string JSON desde la BD si no se parseó antes,
            // pero si el endpoint devuelve objeto ya parseado (gorm serializer:json), será objeto.
            // Asumiremos que es objeto si el backend lo maneja bien, o string si no.
            const config = typeof club.theme_config === 'string'
                ? JSON.parse(club.theme_config)
                : club.theme_config

            theme = { ...theme, ...config }
        } catch (e) {
            console.error("Error parsing theme config", e)
        }
    }

    // Estilos dinámicos
    const style = {
        "--primary": theme.primaryColor,
        "--secondary": theme.secondaryColor,
        "--accent": theme.accentColor,
    } as React.CSSProperties

    return (
        <div className="min-h-screen bg-background" style={style}>
            {/* Header Público */}
            <header className="border-b bg-card">
                <div className="container mx-auto px-4 py-4 flex justify-between items-center">
                    <div className="flex items-center gap-3">
                        {club.logo_url && (
                            <Image src={club.logo_url} alt={club.name} width={40} height={40} className="h-10 w-10 object-contain" />
                        )}
                        <h1 className="text-xl font-bold">{club.name}</h1>
                    </div>
                    <nav className="flex gap-4">
                        <Link href={`/${clubSlug}`} className="hover:text-primary transition-colors">Inicio</Link>
                        <Link href={`/${clubSlug}/championships`} className="hover:text-primary transition-colors">Campeonatos</Link>
                        <Link href={`/${clubSlug}/store`} className="hover:text-primary transition-colors">Tienda</Link>
                        <Link href="/login" className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:opacity-90 transition-opacity">
                            Acceso Socios
                        </Link>
                    </nav>
                </div>
            </header>

            {/* Contenido Principal */}
            <main>
                {children}
            </main>

            {/* Sponsors */}
            <SponsorCarousel clubSlug={clubSlug} />

            {/* Footer Público */}
            <footer className="bg-slate-900 text-white py-8 mt-auto">
                <div className="container mx-auto px-4 text-center">
                    <p>© {new Date().getFullYear()} {club.name}. Todos los derechos reservados.</p>
                    <p className="text-sm text-slate-400 mt-2">Powered by Club Pulse</p>
                </div>
            </footer>
        </div>
    )
}
