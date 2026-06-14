// src/pages/products/WorkflowBuilder/WorkflowNode.tsx
import { Handle, Position } from "reactflow";
import { Settings } from "lucide-react";

interface WorkflowNodeProps {
    data: {
        label: string;
        workflowNode: any;
        adapterConfig: any;
    };
    selected: boolean;
}

export function WorkflowNode({ data, selected }: WorkflowNodeProps) {
    return (
        <div
            className={`
        px-4 py-2 rounded-lg border-2 min-w-[180px] bg-background
        ${selected ? 'border-primary shadow-lg' : 'border-gray-300 shadow'}
        hover:shadow-md transition-shadow cursor-pointer
      `}
        >
            <Handle
                type="target"
                position={Position.Left}
                className="w-3 h-3 bg-primary"
            />

            <div className="flex items-center gap-2">
                <Settings className="h-4 w-4 text-muted-foreground" />
                <div className="font-medium text-sm truncate">{data.label}</div>
            </div>

            <div className="text-xs text-muted-foreground mt-1 truncate">
                {data.adapterConfig?.http_method || 'GET'} {data.adapterConfig?.endpoint_path || ''}
            </div>

            <Handle
                type="source"
                position={Position.Right}
                className="w-3 h-3 bg-primary"
            />
        </div>
    );
}