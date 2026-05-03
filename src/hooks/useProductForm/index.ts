import { useProductFormState } from "./useProductFormState";
import { useProductFormInit } from "./useProductFormInit";
import { useProductFormActions } from "./useProductFormActions";
import { useProductFormHandlers } from "./useProductFormHandlers";

export function useProductForm(isEdit: boolean, productId?: string, initialData?: any) {
    const {
        loading, setLoading,
        testing, setTesting,
        product, setProduct,
    } = useProductFormState();

    // Load data saat edit
    useProductFormInit(isEdit, initialData, setProduct);

    const {
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
    } = useProductFormActions(setProduct);

    const {
        handleTestConnection,
        handleSave,
        handleCancel,
    } = useProductFormHandlers(
        isEdit, productId, product, loading, setLoading,
        testing, setTesting, setProduct
    );

    return {
        product,
        loading,
        testing,
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        handleTestConnection,
        handleSave,
        handleCancel,
    };
}