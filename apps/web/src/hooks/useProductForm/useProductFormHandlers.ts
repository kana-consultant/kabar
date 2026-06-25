import { useNavigate } from "@tanstack/react-router";
import { useToast } from "../use-toast";
import { addProduct, updateProduct, testConnection } from "@/services/product";
import type { Product } from "@/services/product";
import { mapProductToCreateRequest, mapProductToUpdateRequest } from "./productMappers";

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
        const productToSaveRaw = productWithWorkflows || product;

        if (!productToSaveRaw.name || !productToSaveRaw.api_endpoint) {
            toast.error("Isi nama produk dan API endpoint");
            return;
        }

        setLoading(true);

        try {
            let adapterConfig = undefined;

            if (productToSaveRaw.adapter_config) {
                const ac = productToSaveRaw.adapter_config;

                const toJsonString = (val: unknown, fallback = "{}") => {
                    if (!val) return fallback;
                    if (typeof val === "string") return val;
                    return JSON.stringify(val);
                };

                adapterConfig = {
                    id: ac.id || "",
                    custom_headers: toJsonString(ac.custom_headers || (ac as any).customHeaders),
                    meta_config: ac.meta_config ? toJsonString(ac.meta_config) : undefined,
                    sitemap_config: ac.sitemap_config ? toJsonString(ac.sitemap_config) : undefined,
                    timeout_seconds: ac.timeout_seconds || 30,
                };
            }

            // ✅ FIX 2: Clean workflows — strip temp IDs on create, fill required fields
            let workflows = (productToSaveRaw.workflows || []).map(workflow => ({
                ...workflow,
                // Temp ID hanya dipakai di FE, saat create biarkan backend generate
                id: isEdit ? workflow.id : (workflow.id?.startsWith("temp-") ? "" : workflow.id),
                product_id: workflow.product_id || productToSaveRaw.id || "",
                updated_at: new Date().toISOString(),
                nodes: (workflow.nodes || []).map(node => ({
                    ...node,
                    workflow_id:  workflow.id,
                    updated_at: new Date().toISOString(),
                    adapter_config: node.adapter_config
                        ? {
                            ...node.adapter_config,
                            updated_at: new Date().toISOString(),
                            custom_headers: typeof node.adapter_config.custom_headers === "string"
                                ? node.adapter_config.custom_headers
                                : JSON.stringify(node.adapter_config.custom_headers || {}),
                            field_mapping: typeof node.adapter_config.field_mapping === "string"
                                ? node.adapter_config.field_mapping
                                : JSON.stringify(node.adapter_config.field_mapping || {}),
                        }
                        : undefined,
                })),
            }));


            console.log(workflows)

            const now = new Date().toISOString();

            // ✅ FIX 3: created_at fallback ke now jika kosong
            const completeProduct: Partial<Product> = {
                id: productToSaveRaw.id || "",
                name: productToSaveRaw.name,
                platform: productToSaveRaw.platform || "custom",
                api_endpoint: productToSaveRaw.api_endpoint,
                status: productToSaveRaw.status || "pending",
                sync_status: productToSaveRaw.sync_status || "idle",
                created_at: productToSaveRaw.created_at || now,  // ✅ fix empty string
                updated_at: now,

                ...(productToSaveRaw.api_key && { api_key: productToSaveRaw.api_key }),
                ...(productToSaveRaw.last_sync && { last_sync: productToSaveRaw.last_sync }),
                ...(productToSaveRaw.created_by && { created_by: productToSaveRaw.created_by }),
                ...(productToSaveRaw.team_id && { team_id: productToSaveRaw.team_id }),
                ...(productToSaveRaw.user_id && { user_id: productToSaveRaw.user_id }),

                ...(adapterConfig && { adapter_config: adapterConfig }),

                // ✅ FIX 4: Jangan kirim temp workflow_id ke backend
                workflow_id: isEdit
                    ? productToSaveRaw.workflow_id || ""
                    : "",  // backend akan assign setelah create
                workflows,
            };

            console.log(completeProduct)

            if (isEdit && productId) {
                const updateRequest = mapProductToUpdateRequest(completeProduct as Product);
                await updateProduct(productId, updateRequest);
                toast.success("Produk berhasil diperbarui");
            } else {
                const createRequest = mapProductToCreateRequest(completeProduct as Product);
                console.log("Hasil Mapping", createRequest);
                const createdProduct = await addProduct(createRequest);
                console.log("✅ Product created:", createdProduct);
                toast.success("Produk berhasil ditambahkan");
            }

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