import { useState } from "react";
import type { Product } from "@/types/product";

export function useProductFormState() {
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

    return {
        loading, setLoading,
        testing, setTesting,
        product, setProduct,
    };
}