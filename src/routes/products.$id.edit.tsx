import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { ProductForm } from "@/components/products/ProductForm/ProductForm";
import { getProductById } from "@/services/productService";
import { toast } from "sonner";

export const Route = createFileRoute("/products/$id/edit")({
    component: ProductEditWrapper,
});

export default function ProductEditWrapper() {
    const { id } = Route.useParams();
    const [product, setProduct] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadProduct = async () => {
            try {
                const existing = await getProductById(id);
                if (existing) {
                    setProduct(existing as any);
                } else {
                    toast.error("Produk tidak ditemukan");
                }
            } catch (error) {
                toast.error("Gagal memuat produk");
            } finally {
                setLoading(false);
            }
        };
        loadProduct();
    }, [id]);

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="text-slate-500">Memuat...</div>
            </div>
        );
    }

    if (!product) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="text-slate-500">Produk tidak ditemukan</div>
            </div>
        );
    }

    return <ProductForm isEdit={true} productId={id} initialData={product} />;
}