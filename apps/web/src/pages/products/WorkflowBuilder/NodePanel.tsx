// src/pages/products/WorkflowBuilder/NodePanel.tsx
import { useState, useEffect, useMemo } from "react";
import { Trash2, Play, X } from "lucide-react";
import { Button, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@kana-consultant/ui-kit";
import type { Node } from "reactflow";
import { SimpleJsonBuilder } from "../SimpleJsonBuilder";
import { useToast } from "@/hooks/use-toast";

interface NodePanelProps {
    node: Node;
    adapterConfigs: any[];
    allNodes: Node[];
    product: {
        id: string | undefined;
        apiEndpoint: string;
        apiKey?: string;
    };
    onClose: () => void;
    onDelete: (nodeId: string) => void;
    onUpdate: (nodeId: string, updates: Record<string, any>) => void;
}

interface PlaceholderItem {
    value: string;
    label: string;
    group?: string;
}

export function NodePanel({
    node,
    adapterConfigs,
    allNodes,
    product,
    onClose,
    onDelete,
    onUpdate
}: NodePanelProps) {
    const [selectedAdapterId, setSelectedAdapterId] = useState(
        node.data.workflowNode?.adapterConfigId || ""
    );
    const [inputMapping, setInputMapping] = useState<any>(
        node.data.workflowNode?.inputMapping || {}
    );
    const [saving, setSaving] = useState(false);
    const [executing, setExecuting] = useState(false);
    const [executionResult, setExecutionResult] = useState<any>(null);
    const [isModalOpen, setIsModalOpen] = useState(true);
    const [tempSelectedAdapterId, setTempSelectedAdapterId] = useState("");
    const [tempInputMapping, setTempInputMapping] = useState<string>("");

    const toast = useToast();

    const currentNodeStepOrder = node.data.workflowNode?.stepOrder;
    const isFirstNode = currentNodeStepOrder === 1;

    useEffect(() => {
        setSelectedAdapterId(node.data.workflowNode?.adapterConfigId || "");
        setInputMapping(node.data.workflowNode?.inputMapping || {});
        setTempSelectedAdapterId(node.data.workflowNode?.adapterConfigId || "");
        setTempInputMapping(node.data.workflowNode?.inputMapping || {});
    }, [node.id, node.data.workflowNode]);

    const selectedAdapter = adapterConfigs.find((c) => c.id === tempSelectedAdapterId);

    // ==================== GENERATE PLACEHOLDERS FROM PREVIOUS NODES ====================

    const previousNodes = useMemo(() => {
        const previousIds = node.data.workflowNode?.previousNodeIds || [];
        const result = allNodes.filter(n => previousIds.includes(n.id));
        return result;
    }, [allNodes, node.data.workflowNode?.previousNodeIds, node.id]);

    const generateNodePlaceholders = (prevNode: Node): PlaceholderItem[] => {
        const nodeName = prevNode.data.workflowNode?.name || `Node ${prevNode.data.workflowNode?.stepOrder}`;
        const nodeId = prevNode.id;

        const responseExample =
            prevNode.data.workflowNode?.responseExample ||
            prevNode.data.workflowNode?.lastExecution?.data ||
            prevNode.data.workflowNode?.responseData;

        const placeholders: PlaceholderItem[] = [];

        const generateNestedPlaceholders = (obj: any, currentPath: string) => {
            if (!obj || typeof obj !== "object") return;

            Object.entries(obj).forEach(([key, value]) => {
                const fullPath = currentPath ? `${currentPath}.${key}` : key;
                const placeholderValue = `{{node_${nodeId}.response.${fullPath}}}`;

                if (value && typeof value === "object" && !Array.isArray(value)) {
                    generateNestedPlaceholders(value, fullPath);
                } else {
                    placeholders.push({
                        value: placeholderValue,
                        label: fullPath,
                        group: `🔗 ${nodeName}`
                    });
                }
            });
        };

        if (responseExample && typeof responseExample === "object") {
            generateNestedPlaceholders(responseExample, "");
        }

        return placeholders;
    };

    const nodePlaceholders = useMemo(() => {
        if (isFirstNode) {
            return [];
        }

        const placeholders: PlaceholderItem[] = [];
        previousNodes.forEach(prevNode => {
            const generated = generateNodePlaceholders(prevNode);
            if (generated.length > 0) {
                placeholders.push(...generated);
            }
        });
        return placeholders;
    }, [previousNodes, isFirstNode, node.id]);

    const allPlaceholders = useMemo(() => {
        return nodePlaceholders;
    }, [nodePlaceholders]);

    // ==================== EXECUTE NODE ====================
    const handleExecute = async () => {
        if (!selectedAdapter) {
            toast.error("Pilih endpoint terlebih dahulu");
            return;
        }

        setExecuting(true);
        setExecutionResult(null);

        try {
            const baseUrl = product.apiEndpoint?.replace(/\/$/, '') || '';
            const endpointPath = selectedAdapter.endpointPath?.startsWith('/')
                ? selectedAdapter.endpointPath
                : `/${selectedAdapter.endpointPath || ''}`;
            const fullUrl = `${baseUrl}${endpointPath}`;

            if (!fullUrl || fullUrl === '/') {
                toast.error("URL tidak valid. Periksa Base URL dan Endpoint Path");
                setExecuting(false);
                return;
            }

            let headers: Record<string, string> = {
                'Content-Type': 'application/json'
            };

            if (selectedAdapter.customHeaders) {
                if (typeof selectedAdapter.customHeaders === 'string') {
                    try {
                        const parsed = JSON.parse(selectedAdapter.customHeaders);
                        headers = { ...headers, ...parsed };
                    } catch (e) {
                        // ignore
                    }
                } else {
                    headers = { ...headers, ...selectedAdapter.customHeaders };
                }
            }

            if (product.apiKey) {
                headers['X-API-Key'] = product.apiKey;
            }

            let body = null;
            if (selectedAdapter.httpMethod !== 'GET' && selectedAdapter.httpMethod !== 'DELETE') {
                body = JSON.stringify(tempInputMapping);
            }

            const response = await fetch(fullUrl, {
                method: selectedAdapter.httpMethod,
                headers,
                body,
            });

            let data;
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                data = await response.json();
            } else {
                data = await response.text();
            }

            if (response.ok) {
                toast.success(`Success: ${response.status}`);

                onUpdate(node.id, {
                    responseExample: data,
                    lastExecution: {
                        status: response.status,
                        statusText: response.statusText,
                        timestamp: new Date().toISOString(),
                        data
                    }
                });

                setExecutionResult({
                    status: response.status,
                    statusText: response.statusText,
                    data,
                });
            } else {
                toast.error(`Error: ${response.status} ${response.statusText}`);
                setExecutionResult({
                    status: response.status,
                    statusText: response.statusText,
                    data,
                    error: true,
                });
            }
        } catch (error: any) {
            console.error("Execution error:", error);
            toast.error(error.message);
            setExecutionResult({
                error: true,
                message: error.message,
            });
        } finally {
            setExecuting(false);
        }
    };

    const handleSave = async () => {
        setSaving(true);
        try {
            const adapter = adapterConfigs.find((c) => c.id === tempSelectedAdapterId);
            onUpdate(node.id, {
                inputMapping: tempInputMapping,
                adapterConfigId: tempSelectedAdapterId,
                ...(adapter && node.data.label !== adapter.endpointPath && {
                    label: adapter.endpointPath
                })
            });

            setSelectedAdapterId(tempSelectedAdapterId);
            setInputMapping(tempInputMapping);

            toast.success("Node saved");
            onClose();
        } catch (error) {
            toast.error("Failed to save node");
        } finally {
            setSaving(false);
        }
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        onClose();
    };

    const previousNodesCount = previousNodes.length;
    const hasPlaceholders = nodePlaceholders.length > 0;

    return (
        <Dialog open={isModalOpen} onOpenChange={handleCloseModal}>
            <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle className="flex items-center justify-between">
                        <span>Node Configuration - Step {currentNodeStepOrder}</span>
                    </DialogTitle>
                    <DialogDescription>
                        Configure endpoint and input mapping for this workflow node
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 mt-4">
                    {/* Step Order Info */}
                    <div className="bg-muted/50 p-3 rounded-lg">
                        <Label className="text-xs text-muted-foreground">Step Order</Label>
                        <p className="text-lg font-semibold mt-1">{currentNodeStepOrder}</p>
                    </div>

                    {/* Adapter Config Selector */}
                    <div className="space-y-2">
                        <Label className="text-sm font-semibold">Select Endpoint</Label>
                        <Select
                            value={tempSelectedAdapterId}
                            onValueChange={setTempSelectedAdapterId}
                        >
                            <SelectTrigger>
                                <SelectValue placeholder="Select endpoint" />
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

                    {/* Adapter Info */}
                    {selectedAdapter && (
                        <div className="text-xs text-muted-foreground space-y-1 bg-muted p-3 rounded-lg">
                            <p><strong>Method:</strong> {selectedAdapter.httpMethod}</p>
                            <p><strong>Path:</strong> {selectedAdapter.endpointPath}</p>
                            <p><strong>Base URL:</strong> {product.apiEndpoint}</p>
                            <p><strong>Full URL:</strong> {product.apiEndpoint?.replace(/\/$/, '')}{selectedAdapter.endpointPath?.startsWith('/') ? selectedAdapter.endpointPath : `/${selectedAdapter.endpointPath}`}</p>
                            <p><strong>Timeout:</strong> {selectedAdapter.timeoutSeconds || 30}s</p>
                            <p><strong>Retry:</strong> {selectedAdapter.retryCount || 3}x</p>
                        </div>
                    )}

                    {/* Input Mapping */}
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <Label className="text-sm font-semibold">Input Mapping</Label>
                            {!isFirstNode && (
                                <span className="text-xs text-muted-foreground">
                                    {hasPlaceholders
                                        ? `✅ ${nodePlaceholders.length} placeholders available`
                                        : `⚠️ Execute previous node first to get placeholders`}
                                </span>
                            )}
                        </div>
                        <div className="border rounded-lg p-4 bg-muted/30">
                            <SimpleJsonBuilder
                                value={tempInputMapping}
                                onChange={setTempInputMapping}
                                placeholders={allPlaceholders}
                            />
                        </div>
                        <p className="text-xs text-muted-foreground">
                            {isFirstNode
                                ? "First node has no placeholders from previous nodes"
                                : hasPlaceholders
                                    ? `Click 📋 button to select placeholders from ${previousNodesCount} previous node(s)`
                                    : "Execute previous nodes first to get response placeholders"
                            }
                        </p>
                    </div>

                    {/* Execute Button */}
                    <div className="space-y-2 pt-2">
                        <Button
                            type="button"
                            onClick={handleExecute}
                            disabled={executing || !selectedAdapter}
                            className="w-full"
                            variant="outline"
                        >
                            <Play className="h-4 w-4 mr-2" />
                            {executing ? "Executing..." : "Test Execute Node"}
                        </Button>
                    </div>

                    {/* Execution Result */}
                    {executionResult && (
                        <div className={`space-y-2 p-3 rounded-lg ${executionResult.error ? 'bg-red-50 dark:bg-red-900/20' : 'bg-green-50 dark:bg-green-900/20'}`}>
                            <Label className="text-xs font-semibold">Execution Result:</Label>
                            <pre className="text-xs overflow-auto max-h-40 p-2 bg-white dark:bg-slate-900 rounded">
                                {JSON.stringify(executionResult, null, 2)}
                            </pre>
                        </div>
                    )}
                </div>

                <div className="flex justify-end gap-2 mt-6">
                    <Button variant="outline" onClick={handleCloseModal}>
                        Cancel
                    </Button>
                    <Button onClick={handleSave} disabled={saving}>
                        {saving ? "Saving..." : "Save Node"}
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}