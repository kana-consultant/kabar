import { Button } from "@kana-consultant/ui-kit";
import type { ModelFormData } from "../provider.types";
import { BasicInfoSection } from "./sections/BasicInfoSection";
import { RequestTemplateSection } from "./sections/RequestTemplateSection";
import { DefaultParamsSection } from "./sections/DefaultParamsSection";
import { ResponseConfigSection } from "./sections/ResponseConfigSection";
import { StatusSection } from "./sections/StatusSection";
import type { ModelFormProps } from "../provider.types";
import { useState } from "react";
import type { RequestSchema } from "@/types/provider.types";

interface ModelFormLayoutProps extends ModelFormProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    onTestModel: () => void;
    isTestDisabled: boolean;
    isSubmitting?: boolean;
    isLoading?: boolean; 
}

export function ModelFormLayout({
    formData,
    setFormData,
    providers = [],
    families = [], 
    initialData,
    onSubmit,
    onCancel,
    onTestModel,
    isTestDisabled,
    isSubmitting = false,
    isLoading = false,
}: ModelFormLayoutProps) {
    const [errors, setErrors] = useState<Record<string, string>>({});
    console.log("lfasnjmlfsdglkjfsdgkljsdgkjhdgskn=============================")
    console.log(families)
    // ✅ FIX #3 - Add guard for undefined
    const filteredFamilies = families?.filter(f => f.provider_id === formData.provider_id) || [];
    const selectedFamily = families?.find(f => f.id === formData.family_id);
    const selectedSchema = selectedFamily?.schema || null; // ✅ FIX #4
    
    // ✅ FIX #5 - Simplify with nullish coalescing
    const getResponseTextPath = () => 
        formData.response_text_path ?? selectedSchema?.response_text_path ?? "";
    
    const getResponseImagePath = () => 
        formData.response_image_path ?? selectedSchema?.response_image_path ?? "";
    
    // ✅ FIX #7 - Add validation
    const validateForm = (): boolean => {
        const newErrors: Record<string, string> = {};
        
        if (!formData.name?.trim()) newErrors.name = "Name is required";
        if (!formData.display_name?.trim()) newErrors.display_name = "Display name is required";
        if (!formData.provider_id) newErrors.provider_id = "Provider is required";
        
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };
    
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validateForm()) return;
        await onSubmit(formData);
    };
    
    // ✅ FIX #2 - Better loading UI
    if (isLoading) {
        return (
            <div className="flex justify-center items-center p-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
        );
    }
    
    // ✅ FIX #6 - Better validation for submit button
    const isFormValid = formData.name?.trim() && formData.display_name?.trim() && formData.provider_id;
    
    return (
        <form onSubmit={handleSubmit} className="space-y-8 mt-5">
            <BasicInfoSection
                formData={formData}
                setFormData={setFormData}
                providers={providers}
                filteredFamilies={filteredFamilies}
                isEditMode={!!initialData}
            />
            
            {/* <RequestTemplateSection
                formData={formData}
                setFormData={setFormData}
                selectedSchema={selectedSchema as RequestSchema}
            /> */}
            
            <DefaultParamsSection
                formData={formData}
                setFormData={setFormData}
            />
            
            {/* <ResponseConfigSection
                formData={formData}
                setFormData={setFormData}
                selectedSchema={selectedSchema as RequestSchema}
                onTestModel={onTestModel}
                isTestDisabled={isTestDisabled}
                getResponseTextPath={getResponseTextPath}
                getResponseImagePath={getResponseImagePath}
            /> */}
            
            <StatusSection
                formData={formData}
                setFormData={setFormData}
                isEditMode={!!initialData}
                createdBy={initialData?.created_by}
            />
            
            <div className="flex justify-end space-x-2 pt-4 border-t">
                <Button type="button" variant="outline" onClick={onCancel}>
                    Cancel
                </Button>
                <Button 
                    type="submit" 
                    disabled={isSubmitting || !isFormValid}
                >
                    {isSubmitting ? "Saving..." : (initialData ? "Update Model" : "Add Model")}
                </Button>
            </div>
        </form>
    );
}