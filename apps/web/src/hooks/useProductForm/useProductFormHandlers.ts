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
        // Gunakan productWithWorkflows jika ada (dari WorkflowBuilder),否则 pakai product state
        const productToSaveRaw = productWithWorkflows || product;

        if (!productToSaveRaw.name || !productToSaveRaw.api_endpoint) {
            toast.error("Isi nama produk dan API endpoint");
            return;
        }

        setLoading(true);

        try {
            // Prepare field mapping
            const fieldMappingValue = typeof productToSaveRaw.adapter_config?.field_mapping === "string"
                ? productToSaveRaw.adapter_config.field_mapping
                : JSON.stringify(productToSaveRaw.adapter_config?.field_mapping ?? [], null, 2);

            // Prepare custom headers
            const customHeadersValue = typeof productToSaveRaw.adapter_config?.custom_headers === "string"
                ? productToSaveRaw.adapter_config.custom_headers
                : JSON.stringify(productToSaveRaw.adapter_config?.custom_headers ?? {});

            // Prepare response mapping
            const responseMappingValue = productToSaveRaw.adapter_config?.response_mapping
                ? typeof productToSaveRaw.adapter_config.response_mapping === "string"
                    ? productToSaveRaw.adapter_config.response_mapping
                    : JSON.stringify(productToSaveRaw.adapter_config.response_mapping)
                : undefined;

            // Update product with prepared values
            const productToSave = {
                ...productToSaveRaw,
                adapter_config: productToSaveRaw.adapter_config ? {
                    ...productToSaveRaw.adapter_config,
                    field_mapping: fieldMappingValue,
                    custom_headers: customHeadersValue,
                    response_mapping: responseMappingValue,
                } : undefined,
            };

            console.log("💾 Saving product:", productToSave);

            if (isEdit && productId) {
                // Use mapper for update
                const updateRequest = mapProductToUpdateRequest(productToSave);
                await updateProduct(productId, updateRequest);
                toast.success("Produk diperbarui");
            } else {
                // Use mapper for create
                const createRequest = mapProductToCreateRequest(productToSave);
                await addProduct(createRequest);
                toast.success("Produk ditambahkan");
            }

            navigate({ to: "/products" });
        } catch (error) {
            console.error("Save error:", error);
            toast.error("Gagal menyimpan produk");
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