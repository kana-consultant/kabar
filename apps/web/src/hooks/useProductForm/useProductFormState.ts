import { useState } from "react";
import type { Product } from "@/services/product";

const createEmptyProduct = (): Product => ({
    id: "",
    name: "",
    platform: "wordpress",
    api_endpoint: "",
    sync_status : "idle",
    created_at : "", 
    updated_at : "",
    api_key: "",
    status: "pending",
    adapter_config: undefined,
    workflow_id: "",
    workflows: [],
});

export function useProductFormState() {
    const [loading, setLoading] = useState(false);
    const [testing, setTesting] = useState(false);
    const [product, setProduct] = useState<Product>(createEmptyProduct());

    return {
        loading,
        setLoading,
        testing,
        setTesting,
        product,
        setProduct,
    };
}