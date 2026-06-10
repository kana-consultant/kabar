import { useModelForm } from "@/hooks/useModelManagement/useModelForm";
import { TestModelDialog } from "../../TestModelDialog";
import { useState } from "react";
import { Button } from "@kana-consultant/ui-kit";
import { useNavigate, useParams } from "@tanstack/react-router";
import { useEditModel } from "@/hooks/useModelManagement/useEditModel";
import { ModelFormLayout } from "../ModelForm/ModelFormLayout";
import type { AIModelSchema, APIProvider } from "@/types/provider.types";
import { PageTitle } from "@/components/page/page";
import { ArrowLeft } from "lucide-react";

export default function EditModelPage() {
    const navigate = useNavigate();
    const { id } = useParams({ from: '/protected/model/$id/edit' });
    const [showTestModel, setShowTestModel] = useState(false);

    // Use hook for data fetching and mutations
    const { providers, families, schemas, model, isLoading, isSubmitting, updateModel } = useEditModel(id as string);

    const { formData, setFormData, getSaveData } = useModelForm(model as AIModelSchema);

    const isTestDisabled = !formData.name || Object.keys(formData.request_template || {}).length === 0;

    const handleSubmit = async () => {
        try {
            const saveData = getSaveData();
            await updateModel(saveData);
        } catch (error) {
            console.error('Submit failed:', error);
        }
    };

    if (isLoading) {
        return (
            <div className="container mx-auto py-8 max-w-4xl">
                <div className="flex justify-center items-center h-64">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </div>
            </div>
        );
    }

    if (!model && !isLoading) {
        return (
            <div className="container mx-auto py-8 max-w-4xl">
                <div className="text-center py-12">
                    <h2 className="text-2xl font-semibold text-muted-foreground">Model not found</h2>
                    <Button onClick={() => navigate({ to: '/ai-management' })} className="mt-4">
                        Back to Models
                    </Button>
                </div>
            </div>
        );
    }

    const handleBack = () => {
        navigate({ to: "/ai-management" });
    };

    return (
        <div className="container mx-auto">
            <div className="flex items-center space-x-2 mb-4">
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={handleBack}
                >
                    <ArrowLeft className="h-4 w-4 mr-1" />
                </Button>
                <PageTitle>Edit Model</PageTitle>

            </div>
            <p className="text-muted-foreground">Update AI model configuration</p>

            <ModelFormLayout
                formData={formData}
                setFormData={setFormData}
                providers={providers}
                families={families}
                schemas={schemas}
                initialData={model as AIModelSchema}
                onSubmit={handleSubmit}
                onCancel={() => navigate({ to: '/ai-management' })}
                onTestModel={() => setShowTestModel(true)}
                isTestDisabled={isTestDisabled}
                isLoading={isLoading}
                isSubmitting={isSubmitting}
            />

            <TestModelDialog
                open={showTestModel}
                onOpenChange={setShowTestModel}
                family={{
                    id: '',
                    provider_id: formData.provider_id,
                    schema_id: '',
                    name: formData.name,
                    display_name: formData.display_name,
                    description: null,
                    schema: {
                        id: '',
                        provider_id: formData.provider_id,
                        name: 'custom',
                        endpoint_path: '',
                        request_template: JSON.stringify(formData.request_template),
                        response_text_path: "",
                        response_image_path: "",
                        supports_temperature: true,
                        supports_streaming: true,
                    },
                }}
                providerConfig={providers.find((p: any) => p.id === formData.provider_id) as APIProvider}
                onPathsSelected={(textPath, imagePath) => {
                    setFormData((prev: any) => ({
                        ...prev,
                        response_text_path: textPath,
                        response_image_path: imagePath || "",
                    }));
                }}
            />
        </div>
    );
}