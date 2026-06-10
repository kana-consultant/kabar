// pages/admin/ai-management/components/ProviderForm/FamiliesSection/ResponsePaths.tsx

import { useState } from "react";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import type{ Family } from "@/types/provider.types";
import { TestModelDialog } from "@/pages/ai-management/TestModelDialog";
import { Beaker } from "lucide-react";

interface ResponsePathsProps {
    value: Family;
    onChange: (updates: Partial<Family>) => void;
    providerConfig?: {
        base_url: string;
        auth_header: string;
        auth_prefix: string | null;
        default_headers: Record<string, string>;
        display_name?: string;
    };
}

export function ResponsePaths({ value, onChange, providerConfig }: ResponsePathsProps) {
    const [testDialogOpen, setTestDialogOpen] = useState(false);

    const handlePathsSelected = (textPath: string, imagePath?: string) => {
        onChange({
            schema: {
                ...value.schema,
                response_text_path: textPath,
                response_image_path: imagePath || null,
            }
        });
    };

    // Konversi providerConfig ke APIProvider untuk TestModelDialog
    const apiProviderForTest = providerConfig ? {
        base_url: providerConfig.base_url,
        auth_header: providerConfig.auth_header,
        auth_prefix: providerConfig.auth_prefix,
        default_headers: providerConfig.default_headers,
        display_name: providerConfig.display_name || "Provider",
        name: "temp",
        description: null,
        auth_type: "custom",
        is_active: true,
        families: [],
    } : undefined;

    return (
        <>
            <div className="space-y-2">
                <div className="flex items-center justify-between gap-2">
                    <Label className="flex-1">Response Text Path</Label>
                    <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => setTestDialogOpen(true)}
                        disabled={!providerConfig?.base_url}
                    >
                        <Beaker className="h-3 w-3 mr-1" />
                        Test & Auto-Detect
                    </Button>
                </div>
                <Input
                    value={value.schema.response_text_path || ""}
                    onChange={(e) => onChange({ schema : {...value.schema, response_text_path : e.target.value as string } })}
                    placeholder="choices[0].message.content"
                />
                <p className="text-xs text-muted-foreground">
                    JSON path to extract text from response
                </p>
            </div>

            <div className="space-y-2">
                <Label>Response Image Path (Optional)</Label>
                <Input
                    value={value.schema.response_image_path || ""}
                    onChange={(e) => onChange({schema : {...value.schema, response_image_path : e.target.value || null } })}
                    placeholder="data[0].url"
                />
                <p className="text-xs text-muted-foreground">
                    JSON path to extract image URL from response
                </p>
            </div>

            {/* Test Dialog - pass family langsung */}
            {apiProviderForTest && (
                <TestModelDialog
                    open={testDialogOpen}
                    onOpenChange={setTestDialogOpen}
                    family={value}  // ← Langsung pass family
                    providerConfig={apiProviderForTest}
                    onPathsSelected={handlePathsSelected}
                />
            )}
        </>
    );
}