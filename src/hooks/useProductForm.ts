import { useState, useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { addProduct, updateProduct, testConnection } from "@/services/productService";
import type { Product, AdapterConfig } from "@/types/product";

export function useProductForm(isEdit: boolean, productId?: string, initialData?: Product | null) {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [testing, setTesting] = useState(false);
    const [product, setProduct] = useState<Partial<Product>>({
        name: "",
        platform: "wordpress",
        apiEndpoint: "",
        apiKey: "",
        status: "pending",
        lastSync: "-",
        adapterConfig: {
            endpointPath: "",
            httpMethod: "POST",
            customHeaders: {
                "Content-Type": "application/json",
            },
            fieldMapping: JSON.stringify([], null, 2),
        },
    });

    // Load data saat edit
    useEffect(() => {
        if (isEdit && initialData) {
            let fieldMappingValue = initialData.adapterConfig?.fieldMapping || "[]";

            if (typeof fieldMappingValue === 'object') {
                fieldMappingValue = JSON.stringify(fieldMappingValue, null, 2);
            }

            if (typeof fieldMappingValue === 'string' && (fieldMappingValue === "" || fieldMappingValue === "null")) {
                fieldMappingValue = "[]";
            }

            setProduct({
                ...initialData,
                adapterConfig: {
                    ...initialData.adapterConfig,
                    endpointPath: initialData.adapterConfig?.endpointPath || "",
                    httpMethod: initialData.adapterConfig?.httpMethod || "POST",
                    customHeaders: initialData.adapterConfig?.customHeaders || { "Content-Type": "application/json" },
                    fieldMapping: fieldMappingValue,
                },
            });
        }
    }, [isEdit, initialData]);

    const updateProductInfo = (updates: Partial<Product>) => {
        setProduct(prev => ({ ...prev, ...updates }));
    };

    const updateAdapterConfig = (updates: Partial<AdapterConfig>) => {
        setProduct(prev => ({
            ...prev,
            adapterConfig: { ...prev.adapterConfig, ...updates } as AdapterConfig,
        }));
    };

    const updateFieldMapping = (value: string) => {
        console.log("📝 useProductForm.updateFieldMapping called with:", value);
        console.log("📝 type:", typeof value);

        setProduct(prev => {
            const newProduct = {
                ...prev,
                adapterConfig: {
                    ...prev.adapterConfig,
                    fieldMapping: value,
                } as AdapterConfig,
            };
            console.log("📝 New product state:", newProduct);
            return newProduct;
        });
    };

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
                setProduct(prev => ({ ...prev, status: "connected", lastSync: "baru saja" }));
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
        product,
        loading,
        testing,
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        handleTestConnection,
        handleSave,
        handleCancel,
    };
}