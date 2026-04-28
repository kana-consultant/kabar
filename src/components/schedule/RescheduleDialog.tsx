import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Calendar } from "lucide-react";
import type { Draft } from "@/types/draft";

interface RescheduleDialogProps {
    schedule: Draft | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    newDate: string;
    onDateChange: (value: string) => void;
    newTime: string;
    onTimeChange: (value: string) => void;
    onReschedule: () => void;
}

export function RescheduleDialog({
    schedule,
    open,
    onOpenChange,
    newDate,
    onDateChange,
    newTime,
    onTimeChange,
    onReschedule,
}: RescheduleDialogProps) {
    if (!schedule) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Ubah Jadwal Posting</DialogTitle>
                    <DialogDescription>
                        Atur ulang jadwal untuk "{schedule.title}"
                    </DialogDescription>
                </DialogHeader>
                <div className="grid grid-cols-2 gap-3">
                    <div>
                        <label className="text-sm font-medium">Tanggal</label>
                        <Input
                            type="date"
                            value={newDate}
                            onChange={(e) => onDateChange(e.target.value)}
                            className="mt-1"
                        />
                    </div>
                    <div>
                        <label className="text-sm font-medium">Waktu</label>
                        <Input
                            type="time"
                            value={newTime}
                            onChange={(e) => onTimeChange(e.target.value)}
                            className="mt-1"
                        />
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Batal
                    </Button>
                    <Button onClick={onReschedule}>
                        <Calendar className="mr-2 h-4 w-4" />
                        Perbarui Jadwal
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}