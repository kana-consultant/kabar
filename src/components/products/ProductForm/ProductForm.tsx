import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ProductBasicInfo } from "./ProductBasicInfo";
import { ProductApiConfig } from "./ProductApiConfig";
import { ProductFormActions } from "./ProductFormActions";
import { ProductFieldMapping } from "@/components/products/ProductFieldMapping";
import { useProductForm } from "@/hooks/useProductForm";
import type { Product } from "@/types/product";

interface ProductFormProps {
    isEdit: boolean;
    productId?: string;
    initialData?: Product | null;
}

export function ProductForm({ isEdit, productId, initialData }: ProductFormProps) {
    const {
        product,
        loading,
        testing,
        updateProductInfo,
        updateAdapterConfig,
        updateFieldMapping,
        handleTestConnection,
        handleSave,
        handleCancel,
    } = useProductForm(isEdit, productId, initialData);

    return (
        <div className="max-w-7xl mx-auto space-y-6 ">
            {/* Header */}
            <div className="flex items-center gap-4 border-b pb-4">
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={handleCancel}
                    className="shrink-0"
                >
                    <ArrowLeft className="h-5 w-5" />
                </Button>
                <div>
                    <h2 className="text-2xl md:text-3xl font-bold tracking-tight">
                        {isEdit ? "Edit Produk" : "Tambah Produk Baru"}
                    </h2>
                    <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
                        Konfigurasi API endpoint dan field mapping untuk integrasi konten
                    </p>
                </div>
            </div>

            {/* Form Grid */}
            <div className="grid gap-6 lg:grid-cols-2">
                <ProductBasicInfo
                    product={product}
                    onUpdate={updateProductInfo}
                    onTestConnection={handleTestConnection}
                    isTesting={testing}
                />
                <ProductApiConfig
                    config={product.adapterConfig || {}}
                    onUpdate={updateAdapterConfig}
                />
            </div>

            {/* Field Mapping - Full Width */}
            <div className="mt-2">
                <ProductFieldMapping
                    fieldMapping={(() => {
                        try {
                            // Coba parse ke object
                            const parsed = JSON.parse(product.adapterConfig?.fieldMapping || "{}");
                            return parsed;
                        } catch (e) {
                            return {};
                        }
                    })()}
                    onChange={updateFieldMapping}
                />
            </div>

            {/* Action Buttons */}
            <ProductFormActions
                onCancel={handleCancel}
                onSave={handleSave}
                isSaving={loading}
            />
        </div>
    );
}