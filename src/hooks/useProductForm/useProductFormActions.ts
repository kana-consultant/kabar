import type { Product, AdapterConfig } from "@/types/product";

export function useProductFormActions(
    setProduct: (updater: any) => void
) {
    const updateProductInfo = (updates: Partial<Product>) => {
        setProduct((prev : Partial<Product>) => ({ ...prev, ...updates }));
    };

      const updateAdapterConfig = (updates: Partial<AdapterConfig>) => {
        setProduct((prev: Partial<Product>) => ({
            ...prev,
            adapterConfig: { ...prev.adapterConfig, ...updates } as AdapterConfig,
        }));
    };

    const updateFieldMapping = (value: string) => {
        console.log("📝 useProductForm.updateFieldMapping called with:", value);
        console.log("📝 type:", typeof value);

        setProduct((prev : Partial<Product> )=> {
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

    return {
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
    };
}