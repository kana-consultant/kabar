// src/pages/products/WorkflowBuilder/NodePanel.tsx
import { useState, useEffect, useMemo } from "react";
import { X, Trash2 } from "lucide-react";
import { Button, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@kana-consultant/ui-kit";
import type { Node } from "reactflow";
import { SimpleJsonBuilder } from "../SimpleJsonBuilder";

interface NodePanelProps {
    node: Node;
    adapterConfigs: any[];
    allNodes: Node[];
    onClose: () => void;
    onDelete: (nodeId: string) => void;
    onUpdate: (nodeId: string, updates: Record<string, any>) => void;
}

interface PlaceholderItem {
    value: string;
    label: string;
    group?: string;
}

// ==================== DEFAULT PLACEHOLDERS (BUILT-IN) ====================
const DEFAULT_PLACEHOLDERS: PlaceholderItem[] = [
    { value: "{id}", label: "Unique ID", group: "Basic" },
    { value: "{title}", label: "Article Title", group: "Article" },
    { value: "{slug}", label: "Article Slug", group: "Article" },
    { value: "{tags}", label: "Article Tags (array)", group: "Article" },
    { value: "{keywords}", label: "Article Keywords (array)", group: "Article" },
    { value: "{topic}", label: "Article Topic", group: "Article" },
    { value: "{timestamp}", label: "Current Timestamp", group: "Utility" },
    { value: "{content}", label: "Article Content (HTML)", group: "Content" },
    { value: "{content_text}", label: "Article Content (Plain Text)", group: "Content" },
    { value: "{content_with_image}", label: "Article Content + Image after H1 (HTML)", group: "Content" },
    { value: "{excerpt}", label: "Article Excerpt", group: "Content" },
    { value: "{image_url}", label: "Image URL (Plain)", group: "Image" },
    { value: "{image_content_html}", label: "Image HTML Tag", group: "Image" },
    { value: "{scheduled_for}", label: "Scheduled Time", group: "Schedule" },
    { value: "{meta_title}", label: "Meta Title", group: "SEO" },
    { value: "{meta_description}", label: "Meta Description", group: "SEO" },
    { value: "{meta_keywords}", label: "Meta Keywords (JSON)", group: "SEO" },
    { value: "{og_title}", label: "OG Title", group: "Social Media" },
    { value: "{og_description}", label: "OG Description", group: "Social Media" },
    { value: "{og_image}", label: "OG Image URL", group: "Social Media" },
    { value: "{twitter_title}", label: "Twitter Title", group: "Social Media" },
    { value: "{twitter_description}", label: "Twitter Description", group: "Social Media" },
    { value: "{twitter_image}", label: "Twitter Image URL", group: "Social Media" },
    { value: "{sitemap_priority}", label: "Sitemap Priority", group: "SEO" },
    { value: "{sitemap_changefreq}", label: "Sitemap Change Frequency", group: "SEO" },
];

export function NodePanel({ 
    node, 
    adapterConfigs, 
    allNodes,
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
    
    const currentNodeStepOrder = node.data.workflowNode?.stepOrder;
    const isFirstNode = currentNodeStepOrder === 1;

    useEffect(() => {
        setSelectedAdapterId(node.data.workflowNode?.adapterConfigId || "");
        setInputMapping(node.data.workflowNode?.inputMapping || {});
    }, [node.id, node.data.workflowNode]);

    const selectedAdapter = adapterConfigs.find((c) => c.id === selectedAdapterId);

    // ==================== GENERATE PLACEHOLDERS FROM PREVIOUS NODES ====================
    
    // Dapatkan node-node sebelumnya (node yang terhubung ke node ini)
    const previousNodes = useMemo(() => {
        const previousIds = node.data.workflowNode?.previousNodeIds || [];
        return allNodes.filter(n => previousIds.includes(n.id));
    }, [allNodes, node.data.workflowNode]);

    // Generate placeholder dari response node sebelumnya
    const generateNodePlaceholders = (prevNode: Node): PlaceholderItem[] => {
        const nodeName = prevNode.data.workflowNode?.name || prevNode.id;
        const nodeId = prevNode.id;
        
        // Ambil response mapping dari node sebelumnya
        const responseMapping = prevNode.data.adapterConfig?.responseMapping || {};
        const responseExample = prevNode.data.workflowNode?.responseExample;
        
        const placeholders: PlaceholderItem[] = [];
        
        // Function recursive untuk generate nested fields
        const generateNestedPlaceholders = (obj: any, currentPath: string) => {
            if (!obj || typeof obj !== "object") return;
            
            Object.entries(obj).forEach(([key, value]) => {
                const fullPath = currentPath ? `${currentPath}.${key}` : key;
                const placeholderValue = `{{node_${nodeId}.response.${fullPath}}}`;
                
                if (value && typeof value === "object" && !Array.isArray(value)) {
                    generateNestedPlaceholders(value, fullPath);
                } else {
                    let label = fullPath;
                    if (typeof value === 'string') label = `${fullPath} (string)`;
                    else if (typeof value === 'number') label = `${fullPath} (number)`;
                    else if (typeof value === 'boolean') label = `${fullPath} (boolean)`;
                    else if (Array.isArray(value)) label = `${fullPath} (array)`;
                    
                    placeholders.push({
                        value: placeholderValue,
                        label: label,
                        group: `🔗 ${nodeName}`
                    });
                }
            });
        };
        
        // Generate dari responseMapping jika ada
        if (responseMapping && typeof responseMapping === 'object' && Object.keys(responseMapping).length > 0) {
            generateNestedPlaceholders(responseMapping, "");
        } 
        else if (responseExample && typeof responseExample === "object") {
            generateNestedPlaceholders(responseExample, "");
        }
        
        return placeholders;
    };

    // Generate placeholder dari node sebelumnya (HANYA untuk node > 1)
    const nodePlaceholders = useMemo(() => {
        // Jika node pertama, return array kosong
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

    // Gabungkan placeholder:
    // - Node 1: hanya DEFAULT_PLACEHOLDERS
    // - Node > 1: DEFAULT_PLACEHOLDERS + placeholder dari node sebelumnya
    const allPlaceholders = useMemo(() => {
        if (isFirstNode) {
            // Node pertama: hanya default placeholders
            return DEFAULT_PLACEHOLDERS;
        }
        
        // Node selanjutnya: default placeholders + placeholder dari node sebelumnya
        return [...DEFAULT_PLACEHOLDERS, ...nodePlaceholders];
    }, [nodePlaceholders, isFirstNode]);

    const handleSave = async () => {
        setSaving(true);
        try {
            onUpdate(node.id, {
                inputMapping: inputMapping,
                adapterConfigId: selectedAdapterId,
            });
        } finally {
            setSaving(false);
        }
    };

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

            <div className="space-y-4 flex-1 overflow-y-auto">
                {/* Step Order */}
                <div>
                    <Label className="text-xs">Step Order</Label>
                    <p className="text-sm mt-1">{currentNodeStepOrder}</p>
                </div>

                {/* Adapter Config Selector */}
                <div className="space-y-2">
                    <Label className="text-xs">Adapter Config</Label>
                    <Select
                        value={selectedAdapterId}
                        onValueChange={(value) => setSelectedAdapterId(value)}
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

                {/* Adapter Info */}
                {selectedAdapter && (
                    <div className="text-xs text-muted-foreground space-y-1 bg-muted p-2 rounded">
                        <p><strong>Endpoint:</strong> {selectedAdapter.endpointPath}</p>
                        <p><strong>Method:</strong> {selectedAdapter.httpMethod}</p>
                        <p><strong>Timeout:</strong> {selectedAdapter.timeoutSeconds || 30}s</p>
                        <p><strong>Retry:</strong> {selectedAdapter.retryCount || 3}x</p>
                    </div>
                )}

                {/* Input Mapping - Using SimpleJsonBuilder */}
                <div className="space-y-2">
                    <Label className="text-xs">Input Mapping</Label>
                    <div className="border rounded-lg p-4 bg-muted/30">
                        <SimpleJsonBuilder
                            value={inputMapping}
                            onChange={setInputMapping}
                            placeholders={allPlaceholders}
                        />
                    </div>
                    <p className="text-xs text-muted-foreground">
                        {isFirstNode 
                            ? "Klik tombol 📋 untuk memilih placeholder dari daftar bawaan"
                            : `Klik tombol 📋 untuk memilih placeholder dari node sebelumnya atau daftar bawaan`
                        }
                    </p>
                </div>
            </div>

            <Button onClick={handleSave} disabled={saving} className="mt-4">
                {saving ? "Saving..." : "Save"}
            </Button>
        </div>
    );
}