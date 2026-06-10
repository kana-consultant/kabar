// pages/ai-management/model/$id.edit.tsx

import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "@kana-consultant/ui-kit";
import { ArrowLeft } from "lucide-react";
import { ModelFormLayout as ModelForm } from "../ModelForm/ModelFormLayout";
import { useModelManagement } from "@/hooks/useAIManagement/useAIModelManagement";
import { useProviderManagement } from "@/hooks/useAIManagement/useProviderManagement";
import { useEffect } from "react";

export const Route = createFileRoute('/model/$id/edit')({
  component: EditModelPage,
});

function EditModelPage() {
  const params = Route.useParams(); // ✅ PAKE USEPARAMS
  const { id } = params; // ambil id dari params
  const navigate = useNavigate();
  const { selectedModel, loadModel, updateModel, loading } = useModelManagement();
  const { providers } = useProviderManagement();

  useEffect(() => {
    loadModel(id);
  }, [id, loadModel]);

  const handleSubmit = async (data: any) => {
    try {
      await updateModel(id, data);
      navigate({ to: '/model' }); // ✅ FIX: /ai-management/model → /model
    } catch (error) {
      console.error('Failed to update model:', error);
    }
  };

  const handleCancel = () => {
    navigate({ to: '/model' }); // ✅ FIX: /ai-management/model → /model
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
          <p className="mt-2 text-gray-500">Loading model data...</p>
        </div>
      </div>
    );
  }

  if (!selectedModel) {
    return (
      <div className="text-center py-8">
        <p className="text-red-500">Model not found</p>
        <Button variant="link" onClick={handleCancel} className="mt-2">
          Back to Models
        </Button>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-8 px-4">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Button variant="ghost" onClick={handleCancel}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Edit Model</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Update configuration for {selectedModel.display_name}
          </p>
        </div>
      </div>

      {/* Form Card */}
      <div className="bg-white rounded-lg border shadow-sm">
        {/* <div className="p-6">
          <ModelForm
            initialData={selectedModel}
            providers={providers}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
          />
        </div> */}
      </div>
    </div>
  );
}