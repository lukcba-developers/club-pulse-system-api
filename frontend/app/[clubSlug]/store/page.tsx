import PublicStoreClient from "@/components/store/PublicStoreClient"

async function getProducts(slug: string) {
    try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/public/clubs/${slug}/store/products`, {
            next: { revalidate: 60 } // Cache corta
        })
        if (!res.ok) return []
        const data = await res.json()
        return data.data || []
    } catch (error) {
        console.error("Error fetching products", error)
        return []
    }
}

export default async function PublicStorePage({ params }: { params: { clubSlug: string } }) {
    const products = await getProducts(params.clubSlug)

    return <PublicStoreClient products={products} clubSlug={params.clubSlug} />
}
