import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Product } from "@/services/store-service";
import { ShoppingCart } from "lucide-react";

interface ProductCardProps {
    product: Product;
    onAddToCart: (product: Product) => void;
}

export function ProductCard({ product, onAddToCart }: ProductCardProps) {
    const isOutOfStock = product.stock_quantity <= 0;

    return (
        <Card className="flex flex-col h-full">
            <CardHeader className="p-4">
                <div className="flex justify-between items-start">
                    <CardTitle className="text-lg line-clamp-1 truncate" title={product.name}>{product.name}</CardTitle>
                    <Badge variant={isOutOfStock ? "destructive" : "secondary"}>
                        {isOutOfStock ? "Sin Stock" : `$${product.price}`}
                    </Badge>
                </div>
                <CardDescription className="line-clamp-2 text-sm h-10 overflow-hidden">
                    {product.description}
                </CardDescription>
            </CardHeader>
            <CardContent className="flex-grow p-4 pt-0">
                {product.image_url ? (
                    <div className="w-full h-32 bg-gray-100 rounded-md mb-2 overflow-hidden">
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img src={product.image_url} alt={product.name} className="w-full h-full object-cover" />
                    </div>
                ) : (
                    <div className="w-full h-32 bg-gray-100 rounded-md mb-2 flex items-center justify-center text-gray-400 text-sm">
                        Sin Imagen
                    </div>
                )}
                <div className="text-xs text-gray-500 mt-2">
                    SKU: {product.sku} | Stock: {product.stock_quantity}
                </div>
            </CardContent>
            <CardFooter className="p-4 pt-0 mt-auto">
                <Button
                    className="w-full"
                    disabled={isOutOfStock}
                    onClick={() => onAddToCart(product)}
                >
                    <ShoppingCart className="mr-2 h-4 w-4" />
                    Agregar
                </Button>
            </CardFooter>
        </Card>
    );
}
