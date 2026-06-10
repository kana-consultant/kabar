import { Button, Input, Label } from "@kana-consultant/ui-kit";
import type { ModelFormSectionProps } from "../../provider.types";

interface ResponseConfigSectionProps extends ModelFormSectionProps {
    onTestModel: () => void;
    isTestDisabled: boolean;
    getResponseTextPath: () => string;
    getResponseImagePath: () => string;
}

export function ResponseConfigSection({ 
    formData, 
    setFormData, 
    selectedSchema, 
    onTestModel, 
    isTestDisabled,
    getResponseTextPath,
    getResponseImagePath
}: ResponseConfigSectionProps) {
    return (
        <div className="space-y-4">
            <div>
                <h3 className="text-lg font-semibold">Response Configuration</h3>
                <p className="text-sm text-muted-foreground">
                    Define how to extract text and images from API response
                </p>
                {selectedSchema && (
                    <div className="mt-2 p-2 bg-muted/50 rounded-md text-sm">
                        <span className="font-medium">Schema defaults:</span>
                        <br />
                        Text path: {selectedSchema.response_text_path || "Not set"}
                        <br />
                        Image path: {selectedSchema.response_image_path || "Not set"}
                    </div>
                )}
            </div>

            <div className="flex items-center justify-between">
                <Label className="text-base">Response Paths</Label>
                <Button
                    type="button"
                    variant="outline"
                    onClick={onTestModel}
                    disabled={isTestDisabled}
                >
                    Test Model & Auto-Detect
                </Button>
            </div>

            <div className="grid grid-cols-1 gap-4">
                <div className="space-y-2">
                    <Label>Text Response Path</Label>
                    <Input
                        value={getResponseTextPath()}
                        onChange={(e) => setFormData(prev => ({ 
                            ...prev, 
                            response_text_path: e.target.value || null 
                        }))}
                        placeholder="choices[0].message.content"
                    />
                    <p className="text-xs text-muted-foreground">
                        JSONPath to extract text from response. Leave empty to use schema default.
                    </p>
                </div>

                <div className="space-y-2">
                    <Label>Image Response Path</Label>
                    <Input
                        value={getResponseImagePath()}
                        onChange={(e) => setFormData(prev => ({ 
                            ...prev, 
                            response_image_path: e.target.value || null 
                        }))}
                        placeholder="data[0].url"
                    />
                    <p className="text-xs text-muted-foreground">
                        JSONPath to extract image URL from response (optional)
                    </p>
                </div>
            </div>
        </div>
    );
}