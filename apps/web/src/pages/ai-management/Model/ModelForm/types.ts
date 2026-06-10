// components/ModelForm/ModelForm.types.ts

import type { APIProvider,Family, AIModel } from "@/types/provider.types";
import type { ModelFormData } from "../provider.types";
import type { RequestSchema } from "@/types/provider.types";

export interface ModelFormLayoutProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    providers: APIProvider[];
    families: Family[];
    schemas: RequestSchema[];
    initialData?: AIModel | null;
    onSubmit: (data: ModelFormData) => Promise<void>;
    onCancel: () => void;
    onTestModel: () => void;
    isTestDisabled: boolean;
    isLoading?: boolean;
    isSubmitting?: boolean;
}

// ModelFormProps (untuk komponen utama form)
export interface ModelFormProps {
    mode: 'create' | 'edit';
    initialData?: AIModel | null;
    providers: APIProvider[];
    families: Family[];
    schemas: RequestSchema[];
    onSubmit: (data: ModelFormData) => Promise<void>;
    onCancel: () => void;
    onTestModel?: () => void;
    isSubmitting?: boolean;
}

// Section props
export interface BasicInfoSectionProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    providers: APIProvider[];
    filteredFamilies: Family[];
    isEditMode?: boolean;
}

export interface RequestTemplateSectionProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    selectedSchema?: RequestSchema;
}

export interface DefaultParamsSectionProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
}

export interface ResponseConfigSectionProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    selectedSchema?: RequestSchema;
    onTestModel: () => void;
    isTestDisabled: boolean;
    getResponseTextPath: () => string;
    getResponseImagePath: () => string;
}

export interface StatusSectionProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    isEditMode?: boolean;
    createdBy?: string | null;
}

// JsonBuilder props
export interface JsonBuilderProps {
    value: Record<string, any>;
    onChange: (value: Record<string, any>) => void;
    availableVariables?: string[];
}