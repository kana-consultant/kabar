import { useCallback } from "react";
import { getProducts } from "@/services/product";
import type { ToastContextType } from "../use-toast";

export function useProductsData(
    setProducts: (data: any[]) => void,
    setProductNames: (names: string[]) => void,
    setProductsLoading: (val: boolean) => void,
    setProductsError: (val: string | null) => void,
    setLoading: (val : boolean) => void,
    toast : ToastContextType
) {
    const fetchProducts = useCallback(async () => {
        setProductsLoading(true);
        setProductsError(null);
        try {
            const productsData = await getProducts();
            console.log(productsData)
            setProducts(productsData || []);

            // Extract product names for display
            const names = productsData ? productsData.map(p => p.name) : [];
            setProductNames(names);

            console.log("Products loaded:", productsData);
        } catch (error) {
            console.error("Failed to fetch products:", error);
            setProductsError("Gagal memuat data produk");
            toast.error("Gagal memuat data produk", {
                description: "Silahkan refresh halaman",
            });
            // Fallback to default products if API fails
            const defaultProducts = ["TrekkingID", "CampingMart", "OutdoorGear"];
            setProductNames(defaultProducts);
        } finally {
            setProductsLoading(false);
        }
    }, []);

    const loadProducts = useCallback(async () => {
        setProductsLoading(true);
        try {
            const data = await getProducts();
            setProducts(data || []);
            setProductNames(data ? data.map(p => p.name) : []);
            setLoading(false)
        } catch (error) {
            console.error('Failed to load products:', error);
            toast.error('Gagal memuat produk');
        } finally {
            setProductsLoading(false);
        }
    }, []);

    return { fetchProducts, loadProducts };
}