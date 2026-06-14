import { useNavigate } from "@tanstack/react-router";
import { useToast } from "../use-toast";
import { addProduct, updateProduct, testConnection } from "@/services/product";
import type { Product } from "@/services/product";
import  { mapProductToCreateRequest, mapProductToUpdateRequest } from "./productMappers";

export function useProductFormHandlers(
    isEdit: boolean,
    productId: string | undefined,
    product: Partial<Product>,
    loading: boolean,
    setLoading: (val: boolean) => void,
    testing: boolean,
    setTesting: (val: boolean) => void,
    setProduct: (updater: any) => void
) {
    const navigate = useNavigate();
    const toast = useToast()

    const handleTestConnection = async () => {
        if (!product.api_endpoint) {
            toast.error("Isi API endpoint terlebih dahulu");
            return;
        }

        setTesting(true);
        try {
            const result = await testConnection(product as any);
            if (result) {
                toast.success("Koneksi berhasil");
                setProduct((prev: Partial<Product>) => ({
                    ...prev,
                    status: "connected",
                    last_sync: new Date().toISOString()
                }));
            } else {
                toast.error("Koneksi gagal");
            }
        } catch (error) {
            toast.error("Gagal menguji koneksi");
        } finally {
            setTesting(false);
        }
    };

    const handleSave = async (productWithWorkflows?: Product) => {
        // Gunakan productWithWorkflows jika ada (dari WorkflowBuilder), else pakai product state
        const productToSaveRaw = productWithWorkflows || product;

        // Validasi required fields
        if (!productToSaveRaw.name || !productToSaveRaw.api_endpoint) {
            toast.error("Isi nama produk dan API endpoint");
            return;
        }

        setLoading(true);

        try {
            // Prepare adapter config dengan semua field yang diperlukan
            let adapterConfig = undefined;

            if (productToSaveRaw.adapter_config) {
                // Handle field_mapping
                let fieldMappingValue = productToSaveRaw.adapter_config.field_mapping;
                if (fieldMappingValue && typeof fieldMappingValue !== 'string') {
                    fieldMappingValue = JSON.stringify(fieldMappingValue);
                } else if (!fieldMappingValue) {
                    fieldMappingValue = "{}";
                }

                // Handle custom_headers
                let customHeadersValue = productToSaveRaw.adapter_config.custom_headers;
                if (customHeadersValue && typeof customHeadersValue !== 'string') {
                    customHeadersValue = JSON.stringify(customHeadersValue);
                } else if (!customHeadersValue) {
                    customHeadersValue = "{}";
                }

                // Handle response_mapping
                let responseMappingValue = productToSaveRaw.adapter_config.response_mapping;
                if (responseMappingValue && typeof responseMappingValue !== 'string') {
                    responseMappingValue = JSON.stringify(responseMappingValue);
                }

                // Handle meta_config
                let metaConfigValue = productToSaveRaw.adapter_config.meta_config;
                if (metaConfigValue && typeof metaConfigValue !== 'string') {
                    metaConfigValue = JSON.stringify(metaConfigValue);
                }

                // Handle sitemap_config
                let sitemapConfigValue = productToSaveRaw.adapter_config.sitemap_config;
                if (sitemapConfigValue && typeof sitemapConfigValue !== 'string') {
                    sitemapConfigValue = JSON.stringify(sitemapConfigValue);
                }

                adapterConfig = {
                    id: productToSaveRaw.adapter_config.id,
                    product_id: productToSaveRaw.adapter_config.product_id || productToSaveRaw.id || "",
                    endpoint_path: productToSaveRaw.adapter_config.endpoint_path || "",
                    http_method: productToSaveRaw.adapter_config.http_method || "GET",
                    custom_headers: customHeadersValue,
                    field_mapping: fieldMappingValue,
                    response_mapping: responseMappingValue,
                    meta_config: metaConfigValue,
                    sitemap_config: sitemapConfigValue,
                    timeout_seconds: productToSaveRaw.adapter_config.timeout_seconds || 30,
                    retry_count: productToSaveRaw.adapter_config.retry_count || 3,
                    created_at: productToSaveRaw.adapter_config.created_at || new Date().toISOString(),
                    updated_at: new Date().toISOString(),
                };
            }

            // Prepare workflows data
            let workflows = productToSaveRaw.workflows || [];

            // Clean up workflows data (remove temporary ids if needed, ensure all required fields)
            workflows = workflows.map(workflow => ({
                ...workflow,
                product_id: workflow.product_id || productToSaveRaw.id || "",
                updated_at: new Date().toISOString(),
                nodes: workflow.nodes?.map(node => ({
                    ...node,
                    workflow_id: node.workflow_id || workflow.id,
                    updated_at: new Date().toISOString(),
                    // Ensure adapter_config in node is properly formatted
                    adapter_config: node.adapter_config ? {
                        ...node.adapter_config,
                        updated_at: new Date().toISOString(),
                        custom_headers: typeof node.adapter_config.custom_headers === 'string'
                            ? node.adapter_config.custom_headers
                            : JSON.stringify(node.adapter_config.custom_headers || {}),
                        field_mapping: typeof node.adapter_config.field_mapping === 'string'
                            ? node.adapter_config.field_mapping
                            : JSON.stringify(node.adapter_config.field_mapping || {}),
                    } : undefined,
                })) || [],
            }));

            // Build complete product object
            const completeProduct: Partial<Product> = {
                // Required fields
                id: productToSaveRaw.id,
                name: productToSaveRaw.name,
                platform: productToSaveRaw.platform || "custom",
                api_endpoint: productToSaveRaw.api_endpoint,
                status: productToSaveRaw.status || "pending",
                sync_status: productToSaveRaw.sync_status || "idle",
                created_at: productToSaveRaw.created_at || new Date().toISOString(),
                updated_at: new Date().toISOString(),

                // Optional fields
                ...(productToSaveRaw.api_key && { api_key: productToSaveRaw.api_key }),
                ...(productToSaveRaw.last_sync && { last_sync: productToSaveRaw.last_sync }),
                ...(productToSaveRaw.created_by && { created_by: productToSaveRaw.created_by }),
                ...(productToSaveRaw.team_id && { team_id: productToSaveRaw.team_id }),
                ...(productToSaveRaw.user_id && { user_id: productToSaveRaw.user_id }),

                // Adapter config
                ...(adapterConfig && { adapter_config: adapterConfig }),

                // Workflow related
                workflow_id: productToSaveRaw.workflow_id || (workflows[0]?.id || ""),
                workflows: workflows,
            };

            console.log("💾 Saving complete product:", {
                ...completeProduct,
                workflows_count: workflows.length,
                nodes_count: workflows.reduce((acc, w) => acc + (w.nodes?.length || 0), 0)
            });

            if (isEdit && productId) {
                // Update existing product
                const updateRequest = mapProductToUpdateRequest(completeProduct as Product);
                await updateProduct(productId, updateRequest);
                toast.success("Produk berhasil diperbarui");
            } else {
                // Create new product
                const createRequest = mapProductToCreateRequest(completeProduct as Product);
                console.log("Hasil Mapping",createRequest)
                const createdProduct = await addProduct(createRequest);
                console.log("✅ Product created:", createdProduct);
                toast.success("Produk berhasil ditambahkan");
            }

            // Navigate back to products list
            navigate({ to: "/products" });
        } catch (error) {
            console.error("Save error:", error);
            toast.error("Gagal menyimpan produk: " + (error as Error).message);
        } finally {
            setLoading(false);
        }
    };

    const handleCancel = () => {
        navigate({ to: "/products" });
    };

    return {
        handleTestConnection,
        handleSave,
        handleCancel,
    };
}