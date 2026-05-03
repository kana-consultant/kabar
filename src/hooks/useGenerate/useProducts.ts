import { useState, useEffect, useCallback } from "react";
import { getProducts } from "@/services/product";
import type { Product } from "@/types/product";

export function useProducts() {
    const [products, setProducts] = useState<Product[]>([]);
    const [productNames, setProductNames] = useState<string[]>([]);
    const [productsLoading, setProductsLoading] = useState(true);
    const [productsError, setProductsError] = useState<string | null>(null);

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

    useEffect(() => {
        fetchProducts();
    }, [fetchProducts]);

    return {
        products,
        productNames,
        productsLoading,
        productsError,
        refetchProducts: fetchProducts
    };
}