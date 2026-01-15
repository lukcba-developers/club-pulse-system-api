"use client"

import { useState } from "react"
import { ShoppingCart, Trash2, CreditCard, Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"

interface Product {
    id: string
    name: string
    description: string
    price: string
    image_url?: string
    stock_quantity: number
    is_active: boolean
    created_at: string
    updated_at: string
}

interface CartItem extends Product {
    quantity: number
}

export default function PublicStoreClient({ products, clubSlug }: { products: Product[], clubSlug: string }) {
    const [cart, setCart] = useState<CartItem[]>([])
    const [isCheckoutOpen, setIsCheckoutOpen] = useState(false)

    // Checkout Form State
    const [guestName, setGuestName] = useState("")
    const [guestEmail, setGuestEmail] = useState("")
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [orderSuccess, setOrderSuccess] = useState(false)

    const addToCart = (product: Product) => {
        setCart(prev => {
            const existing = prev.find(item => item.id === product.id)
            if (existing) {
                // Validación de stock: No permitir agregar más del stock disponible
                if (existing.quantity >= product.stock_quantity) {
                    return prev
                }
                return prev.map(item =>
                    item.id === product.id ? { ...item, quantity: item.quantity + 1 } : item
                )
            }
            // Validación de stock inicial: Asegurar que hay al menos 1 disponible
            if (product.stock_quantity <= 0) {
                return prev
            }
            return [...prev, { ...product, quantity: 1 }]
        })
    }

    const removeFromCart = (productId: string) => {
        setCart(prev => prev.filter(item => item.id !== productId))
    }

    const totalAmount = cart.reduce((sum, item) => sum + (parseFloat(item.price) * item.quantity), 0)

    const handleCheckout = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!guestName || !guestEmail || cart.length === 0) return

        setIsSubmitting(true)
        try {
            const itemsPayload = cart.map(item => ({
                product_id: item.id,
                quantity: item.quantity
            }))

            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/public/clubs/${clubSlug}/store/purchase`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    guest_name: guestName,
                    guest_email: guestEmail,
                    items: itemsPayload
                })
            })

            if (res.ok) {
                setOrderSuccess(true)
                setCart([])
            } else {
                const err = await res.json()
                alert("Error al procesar la orden: " + (err.error || "Desconocido"))
            }
        } catch (error) {
            console.error(error)
            alert("Error de conexión")
        } finally {
            setIsSubmitting(false)
        }
    }

    if (orderSuccess) {
        return (
            <div className="container mx-auto px-4 py-12 text-center">
                <Card className="max-w-md mx-auto">
                    <CardHeader>
                        <div className="mx-auto w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mb-4">
                            <ShoppingCart className="text-green-600 w-6 h-6" />
                        </div>
                        <CardTitle className="text-2xl text-green-700">¡Compra Exitosa!</CardTitle>
                        <CardDescription>
                            Gracias por tu compra, {guestName}. Hemos enviado un comprobante a {guestEmail}.
                        </CardDescription>
                    </CardHeader>
                    <CardFooter className="justify-center">
                        <Button onClick={() => setOrderSuccess(false)}>Seguir Comprando</Button>
                    </CardFooter>
                </Card>
            </div>
        )
    }

    return (
        <div className="container mx-auto px-4 py-8 relative min-h-screen">
            <div className="flex justify-between items-center mb-8">
                <h1 className="text-3xl font-bold">Tienda Oficial</h1>

                {/* Cart Trigger */}
                <Dialog open={isCheckoutOpen} onOpenChange={setIsCheckoutOpen}>
                    <DialogTrigger asChild>
                        <Button size="lg" className="relative">
                            <ShoppingCart className="mr-2 h-5 w-5" />
                            Ver Carrito
                            {cart.length > 0 && (
                                <Badge variant="destructive" className="absolute -top-2 -right-2">
                                    {cart.reduce((a, b) => a + b.quantity, 0)}
                                </Badge>
                            )}
                        </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-lg">
                        <DialogHeader>
                            <DialogTitle>Tu Carrito de Compras</DialogTitle>
                            <DialogDescription>
                                Revisa tus productos antes de finalizar la compra.
                            </DialogDescription>
                        </DialogHeader>

                        {cart.length === 0 ? (
                            <div className="py-8 text-center text-muted-foreground">
                                Tu carrito está vacío.
                            </div>
                        ) : (
                            <div className="space-y-4 max-h-[60vh] overflow-y-auto pr-2">
                                {cart.map(item => (
                                    <div key={item.id} className="flex justify-between items-center bg-muted/50 p-3 rounded-lg">
                                        <div className="flex-1">
                                            <h4 className="font-semibold">{item.name}</h4>
                                            <div className="text-sm text-muted-foreground">
                                                {item.quantity} x ${parseFloat(item.price).toFixed(2)}
                                            </div>
                                        </div>
                                        <div className="font-bold mr-4">
                                            ${(item.quantity * parseFloat(item.price)).toFixed(2)}
                                        </div>
                                        <Button variant="ghost" size="icon" onClick={() => removeFromCart(item.id)}>
                                            <Trash2 className="h-4 w-4 text-destructive" />
                                        </Button>
                                    </div>
                                ))}

                                <div className="border-t pt-4 mt-4 flex justify-between items-center">
                                    <span className="font-bold text-lg">Total a Pagar:</span>
                                    <span className="font-bold text-2xl text-primary">${totalAmount.toFixed(2)}</span>
                                </div>

                                <form id="checkout-form" onSubmit={handleCheckout} className="space-y-4 pt-4 border-t mt-4">
                                    <div className="space-y-2">
                                        <Label htmlFor="guestName">Nombre Completo</Label>
                                        <Input
                                            id="guestName"
                                            required
                                            placeholder="Juan Pérez"
                                            value={guestName}
                                            onChange={(e) => setGuestName(e.target.value)}
                                        />
                                    </div>
                                    <div className="space-y-2">
                                        <Label htmlFor="guestEmail">Correo Electrónico</Label>
                                        <Input
                                            id="guestEmail"
                                            type="email"
                                            required
                                            placeholder="juan@ejemplo.com"
                                            value={guestEmail}
                                            onChange={(e) => setGuestEmail(e.target.value)}
                                        />
                                    </div>
                                </form>
                            </div>
                        )}

                        <DialogFooter className="mt-4">
                            {cart.length > 0 && (
                                <Button type="submit" form="checkout-form" className="w-full" disabled={isSubmitting}>
                                    {isSubmitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <CreditCard className="mr-2 h-4 w-4" />}
                                    {isSubmitting ? "Procesando..." : "Confirmar Compra"}
                                </Button>
                            )}
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            {products.length === 0 ? (
                <div className="text-center py-12 bg-muted rounded-lg">
                    <p className="text-muted-foreground">Próximamente productos disponibles.</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
                    {products.map((product) => (
                        <Card key={product.id} className="overflow-hidden flex flex-col hover:shadow-lg transition-shadow">
                            <div className="aspect-square bg-muted relative">
                                {product.image_url ? (
                                    /* eslint-disable @next/next/no-img-element */
                                    <img src={product.image_url} alt={product.name} className="object-cover w-full h-full" />
                                ) : (
                                    <div className="flex items-center justify-center h-full text-muted-foreground bg-slate-100">
                                        Sin Imagen
                                    </div>
                                )}
                                {product.stock_quantity <= 0 && (
                                    <div className="absolute inset-0 bg-black/60 flex items-center justify-center text-white font-bold backdrop-blur-sm">
                                        AGOTADO
                                    </div>
                                )}
                            </div>
                            <CardHeader>
                                <CardTitle className="text-lg">{product.name}</CardTitle>
                                <CardDescription className="line-clamp-2">{product.description}</CardDescription>
                            </CardHeader>
                            <CardContent className="flex-grow">
                                <div className="font-bold text-xl text-primary">${parseFloat(product.price).toFixed(2)}</div>
                                <div className="text-xs text-muted-foreground mt-1">Stock: {product.stock_quantity}</div>
                            </CardContent>
                            <CardFooter>
                                <Button
                                    className="w-full gap-2"
                                    onClick={() => addToCart(product)}
                                    disabled={product.stock_quantity <= 0}
                                >
                                    <ShoppingCart className="h-4 w-4" />
                                    Agregar
                                </Button>
                            </CardFooter>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    )
}
