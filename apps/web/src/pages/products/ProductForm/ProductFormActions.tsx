import { Button } from "@/components/ui/button";
import { Save, Loader2 } from "lucide-react";

interface ProductFormActionsProps {
    onCancel: () => void;
    onSave: () => void;
    isSaving: boolean;
}

export function ProductFormActions({ onCancel, onSave, isSaving }: ProductFormActionsProps) {
    return (
        <div className="flex justify-end gap-3 pt-4 border-t">
            <Button
                variant="outline"
                onClick={onCancel}
                className="h-10 px-6"
            >
                Batal
            </Button>
            <Button
                onClick={onSave}
                disabled={isSaving}
                className="h-10 px-6 bg-blue-600 hover:bg-blue-700"
            >
                {isSaving ? (
                    <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Menyimpan...
                    </>
                ) : (
                    <>
                        <Save className="mr-2 h-4 w-4" />
                        Simpan Produk
                    </>
                )}
            </Button>
        </div>
    );
}