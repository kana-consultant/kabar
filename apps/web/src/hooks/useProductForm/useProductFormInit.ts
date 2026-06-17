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
            });
        }
    }, [isEdit, initialData]);
}