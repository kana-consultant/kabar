import { Input, Label, Switch } from "@kana-consultant/ui-kit";
import type { ModelFormSectionProps } from "../../provider.types";

interface StatusSectionProps extends ModelFormSectionProps {
    isEditMode?: boolean;
    createdBy?: string | null;
}

export function StatusSection({ formData, setFormData, isEditMode, createdBy }: StatusSectionProps) {
    return (
        <div className="space-y-4">
            <h3 className="text-lg font-semibold">Status</h3>

            <div className="flex items-center justify-between">
                <div className="space-y-1">
                    <Label>Active</Label>
                    <p className="text-sm text-muted-foreground">
                        Model can be used when active
                    </p>
                </div>
                <Switch
                    checked={formData.is_active}
                    onCheckedChange={(checked) => setFormData((prev : any) => ({ ...prev, is_active: checked }))}
                />
            </div>

            <div className="flex items-center justify-between">
                <div className="space-y-1">
                    <Label>Default Model</Label>
                    <p className="text-sm text-muted-foreground">
                        Set as default for this provider
                    </p>
                </div>
                <Switch
                    checked={formData.is_default}
                    onCheckedChange={(checked) => setFormData((prev : any) => ({ ...prev, is_default: checked }))}
                />
            </div>

            {isEditMode && createdBy && (
                <div className="space-y-2">
                    <Label>Created By</Label>
                    <Input value={createdBy} disabled className="bg-muted" />
                </div>
            )}
        </div>
    );
}