import { useState } from "react";
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

interface DeleteScheduleDialogProps {
    schedule: Draft | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onConfirm: () => void;
}

export function DeleteScheduleDialog({ schedule, open, onOpenChange, onConfirm }: DeleteScheduleDialogProps) {
    const [isDeleting, setIsDeleting] = useState(false);
    
    if (!schedule) return null;

    const handleConfirm = async () => {
        setIsDeleting(true);
        await onConfirm();
        setIsDeleting(false);
    };

    return (
        <Dialog open={open} onOpenChange={(newOpen : boolean) => {
            if (!isDeleting) onOpenChange(newOpen);
        }}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Hapus Jadwal?</DialogTitle>
                    <DialogDescription>
                        Apakah Anda yakin ingin menghapus jadwal "{schedule.title}"?
                        Tindakan ini tidak dapat dibatalkan.
                    </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                    <DialogClose asChild>
                        <Button variant="outline" disabled={isDeleting}>Batal</Button>
                    </DialogClose>
                    <Button 
                        onClick={handleConfirm} 
                        variant="destructive"
                        disabled={isDeleting}
                    >
                        {isDeleting ? "Menghapus..." : "Hapus"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}