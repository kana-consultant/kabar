// pages/admin/ai-management/components/ProviderForm/ConnectionSection.tsx

import { useState } from "react";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { KeyValueEditor } from "@/components/KeyValueEditor";
import { type AuthPreset, AUTH_PRESETS } from "@/constants/auth-presets";
import type { ProviderFormData } from "@/types/provider.types";
interface ConnectionSectionProps {
    value: ProviderFormData;
    onChange: (updates: Partial<ProviderFormData>) => void;
    errors?: Record<string, string>;
    onTestConnection?: () => void;
}

export function ConnectionSection({ value, onChange, errors, onTestConnection }: ConnectionSectionProps) {
    const [authPreset, setAuthPreset] = useState<AuthPreset>(() => {
        if (value.auth_type === "bearer" && value.auth_header === "Authorization" && value.auth_prefix === "Bearer") {
            return "bearer_openai";
        }
        if (value.auth_type === "api_key" && value.auth_header === "x-goog-api-key") {
            return "api_key_google";
        }
        return "custom";
    });

    const handleAuthPresetChange = (preset: AuthPreset) => {
        setAuthPreset(preset);
        const config = AUTH_PRESETS[preset];
        onChange({
            auth_type: config.auth_type,
            auth_header: config.auth_header,
            auth_prefix: config.auth_prefix,
        });
    };

    return (
        <div className="space-y-4">
            <div className="space-y-2">
                <Label htmlFor="base_url">
                    Base URL <span className="text-red-500">*</span>
                </Label>
                <Input
                    id="base_url"
                    value={value.base_url}
                    onChange={(e) => onChange({ base_url: e.target.value })}
                    placeholder="https://api.openai.com/v1"
                    className={errors?.base_url ? "border-red-500" : ""}
                />
                {errors?.base_url && (
                    <p className="text-sm text-red-500">{errors.base_url}</p>
                )}
                <p className="text-xs text-muted-foreground">
                    The base endpoint URL for all API calls
                </p>
            </div>

            <div className="space-y-2">
                <Label>Authentication Preset</Label>
                <Select value={authPreset} onValueChange={handleAuthPresetChange}>
                    <SelectTrigger>
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                        {Object.entries(AUTH_PRESETS).map(([key, config]) => (
                            <SelectItem key={key} value={key}>
                                {config.label}
                            </SelectItem>
                        ))}
                    </SelectContent>
                </Select>
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                    <Label htmlFor="auth_type">Auth Type</Label>
                    <Input
                        id="auth_type"
                        value={value.auth_type}
                        onChange={(e) => onChange({ auth_type: e.target.value })}
                        placeholder="bearer"
                    />
                </div>

                <div className="space-y-2">
                    <Label htmlFor="auth_header">Auth Header</Label>
                    <Input
                        id="auth_header"
                        value={value.auth_header}
                        onChange={(e) => onChange({ auth_header: e.target.value })}
                        placeholder="Authorization"
                    />
                </div>
            </div>

            <div className="space-y-2">
                <Label htmlFor="auth_prefix">Auth Prefix (Optional)</Label>
                <Input
                    id="auth_prefix"
                    value={value.auth_prefix || ""}
                    onChange={(e) => onChange({ auth_prefix: e.target.value || null })}
                    placeholder="Bearer"
                />
            </div>

            <div className="space-y-2">
                <Label>Default Headers</Label>
                <KeyValueEditor
                    value={value.default_headers}
                    onChange={(headers) => onChange({ default_headers: headers })}
                />
                <p className="text-xs text-muted-foreground">
                    Headers that will be sent with every request
                </p>
            </div>

            {onTestConnection && (
                <Button type="button" variant="outline" onClick={onTestConnection}>
                    Test Connection
                </Button>
            )}
        </div>
    );
}