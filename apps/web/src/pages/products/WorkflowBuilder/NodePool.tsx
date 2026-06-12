// src/pages/products/WorkflowBuilder/NodePanel.tsx (REVISED)
import { useState, useEffect } from "react";
import { X, Trash2 } from "lucide-react";
import { Button, Input, Label, Textarea, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@kana-consultant/ui-kit";
import type { Node } from "reactflow";

interface NodePanelProps {
  node: Node;
  adapterConfigs: any[];
  onClose: () => void;
  onDelete: (nodeId: string) => void;
  onUpdate: (nodeId: string, updates: Record<string, any>) => void;
}

export function NodePanel({ node, adapterConfigs, onClose, onDelete, onUpdate }: NodePanelProps) {
  const [selectedAdapterId, setSelectedAdapterId] = useState(
    node.data.workflowNode?.adapterConfigId || ""
  );
  const [inputMapping, setInputMapping] = useState(
    JSON.stringify(node.data.workflowNode?.inputMapping || {}, null, 2)
  );

  useEffect(() => {
    setSelectedAdapterId(node.data.workflowNode?.adapterConfigId || "");
    setInputMapping(JSON.stringify(node.data.workflowNode?.inputMapping || {}, null, 2));
  }, [node.id]);

  const selectedAdapter = adapterConfigs.find((c) => c.id === selectedAdapterId);

  return (
    <div className="flex flex-col h-full py-4">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-semibold text-sm">Node Properties</h3>
        <div className="flex gap-1">
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6 text-destructive"
            onClick={() => onDelete(node.id)}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="icon" className="h-6 w-6" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className="space-y-4 flex-1">
        <div>
          <Label className="text-xs">Step Order</Label>
          <p className="text-sm mt-1">{node.data.workflowNode?.stepOrder}</p>
        </div>

        <div className="space-y-2">
          <Label className="text-xs">Adapter Config</Label>
          <Select
            value={selectedAdapterId}
            onValueChange={(value) => {
              setSelectedAdapterId(value);
              onUpdate(node.id, { adapterConfigId: value });
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder="Pilih adapter" />
            </SelectTrigger>
            <SelectContent>
              {adapterConfigs.map((config) => (
                <SelectItem key={config.id} value={config.id}>
                  {config.httpMethod} {config.endpointPath}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {selectedAdapter && (
          <div className="text-xs text-muted-foreground space-y-1 bg-muted p-2 rounded">
            <p><strong>Endpoint:</strong> {selectedAdapter.endpointPath}</p>
            <p><strong>Method:</strong> {selectedAdapter.httpMethod}</p>
            <p><strong>Timeout:</strong> {selectedAdapter.timeoutSeconds || 30}s</p>
            <p><strong>Retry:</strong> {selectedAdapter.retryCount || 3}x</p>
          </div>
        )}

        <div className="space-y-2">
          <Label className="text-xs">Input Mapping</Label>
          <Textarea
            value={inputMapping}
            onChange={(e) => setInputMapping(e.target.value)}
            rows={8}
            className="text-xs font-mono"
            placeholder='{"key": "value"}'
            onBlur={() => {
              try {
                const parsed = JSON.parse(inputMapping);
                onUpdate(node.id, { inputMapping: parsed });
              } catch {}
            }}
          />
          <p className="text-[10px] text-muted-foreground">
            Gunakan {"{{prev.response.field}}"} untuk referensi node sebelumnya
          </p>
        </div>
      </div>
    </div>
  );
}