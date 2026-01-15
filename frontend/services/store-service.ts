import api from "@/lib/axios";

export interface Product {
    id: string;
    club_id: string;
    name: string;
    description: string;
    price: string;
    stock_quantity: number;
    sku: string;
    category: string;
    image_url?: string;
    status: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface CartItem {
    product_id: string;
    quantity: number;
}

export interface PurchaseRequest {
    items: CartItem[];
    payment_method: string; // 'CASH', 'MERCADO_PAGO', etc.
}

export const storeService = {
    getProducts: async (category?: string): Promise<Product[]> => {
        const params = category ? { category } : {};
        const response = await api.get("/store/products", { params });
        // Handler returns { data: products }
        return response.data.data;
    },

    purchaseItems: async (items: CartItem[], paymentMethod: string = 'CASH') => {
        const payload: PurchaseRequest = {
            items,
            payment_method: paymentMethod
        };
        const response = await api.post("/store/purchase", payload);
        return response.data;
    }
};
