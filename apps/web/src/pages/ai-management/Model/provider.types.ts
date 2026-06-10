
import type { APIProvider,AIModelSchema,Family } from "@/types/provider.types";
import type { RequestSchema } from "@/types/provider.types";

export interface ModelFormData {
    // Relations
    provider_id: string;
    family_id: string;

    // Basic Info
    name: string;
    display_name: string;
    description: string | null;

    // Request Template
    request_template: Record<string, any> | string | null;

    // Response Paths
    response_text_path: string | null;
    response_image_path: string | null;

    // Default Parameters
    max_tokens: number | null;
    temperature: number | null;

    // Status
    is_active: boolean;
    is_default: boolean;

    // Ownership
    team_id: string | null;
    created_by: string | null;
}

export interface ModelFormProps {
    initialData?: AIModelSchema;
    providers: APIProvider[];
    families: Family[];
    schemas: RequestSchema[];
    currentUserId?: string;
    onSubmit: (data: ModelFormData) => Promise<void>;
    onCancel: () => void;
    isLoading?: boolean;
}

export interface ModelFormSectionProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    providers?: APIProvider[];
    families?: Family[];
    schemas?: RequestSchema[];
    selectedFamily?: Family;
    selectedSchema?: RequestSchema;
    filteredFamilies?: Family[];
    onTestModel?: () => void;
    isTestDisabled?: boolean;
    isEditMode?: boolean;
}

export interface TemplateEditorProps {
    formData: ModelFormData;
    setFormData: React.Dispatch<React.SetStateAction<ModelFormData>>;
    selectedSchema?: RequestSchema;
}