// pages/admin/ai-management/create.tsx
import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import {
    Page, PageHeader,
    PageTitle,
    PageDescription,
    PageContent,
} from "@/components/page/page";
import { Button } from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";
import { ProviderForm } from "./ProviderForm";
import { createProvider } from "@/services/provider/provider-api";
import { testConnection } from "@/services/provider/provider-api";
import type { ProviderFormData, TestConnectionConfig } from "@/types/provider.types";

export default function CreateProviderPage() {
    const toast = useToast()
    const navigate = useNavigate();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isTesting, setIsTesting] = useState(false);

    const handleSubmit = async (data: ProviderFormData) => {
        setIsSubmitting(true);
        try {
            const newProvider = await createProvider(data);
            toast.success("Provider created successfully");
            navigate({
                to: "/ai-management",
                params: { id: newProvider.id! }
            });
        } catch (error) {
            toast.error("Failed to create provider");
            console.error(error);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleTestConnection = async (config: Partial<ProviderFormData>) => {
        setIsTesting(true);
        try {
            // Convert Partial<ProviderFormData> ke TestConnectionConfig
            const testConfig: TestConnectionConfig = {
                base_url: config.base_url || "",
                auth_type: config.auth_type || "bearer",
                auth_header: config.auth_header || "Authorization",
                auth_prefix: config.auth_prefix as string,
            };

            const result = await testConnection(testConfig);
            if (result.success) {
                toast.success(`Connection successful! (${result.latency}ms)`);
            } else {
                toast.error(`Connection failed: ${result.message}`);
            }
        } catch (error) {
            toast.error("Failed to test connection");
        } finally {
            setIsTesting(false);
        }
    };
    const handleBack = () => {
        navigate({ to: "/ai-management" });
    };

    return (
        <Page>
            <PageHeader>
                <div className="flex items-center space-x-4">
                    <Button
                        variant="ghost"
                        size="sm"
                        onClick={handleBack}
                    >
                        <ArrowLeft className="h-4 w-4 mr-2" />
                    </Button>
                    <div>
                        <PageTitle>Create New Provider</PageTitle>
                        <PageDescription>
                            Configure a new AI API provider with its authentication and API families
                        </PageDescription>
                    </div>
                </div>
            </PageHeader>

            <PageContent>
                <ProviderForm
                    mode="create"
                    onSubmit={handleSubmit}
                    onTestConnection={handleTestConnection}
                    isSubmitting={isSubmitting}
                    isTesting={isTesting}
                />
            </PageContent>
        </Page>
    );
}