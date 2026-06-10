// pages/admin/ai-management/components/ProviderForm/BasicInfoSection.tsx

import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Textarea } from "@kana-consultant/ui-kit";
import type { ProviderFormData } from "@/types/provider.types";

interface BasicInfoSectionProps {
    value: ProviderFormData;
    onChange: (updates: Partial<ProviderFormData>) => void;
    errors?: Record<string, string>;
}

export function BasicInfoSection({ value, onChange, errors }: BasicInfoSectionProps) {
    return (
        <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                    <Label htmlFor="name">
                        Provider Name <span className="text-red-500">*</span>
                    </Label>
                    <Input
                        id="name"
                        value={value.name}
                        onChange={(e) => onChange({ name: e.target.value })}
                        placeholder="openrouter"
                        className={errors?.name ? "border-red-500" : ""}
                    />
                    {errors?.name && (
                        <p className="text-sm text-red-500">{errors.name}</p>
                    )}
                    <p className="text-xs text-muted-foreground">
                        Unique identifier, URL-friendly (e.g., "openrouter", "anthropic")
                    </p>
                </div>

                <div className="space-y-2">
                    <Label htmlFor="display_name">
                        Display Name <span className="text-red-500">*</span>
                    </Label>
                    <Input
                        id="display_name"
                        value={value.display_name}
                        onChange={(e) => onChange({ display_name: e.target.value })}
                        placeholder="OpenRouter"
                        className={errors?.display_name ? "border-red-500" : ""}
                    />
                    {errors?.display_name && (
                        <p className="text-sm text-red-500">{errors.display_name}</p>
                    )}
                </div>
            </div>

            <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                    id="description"
                    value={value.description || ""}
                    onChange={(e) => onChange({ description: e.target.value || null })}
                    placeholder="Unified API for multiple AI models"
                    rows={3}
                />
                <p className="text-xs text-muted-foreground">
                    Optional description of this provider
                </p>
            </div>
        </div>
    );
}