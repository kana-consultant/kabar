// pages/admin/ai-management/ModelsTab.tsx (updated)
import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { Button } from "@kana-consultant/ui-kit";
import { Plus } from "lucide-react";
import { ModelTable } from "./Model/ModelList/ModelTable";
import { DeleteAlertDialog } from "./DeleteAlertDialog";
import { ModelDetailModal } from "./ModelDetailModal";
import type { AIModel } from "@/types/provider.types";

interface ModelsTabProps {
  models: AIModel[];
  loading: boolean;
  onDelete: (id: string) => void;
}

export function ModelsTab({ 
  models, 
  loading, 
  onDelete, 
}: ModelsTabProps) {
  const navigate = useNavigate();
  const [deleteModel, setDeleteModel] = useState<AIModel | null>(null);
  const [viewModel, setViewModel] = useState<AIModel | null>(null);

  const handleAdd = () => {
    navigate({ to: "/model/add" });
  };

  const handleEdit = (id: string) => {
    navigate({ 
      to: "/model/$id/edit",
      params: { id: id }
    });
  };

  const handleViewDetail = (model: AIModel) => {
    setViewModel(model);
  };

  return (
    <>
      <div className="flex justify-end">
        <Button onClick={handleAdd}>
          <Plus className="w-4 h-4 mr-2" />
          Add Model
        </Button>
      </div>

      <ModelTable
        models={models}
        isLoading={loading}
        onEdit={handleEdit}
        onView={handleViewDetail}
        onDelete={(model) => setDeleteModel(model)}
      />

      <DeleteAlertDialog
        open={!!deleteModel}
        onOpenChange={() => setDeleteModel(null)}
        itemName={deleteModel?.display_name || ""}
        itemType="model"
        onConfirm={() => {
          if (deleteModel) onDelete(deleteModel.id);
          setDeleteModel(null);
        }}
      />

      <ModelDetailModal
        open={!!viewModel}
        onOpenChange={() => setViewModel(null)}
        model={viewModel}
      />
    </>
  );
}