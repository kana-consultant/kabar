import { useState } from "react";
import { ArrowLeft } from "lucide-react";
import { Button } from "@kana-consultant/ui-kit";
import { ProductBasicInfo } from "./ProductBasicInfo";
import { ProductApiConfig } from "./ProductApiConfig";
import { ProductFormActions } from "./ProductFormActions";
import { ProductFieldMapping } from "@/pages/products/ProductFieldMapping";
import { useProductForm } from "@/hooks/useProductForm";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    Input,
    Label,
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";
import { WorkflowBuilder } from "../WorkflowBuilder/WorkflowBuilder";
import type { AdapterConfig } from "@/types/product";
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
        updateMetaConfig,
        updateSitemapConfig,
        updateWorkflowId,
        handleSave,
        handleCancel,
    } = useProductForm(isEdit, productId, initialData);

    const toast = useToast();

    // Modal test state
    const [showModal, setShowModal] = useState(false);
    const [testUrl, setTestUrl] = useState("");
    const [testMethod, setTestMethod] = useState("GET");
    const [testAuthType, setTestAuthType] = useState<"none" | "apiKey" | "bearer">("none");
    const [testApiKey, setTestApiKey] = useState("");
    const [testBearerToken, setTestBearerToken] = useState("");
    const [testCustomHeaders, setTestCustomHeaders] = useState<{ key: string; value: string }[]>([
        { key: "", value: "" }
    ]);
    const [isLoading, setIsLoading] = useState(false);

    const addTestHeader = () => {
        setTestCustomHeaders([...testCustomHeaders, { key: "", value: "" }]);
    };

    const updateTestHeader = (index: number, field: "key" | "value", value: string) => {
        const updated = [...testCustomHeaders];
        updated[index][field] = value;
        setTestCustomHeaders(updated);
    };

    const removeTestHeader = (index: number) => {
        const updated = testCustomHeaders.filter((_, i) => i !== index);
        setTestCustomHeaders(updated.length ? updated : [{ key: "", value: "" }]);
    };

    const handleTest = async () => {
        if (!testUrl) {
            toast.error("URL tidak boleh kosong");
            return;
        }

        if (testAuthType === "apiKey" && !testApiKey) {
            toast.error("Mohon masukkan API Key");
            return;
        }

        if (testAuthType === "bearer" && !testBearerToken) {
            toast.error("Mohon masukkan Bearer Token");
            return;
        }

        setIsLoading(true);

        try {
            const headers: Record<string, string> = {};

            testCustomHeaders.forEach(header => {
                if (header.key.trim()) {
                    headers[header.key.trim()] = header.value;
                }
            });

            if (testAuthType === "apiKey") {
                headers["X-API-Key"] = testApiKey;
            } else if (testAuthType === "bearer") {
                headers["Authorization"] = `Bearer ${testBearerToken}`;
            }

            const response = await fetch(testUrl, {
                method: testMethod,
                headers,
            });

            let data;
            const contentType = response.headers.get("content-type");
            if (contentType && contentType.includes("application/json")) {
                data = await response.json();
            } else {
                data = await response.text();
            }

            if (response.ok) {
                toast.success(`Status: ${response.status} ${response.statusText}`);
                console.log("Response data:", data);
            } else {
                toast.error(`Status: ${response.status} ${response.statusText}`);
            }
        } catch (error: any) {
            toast.error(error.message);
        } finally {
            setIsLoading(false);
        }
    };

    const adapterConfigs: AdapterConfig[] = product.adapterConfigs
        ? product.adapterConfigs.map(config => ({
            id: config.id,
            productId: config.productId,
            endpointPath: config.endpointPath,
            httpMethod: config.httpMethod,
            customHeaders: config.customHeaders,
            fieldMapping: config.fieldMapping,
            responseMapping: config.responseMapping,
            metaConfig: config.metaConfig,
            sitemapConfig: config.sitemapConfig,
            timeoutSeconds: config.timeoutSeconds,
            retryCount: config.retryCount,
            createdAt: config.createdAt,
            updatedAt: config.updatedAt,
        }))
        : [];

    console.log("HASIL AKHIR")
    console.log(product)

    return (
        <div className="max-w-full mx-auto space-y-6">
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

            {/* Single Endpoint Section */}
            <div className="space-y-6">
                <div className="grid gap-6 lg:grid-cols-2">
                    <ProductBasicInfo
                        product={product}
                        onUpdate={updateProductInfo}
                        onTestConnection={() => setShowModal(true)}
                        isTesting={testing}
                    />
                    <ProductApiConfig
                        config={product.adapterConfig || {}}
                        onUpdate={updateAdapterConfig}
                    />
                </div>

                <div>
                    <ProductFieldMapping
                        domain={product?.apiEndpoint as string || ""}
                        fieldMapping={(() => {
                            try {
                                const parsed = JSON.parse(product.adapterConfig?.fieldMapping || "{}");
                                return parsed;
                            } catch (e) {
                                return {};
                            }
                        })()}
                        metaConfig={(() => {
                            try {
                                return JSON.parse(product.adapterConfig?.metaConfig || "{}");
                            } catch (e) {
                                return {};
                            }
                        })()}
                        sitemapConfig={(() => {
                            try {
                                return JSON.parse(product.adapterConfig?.sitemapConfig || "{}");
                            } catch (e) {
                                return {};
                            }
                        })()}
                        onChange={(fieldMapping, metaConfig, sitemapConfig) => {
                            updateFieldMapping(fieldMapping);
                            if (metaConfig) updateMetaConfig(metaConfig);
                            if (sitemapConfig) updateSitemapConfig(sitemapConfig);
                        }}
                    />
                </div>
            </div>

            {/* Workflow Builder Section */}
            <div className="border-t pt-6">
                <div className="mb-4">
                    <h3 className="text-lg font-semibold">Workflow Builder</h3>
                    <p className="text-sm text-slate-500 dark:text-slate-400">
                        Buat workflow multi-step dengan menghubungkan beberapa adapter
                    </p>
                </div>

                <WorkflowBuilder
                    productId={product.id || ""}
                    adapterConfigs={adapterConfigs}
                    product={product}  // ← tambahkan ini
                    onUpdateAdapterConfigs={(newAdapterConfigs) => {
                        updateProductInfo({ adapterConfigs: newAdapterConfigs });
                    }}
                    onChange={(workflowId) => {
                        updateWorkflowId(workflowId);
                    }}
                />
            </div>
            {/* Save Button */}
            <ProductFormActions
                onCancel={handleCancel}
                onSave={handleSave}
                isSaving={loading}
            />

            {/* Modal Test */}
            <Dialog open={showModal} onOpenChange={setShowModal}>
                <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
                    <DialogHeader>
                        <DialogTitle>Test API Connection</DialogTitle>
                        <DialogDescription>
                            Masukkan semua konfigurasi API untuk testing
                        </DialogDescription>
                    </DialogHeader>

                    <div className="space-y-4 py-4">
                        <div className="space-y-2">
                            <Label>URL Endpoint</Label>
                            <Input
                                placeholder="https://api.example.com/v1/products"
                                value={testUrl}
                                onChange={(e: any) => setTestUrl(e.target.value)}
                            />
                        </div>

                        <div className="space-y-2">
                            <Label>HTTP Method</Label>
                            <Select value={testMethod} onValueChange={setTestMethod}>
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="GET">GET</SelectItem>
                                    <SelectItem value="POST">POST</SelectItem>
                                    <SelectItem value="PUT">PUT</SelectItem>
                                    <SelectItem value="PATCH">PATCH</SelectItem>
                                    <SelectItem value="DELETE">DELETE</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-2">
                            <Label>Autentikasi</Label>
                            <Select
                                value={testAuthType}
                                onValueChange={(v: any) => setTestAuthType(v as any)}
                            >
                                <SelectTrigger>
                                    <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="none">Tanpa Autentikasi</SelectItem>
                                    <SelectItem value="apiKey">API Key</SelectItem>
                                    <SelectItem value="bearer">Bearer Token</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        {testAuthType === "apiKey" && (
                            <div className="space-y-2">
                                <Label>API Key</Label>
                                <Input
                                    type="password"
                                    placeholder="Masukkan API Key"
                                    value={testApiKey}
                                    onChange={(e: any) => setTestApiKey(e.target.value)}
                                />
                            </div>
                        )}

                        {testAuthType === "bearer" && (
                            <div className="space-y-2">
                                <Label>Bearer Token</Label>
                                <Input
                                    type="password"
                                    placeholder="Masukkan Bearer Token"
                                    value={testBearerToken}
                                    onChange={(e: any) => setTestBearerToken(e.target.value)}
                                />
                            </div>
                        )}

                        <div className="space-y-2">
                            <div className="flex items-center justify-between">
                                <Label>Custom Headers</Label>
                                <Button type="button" variant="outline" size="sm" onClick={addTestHeader}>
                                    Tambah Header
                                </Button>
                            </div>
                            <div className="space-y-2">
                                {testCustomHeaders.map((header, index) => (
                                    <div key={index} className="flex gap-2">
                                        <Input
                                            placeholder="Header Key"
                                            value={header.key}
                                            onChange={(e: any) => updateTestHeader(index, "key", e.target.value)}
                                            className="flex-1"
                                        />
                                        <Input
                                            placeholder="Header Value"
                                            value={header.value}
                                            onChange={(e: any) => updateTestHeader(index, "value", e.target.value)}
                                            className="flex-1"
                                        />
                                        <Button
                                            type="button"
                                            variant="destructive"
                                            size="sm"
                                            onClick={() => removeTestHeader(index)}
                                        >
                                            Hapus
                                        </Button>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>

                    <div className="flex gap-3">
                        <Button variant="outline" onClick={() => setShowModal(false)} className="flex-1">
                            Tutup
                        </Button>
                        <Button onClick={handleTest} disabled={isLoading} className="flex-1">
                            {isLoading ? "Menguji..." : "Test"}
                        </Button>
                    </div>
                </DialogContent>
            </Dialog>
        </div>
    );
}