// pages/admin/ai-management/components/FamiliesList.tsx (Simplified)

import { useState } from "react";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Copy, Check, Zap, Thermometer, MessagesSquare } from "lucide-react";
import type { Family } from "@/types/provider.types";

interface FamiliesListProps {
    families: Family[];
}

export function FamiliesList({ families }: FamiliesListProps) {
    const [copiedId, setCopiedId] = useState<string | null>(null);

    const handleCopyTemplate = async (template: string, familyId: string) => {
        try {
            await navigator.clipboard.writeText(template);
            setCopiedId(familyId);
            setTimeout(() => setCopiedId(null), 2000);
        } catch (err) {
            console.error("Failed to copy:", err);
        }
    };

    const formatJSON = (jsonString: string) => {
        try {
            const parsed = JSON.parse(jsonString);
            return JSON.stringify(parsed, null, 2);
        } catch {
            return jsonString;
        }
    };

    if (families.length === 0) {
        return (
            <Card>
                <CardContent className="py-8 text-center">
                    <p className="text-muted-foreground">No families configured for this provider.</p>
                </CardContent>
            </Card>
        );
    }

    return (
        <div className="space-y-4">
            {families.map((family) => (
                <Card key={family.id}>
                    <CardHeader>
                        <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-3">
                                <MessagesSquare className="h-5 w-5 text-muted-foreground" />
                                <div>
                                    <CardTitle>{family.display_name}</CardTitle>
                                    <p className="text-sm text-muted-foreground">{family.name}</p>
                                </div>
                            </div>
                            <div className="flex items-center space-x-2">
                                {family.supports_streaming && (
                                    <Badge tone="outline">Streaming</Badge>
                                )}
                                {family.supports_temperature && (
                                    <Badge tone="outline">Temperature</Badge>
                                )}
                            </div>
                        </div>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        {/* Endpoint */}
                        <div>
                            <h4 className="text-sm font-medium mb-1">Endpoint</h4>
                            <code className="text-sm bg-muted px-2 py-1 rounded block">
                                {family.endpoint_path}
                            </code>
                        </div>

                        {/* Response Paths */}
                        <div className="grid grid-cols-2 gap-4">
                            {family.response_text_path && (
                                <div>
                                    <h4 className="text-sm font-medium mb-1">Response Text Path</h4>
                                    <code className="text-xs bg-muted px-2 py-1 rounded">
                                        {family.response_text_path}
                                    </code>
                                </div>
                            )}
                            {family.response_image_path && (
                                <div>
                                    <h4 className="text-sm font-medium mb-1">Response Image Path</h4>
                                    <code className="text-xs bg-muted px-2 py-1 rounded">
                                        {family.response_image_path}
                                    </code>
                                </div>
                            )}
                        </div>

                        {/* Request Template */}
                        <div>
                            <div className="flex items-center justify-between mb-2">
                                <h4 className="text-sm font-medium">Request Template</h4>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => handleCopyTemplate(family.request_template, family.id)}
                                >
                                    {copiedId === family.id ? (
                                        <>
                                            <Check className="h-3 w-3 mr-1" />
                                            Copied
                                        </>
                                    ) : (
                                        <>
                                            <Copy className="h-3 w-3 mr-1" />
                                            Copy
                                        </>
                                    )}
                                </Button>
                            </div>
                            <pre className="text-xs bg-muted p-3 rounded-lg overflow-x-auto">
                                {formatJSON(family.request_template)}
                            </pre>
                        </div>

                        {/* Features */}
                        <div className="flex gap-4 pt-2">
                            <div className="flex items-center space-x-1">
                                <Zap className="h-4 w-4" />
                                <span className="text-sm">
                                    {family.supports_streaming ? "Streaming" : "No Streaming"}
                                </span>
                            </div>
                            <div className="flex items-center space-x-1">
                                <Thermometer className="h-4 w-4" />
                                <span className="text-sm">
                                    {family.supports_temperature ? "Temperature" : "No Temperature"}
                                </span>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}