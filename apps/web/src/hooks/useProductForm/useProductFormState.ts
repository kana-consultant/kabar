import { useState } from "react";
import type { Product, AdapterConfig } from "@/services/product";

const createEmptyAdapterConfig = (): AdapterConfig => ({
    custom_headers: JSON.stringify({
        "Content-Type": "application/json",
    }),
});

const createEmptyProduct = (): Product => ({
    id: "",
    name: "",
    platform: "wordpress",
    api_endpoint: "",
    sync_status: "idle",
    created_at: "",
    updated_at: "",
    api_key: "",
    status: "pending",
    adapter_config: createEmptyAdapterConfig(),
    workflow_id: "",
    workflows: [],
});

export function useProductFormState() {
    const [loading, setLoading] = useState(false);
    const [testing, setTesting] = useState(false);
    const [product, setProduct] = useState<Product | null>(null);

    return {
        loading,
        setLoading,
        testing,
        setTesting,
        product,
        setProduct,
    };
}