import { useEffect } from "react";
import type { Product } from "@/types/product";

export function useProductFormInit(
    isEdit: boolean,
    initialData: Product | null | undefined,
    setProduct: (updater: any) => void
) {
    useEffect(() => {
        if (isEdit && initialData) {
            let fieldMappingValue = initialData.adapterConfig?.fieldMapping || "[]";

            if (typeof fieldMappingValue === 'object') {
                fieldMappingValue = JSON.stringify(fieldMappingValue, null, 2);
            }

            if (typeof fieldMappingValue === 'string' && (fieldMappingValue === "" || fieldMappingValue === "null")) {
                fieldMappingValue = "[]";
            }

            setProduct({
                ...initialData,
                adapterConfig: {
                    ...initialData.adapterConfig,
                    endpointPath: initialData.adapterConfig?.endpointPath || "",
                    httpMethod: initialData.adapterConfig?.httpMethod || "POST",
                    customHeaders: initialData.adapterConfig?.customHeaders || { "Content-Type": "application/json" },
                    fieldMapping: fieldMappingValue,
                },
            });
        }
    }, [isEdit, initialData]);
}