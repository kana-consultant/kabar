import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { type Field } from "./types";

interface JsonBuilderHeaderProps {
    fields: Field[];
    onAddField: () => void;
    onAddObject: () => void;
    nextFieldNumber: number;
    nextObjectNumber: number;
}

export function JsonBuilderHeader({
    onAddField,
    onAddObject,
    nextFieldNumber,
    nextObjectNumber
}: JsonBuilderHeaderProps) {
    return (
        <div className="flex gap-2 mb-4">
            <Button size="sm" onClick={onAddField}>
                <Plus className="h-3 w-3 mr-1" />
                + Tambah Field (field{nextFieldNumber})
            </Button>
            <Button size="sm" variant="outline" onClick={onAddObject}>
                <Plus className="h-3 w-3 mr-1" />
                + Tambah Object (object{nextObjectNumber})
            </Button>
        </div>
    );
}