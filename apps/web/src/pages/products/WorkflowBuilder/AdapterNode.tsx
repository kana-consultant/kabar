// src/pages/products/WorkflowBuilder/AdapterNode.tsx
import { memo } from "react";
import { Handle, Position, type NodeProps } from "reactflow";

export const AdapterNode = memo(({ data, selected }: NodeProps) => {
  return (
    <div
      className={`px-4 py-2 shadow-md rounded-md border-2 bg-card ${
        selected ? "border-primary" : "border-muted"
      }`}
    >
      <Handle type="target" position={Position.Top} className="!bg-blue-500" />
      
      <div className="flex flex-col">
        <div className="flex items-center gap-2">
          <span className="text-xs font-bold text-primary bg-primary/10 px-1.5 py-0.5 rounded">
            {data.adapterConfig?.httpMethod || "POST"}
          </span>
          <span className="text-xs font-medium truncate max-w-[120px]">
            {data.label}
          </span>
        </div>
        <p className="text-[10px] text-muted-foreground mt-1">
          Step {data.workflowNode?.stepOrder || "?"}
        </p>
      </div>

      <Handle type="source" position={Position.Bottom} className="!bg-green-500" />
    </div>
  );
});

AdapterNode.displayName = "AdapterNode";