import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogClose,
} from "@kana-consultant/ui-kit";
import type { HistoryItem } from "@/services/history";

interface DeleteHistoryDialogProps {
    item: HistoryItem | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onConfirm: () => void;
}

export function DeleteHistoryDialog({ item, open, onOpenChange, onConfirm }: DeleteHistoryDialogProps) {
    if (!item) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Hapus Riwayat?</DialogTitle>
                    <DialogDescription>
                        Apakah Anda yakin ingin menghapus riwayat "{item.title}"?
                        Tindakan ini tidak dapat dibatalkan.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <DialogClose asChild>
                        <button className="px-4 py-2 text-sm font-medium rounded-md border bg-background hover:bg-muted">
                            Batal
                        </button>
                    </DialogClose>
                    <button 
                        onClick={onConfirm} 
                        className="px-4 py-2 text-sm font-medium rounded-md bg-red-600 text-white hover:bg-red-700"
                    >
                        Hapus
                    </button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}