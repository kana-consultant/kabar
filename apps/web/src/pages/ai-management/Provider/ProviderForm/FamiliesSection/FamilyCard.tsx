// pages/admin/ai-management/components/ProviderForm/FamiliesSection/FamilyCard.tsx

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Trash2, Copy, ChevronDown, ChevronUp } from "lucide-react";
import { FamilyBasicInfo } from "./FamilyBasicInfo";
import { ResponsePaths } from "./ResponsePaths";
import { RequestTemplateEditor } from "./RequestTemplateEditor";
import type { Family } from "@/types/provider.types";

interface FamilyCardProps {
    family: Family;
    index: number;
    isOnly: boolean;
    onUpdate: (updates: Partial<Family>) => void;
    onRemove: () => void;
    onDuplicate: () => void;
    errors?: Record<string, string>;
    providerConfig?: {
        base_url: string;
        auth_header: string;
        auth_prefix: string | null;
        default_headers: Record<string, string>;
        display_name?: string;
    };
}

export function FamilyCard({
    family,
    index,
    isOnly,
    onUpdate,
    onRemove,
    onDuplicate,
    errors,
    providerConfig,
}: FamilyCardProps) {
    const [isExpanded, setIsExpanded] = useState(false);

    // Handler untuk update parameter family
    const handleFamilyParamChange = (field: string, value: number | string | undefined) => {
        onUpdate({ [field]: value });
    };

    return (
        <Card className="relative">
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                        <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            onClick={() => setIsExpanded(!isExpanded)}
                        >
                            {isExpanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                        </Button>
                        <CardTitle className="text-lg">
                            Family {index + 1}: {family.display_name || family.name || "New Family"}
                        </CardTitle>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Button type="button" variant="ghost" size="sm" onClick={onDuplicate}>
                            <Copy className="h-4 w-4" />
                        </Button>
                        {!isOnly && (
                            <Button type="button" variant="ghost" size="sm" onClick={onRemove} className="text-red-500 hover:text-red-700">
                                <Trash2 className="h-4 w-4" />
                            </Button>
                        )}
                    </div>
                </div>
            </CardHeader>

            {isExpanded && (
                <CardContent className="space-y-6 pt-0">
                    <FamilyBasicInfo
                        value={family}
                        onChange={onUpdate}
                        errors={errors}
                    />

                    {/* Family Parameters Section */}
                    <div className="space-y-4 p-4 border rounded-lg">
                        <div>
                            <h5 className="font-medium text-sm">Family Parameters</h5>
                            <p className="text-xs text-muted-foreground">
                                Override global parameters for this family
                            </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Max Token */}
                            <div className="space-y-2">
                                <label className="text-sm font-medium">
                                    Max Token
                                    <span className="text-xs text-muted-foreground ml-2">
                                        (Override global)
                                    </span>
                                </label>
                                <input
                                    type="number"
                                    min="1"
                                    max="32000"
                                    value={family.max_token ?? ''}
                                    onChange={(e) => handleFamilyParamChange('max_token', parseInt(e.target.value) || undefined)}
                                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    placeholder="Use global default"
                                />
                                {errors?.[`family_${family.id}_max_token`] && (
                                    <p className="text-xs text-red-500">{errors[`family_${family.id}_max_token`]}</p>
                                )}
                                <p className="text-xs text-muted-foreground">
                                    Maximum tokens for this family (leave empty to use global)
                                </p>
                            </div>

                            {/* Temperature */}
                            <div className="space-y-2">
                                <label className="text-sm font-medium">
                                    Temperature
                                    <span className="text-xs text-muted-foreground ml-2">
                                        (Override global)
                                    </span>
                                </label>
                                <input
                                    type="number"
                                    step="0.1"
                                    min="0"
                                    max="2"
                                    value={family.temperature ?? ''}
                                    onChange={(e) => {
                                        const value = e.target.value;
                                        handleFamilyParamChange('temperature', value === "" ? undefined : parseFloat(value));
                                    }}
                                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    placeholder="Use global default"
                                />
                                {errors?.[`family_${family.id}_temperature`] && (
                                    <p className="text-xs text-red-500">{errors[`family_${family.id}_temperature`]}</p>
                                )}
                                <p className="text-xs text-muted-foreground">
                                    Randomness level for this family (0-2, leave empty to use global)
                                </p>
                            </div>

                            {/* System Prompt */}
                            <div className="md:col-span-2 space-y-2">
                                <label className="text-sm font-medium">
                                    System Prompt
                                    <span className="text-xs text-muted-foreground ml-2">
                                        (Override global)
                                    </span>
                                </label>
                                <textarea
                                    value={family.system_prompt || ''}
                                    onChange={(e) => handleFamilyParamChange('system_prompt', e.target.value)}
                                    rows={3}
                                    className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
                                    placeholder="Use global system prompt"
                                />
                                {errors?.[`family_${family.id}_system_prompt`] && (
                                    <p className="text-xs text-red-500">{errors[`family_${family.id}_system_prompt`]}</p>
                                )}
                                <p className="text-xs text-muted-foreground">
                                    Custom system prompt for this family (overrides global)
                                </p>
                            </div>
                        </div>
                    </div>

                    <RequestTemplateEditor
                        familyId={family.id}
                        value={family.schema?.request_template || "{}"}
                        onChange={(template) => onUpdate({ schema: { ...family.schema, request_template: template } })}
                        error={errors?.[`family_${family.id}_template`]}
                    />

                    <ResponsePaths
                        value={family}
                        onChange={onUpdate}
                        providerConfig={providerConfig}
                    />
                </CardContent>
            )}
        </Card>
    );
}