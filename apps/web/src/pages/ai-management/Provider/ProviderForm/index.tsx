// pages/admin/ai-management/components/ProviderForm/index.tsx

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@kana-consultant/ui-kit";
import { BasicInfoSection } from "./BasicInfoSection";
import { ConnectionSection } from "./ConnectionSection";
import { FamiliesSection } from "./FamiliesSection";
import { StatusSection } from "./StatusSection";
import { FormActions } from "./FormActions";
import { useProviderForm } from "@/hooks/useAIManagement/useProviderForm";
import type { APIProvider, ProviderFormData, TestConnectionConfig } from "@/types/provider.types";

interface ProviderFormProps {
    mode: "create" | "edit";
    initialData?: APIProvider;
    onSubmit: (data: ProviderFormData) => Promise<void>;
    onTestConnection?: (config: TestConnectionConfig) => void;
    isSubmitting?: boolean;
    isTesting?: boolean;
}

export function ProviderForm({
    mode,
    initialData,
    onSubmit,
    onTestConnection,
    isSubmitting = false,
   
}: ProviderFormProps) {
    const {
        formData,
        updateFormData,
        validateForm,
        getTestConnectionConfig,
    } = useProviderForm(initialData || null);

    const [errors, setErrors] = useState<Record<string, string>>({});

    const handleSubmit = async () => {
        const validationErrors = validateForm();
        // if (Object.keys(validationErrors).length > 0) {
        //     setErrors(validationErrors);
        //     return;
        // }
        setErrors({});
        await onSubmit(formData);
    };

    const handleTestConnection = () => {
        if (onTestConnection) {
            const config = getTestConnectionConfig();
            onTestConnection(config);
        }
    };

    console.log("formData : ",formData)

    return (
        <form onSubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
            <div className="space-y-6">
                <Card>
                    <CardHeader>
                        <CardTitle>Basic Information</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <BasicInfoSection
                            value={formData}
                            onChange={updateFormData}
                            errors={errors}
                        />
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>Connection Settings</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <ConnectionSection
                            value={formData}
                            onChange={updateFormData}
                            errors={errors}
                            onTestConnection={handleTestConnection}
                        />
                    </CardContent>
                </Card>

                <Card>
                    <CardContent>
                        <CardHeader>
                            <CardTitle>API Families</CardTitle>
                        </CardHeader>
                        <FamiliesSection
                            value={formData}
                            onChange={updateFormData}
                            errors={errors}
                            providerConfig={{                      // ← Tambahkan ini
                                base_url: formData.base_url,
                                auth_header: formData.auth_header,
                                auth_prefix: formData.auth_prefix,
                                default_headers: formData.default_headers,
                                display_name: formData.display_name,
                            }}
                        />
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader>
                        <CardTitle>Status</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <StatusSection
                            value={formData.is_active}
                            onChange={(isActive) => updateFormData({ is_active: isActive })}
                        />
                    </CardContent>
                </Card>

                <FormActions
                    mode={mode}
                    isSubmitting={isSubmitting}
                    onCancel={() => window.history.back()}
                />
            </div>
        </form >
    );
}