// pages/admin/ai-management/create.tsx

import { useModelForm } from "@/hooks/useModelManagement/useModelForm";
import { ModelFormLayout } from "./ModelForm/ModelFormLayout";
import { TestModelDialog } from "../TestModelDialog";
import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useCreateModel } from "@/hooks/useModelManagement/useCreateModel";
import { useToast } from "@/hooks/use-toast";
import type { APIProvider } from "@/types/provider.types";

export default function CreateModelPage() {
    const navigate = useNavigate();
    const toast = useToast()
    const [showTestModel, setShowTestModel] = useState(false);
    
    // Use hook for data fetching and mutations
    const { providers, families, schemas, isLoading, isSubmitting } = useCreateModel();
    
    // Use form hook
    const { formData, setFormData, handleSubmitModel } = useModelForm();
    
    const isTestDisabled = !formData.name || Object.keys(formData.request_template || {}).length === 0;
    
    const handleSubmit = async () => {
        try {
            const saveData = handleSubmitModel();
          
            toast.success('Model created successfully');
            navigate({ to: '/ai-management' });
        } catch (error: any) {
            toast.error(error.message || 'Failed to create model');
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
    
    return (
        <div className="container mx-auto ">
            <div className="mb-6">
                <h1 className="text-2xl font-bold">Create New Model</h1>
                <p className="text-muted-foreground">Configure a new AI model for your application</p>
            </div>
            
            <ModelFormLayout
                formData={formData}
                setFormData={setFormData}
                providers={providers}
                families={families}
                schemas={schemas}
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
                    id: "",
                    provider_id: formData.provider_id,
                    schema_id: "",
                    name: formData.name,
                    display_name: formData.display_name,
                    description: null,
                    schema: {
                        id: "",
                        provider_id: formData.provider_id,
                        name: "custom",
                        endpoint_path: "",
                        request_template: JSON.stringify(formData.request_template),
                        response_text_path: "",
                        response_image_path: "",
                        supports_temperature: true,
                        supports_streaming: true,
                    },
                }}
                providerConfig={providers.find(p => p.id === formData.provider_id) as APIProvider}
                onPathsSelected={(textPath, imagePath) => {
                    setFormData(prev => ({
                        ...prev,
                        response_text_path: textPath,
                        response_image_path: imagePath || null,
                    }));
                }}
            />
        </div>
    );
}