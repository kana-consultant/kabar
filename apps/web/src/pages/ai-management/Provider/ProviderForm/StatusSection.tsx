// pages/admin/ai-management/components/ProviderForm/StatusSection.tsx

import { Label } from "@kana-consultant/ui-kit";
import { Switch } from "@kana-consultant/ui-kit";

interface StatusSectionProps {
    value: boolean;
    onChange: (value: boolean) => void;
}

export function StatusSection({ value, onChange }: StatusSectionProps) {
    return (
        <div className="flex items-center justify-between p-4 bg-muted/30 rounded-lg">
            <div className="space-y-1">
                <Label className="text-base">Active Status</Label>
                <p className="text-sm text-muted-foreground">
                    Enable or disable this provider. Disabled providers won't be available for use.
                </p>
            </div>
            <Switch
                checked={value}
                onCheckedChange={onChange}
            />
        </div>
    );
}