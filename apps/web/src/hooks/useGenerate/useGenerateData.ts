import { useEffect, useCallback } from "react";
// Hapus import Toast langsung
// import { Toast } from "@kana-consultant/ui-kit";
import type { ToastContextType } from "@/hooks/use-toast"; // Import type untuk typing
import { getProducts } from "@/services/product";
import { getModelsFromAPIKeys } from "@/services/model";

export function useGenerateData(
    setProducts: (data: any[]) => void,
    setProductNames: (names: string[]) => void,
    setProductsLoading: (loading: boolean) => void,
    setProductsError: (error: string | null) => void,
    setModels: (models: any[]) => void,
    setLoadingModels: (loading: boolean) => void,
    setSelectedModelId: (id: string) => void,
    toast: ToastContextType //   Tambahkan parameter toast
) {
    const fetchProducts = useCallback(async () => {
        setProductsLoading(true);
        setProductsError(null);
        try {
            const productsData = await getProducts();
            console.log(productsData);
            setProducts(productsData || []);
            setProductNames(productsData ? productsData.map(p => p.name) : []);
        } catch (error) {
            console.error("Failed to fetch products:", error);
            setProductsError("Gagal memuat data produk");
            toast.error("Gagal memuat data produk"); //   Ganti toast.error dengan toast.error
            setProductNames(["TrekkingID", "CampingMart", "OutdoorGear"]);
        } finally {
            setProductsLoading(false);
        }
    }, [setProducts, setProductNames, setProductsLoading, setProductsError, toast]);

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
            toast.error("Gagal memuat model AI"); //   Ganti toast.error dengan toast.error
        } finally {
            setLoadingModels(false);
        }
    }, [setModels, setLoadingModels, setSelectedModelId, toast]);

    useEffect(() => {
        fetchProducts();
        loadModels();
    }, [fetchProducts, loadModels]); //   Tambahkan dependencies

    return { fetchProducts, loadModels };
}