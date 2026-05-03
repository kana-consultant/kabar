import { useEffect, useCallback } from "react";
import { toast } from "sonner";
import { getProducts } from "@/services/product";
import { getModelsFromAPIKeys } from "@/services/model";

export function useGenerateData(
    setProducts: (data: any[]) => void,
    setProductNames: (names: string[]) => void,
    setProductsLoading: (loading: boolean) => void,
    setProductsError: (error: string | null) => void,
    setModels: (models: any[]) => void,
    setLoadingModels: (loading: boolean) => void,
    setSelectedModelId: (id: string) => void
) {
    const fetchProducts = useCallback(async () => {
        setProductsLoading(true);
        setProductsError(null);
        try {
            const productsData = await getProducts();
            setProducts(productsData || []);
            setProductNames(productsData ? productsData.map(p => p.name) : []);
        } catch (error) {
            console.error("Failed to fetch products:", error);
            setProductsError("Gagal memuat data produk");
            setProductNames(["TrekkingID", "CampingMart", "OutdoorGear"]);
        } finally {
            setProductsLoading(false);
        }
    }, []);

    const loadModels = useCallback(async () => {
        setLoadingModels(true);
        try {
            const modelsData = await getModelsFromAPIKeys();
            setModels(modelsData as any);
            if (modelsData.length > 0 && !setSelectedModelId) {
                // selectedModelId logic will be handled in useEffect
            }
        } catch (error) {
            console.error("Failed to load models:", error);
            toast.error("Gagal memuat model AI");
        } finally {
            setLoadingModels(false);
        }
    }, []);

    useEffect(() => {
        fetchProducts();
        loadModels();
    }, []);

    return { fetchProducts, loadModels };
}