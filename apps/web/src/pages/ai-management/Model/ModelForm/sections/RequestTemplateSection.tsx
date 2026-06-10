// components/ModelForm/sections/RequestTemplateSection.tsx

import { Tabs, TabsContent, TabsList, TabsTrigger, Button, Label, Textarea } from "@kana-consultant/ui-kit";
import { JsonBuilder } from "@/pages/ai-management/JsonBuilder";
import { Code2, Blocks } from "lucide-react";
import { useTemplateEditor } from "@/hooks/useModelManagement/useTemplateEditor";
import type { RequestTemplateSectionProps } from "../types";

export function RequestTemplateSection({ 
    formData, 
    setFormData, 
    selectedSchema 
}: RequestTemplateSectionProps) {
    const {
        mode,
        templateString,
        templateObject,
        updateFromObject,
        updateFromString,
        formatJSON,
        setMode,
    } = useTemplateEditor({
        initialTemplate: typeof formData.request_template === 'object' && formData.request_template !== null
            ? formData.request_template
            : {},
        onTemplateChange: (templateString) => {
            setFormData((prev : any) => ({ ...prev, request_template: templateString }));
        }
    });

    const handleObjectChange = (jsonObject: Record<string, any>) => {
        updateFromObject(jsonObject);
    };

    const getTemplateObject = (): Record<string, any> => {
        if (formData.request_template && typeof formData.request_template === 'object') {
            return formData.request_template as Record<string, any>;
        }
        if (templateObject && Object.keys(templateObject).length > 0) {
            return templateObject;
        }
        return {};
    };

    return (
        <div className="space-y-4">
            <div>
                <h3 className="text-lg font-semibold">Request Template</h3>
                <p className="text-sm text-muted-foreground">
                    Define the request body structure. Use {"{variable}"} placeholders that will be replaced at runtime.
                </p>
                {selectedSchema && (
                    <div className="mt-2 p-2 bg-muted/50 rounded-md text-sm">
                        <span className="font-medium">Inherited from schema:</span> {selectedSchema.name}
                        <br />
                        <span className="font-medium">Endpoint:</span> {selectedSchema.endpoint_path}
                    </div>
                )}
            </div>

            <Tabs value={mode} onValueChange={(v) => setMode(v as "visual" | "code")}>
                <TabsList className="w-full">
                    <TabsTrigger value="visual" className="flex-1">
                        <Blocks className="h-4 w-4 mr-2" />
                        Visual Builder
                    </TabsTrigger>
                    <TabsTrigger value="code" className="flex-1">
                        <Code2 className="h-4 w-4 mr-2" />
                        Code Editor
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="visual" className="mt-2">
                    <JsonBuilder
                        value={getTemplateObject()}
                        onChange={handleObjectChange}
                        availableVariables={["{model}", "{prompt}", "{system_prompt}", "{temperature}", "{max_tokens}"]}
                    />
                </TabsContent>

                <TabsContent value="code" className="mt-2">
                    <div className="space-y-2">
                        <div className="flex items-center justify-between">
                            <Label>Request Template JSON</Label>
                            <Button
                                type="button"
                                variant="outline"
                                size="sm"
                                onClick={formatJSON}
                                className="text-xs h-7"
                            >
                                Format JSON
                            </Button>
                        </div>
                        <Textarea
                            value={templateString}
                            onChange={(e) => updateFromString(e.target.value)}
                            rows={10}
                            className="font-mono text-sm"
                            placeholder={`{
  "model": "{model}",
  "max_tokens": "{max_tokens}",
  "temperature": "{temperature}",
  "system": "{system_prompt}",
  "messages": "{prompt}"
}`}
                        />
                        <p className="text-xs text-muted-foreground">
                            Available variables: {"{model}"}, {"{prompt}"}, {"{system_prompt}"}, {"{temperature}"}, {"{max_tokens}"}
                        </p>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    );
}