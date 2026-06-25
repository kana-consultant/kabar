import { useEffect } from "react";
import type { Product } from "@/services/product";

export function useProductFormInit(
    isEdit: boolean,
    initialData: Product | null | undefined,
    setProduct: (updater: any) => void
) {
    useEffect(() => {
        if (isEdit && initialData) {

            setProduct({
                ...initialData,
            });
        }
    }, [isEdit, initialData]);
}