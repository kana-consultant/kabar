import { useState } from "react";
import type { Product, AdapterConfig } from "@/services/product";

const createEmptyAdapterConfig = (): Omit<AdapterConfig, 'id' | 'productId'> => ({
    endpointPath: "",
    httpMethod: "POST",
    customHeaders: { "Content-Type": "application/json" },
    fieldMapping: "{}",
    responseMapping: {},
    metaConfig: "{}",
    sitemapConfig: "{}",
    timeoutSeconds: 30,
    retryCount: 3,
});

const createEmptyProduct = (): Partial<Product> => ({
    name: "",
    platform: "wordpress",
    apiEndpoint: "",
    apiKey: "",
    status: "pending",
    syncStatus: "idle",
    lastSync: undefined,
    adapterConfigs: [createEmptyAdapterConfig() as AdapterConfig],
});

export function useProductFormState() {
    const [loading, setLoading] = useState(false);
    const [testing, setTesting] = useState(false);
    const [product, setProduct] = useState<Partial<Product>>(createEmptyProduct());

    return {
        loading,
        setLoading,
        testing,
        setTesting,
        product,
        setProduct,
    };
}