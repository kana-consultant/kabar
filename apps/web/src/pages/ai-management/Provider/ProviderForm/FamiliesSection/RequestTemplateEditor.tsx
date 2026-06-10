// pages/admin/ai-management/components/ProviderForm/FamiliesSection/RequestTemplateEditor.tsx

import { Tabs, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Textarea } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Code2, Blocks, Braces } from "lucide-react";
import { JsonBuilder } from "@/pages/ai-management/JsonBuilder";
import { useTemplateEditor } from "@/hooks/useAIManagement/useTemplateEditor";

interface RequestTemplateEditorProps {
    familyId: string;
    value: string;
    onChange: (value: string) => void;
    error?: string;
}

const AVAILABLE_VARIABLES =["{model}", "{prompt}", "{temperature}", "{max_token}", "{system_prompt}"];

const PLACEHOLDER_JSON = `{
  "model": "{model}",
  "messages": [
    {
      "role": "system",
      "content": "{system_prompt}"
    },
    {
      "role": "user",
      "content": "{prompt}"
    }
  ],
  "temperature": {temperature},
  "max_tokens": {max_tokens},
  "stream": true
}`;

export function RequestTemplateEditor({ familyId, value, onChange, error }: RequestTemplateEditorProps) {
    const {
        mode,
        templateString,
        templateObject,
        isValid,
        updateFromObject,
        updateFromString,
        formatJSON,
        setMode,
    } = useTemplateEditor({
        familyId,
        initialTemplate: value,
        onTemplateChange: onChange,
    });

    console.log("JSON TEMPLATE ==", templateObject)

    return (
        <div className="space-y-2">
            <div className="flex items-center justify-between">
                <Label className="text-base font-semibold">
                    Request Template <span className="text-red-500">*</span>
                </Label>
                <Tabs value={mode} onValueChange={(v) => setMode(v as any)} className="w-auto">
                    <TabsList className="h-8">
                        <TabsTrigger value="visual" className="text-xs px-3">
                            <Blocks className="h-3 w-3 mr-1" />
                            Visual Builder
                        </TabsTrigger>
                        <TabsTrigger value="code" className="text-xs px-3">
                            <Code2 className="h-3 w-3 mr-1" />
                            Code Editor
                        </TabsTrigger>
                    </TabsList>
                </Tabs>
            </div>

            {mode === "visual" ? (
                <div className="border rounded-lg p-4">
                    <JsonBuilder
                        value={templateObject}
                        onChange={updateFromObject}
                        availableVariables={AVAILABLE_VARIABLES}
                    />
                </div>
            ) : (
                <div className="space-y-2">
                    <div className="flex items-center justify-between">
                        <Label>JSON Editor</Label>
                        <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={formatJSON}
                            disabled={!isValid}
                            className="text-xs h-7"
                        >
                            <Braces className="h-3 w-3 mr-1" />
                            Format JSON
                        </Button>
                    </div>
                    <Textarea
                        value={templateString}
                        onChange={(e) => updateFromString(e.target.value)}
                        rows={20}
                        className={`font-mono text-sm ${error ? "border-red-500" : ""}`}
                        placeholder={PLACEHOLDER_JSON}
                    />
                    <div className="flex items-center justify-between flex-wrap gap-2">
                        <div className="flex flex-wrap gap-1">
                            {AVAILABLE_VARIABLES.map(variable => (
                                <Badge key={variable} tone="outline" className="text-xs">
                                    {variable}
                                </Badge>
                            ))}
                        </div>
                        {templateString && templateString.trim() !== "" ? (
                            isValid ? (
                                <Badge tone="outline" className="text-xs text-green-600 border-green-600">
                                    ✓ Valid JSON
                                </Badge>
                            ) : (
                                <Badge tone="outline" className="text-xs text-red-600 border-red-600">
                                    ✗ Invalid JSON
                                </Badge>
                            )
                        ) : (
                            <Badge tone="outline" className="text-xs text-gray-400 border-gray-400">
                                Empty
                            </Badge>
                        )}
                    </div>
                    {error && <p className="text-sm text-red-500">{error}</p>}
                </div>
            )}
            <p className="text-xs text-muted-foreground">
                Define the request body template. Use variables that will be replaced at runtime.
            </p>
        </div>
    );
}