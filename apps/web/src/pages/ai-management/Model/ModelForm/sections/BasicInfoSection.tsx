import { Input, Textarea, Label, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@kana-consultant/ui-kit";
import type { ModelFormSectionProps } from "../../provider.types";

export function BasicInfoSection({ formData, setFormData, providers, filteredFamilies, onTestModel, isEditMode }: ModelFormSectionProps) {
    const handleProviderChange = (providerId: string) => {
        setFormData(prev => ({
            ...prev,
            provider_id: providerId,
            family_id: "",
        }));
    };

    const handleFamilyChange = (familyId: string | null) => {
        setFormData(prev => ({ ...prev, family_id: familyId as string }));
    };

    console.log("jfaslomfalsml")
    console.log(formData, filteredFamilies)

    return (
        <div className="space-y-4">
            <h3 className="text-lg font-semibold">Basic Information</h3>

            <div className="">
                <div className="space-y-2">
                    <Label>Provider *</Label>
                    <Select value={formData.provider_id} onValueChange={handleProviderChange}>
                        <SelectTrigger>
                            <SelectValue placeholder="Select provider" />
                        </SelectTrigger>
                        <SelectContent>
                            {providers?.map((provider) => (
                                <SelectItem key={provider.id} value={provider.id as string}>
                                    {provider.display_name}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>

                
            </div>
            <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                    <Label>Model Name *</Label>
                    <Input
                        value={formData.name}
                        onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                        placeholder="gpt-4-turbo"
                     
                    />
                    <p className="text-xs text-muted-foreground">
                        API model identifier (must match provider's model name)
                    </p>
                </div>

                <div className="space-y-2">
                    <Label>Display Name *</Label>
                    <Input
                        value={formData.display_name}
                        onChange={(e) => setFormData(prev => ({ ...prev, display_name: e.target.value }))}
                        placeholder="GPT-4 Turbo"
                    />
                </div>
            </div>

            <div className="space-y-2">
                <Label>Description</Label>
                <Textarea
                    value={formData.description || ""}
                    onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value || null }))}
                    placeholder="Model description..."
                    rows={2}
                />
            </div>

            {/* <div className="space-y-2">
                <Label>Team ID (Optional)</Label>
                <Input
                    value={formData.team_id || ""}
                    onChange={(e) => setFormData(prev => ({ ...prev, team_id: e.target.value || null }))}
                    placeholder="Leave empty for global model"
                />
                <p className="text-xs text-muted-foreground">
                    Assign to specific team for multi-tenant isolation
                </p>
            </div> */}
        </div>
    );
}