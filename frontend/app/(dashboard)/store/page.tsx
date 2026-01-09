"use client";

import { useEffect, useState } from "react";
import { storeService, Product } from "@/services/store-service";
import { ProductCard } from "@/components/store/product-card";
import { DashboardHeader } from "@/components/dashboard/header";
import { Shell } from "@/components/layout/shell";
import { useToast } from "@/components/ui/use-toast";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function StorePage() {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);
    const { toast } = useToast();
    const [cart, setCart] = useState<{ product: Product, quantity: number }[]>([]);

    useEffect(() => {
        fetchProducts();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    const fetchProducts = async () => {
        try {
            setLoading(true);
            const data = await storeService.getProducts();
            setProducts(data || []);
        } catch {
            toast({
                title: "Error",
                description: "No se pudieron cargar los productos.",
                variant: "destructive",
            });
        } finally {
            setLoading(false);
        }
    };

    const handleAddToCart = (product: Product) => {
        setCart(prev => {
            const existing = prev.find(item => item.product.id === product.id);
            if (existing) {
                return prev.map(item =>
                    item.product.id === product.id
                        ? { ...item, quantity: item.quantity + 1 }
                        : item
                );
            }
            return [...prev, { product, quantity: 1 }];
        });
        toast({
            title: "Agregado",
            description: `${product.name} agregado al carrito.`,
        });
    };

    const handleCheckout = async () => {
        if (cart.length === 0) return;
        try {
            const items = cart.map(item => ({ product_id: item.product.id, quantity: item.quantity }));
            await storeService.purchaseItems(items, 'CASH'); // Mocking CASH payment for now
            toast({
                title: "Compra Exitosa",
                description: "Tu pedido ha sido registrado.",
            });
            setCart([]);
            fetchProducts(); // Refresh stock
        } catch {
            toast({
                title: "Error",
                description: "No se pudo procesar la compra.",
                variant: "destructive",
            });
        }
    };

    return (
        <Shell>
            <DashboardHeader heading="Tienda del Club" text="Adquiere indumentaria y accesorios oficiales.">
                {cart.length > 0 && (
                    <Button onClick={handleCheckout}>
                        Pagar Carrito ({cart.reduce((acc, item) => acc + item.quantity, 0)})
                    </Button>
                )}
            </DashboardHeader>

            {loading ? (
                <div className="flex justify-center p-8">
                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                </div>
            ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                    {products.length > 0 ? (
                        products.map(product => (
                            <ProductCard
                                key={product.id}
                                product={product}
                                onAddToCart={handleAddToCart}
                            />
                        ))
                    ) : (
                        <div className="col-span-full text-center text-muted-foreground p-8">
                            No hay productos disponibles en este momento.
                        </div>
                    )}
                </div>
            )}
        </Shell>
    );
}
