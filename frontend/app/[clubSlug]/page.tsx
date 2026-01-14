import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Calendar, Trophy, Users, ShoppingBag } from "lucide-react"
import Link from "next/link"
import Image from "next/image"

// Función auxiliar para obtener noticias (client o server side fetch)
async function getNews(slug: string) {
    try {
        const res = await fetch(`${process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/public/clubs/${slug}/news`, { next: { revalidate: 300 } })
        if (!res.ok) return []
        const data = await res.json()
        return data.data || []
    } catch (error) {
        console.error("Error fetching news:", error)
        return []
    }
}
interface NewsItem {
    id: string
    title: string
    content: string
    image_url?: string
    created_at: string
}

export default async function PublicClubHome({ params }: { params: Promise<{ clubSlug: string }> }) {
    const { clubSlug } = await params
    const news = await getNews(clubSlug)

    return (
        <div className="space-y-12 pb-12">
            {/* Hero Section */}
            <section className="relative py-24 bg-gradient-to-br from-primary/10 to-accent/5">
                <div className="container mx-auto px-4 text-center">
                    <h1 className="text-4xl md:text-6xl font-extrabold tracking-tight mb-6 text-foreground">
                        Bienvenido al Club
                    </h1>
                    <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-8">
                        Pasión, deporte y comunidad. Únete a nosotros y forma parte de la historia.
                    </p>
                    <div className="flex justify-center gap-4">
                        <Link href={`/${clubSlug}/store`}>
                            <Button size="lg" className="gap-2">
                                <ShoppingBag className="h-5 w-5" />
                                Ir a la Tienda
                            </Button>
                        </Link>
                        <Link href="/register">
                            <Button size="lg" variant="outline" className="gap-2">
                                <Users className="h-5 w-5" />
                                Hacerme Socio
                            </Button>
                        </Link>
                    </div>
                </div>
            </section>

            {/* Featured Sections Grid */}
            <section className="container mx-auto px-4">
                <div className="grid md:grid-cols-3 gap-6">
                    <Card>
                        <CardHeader>
                            <Trophy className="h-10 w-10 text-primary mb-2" />
                            <CardTitle>Torneos Activos</CardTitle>
                            <CardDescription>Sigue el desempeño de nuestros equipos</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <p className="mb-4">Consulta fixtures, resultados y tablas de posiciones de todas las categorías.</p>
                            <Button variant="link" className="p-0">Ver Torneos &rarr;</Button>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <Calendar className="h-10 w-10 text-primary mb-2" />
                            <CardTitle>Próximos Eventos</CardTitle>
                            <CardDescription>No te pierdas ninguna actividad</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <ul className="space-y-2 mb-4">
                                <li className="text-sm">• Cena de fin de año - 20 Dic</li>
                                <li className="text-sm">• Torneo de Verano - 15 Ene</li>
                            </ul>
                            <Button variant="link" className="p-0">Ver Calendario &rarr;</Button>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <ShoppingBag className="h-10 w-10 text-primary mb-2" />
                            <CardTitle>Tienda Oficial</CardTitle>
                            <CardDescription>Merchandising y equipamiento</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <p className="mb-4">Adquiere la camiseta oficial, ropa de entrenamiento y accesorios del club.</p>
                            <Link href={`/${clubSlug}/store`}>
                                <Button variant="link" className="p-0">Visitar Tienda &rarr;</Button>
                            </Link>
                        </CardContent>
                    </Card>
                </div>
            </section>

            {/* News Section */}
            <section className="container mx-auto px-4">
                <h2 className="text-3xl font-bold mb-8 text-center">Últimas Novedades</h2>
                {news.length > 0 ? (
                    <div className="grid md:grid-cols-2 gap-8">
                        {news.map((item: NewsItem) => (
                            <div key={item.id} className="bg-card rounded-lg overflow-hidden border shadow-sm flex flex-col">
                                <div className="h-48 bg-muted flex items-center justify-center text-muted-foreground relative">
                                    {item.image_url ? (
                                        <Image src={item.image_url} alt={item.title} fill className="object-cover" />
                                    ) : (
                                        <span>Sin Imagen</span>
                                    )}
                                </div>
                                <div className="p-6 flex-grow">
                                    <h3 className="text-xl font-bold mb-2">{item.title}</h3>
                                    <p className="text-muted-foreground mb-4 line-clamp-3">{item.content}</p>
                                    <span className="text-xs text-slate-400">{new Date(item.created_at).toLocaleDateString()}</span>
                                </div>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="text-center text-muted-foreground p-8">No hay novedades recientes.</div>
                )}
            </section>
        </div>
    )
}
