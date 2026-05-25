import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogClose,
  Button,
} from "@kana-consultant/ui-kit";
import type { Draft } from "@/services/draft";

interface DeleteAlertDialogProps {
    draft: Draft | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onConfirm: () => void;
}

export function DeleteAlertDialog({ draft, open, onOpenChange, onConfirm }: DeleteAlertDialogProps) {
    if (!draft) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Hapus Draft?</DialogTitle>
                    <DialogDescription>
                        Apakah Anda yakin ingin menghapus "{draft.title}"?
                        Tindakan ini tidak dapat dibatalkan.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <DialogClose asChild>
                        <Button variant="outline">Batal</Button>
                    </DialogClose>
                    <Button 
                        onClick={onConfirm} 
                        variant="destructive"
                    >
                        Hapus
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}