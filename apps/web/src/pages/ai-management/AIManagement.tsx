// pages/admin/ai-management/AIManagement.tsx
import { useState } from "react";
import { ProvidersTab } from "./ProvidersTab";
import { useAIManagement } from "@/hooks/useAIManagement/useAIManagement";

export default function AIManagement() {
  const {
    providers,
    loading,
    deleteProvider,
  } = useAIManagement();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">AI Provider Management</h1>
          <p className="text-muted-foreground">
            Manage API providers configuration
          </p>
        </div>
      </div>

      <ProvidersTab
        providers={providers}
        loading={loading}
        onDelete={deleteProvider}
      />
    </div>
  );
}