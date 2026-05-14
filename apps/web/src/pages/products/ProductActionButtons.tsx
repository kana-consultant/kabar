import { Button } from "@/components/ui/button";
import { Save } from "lucide-react";

interface ProductActionButtonsProps {
    onCancel: () => void;
    onSave: () => void;
    isSaving?: boolean;
}

export function ProductActionButtons({ onCancel, onSave, isSaving = false }: ProductActionButtonsProps) {
    return (
        <div className="flex justify-end gap-3">
            <Button variant="outline" onClick={onCancel}>
                Batal
            </Button>
            <Button onClick={onSave} disabled={isSaving}>
                <Save className="mr-2 h-4 w-4" />
                {isSaving ? "Menyimpan..." : "Simpan Produk"}
            </Button>
        </div>
    );
}