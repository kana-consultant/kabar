// pages/admin/ai-management/components/ProviderForm/FamiliesSection/index.tsx

import { Button } from "@kana-consultant/ui-kit";
import { Plus } from "lucide-react";
import { FamilyCard } from "./FamilyCard";
import { useFamilyManagement } from "@/hooks/useAIManagement/useFamilyManagement";
import type { ProviderFormData, Family } from "@/types/provider.types";

interface FamiliesSectionProps {
    value: ProviderFormData;
    onChange: (updates: Partial<ProviderFormData>) => void;
    errors?: Record<string, string>;
    providerConfig?: {
        base_url: string;
        auth_header: string;
        auth_prefix: string | null;
        default_headers : Record<string, string>;
        display_name?: string;
    };
}

export function FamiliesSection({ value, onChange, errors, providerConfig }: FamiliesSectionProps) {
    const {
        families,
        addFamily,
        removeFamily,
        updateFamily,
        updateFamilySchema,
        duplicateFamily,
    } = useFamilyManagement(value.families);
    console.log("Families Section : ")
    console.log(value)

    // Handle update seluruh family (termasuk schema)
    const handleUpdateFamily = (id: string, updates: Partial<Family>) => {
        // Pisahkan update untuk family fields dan schema fields
        const { schema, ...familyUpdates } = updates;
        
        if (schema) {
            // Jika ada update schema, gunakan updateFamilySchema
            updateFamilySchema(id, schema);
        }
        
        if (Object.keys(familyUpdates).length > 0) {
            // Update family fields
            updateFamily(id, familyUpdates);
        }
        
        // Sync ke parent
        const updatedFamilies = families.map(f => {
            if (f.id === id) {
                let updated = { ...f, ...familyUpdates };
                if (schema) {
                    updated = { ...updated, schema: { ...updated.schema, ...schema } };
                }
                return updated;
            }
            return f;
        });
        
        onChange({ families: updatedFamilies });
    };

    const handleAddFamily = () => {
        const newFamily = addFamily();
        onChange({ families: [...families, newFamily] });
    };

    const handleRemoveFamily = (id: string) => {
        removeFamily(id);
        const updatedFamilies = families.filter(f => f.id !== id);
        onChange({ families: updatedFamilies });
    };

    const handleDuplicateFamily = (family: Family) => {
        const duplicated = duplicateFamily(family);
        onChange({ families: [...families, duplicated] });
    };

    console.log("source :",value.families)
    console.log("Families with nested schema:", families);

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <div>
                    <h4 className="font-medium">API Families</h4>
                    <p className="text-sm text-muted-foreground">
                        Define different API endpoints and their configurations
                    </p>
                </div>
                <Button type="button" variant="outline" size="sm" onClick={handleAddFamily}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Family
                </Button>
            </div>

            {errors?.families && (
                <p className="text-sm text-red-500">{errors.families}</p>
            )}

            <div className="space-y-4">
                {families.map((family, index) => (
                    <FamilyCard
                        key={family.id}
                        family={family}
                        index={index}
                        isOnly={families.length === 1}
                        onUpdate={(updates) => handleUpdateFamily(family.id, updates)}
                        onRemove={() => handleRemoveFamily(family.id)}
                        onDuplicate={() => handleDuplicateFamily(family)}
                        errors={errors}
                        providerConfig={providerConfig}
                    />
                ))}
            </div>

            {/* Empty state */}
            {families.length === 0 && (
                <div className="text-center py-8 border rounded-lg bg-gray-50">
                    <p className="text-muted-foreground">No families configured</p>
                    <Button 
                        type="button" 
                        variant="link" 
                        onClick={handleAddFamily}
                        className="mt-2"
                    >
                        Add your first family
                    </Button>
                </div>
            )}
        </div>
    );
}