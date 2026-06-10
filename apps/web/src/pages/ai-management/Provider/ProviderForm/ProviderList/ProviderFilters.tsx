// pages/admin/ai-management/components/ProviderList/ProviderFilters.tsx

import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";

interface Filters {
    status: string;
    family: string;
}

interface ProviderFiltersProps {
    filters: Filters;
    onFiltersChange: (filters: Filters) => void;
}

export function ProviderFilters({ filters, onFiltersChange }: ProviderFiltersProps) {
    return (
        <div className="flex gap-4">
            <div className="space-y-1">
                <Label className="text-xs">Status</Label>
                <Select
                    value={filters.status}
                    onValueChange={(value) => onFiltersChange({ ...filters, status: value })}
                >
                    <SelectTrigger className="w-32">
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">All</SelectItem>
                        <SelectItem value="active">Active</SelectItem>
                        <SelectItem value="inactive">Inactive</SelectItem>
                    </SelectContent>
                </Select>
            </div>

            <div className="space-y-1">
                <Label className="text-xs">Family</Label>
                <Select
                    value={filters.family}
                    onValueChange={(value) => onFiltersChange({ ...filters, family: value })}
                >
                    <SelectTrigger className="w-40">
                        <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">All Families</SelectItem>
                        <SelectItem value="chat">Chat</SelectItem>
                        <SelectItem value="completion">Completion</SelectItem>
                        <SelectItem value="embedding">Embedding</SelectItem>
                        <SelectItem value="image">Image Generation</SelectItem>
                    </SelectContent>
                </Select>
            </div>
        </div>
    );
}