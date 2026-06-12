// src/pages/products/WorkflowBuilder/WorkflowList.tsx
import { Trash2 } from "lucide-react";
import { Button } from "@kana-consultant/ui-kit";
import type { WorkflowDefinition } from "@/types/workflow";

interface WorkflowListProps {
    workflows: WorkflowDefinition[];
    selectedId: string | null;
    onSelect: (id: string) => void;
    onDelete: (id: string) => void;
    loading: boolean;
}

export function WorkflowList({
    workflows,
    selectedId,
    onSelect,
    onDelete,
    loading,
}: WorkflowListProps) {
    if (loading) {
        return <div className="text-sm text-muted-foreground">Loading...</div>;
    }

    if (workflows.length === 0) {
        return (
            <div className="text-sm text-muted-foreground">
                Belum ada workflow
            </div>
        );
    }

    return (
        <div className="space-y-1">
            {workflows.map((wf) => (
                <div
                    key={wf.id}
                    className={`flex items-center justify-between p-2 rounded cursor-pointer text-sm ${selectedId === wf.id
                            ? "bg-primary/10 text-primary"
                            : "hover:bg-muted"
                        }`}
                    onClick={() => onSelect(wf.id)}
                >
                    <span className="truncate">{wf.name}</span>
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-6 w-6 opacity-0 group-hover:opacity-100"
                        onClick={(e) => {
                            e.stopPropagation();
                            onDelete(wf.id);
                        }}
                    >
                        <Trash2 className="h-3 w-3" />
                    </Button>
                </div>
            ))}
        </div>
    );
}