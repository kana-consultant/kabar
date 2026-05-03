
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { addProduct, updateProduct, testConnection } from "@/services/product";
import type { Product } from "@/types/product";

export function useProductFormHandlers(
    isEdit: boolean,
    productId: string | undefined,
    product: Partial<Product>,
    loading: boolean,
    setLoading: (val: boolean) => void,
    testing: boolean,
    setTesting: (val: boolean) => void,
    setProduct: (updater: any) => void
) {
    const navigate = useNavigate();

    const handleTestConnection = async () => {
        if (!product.apiEndpoint) {
            toast.error("Isi API endpoint terlebih dahulu");
            return;
        }

        setTesting(true);
        try {
            const result = await testConnection(product as any);
            if (result) {
                toast.success("Koneksi berhasil");
                setProduct((prev: Partial<Product>) => ({ ...prev, status: "connected", lastSync: "baru saja" }));
            } else {
                toast.error("Koneksi gagal");
            }
        } catch (error) {
            toast.error("Gagal menguji koneksi");
        } finally {
            setTesting(false);
        }
    };

    const handleSave = async () => {
        if (!product.name || !product.apiEndpoint) {
            toast.error("Isi nama produk dan API endpoint");
            return;
        }

        setLoading(true);
        try {
            // Pastikan fieldMapping adalah string
            let fieldMappingValue = product.adapterConfig?.fieldMapping;
            if (typeof fieldMappingValue !== 'string') {
                fieldMappingValue = JSON.stringify(fieldMappingValue || [], null, 2);
            }

            const productToSave = {
                ...product,
                adapterConfig: {
                    ...product.adapterConfig,
                    endpoint: product.adapterConfig?.endpointPath || "",
                    method: product.adapterConfig?.httpMethod || "POST",
                    headers: product.adapterConfig?.customHeaders || { "Content-Type": "application/json" },
                    fieldMapping: fieldMappingValue,
                },
            };

            console.log("💾 Saving product:", productToSave);
            console.log("💾 fieldMapping type:", typeof fieldMappingValue);

            if (isEdit && productId) {
                await updateProduct(productId, productToSave);
                toast.success("Produk diperbarui");
            } else {
                await addProduct(productToSave as any);
                toast.success("Produk ditambahkan");
            }
            navigate({ to: "/products" });
        } catch (error) {
            console.error("Save error:", error);
            toast.error("Gagal menyimpan produk");
        } finally {
            setLoading(false);
        }
    };

    const handleCancel = () => {
        navigate({ to: "/products" });
    };

    return {
        handleTestConnection,
        handleSave,
        handleCancel,
    };
}