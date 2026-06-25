// src/pages/products/WorkflowBuilder/NodePanel.tsx
import { useState, useEffect, useMemo } from "react";
import { Trash2, Play, Save } from "lucide-react";
import { Button, Label, Input, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, Textarea } from "@kana-consultant/ui-kit";
import type { Node } from "reactflow";
import { SimpleJsonBuilder } from "../SimpleJsonBuilder";
import { useToast } from "@/hooks/use-toast";
import type { AdapterConfig, AdapterConfigNode, WorkflowNode } from "@/types/product";

interface NodePanelProps {
    node: Node;
    allNodes: Node[];
    product: {
        id: string | undefined;
        api_endpoint: string;
        api_key?: string;
    };
    onClose: () => void;
    onDelete: (nodeId: string) => void;
    onUpdate: (nodeId: string, updates: Partial<WorkflowNode> & { adapter_config?: Partial<AdapterConfig> }) => void;
}

interface PlaceholderItem {
    value: string;
    label: string;
    group?: string;
}

export function NodePanel({
    node,
    allNodes,
    product,
    onClose,
    onDelete,
    onUpdate
}: NodePanelProps) {
    const [endpointPath, setEndpointPath] = useState<string>("");
    const [httpMethod, setHttpMethod] = useState<AdapterConfigNode['http_method']>("GET");
    const [customHeaders, setCustomHeaders] = useState<string>("{}");
    const [timeoutSeconds, setTimeoutSeconds] = useState<number>(30);
    const [retryCount, setRetryCount] = useState<number>(3);
    const [inputMapping, setInputMapping] = useState<Record<string, any>>({});
    const [saving, setSaving] = useState(false);
    const [executing, setExecuting] = useState(false);
    const [executionResult, setExecutionResult] = useState<any>(null);
    const [isModalOpen, setIsModalOpen] = useState(true);
    const [headersError, setHeadersError] = useState<string | null>(null);

    const toast = useToast();

    const workflowNode = node.data.workflowNode as WorkflowNode;
    const adapterConfig = node.data.adapterConfig as AdapterConfigNode;
    const currentNodeStepOrder = workflowNode?.step_order;
    const isFirstNode = currentNodeStepOrder === 1;

    // Initialize state from node data
    useEffect(() => {
        if (adapterConfig) {
            setEndpointPath(adapterConfig.endpoint_path || "");
            setHttpMethod(adapterConfig.http_method || "GET");

            // Handle custom_headers (bisa string atau object)
            let headersStr = "{}";
            if (typeof adapterConfig.custom_headers === 'string') {
                headersStr = adapterConfig.custom_headers;
            } else if (adapterConfig.custom_headers && typeof adapterConfig.custom_headers === 'object') {
                headersStr = JSON.stringify(adapterConfig.custom_headers, null, 2);
            }
            setCustomHeaders(headersStr);

            setTimeoutSeconds(adapterConfig.timeout_seconds || 30);
            setRetryCount(adapterConfig.retry_count || 3);
        }

        if (workflowNode) {
            setInputMapping(JSON.parse(workflowNode.adapter_config?.field_mapping || "{}"));
        }
    }, [adapterConfig, workflowNode]);


    console.log("json")
    console.log(inputMapping)


    // Get previous nodes
    const previousNodes = useMemo(() => {
        const previousIds = workflowNode?.previous_node_ids || [];
        return allNodes.filter(n => previousIds.includes(n.id));
    }, [allNodes, workflowNode?.previous_node_ids]);

    // Generate placeholders from previous nodes
    const generateNodePlaceholders = (prevNode: Node): PlaceholderItem[] => {
        const prevWorkflowNode = prevNode.data.workflowNode as WorkflowNode;
        const nodeName = prevWorkflowNode?.step_order
            ? `Step ${prevWorkflowNode.step_order}`
            : `Node ${prevNode.id.substring(0, 6)}`;
        const nodeId = prevNode.id;

        // Get response data from last execution or response example
        const responseData =
            prevNode.data.responseExample ||
            prevNode.data.lastExecution?.data ||
            prevNode.data.workflowNode?.last_execution?.data;

        const placeholders: PlaceholderItem[] = [];

        const generateNestedPlaceholders = (obj: any, currentPath: string) => {
            if (!obj || typeof obj !== "object") return;

            Object.entries(obj).forEach(([key, value]) => {
                const fullPath = currentPath ? `${key}` : key;
                const placeholderValue = `{{response.${fullPath}}}`;

                if (value && typeof value === "object" && !Array.isArray(value)) {
                    generateNestedPlaceholders(value, fullPath);
                } else {
                    let label = fullPath;
                    let group = `📦 ${nodeName}`;

                    // Format array access
                    if (Array.isArray(value)) {
                        label = `${fullPath}[]`;
                    }

                    // Add value preview
                    const valuePreview = typeof value === 'string' && value.length > 30
                        ? value.substring(0, 30) + '...'
                        : String(value);
                    if (valuePreview && valuePreview !== '[object Object]') {
                        label = `${fullPath} (${valuePreview})`;
                    }

                    placeholders.push({
                        value: placeholderValue,
                        label: label,
                        group: group
                    });
                }
            });
        };

        if (responseData && typeof responseData === "object") {
            generateNestedPlaceholders(responseData, "");
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
    }, [previousNodes, isFirstNode]);

    // Validate custom headers JSON
    const validateHeaders = (headersStr: string): boolean => {
        try {
            if (headersStr.trim() === "") {
                setHeadersError(null);
                return true;
            }
            const parsed = JSON.parse(headersStr);
            if (typeof parsed !== 'object' || Array.isArray(parsed)) {
                setHeadersError("Headers must be a JSON object");
                return false;
            }
            setHeadersError(null);
            return true;
        } catch (e) {
            setHeadersError("Invalid JSON format");
            return false;
        }
    };

    // Execute node
    // Execute node - MODIFIED VERSION


    // Save node changes - FIXED VERSION with proper typing
    // Save node changes - FIXED VERSION dengan type yang benar
    const handleSave = async () => {
        if (!endpointPath) {
            toast.error("Please enter an endpoint path");
            return;
        }

        setSaving(true);
        try {
            // Parse custom headers - langsung sebagai string, tidak perlu JSON.stringify lagi
            const adapterConfigUpdates: Partial<AdapterConfig> = {
                timeout_seconds: timeoutSeconds,
                retry_count: retryCount,
                updated_at: new Date().toISOString(),
            };

           
            console.log("📦 Saving node with updates:", {
                adapterConfigUpdates
            });

            // Combine updates dan kirim ke onUpdate
            onUpdate(node.id, {
                adapter_config: adapterConfigUpdates as any
            });

            toast.success("Node configuration saved");
            onClose();
        } catch (error) {
            console.error("Save error:", error);
            toast.error("Failed to save node configuration");
        } finally {
            setSaving(false);
        }
    };

    // Execute node - dengan menyimpan response ke node
    const handleExecute = async () => {
        setExecuting(true);
        setExecutionResult(null);

        try {
            const baseUrl = product.api_endpoint?.replace(/\/$/, '') || '';
            const fullPath = endpointPath.startsWith('/') ? endpointPath : `/${endpointPath}`;
            const fullUrl = `${baseUrl}${fullPath}`;

            console.group("🚀 EXECUTE REQUEST");
            console.log("Base URL:", baseUrl);
            console.log("Endpoint Path:", endpointPath);
            console.log("Full URL:", fullUrl);
            console.log("Method:", httpMethod);

            let headers: Record<string, string> = {
                'Content-Type': 'application/json'
            };

            // Parse custom headers dari string JSON
            if (customHeaders && customHeaders.trim() !== "" && customHeaders !== "{}") {
                try {
                    const customHeadersObj = JSON.parse(customHeaders);
                    headers = { ...headers, ...customHeadersObj };
                } catch (e) {
                    console.error("Failed to parse custom headers", e);
                }
            }

            if (product.api_key) {
                headers['X-API-Key'] = product.api_key;
            }

            console.log("Headers:", headers);

            let body = null;

            if (httpMethod !== 'GET' && httpMethod !== 'DELETE') {
                body = JSON.stringify(inputMapping);
            }

            console.log("Request Body:", body);

            const controller = new AbortController();
            const timeoutId = setTimeout(
                () => controller.abort(),
                timeoutSeconds * 1000
            );

            console.log("Sending request...");

            const response = await fetch(fullUrl, {
                method: httpMethod,
                headers,
                body,
                signal: controller.signal,
            });

            clearTimeout(timeoutId);

            console.log("Response Status:", response.status);
            console.log("Response Status Text:", response.statusText);

            const contentType = response.headers.get('content-type');
            console.log("Content Type:", contentType);

            let data;
            if (contentType?.includes('application/json')) {
                data = await response.json();
            } else {
                data = await response.text();
            }

            console.log("Response Data:", data);
            console.groupEnd();

            // SAVE RESPONSE DATA TO NODE untuk placeholder
            if (response.ok && data) {
                console.log("💾 Saving response data to node for placeholders...");

                // Update node dengan response data
                onUpdate(node.id, {
                    last_execution: {
                        data: data,
                        status: response.status,
                        timestamp: new Date().toISOString()
                    }
                } as any);

                toast.success(`Response data saved! You can now use it as placeholder in next nodes.`);
            }

            setExecutionResult({
                success: response.ok,
                status: response.status,
                statusText: response.statusText,
                data: data,
                error: !response.ok
            });

            if (response.ok) {
                toast.success(`Request successful: ${response.status} ${response.statusText}`);
            } else {
                toast.error(`Request failed: ${response.status} ${response.statusText}`);
            }

        } catch (error: any) {
            console.group("❌ EXECUTION ERROR");
            console.error("Error Message:", error?.message);
            console.error("Full Error:", error);
            console.groupEnd();

            setExecutionResult({
                success: false,
                error: true,
                message: error?.message || "Unknown error occurred",
                data: null
            });

            toast.error(error?.message || "Failed to execute request");
        } finally {
            setExecuting(false);
        }
    };
    // Delete node
    const handleDelete = () => {
        if (confirm("Are you sure you want to delete this node?")) {
            onDelete(node.id);
            onClose();
        }
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        onClose();
    };

    const hasPlaceholders = nodePlaceholders.length > 0;
    const hasPreviousNodes = previousNodes.length > 0;

    return (
        <Dialog open={isModalOpen} onOpenChange={handleCloseModal}>
            <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle className="flex items-center justify-between">
                        <span>Configure Node - Step {currentNodeStepOrder}</span>
                        {/* <Button
                            variant="ghost"
                            size="sm"
                            onClick={handleDelete}
                            className="text-red-500 hover:text-red-700"
                        >
                            <Trash2 className="h-4 w-4" />
                        </Button> */}
                    </DialogTitle>
                    <DialogDescription>
                        Configure the endpoint and input mapping for this workflow step
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-6 mt-4">
                    {/* Step Order Info */}
                    <div className="bg-muted/50 p-3 rounded-lg">
                        <Label className="text-xs text-muted-foreground uppercase">Step Order</Label>
                        <p className="text-lg font-semibold mt-1">{currentNodeStepOrder}</p>
                        {isFirstNode && (
                            <p className="text-xs text-muted-foreground mt-1">
                                This is the first node in the workflow
                            </p>
                        )}
                    </div>

                    {/* Endpoint Configuration */}
                    <div className="space-y-4">
                        <Label className="text-sm font-semibold">Endpoint Configuration</Label>

                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label>HTTP Method</Label>
                                <Select value={httpMethod} onValueChange={(v) => setHttpMethod(v as AdapterConfigNode['http_method'])}>
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="GET">GET</SelectItem>
                                        <SelectItem value="POST">POST</SelectItem>
                                        <SelectItem value="PUT">PUT</SelectItem>
                                        <SelectItem value="PATCH">PATCH</SelectItem>
                                        <SelectItem value="DELETE">DELETE</SelectItem>
                                    </SelectContent>
                                </Select>
                            </div>

                            <div className="space-y-2">
                                <Label>Endpoint Path</Label>
                                <Input
                                    placeholder="/api/posts"
                                    value={endpointPath}
                                    onChange={(e) => setEndpointPath(e.target.value)}
                                />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label>Timeout (seconds)</Label>
                                <Input
                                    type="number"
                                    min={1}
                                    max={300}
                                    value={timeoutSeconds}
                                    onChange={(e) => setTimeoutSeconds(parseInt(e.target.value) || 30)}
                                />
                            </div>

                            <div className="space-y-2">
                                <Label>Retry Count</Label>
                                <Input
                                    type="number"
                                    min={0}
                                    max={10}
                                    value={retryCount}
                                    onChange={(e) => setRetryCount(parseInt(e.target.value) || 3)}
                                />
                            </div>
                        </div>

                        {/* URL Preview */}
                        <div className="bg-muted p-2 rounded text-xs">
                            <span className="text-muted-foreground">Full URL: </span>
                            <span className="font-mono">
                                {product.api_endpoint?.replace(/\/$/, '') || ''}
                                {endpointPath.startsWith('/') ? endpointPath : `/${endpointPath}`}
                            </span>
                        </div>
                    </div>

                    {/* Input Mapping */}
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <Label className="text-sm font-semibold">Input Mapping (Request Body)</Label>
                            {!isFirstNode && (
                                <span className="text-xs text-muted-foreground">
                                    {hasPlaceholders
                                        ? `✓ ${nodePlaceholders.length} placeholder(s) available from ${previousNodes.length} previous node(s)`
                                        : hasPreviousNodes
                                            ? `⚠️ Execute previous node first to get placeholders`
                                            : `No previous nodes available`}
                                </span>
                            )}
                        </div>

                        <div className="border rounded-lg p-4 bg-muted/30">
                            <SimpleJsonBuilder
                                value={inputMapping}
                                onChange={setInputMapping}
                                placeholders={nodePlaceholders}
                            />
                        </div>

                        <p className="text-xs text-muted-foreground">
                            {isFirstNode
                                ? "First node doesn't have access to previous node responses"
                                : hasPlaceholders
                                    ? "Click the 📋 button to insert response data from previous nodes"
                                    : hasPreviousNodes
                                        ? "Execute the previous node(s) first to make their response data available as placeholders"
                                        : "Add previous nodes to this workflow to access their responses"}
                        </p>
                    </div>

                    {/* Execute Button */}
                    <div className="space-y-2 pt-2">
                        <Button
                            type="button"
                            onClick={handleExecute}
                            disabled={executing}
                            className="w-full"
                            variant="outline"
                        >
                            <Play className="h-4 w-4 mr-2" />
                            {executing ? "Executing..." : "Test Execute Endpoint"}
                        </Button>
                        <p className="text-xs text-muted-foreground text-center">
                            Test the endpoint with current input mapping
                        </p>
                    </div>

                    {/* Execution Result */}
                    {executionResult && (
                        <div className={`space-y-2 p-3 rounded-lg ${!executionResult.success
                            ? 'bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800'
                            : 'bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800'
                            }`}>
                            <Label className="text-xs font-semibold flex items-center gap-2">
                                <span>Execution Result:</span>
                                {executionResult.success && (
                                    <span className="text-green-600 dark:text-green-400">
                                        {executionResult.status} {executionResult.statusText}
                                    </span>
                                )}
                                {!executionResult.success && executionResult.status && (
                                    <span className="text-red-600 dark:text-red-400">
                                        {executionResult.status} {executionResult.statusText}
                                    </span>
                                )}
                            </Label>
                            <pre className="text-xs overflow-auto max-h-60 p-3 bg-background rounded">
                                {JSON.stringify(executionResult.data || executionResult, null, 2)}
                            </pre>
                        </div>
                    )}
                </div>

                {/* Actions */}
                <div className="flex justify-end gap-2 mt-6 pt-4 border-t">
                    <Button variant="outline" onClick={handleCloseModal}>
                        Cancel
                    </Button>
                    <Button onClick={handleSave} disabled={saving}>
                        <Save className="h-4 w-4 mr-2" />
                        {saving ? "Saving..." : "Save Node"}
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}