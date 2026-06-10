import { Input, Label } from "@kana-consultant/ui-kit";
import type { ModelFormSectionProps } from "../../provider.types";

export function DefaultParamsSection({ formData, setFormData }: ModelFormSectionProps) {
    return (
        <div className="space-y-4">
            <h3 className="text-lg font-semibold">Default Parameters</h3>

            <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                    <Label>Max Tokens</Label>
                    <Input
                        type="number"
                        value={formData.max_tokens ?? ""}
                        onChange={(e) => setFormData(prev => ({ 
                            ...prev, 
                            max_tokens: e.target.value ? parseInt(e.target.value) : null 
                        }))}
                        placeholder="4096"
                    />
                    <p className="text-xs text-muted-foreground">
                        Maximum tokens per request
                    </p>
                </div>

                <div className="space-y-2">
                    <Label>Temperature</Label>
                    <Input
                        type="number"
                        min={0}
                        max={2}
                        step={0.1}
                        value={formData.temperature ?? ""}
                        onChange={(e) => setFormData(prev => ({ 
                            ...prev, 
                            temperature: e.target.value ? parseFloat(e.target.value) : null 
                        }))}
                        placeholder="0.7"
                    />
                    <p className="text-xs text-muted-foreground">
                        0 = Deterministic, 1 = Creative, 2 = Very creative
                    </p>
                </div>
            </div>
        </div>
    );
}