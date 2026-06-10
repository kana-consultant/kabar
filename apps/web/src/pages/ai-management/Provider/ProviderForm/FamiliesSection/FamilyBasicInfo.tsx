// pages/admin/ai-management/components/ProviderForm/FamiliesSection/FamilyBasicInfo.tsx

import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { type Family } from "@/types/provider.types";

interface FamilyBasicInfoProps {
    value: Family;
    onChange: (updates: Partial<Family>) => void;
    errors?: Record<string, string>;
}

export function FamilyBasicInfo({ value, onChange, errors }: FamilyBasicInfoProps) {
    return (
        <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
                <Label>
                    Family Name <span className="text-red-500">*</span>
                </Label>
                <Input
                    value={value.name}
                    onChange={(e) => onChange({ name: e.target.value })}
                    placeholder="chat"
                />
                {errors?.[`family_${value.id}_name`] && (
                    <p className="text-sm text-red-500">{errors[`family_${value.id}_name`]}</p>
                )}
                <p className="text-xs text-muted-foreground">
                    Unique identifier for this family (e.g., "chat", "completion")
                </p>
            </div>

            <div className="space-y-2">
                <Label>
                    Display Name <span className="text-red-500">*</span>
                </Label>
                <Input
                    value={value.display_name}
                    onChange={(e) => onChange({ display_name: e.target.value })}
                    placeholder="Chat Completion"
                />
                {errors?.[`family_${value.id}_display_name`] && (
                    <p className="text-sm text-red-500">{errors[`family_${value.id}_display_name`]}</p>
                )}
            </div>

            <div className="col-span-2 space-y-2">
                <Label>
                    Endpoint Path <span className="text-red-500">*</span>
                </Label>
                <Input
                    value={value.schema.endpoint_path as any}
                    onChange={(e) => onChange({ schema: { ...value.schema, endpoint_path: e.target.value } })}
                    placeholder="/chat/completions"
                />
                {errors?.[`family_${value.id}_endpoint_path`] && (
                    <p className="text-sm text-red-500">{errors[`family_${value.id}_endpoint_path`]}</p>
                )}
                <p className="text-xs text-muted-foreground">
                    The endpoint path appended to the base URL
                </p>
            </div>
        </div>
    );
}