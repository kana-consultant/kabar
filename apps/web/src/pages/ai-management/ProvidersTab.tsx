// pages/admin/ai-management/ProvidersTab.tsx
import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { Button } from "@kana-consultant/ui-kit";
import { Plus } from "lucide-react";
import { ProviderList } from "./ProviderList";
import { DeleteAlertDialog } from "./DeleteAlertDialog";
import { type APIProvider } from "@/types/provider.types";

interface ProvidersTabProps {
  providers: APIProvider[];
  loading: boolean;
  onDelete: (id: string) => void;
}

export function ProvidersTab({ 
  providers, 
  loading, 
  onDelete
}: ProvidersTabProps) {
  const navigate = useNavigate();
  const [deleteProvider, setDeleteProvider] = useState<APIProvider | null>(null);

  const handleAdd = () => {
    navigate({ to: "/provider/add" });
  };

  const handleEdit = (provider: APIProvider) => {
    navigate({ 
      to: "/provider/$id/edit",
      params: { id: provider.id as string}
    });
  };

  return (
    <>
      <div className="flex justify-end">
        <Button onClick={handleAdd}>
          <Plus className="w-4 h-4 mr-2" />
          Add Provider
        </Button>
      </div>

      <ProviderList
        providers={providers}
        loading={loading}
        onEdit={handleEdit}
        onDelete={(provider) => setDeleteProvider(provider)}
      />

      <DeleteAlertDialog
        open={!!deleteProvider}
        onOpenChange={() => setDeleteProvider(null)}
        itemName={deleteProvider?.display_name || ""}
        itemType="provider"
        onConfirm={() => {
          if (deleteProvider) onDelete(deleteProvider.id as string);
          setDeleteProvider(null);
        }}
      />
    </>
  );
}