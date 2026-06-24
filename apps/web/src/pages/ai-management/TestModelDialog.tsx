// pages/admin/ai-management/TestModelDialog.tsx

import { useState, useEffect } from "react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Checkbox } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Loader, CheckCircle, XCircle, Image, Type, Hash, ToggleLeft, Key, MessageSquare, Globe, Sparkles, Thermometer, Sigma, AlignLeft, Info } from "lucide-react";
import { cn } from "@/lib/utils";
import type { APIProvider, Family } from '@/types/provider.types';

interface TestModelDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    family: Family;
    providerConfig: APIProvider;
    onPathsSelected: (textPath: string, imagePath?: string) => void;
    globalConfig?: {
        max_token?: number;
        temperature?: number;
        system_prompt?: string;
    };
}

interface PathNode {
    path: string;
    value: any;
    type: "text" | "image" | "number" | "boolean" | "null";
}

export function TestModelDialog({
    open,
    onOpenChange,
    family,
    providerConfig,
    onPathsSelected,
    globalConfig
}: TestModelDialogProps) {
    const [apiKey, setApiKey] = useState("");
    const [testPrompt, setTestPrompt] = useState("Hello, this is a test message");

    // SEMUA nilai diambil dari family, tidak bisa diedit
    const maxTokenValue = (() => {
        if (family.max_token !== undefined && family.max_token !== null) return family.max_token;
        if (globalConfig?.max_token !== undefined) return globalConfig.max_token;
        return 1000;
    })();

    const temperatureValue = (() => {
        if (family.temperature !== undefined && family.temperature !== null) return family.temperature;
        if (globalConfig?.temperature !== undefined) return globalConfig.temperature;
        return 0.7;
    })();

    const systemPromptValue = (() => {
        if (family.system_prompt !== undefined && family.system_prompt !== null) return family.system_prompt;
        if (globalConfig?.system_prompt !== undefined) return globalConfig.system_prompt;
        return "";
    })();

    const [loading, setLoading] = useState(false);
    const [rawResponse, setRawResponse] = useState<any>(null);
    const [error, setError] = useState<string | null>(null);
    const [selectedTextPath, setSelectedTextPath] = useState<string>("");
    const [selectedImagePath, setSelectedImagePath] = useState<string>("");
    const [leafNodes, setLeafNodes] = useState<PathNode[]>([]);

    // Reset state ketika dialog ditutup
    const handleOpenChange = (open: boolean) => {
        if (!open) {
            setApiKey("");
            setRawResponse(null);
            setError(null);
            setSelectedTextPath("");
            setSelectedImagePath("");
            setLeafNodes([]);
            setTestPrompt("Hello, this is a test message");
        }
        onOpenChange(open);
    };

    // Auto-select paths ketika response diterima
    useEffect(() => {
        if (rawResponse) {
            const nodes = flattenJSON(rawResponse);
            setLeafNodes(nodes);

            // Auto-select first text field
            const firstText = nodes.find(n => n.type === "text");
            if (firstText && !selectedTextPath) {
                setSelectedTextPath(firstText.path);
            }

            // Auto-select first image field
            const firstImage = nodes.find(n => n.type === "image");
            if (firstImage && !selectedImagePath) {
                setSelectedImagePath(firstImage.path);
            }
        }
    }, [rawResponse]);

    const getTemplateString = (template: any): string => {
        if (!template) return "{}";
        if (typeof template === 'string') return template;
        return JSON.stringify(template);
    };

    const replaceVariables = (str: string): string => {
        const familyName = family.name || "";
        return str
            .replace(/\{model\}/g, familyName)
            .replace(/\{prompt\}/g, testPrompt)
            .replace(/\{max_token\}/g, String(maxTokenValue))
            .replace(/\{temperature\}/g, temperatureValue.toString())
            .replace(/\{system_prompt\}/g, systemPromptValue);
    };

    // IMPROVED: Better image detection
    const isImageValue = (value: string): boolean => {
        if (typeof value !== "string") return false;
        
        const imagePatterns = [
            /^https?:\/\/[^\s]+\.(jpg|jpeg|png|gif|webp|svg|bmp|ico|tiff?)/i,
            /^https?:\/\/[^\s]+\/(image|img|photo|picture|avatar|thumbnail|cover)/i,
            /^data:image\//,
            /^\/9j\//, // JPEG base64
            /^iVBOR/,  // PNG base64
            /^R0lGOD/, // GIF base64
            /^data:image\/svg\+xml/,
            /^blob:https?:\/\//,
        ];
        return imagePatterns.some(pattern => pattern.test(value));
    };

    const flattenJSON = (obj: any, prefix = ""): PathNode[] => {
        const nodes: PathNode[] = [];

        if (obj === null || obj === undefined) {
            nodes.push({ path: prefix, value: null, type: "null" });
            return nodes;
        }

        if (Array.isArray(obj)) {
            obj.forEach((item, index) => {
                const arrayPath = prefix ? `${prefix}[${index}]` : `[${index}]`;
                nodes.push(...flattenJSON(item, arrayPath));
            });
            return nodes;
        }

        if (typeof obj === "object") {
            Object.keys(obj).forEach((key) => {
                const newPath = prefix ? `${prefix}.${key}` : key;
                const value = obj[key];

                if (typeof value === "object" && value !== null && !Array.isArray(value)) {
                    nodes.push(...flattenJSON(value, newPath));
                } else if (Array.isArray(value)) {
                    nodes.push(...flattenJSON(value, newPath));
                } else {
                    let type: PathNode["type"] = "text";

                    if (typeof value === "string") {
                        if (isImageValue(value)) {
                            type = "image";
                        }
                    } else if (typeof value === "number") {
                        type = "number";
                    } else if (typeof value === "boolean") {
                        type = "boolean";
                    } else if (value === null) {
                        type = "null";
                    }

                    nodes.push({ path: newPath, value, type });
                }
            });
        }

        return nodes;
    };

    const handlePathToggle = (node: PathNode, checked: boolean) => {
        if (!checked) {
            if (node.type === "image" && selectedImagePath === node.path) {
                setSelectedImagePath("");
            } else if (selectedTextPath === node.path) {
                setSelectedTextPath("");
            }
        } else {
            if (node.type === "image") {
                setSelectedImagePath(node.path);
            } else if (node.type === "text" || node.type === "number" || node.type === "boolean") {
                setSelectedTextPath(node.path);
            }
        }
    };

    const handleConfirm = () => {
        if (!selectedTextPath) return;
        onPathsSelected(selectedTextPath, selectedImagePath || undefined);
        handleOpenChange(false);
    };

    const getValuePreview = (value: any): string => {
        if (value === null) return "null";
        if (value === undefined) return "undefined";
        if (typeof value === "string") {
            if (isImageValue(value)) {
                if (value.length > 60) return value.substring(0, 60) + "...";
                return value;
            }
            if (value.length > 100) return value.substring(0, 100) + "...";
            return value;
        }
        return String(value);
    };

    const getTypeBadge = (type: PathNode["type"]) => {
        switch (type) {
            case "image":
                return <Badge tone="info" className="text-[10px] h-4 px-1.5 shrink-0"><Image className="h-2.5 w-2.5 mr-0.5" />IMAGE</Badge>;
            case "number":
                return <Badge tone="outline" className="text-[10px] h-4 px-1.5 shrink-0"><Hash className="h-2.5 w-2.5 mr-0.5" />NUM</Badge>;
            case "boolean":
                return <Badge tone="outline" className="text-[10px] h-4 px-1.5 shrink-0"><ToggleLeft className="h-2.5 w-2.5 mr-0.5" />BOOL</Badge>;
            case "null":
                return <Badge tone="outline" className="text-[10px] h-4 px-1.5 shrink-0">NULL</Badge>;
            default:
                return <Badge tone="info" className="text-[10px] h-4 px-1.5 shrink-0"><Type className="h-2.5 w-2.5 mr-0.5" />TEXT</Badge>;
        }
    };

    const getPreviewEndpoint = () => {
        if (!providerConfig) return "Not selected";
        const endpoint = family.schema?.endpoint_path || "";
        const baseUrl = providerConfig.base_url || "";
        const familyName = family.name || "{model}";

        let previewEndpoint = endpoint.replace(/\{model\}/g, familyName);
        let previewBaseUrl = baseUrl.replace(/\{model\}/g, familyName);

        return `${previewBaseUrl}${previewEndpoint}`;
    };

    const copyToClipboard = () => {
        if (rawResponse) {
            navigator.clipboard.writeText(JSON.stringify(rawResponse, null, 2));
        }
    };

    const handleTest = async () => {
        if (!apiKey.trim()) {
            setError("API Key is required");
            return;
        }

        if (!family.schema?.request_template) {
            setError("Request template is required. Please setup request template first.");
            return;
        }

        setLoading(true);
        setError(null);
        setRawResponse(null);
        setSelectedTextPath("");
        setSelectedImagePath("");

        try {
            const templateString = getTemplateString(family.schema.request_template);
            const requestBodyString = replaceVariables(templateString);
            const parsedBody = JSON.parse(requestBodyString);

            let defaultHeaders: Record<string, string> = {};
            if (providerConfig?.default_headers) {
                if (typeof providerConfig.default_headers === "string") {
                    try {
                        defaultHeaders = JSON.parse(providerConfig.default_headers);
                    } catch (err) {
                        defaultHeaders = {};
                    }
                } else {
                    defaultHeaders = providerConfig.default_headers;
                }
            }

            const headers: Record<string, string> = {
                "Content-Type": "application/json",
                ...defaultHeaders,
            };

            if (providerConfig?.auth_header && apiKey) {
                const authValue = providerConfig.auth_prefix
                    ? `${providerConfig.auth_prefix} ${apiKey}`.trim()
                    : apiKey;
                headers[providerConfig.auth_header] = authValue;
            }

            let endpoint = family.schema.endpoint_path || "";
            let baseUrl = providerConfig?.base_url || "";

            endpoint = replaceVariables(endpoint);
            baseUrl = replaceVariables(baseUrl);

            const url = `${baseUrl}${endpoint}`;

            const response = await fetch(url, {
                method: "POST",
                headers,
                body: JSON.stringify(parsedBody),
            });

            const responseText = await response.text();
            let data;

            try {
                data = JSON.parse(responseText);
            } catch (err) {
                data = responseText;
            }

            if (!response.ok) {
                throw new Error(
                    data?.error?.message ||
                    data?.error ||
                    `HTTP ${response.status}: Request failed`
                );
            }

            setRawResponse(data);
        } catch (err: any) {
            setError(
                err?.message ||
                "Failed to test. Please check your configuration."
            );
        } finally {
            setLoading(false);
        }
    };

    // Cek sumber nilai
    const getValueSource = (value: any, familyValue: any, globalValue: any) => {
        if (familyValue !== undefined && familyValue !== null) return "family";
        if (globalValue !== undefined) return "global";
        return "default";
    };

    const maxTokenSource = getValueSource(maxTokenValue, family.max_token, globalConfig?.max_token);
    const temperatureSource = getValueSource(temperatureValue, family.temperature, globalConfig?.temperature);
    const systemPromptSource = getValueSource(systemPromptValue, family.system_prompt, globalConfig?.system_prompt);

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogContent className="max-w-6xl max-h-[90vh] overflow-hidden flex flex-col p-0">
                <DialogHeader className="px-6 pt-6 pb-4 border-b">
                    <DialogTitle className="flex items-center gap-2">
                        <Sparkles className="h-5 w-5 text-purple-500" />
                        Test {family.display_name || family.name} & Auto-Detect Response Paths
                    </DialogTitle>
                    <DialogDescription>
                        Send a test request to see the actual API response and select which fields to extract
                    </DialogDescription>
                </DialogHeader>

                <div className="flex-1 overflow-y-auto px-6 py-4 space-y-5">
                    {/* Test Configuration */}
                    <div className="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-950/20 dark:to-indigo-950/20 rounded-xl p-4 border">
                        <h4 className="text-sm font-semibold mb-3 flex items-center gap-2">
                            <Key className="h-4 w-4 text-blue-500" />
                            Test Configuration
                        </h4>
                        <div className="grid grid-cols-2 gap-4">
                            <div className="space-y-2">
                                <Label className="text-xs font-medium">API Key *</Label>
                                <Input
                                    type="password"
                                    value={apiKey}
                                    onChange={(e) => setApiKey(e.target.value)}
                                    placeholder="Enter temporary API key"
                                    className="bg-white dark:bg-gray-900"
                                />
                                <p className="text-xs text-muted-foreground">
                                    ⚠️ This key is only used for testing and will NOT be saved
                                </p>
                            </div>
                            <div className="space-y-2">
                                <Label className="text-xs font-medium">Test Prompt</Label>
                                <Input
                                    value={testPrompt}
                                    onChange={(e) => setTestPrompt(e.target.value)}
                                    placeholder="Enter test prompt"
                                    className="bg-white dark:bg-gray-900"
                                />
                            </div>
                        </div>
                    </div>

                    {/* AI Parameters Section - READ ONLY FROM FAMILY */}
                    <div className="bg-gradient-to-r from-purple-50 to-pink-50 dark:from-purple-950/20 dark:to-pink-950/20 rounded-xl p-4 border">
                        <h4 className="text-sm font-semibold mb-3 flex items-center gap-2">
                            <Sparkles className="h-4 w-4 text-purple-500" />
                            Family Parameters (Read Only)
                            <Badge tone="warning" className="text-[10px]">
                                From Family Configuration
                            </Badge>
                        </h4>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Max Token - Read Only */}
                            <div className="space-y-2">
                                <Label className="text-xs font-medium flex items-center gap-1">
                                    <Sigma className="h-3 w-3" />
                                    Max Token
                                    <Badge tone="info" className="text-[9px] ml-1">
                                        {maxTokenSource === "family" ? "From Family" : maxTokenSource === "global" ? "From Global" : "Default"}
                                    </Badge>
                                </Label>
                                <div className="relative">
                                    <Input
                                        type="number"
                                        value={maxTokenValue}
                                        disabled
                                        className="w-full bg-gray-100 dark:bg-gray-800 cursor-not-allowed opacity-75"
                                    />
                                </div>
                                <p className="text-xs text-muted-foreground">
                                    Maximum tokens to generate
                                    {maxTokenSource === "family" && family.max_token && ` (Value: ${family.max_token})`}
                                </p>
                            </div>

                            {/* Temperature - Read Only */}
                            <div className="space-y-2">
                                <Label className="text-xs font-medium flex items-center gap-1">
                                    <Thermometer className="h-3 w-3" />
                                    Temperature
                                    <Badge tone="info" className="text-[9px] ml-1">
                                        {temperatureSource === "family" ? "From Family" : temperatureSource === "global" ? "From Global" : "Default"}
                                    </Badge>
                                </Label>
                                <div className="relative">
                                    <Input
                                        type="number"
                                        value={temperatureValue}
                                        disabled
                                        className="w-full bg-gray-100 dark:bg-gray-800 cursor-not-allowed opacity-75"
                                        step={0.1}
                                    />
                                    <Info className="absolute right-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                                </div>
                                <p className="text-xs text-muted-foreground">
                                    Randomness level (0-2)
                                    {temperatureSource === "family" && family.temperature && ` (Value: ${family.temperature})`}
                                </p>
                            </div>

                            {/* System Prompt - Read Only */}
                            <div className="md:col-span-2 space-y-2">
                                <Label className="text-xs font-medium flex items-center gap-1">
                                    <AlignLeft className="h-3 w-3" />
                                    System Prompt
                                    <Badge tone="info" className="text-[9px] ml-1">
                                        {systemPromptSource === "family" ? "From Family" : systemPromptSource === "global" ? "From Global" : "Default"}
                                    </Badge>
                                </Label>
                                <textarea
                                    value={systemPromptValue}
                                    disabled
                                    rows={2}
                                    className="w-full px-3 py-2 border rounded-md bg-gray-100 dark:bg-gray-800 cursor-not-allowed opacity-75 font-mono text-sm"
                                    placeholder="No system prompt set"
                                />
                                <p className="text-xs text-muted-foreground">
                                    System instructions for the AI model (read only from family configuration)
                                </p>
                            </div>
                        </div>
                    </div>

                    {/* Request Configuration Summary */}
                    <div className="bg-muted/30 rounded-xl p-4 border">
                        <h4 className="text-sm font-semibold mb-3 flex items-center gap-2">
                            <Globe className="h-4 w-4 text-green-500" />
                            Request Configuration
                        </h4>
                        <div className="space-y-2 text-sm">
                            <div className="flex items-start gap-2">
                                <span className="font-medium w-20 shrink-0">Provider:</span>
                                <span>{providerConfig?.display_name || "Not selected"}</span>
                            </div>
                            <div className="flex items-start gap-2">
                                <span className="font-medium w-20 shrink-0">Family:</span>
                                <code className="text-xs bg-muted px-1.5 py-0.5 rounded font-mono">
                                    {family.display_name || family.name} ({family.name})
                                </code>
                            </div>
                            <div className="flex items-start gap-2">
                                <span className="font-medium w-20 shrink-0">Endpoint:</span>
                                <code className="text-xs text-green-600 break-all font-mono bg-green-50 dark:bg-green-950/30 px-1.5 py-0.5 rounded">
                                    {getPreviewEndpoint()}
                                </code>
                            </div>
                        </div>
                    </div>

                    {/* Send Button */}
                    <Button
                        onClick={handleTest}
                        disabled={!apiKey.trim() || loading || !providerConfig?.base_url}
                        className="w-full"
                    >
                        {loading ? (
                            <>
                                <Loader className="h-4 w-4 mr-2 animate-spin" />
                                Sending request...
                            </>
                        ) : (
                            <>
                                <Sparkles className="h-4 w-4 mr-2" />
                                Send Test Request
                            </>
                        )}
                    </Button>

                    {/* Error */}
                    {error && (
                        <div className="p-4 rounded-xl bg-red-50 dark:bg-red-500/10 border border-red-200 flex items-start space-x-3">
                            <XCircle className="h-5 w-5 text-red-500 mt-0.5 shrink-0" />
                            <div className="flex-1">
                                <div className="font-semibold text-sm text-red-700">Test Failed</div>
                                <p className="text-sm text-red-600 mt-0.5">{error}</p>
                            </div>
                        </div>
                    )}

                    {/* Response */}
                    {rawResponse && !error && (
                        <div className="space-y-4">
                            <div className="flex items-center gap-2">
                                <div className="h-0.5 flex-1 bg-gradient-to-r from-green-500 to-emerald-500 rounded-full" />
                                <Badge tone="success" className="gap-1">
                                    <CheckCircle className="h-3 w-3" />
                                    Response Received
                                </Badge>
                                <div className="h-0.5 flex-1 bg-gradient-to-r from-emerald-500 to-green-500 rounded-full" />
                            </div>

                            <div className="flex flex-col lg:grid lg:grid-cols-2 gap-5">
                                {/* Raw Response */}
                                <div className="space-y-2">
                                    <Label className="text-xs font-semibold uppercase flex items-center justify-between">
                                        <span>Raw JSON Response</span>
                                        <Button variant="ghost" size="sm" className="h-6 px-2 text-xs" onClick={copyToClipboard}>
                                            Copy
                                        </Button>
                                    </Label>
                                    <div className="bg-slate-900 text-slate-100 p-3 rounded-xl max-h-[420px] overflow-y-auto">
                                        <pre className="text-xs font-mono whitespace-pre-wrap break-words">
                                            {JSON.stringify(rawResponse, null, 2)}
                                        </pre>
                                    </div>
                                </div>

                                {/* Path Picker */}
                                <div className="space-y-2">
                                    <Label className="text-xs font-semibold uppercase block">
                                        Select Response Paths
                                        <span className="text-muted-foreground ml-2">({leafNodes.length} fields)</span>
                                        {leafNodes.length > 0 && (
                                            <span className="text-xs font-normal text-green-600 ml-2">
                                                ✨ Auto-selected first text & image
                                            </span>
                                        )}
                                    </Label>
                                    <div className="border rounded-xl max-h-[420px] overflow-y-auto bg-white dark:bg-gray-950">
                                        <div className="divide-y">
                                            {leafNodes.map((node, index) => {
                                                const isTextSelected = node.type !== "image" && selectedTextPath === node.path;
                                                const isImageSelected = node.type === "image" && selectedImagePath === node.path;
                                                const isSelected = isTextSelected || isImageSelected;

                                                return (
                                                    <div
                                                        key={index}
                                                        className={cn(
                                                            "flex items-start space-x-3 p-3 transition-all cursor-pointer hover:bg-muted/50",
                                                            isSelected ? "bg-blue-50 dark:bg-blue-500/10 border-l-4 border-l-blue-500" : ""
                                                        )}
                                                        onClick={() => handlePathToggle(node, !isSelected)}
                                                    >
                                                        <Checkbox
                                                            checked={isSelected}
                                                            onCheckedChange={(checked) => handlePathToggle(node, checked as boolean)}
                                                            className="mt-0.5 shrink-0"
                                                        />
                                                        <div className="flex-1 min-w-0">
                                                            <code className="text-xs font-mono font-semibold break-all">
                                                                {node.path}
                                                            </code>
                                                            <div className="mt-1 flex items-center gap-2">
                                                                {getTypeBadge(node.type)}
                                                                {node.type === "image" && (
                                                                    <span className="text-[10px] text-blue-600">
                                                                        🖼️ Image detected
                                                                    </span>
                                                                )}
                                                            </div>
                                                            <p className="text-xs text-muted-foreground mt-1 truncate font-mono">
                                                                {getValuePreview(node.value)}
                                                            </p>
                                                        </div>
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            {/* Selected Summary */}
                            {(selectedTextPath || selectedImagePath) && (
                                <div className="space-y-3 border-t pt-4">
                                    <Label className="text-xs font-semibold uppercase">Selected Paths Summary</Label>
                                    <div className="grid gap-2">
                                        {selectedTextPath && (
                                            <div className="p-3 rounded-xl border border-green-200 bg-green-50 dark:bg-green-950/20">
                                                <div className="flex items-center gap-2 mb-1">
                                                    <Badge tone="success" className="gap-1">
                                                        <Type className="h-3 w-3" />
                                                        TEXT
                                                    </Badge>
                                                    <span className="text-xs text-muted-foreground">Path:</span>
                                                </div>
                                                <code className="text-sm font-mono text-green-700 dark:text-green-400 block break-all">
                                                    {selectedTextPath}
                                                </code>
                                                <div className="mt-2 text-xs text-muted-foreground">
                                                    <span className="font-medium">Preview value:</span>{' '}
                                                    <span className="font-mono bg-white/50 dark:bg-black/20 px-2 py-0.5 rounded">
                                                        {getValuePreview(
                                                            leafNodes.find(n => n.path === selectedTextPath)?.value
                                                        )}
                                                    </span>
                                                </div>
                                            </div>
                                        )}
                                        {selectedImagePath && (
                                            <div className="p-3 rounded-xl border border-blue-200 bg-blue-50 dark:bg-blue-950/20">
                                                <div className="flex items-center gap-2 mb-1">
                                                    <Badge tone="info" className="gap-1">
                                                        <Image className="h-3 w-3" />
                                                        IMAGE
                                                    </Badge>
                                                    <span className="text-xs text-muted-foreground">Path:</span>
                                                </div>
                                                <code className="text-sm font-mono text-blue-700 dark:text-blue-400 block break-all">
                                                    {selectedImagePath}
                                                </code>
                                                <div className="mt-2 text-xs text-muted-foreground">
                                                    <span className="font-medium">Preview value:</span>{' '}
                                                    <span className="font-mono bg-white/50 dark:bg-black/20 px-2 py-0.5 rounded truncate block">
                                                        {getValuePreview(
                                                            leafNodes.find(n => n.path === selectedImagePath)?.value
                                                        )}
                                                    </span>
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                </div>

                {/* Footer */}
                <div className="border-t px-6 py-4 bg-muted/20 flex justify-between items-center">
                    <p className="text-xs text-muted-foreground">
                        {rawResponse
                            ? "✓ Select one text field and optionally one image field (auto-selected for you)"
                            : "🔑 Enter API key and click test to see response structure"}
                    </p>
                    <div className="flex gap-2">
                        <Button variant="outline" onClick={() => handleOpenChange(false)}>
                            Cancel
                        </Button>
                        <Button
                            onClick={handleConfirm}
                            disabled={!selectedTextPath || !rawResponse}
                        >
                            <CheckCircle className="h-4 w-4 mr-2" />
                            Confirm & Fill Paths
                        </Button>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}