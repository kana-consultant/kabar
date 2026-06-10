import { useState, useEffect } from "react";
import { useNavigate, useParams } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { Page, PageHeader, PageTitle, PageDescription, PageContent } from "@/components/page/page";
import { Button } from "@kana-consultant/ui-kit";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";
import { ProviderForm } from "../ProviderForm";
import { DangerZone } from "../DangerZone";
import type { ProviderFormData, Family } from "@/types/provider.types";
import { useProviderManagement } from "@/hooks/useAIManagement/useProviderManagement";

export default function EditProviderPage() {
    const toast = useToast();
    const navigate = useNavigate();
    const { id } = useParams({ from: '/protected/provider/$id/edit' });

    const {
        selectedProvider: provider,
        loading: isLoading,
        loadProviderById,
        updateProvider,
        deleteProvider,
    } = useProviderManagement({ autoLoad: false });

    const [isSubmitting, setIsSubmitting] = useState(false);
    const [activeTab, setActiveTab] = useState("edit");

    useEffect(() => {
        if (id) loadProvider();
    }, [id, loadProviderById]);

    const loadProvider = async () => {
        try {
            await loadProviderById(id);
        } catch {
            navigate({ to: "/ai-management" });
        }
    };

    const handleSubmit = async (data: ProviderFormData) => {
        setIsSubmitting(true);
        try {
            await updateProvider(id, {
                ...data,
                families: (data.families || []).map(family => ({
                    id: family.id,
                    provider_id: family.provider_id,
                    schema_id: family.schema_id,
                    name: family.name,
                    display_name: family.display_name,
                    description: family.description ?? null,
                    max_token: family.max_token ?? 1024,        
                    temperature: family.temperature ?? 1.0,       
                    system_prompt: family.system_prompt ?? "You are a helpful assistant.", 
                    schema: family.schema ? {
                        id: family.schema.id,
                        provider_id: family.schema.provider_id,
                        name: family.schema.name,
                        endpoint_path: family.schema.endpoint_path,
                        request_template: family.schema.request_template ?? null,
                        response_text_path: family.schema.response_text_path,
                        response_image_path: family.schema.response_image_path ?? null,
                        supports_temperature: family.schema.supports_temperature,
                        supports_streaming: family.schema.supports_streaming,
                    } : undefined,
                } satisfies Family))
            });

            toast.success("Provider updated successfully");
            navigate({ to: "/ai-management" });

        } catch (error) {
            toast.error("Failed to update provider");
            console.error(error);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleDelete = async () => {
        try {
            await deleteProvider(id);
            toast.success("Provider deleted successfully");
            navigate({ to: "/ai-management" });
        } catch (error) {
            toast.error("Failed to delete provider");
            console.error(error);
        }
    };

    if (isLoading) {
        return (
            <Page>
                <PageContent>
                    <div className="flex items-center justify-center h-64">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                    </div>
                </PageContent>
            </Page>
        );
    }

    if (!provider) return null;

    return (
        <Page>
            <PageHeader>
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                        <Button variant="ghost" size="sm" onClick={() => navigate({ to: "/ai-management" })}>
                            <ArrowLeft className="h-4 w-4 mr-2" />
                        </Button>
                        <div>
                            <PageTitle>Edit Provider: {provider.display_name}</PageTitle>
                            <PageDescription>Manage provider configuration and settings</PageDescription>
                        </div>
                    </div>
                </div>
            </PageHeader>

            <PageContent>
                <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
                    <TabsList>
                        <TabsTrigger value="edit">Edit Configuration</TabsTrigger>
                        <TabsTrigger value="danger">Danger Zone</TabsTrigger>
                    </TabsList>

                    <TabsContent value="edit" className="space-y-6">
                        <div className="gap-6">
                            <div className="lg:col-span-2">
                                <ProviderForm
                                    mode="edit"
                                    initialData={provider}
                                    onSubmit={handleSubmit}
                                    isSubmitting={isSubmitting}
                                />
                            </div>
                        </div>
                    </TabsContent>

                    <TabsContent value="danger">
                        <DangerZone onDelete={handleDelete} providerName={provider.display_name} />
                    </TabsContent>
                </Tabs>
            </PageContent>
        </Page>
    );
}