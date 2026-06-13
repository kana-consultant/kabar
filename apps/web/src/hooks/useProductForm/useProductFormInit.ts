import { useEffect } from "react";
import type { Product } from "@/services/product";

export function useProductFormInit(
    isEdit: boolean,
    initialData: Product | null | undefined,
    setProduct: (updater: any) => void
) {
    useEffect(() => {
        if (isEdit && initialData) {
            let fieldMappingValue = initialData.adapter_config?.field_mapping || "[]";

            if (typeof fieldMappingValue === 'object') {
                fieldMappingValue = JSON.stringify(fieldMappingValue, null, 2);
            }

            if (typeof fieldMappingValue === 'string' && (fieldMappingValue === "" || fieldMappingValue === "null")) {
                fieldMappingValue = "[]";
            }

            setProduct({
                ...initialData,
                adapterConfig: {
                    ...initialData.adapter_config,
                    endpointPath: initialData.adapter_config?.endpoint_path || "",
                    httpMethod: initialData.adapter_config?.http_method || "POST",
                    customHeaders: initialData.adapter_config?.custom_headers || { "Content-Type": "application/json" },
                    fieldMapping: fieldMappingValue,
                },
            });
        }
    }, [isEdit, initialData]);
}