import { useProductFormState } from "./useProductFormState";
import { useProductFormInit } from "./useProductFormInit";
import { useProductFormActions } from "./useProductFormActions";
import { useProductFormHandlers } from "./useProductFormHandlers";
import type { Product } from "@/services/product/types";

export function useProductForm(isEdit: boolean, productId?: string, initialData?: any) {
    const {
        loading, setLoading,
        testing, setTesting,
        product, setProduct,
    } = useProductFormState();

    const {
        handleTestConnection,
        handleSave,
        handleCancel,
    } = useProductFormHandlers(
        isEdit, productId, product as Product
        , loading, setLoading,
        testing, setTesting, setProduct
    );

    // Load data saat edit
    useProductFormInit(isEdit, initialData, setProduct);

    const {
        // Product info
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        updateMetaConfig,
        updateSitemapConfig,

        // Workflow actions
        updateWorkflowId,
        addWorkflow,
        updateWorkflow,
        deleteWorkflow,
        setActiveWorkflow,
        setWorkflows,
        clearWorkflows,

        // Node actions
        addNodeToWorkflow,
        updateNodeInWorkflow,
        deleteNodeFromWorkflow,
        updateNodeConnections,
        reorderWorkflowNodes,


        setAdapterConfigs,
    } = useProductFormActions(setProduct);



    return {
        // State
        product,
        loading,
        testing,

        // Product info actions
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        updateMetaConfig,
        updateSitemapConfig,

        // Workflow actions
        updateWorkflowId,
        addWorkflow,
        updateWorkflow,
        deleteWorkflow,
        setActiveWorkflow,
        setWorkflows,
        clearWorkflows,

        // Node actions
        addNodeToWorkflow,
        updateNodeInWorkflow,
        deleteNodeFromWorkflow,
        updateNodeConnections,
        reorderWorkflowNodes,

        setAdapterConfigs,

        // Handlers
        handleTestConnection,
        handleSave,
        handleCancel,
    };
}