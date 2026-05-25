import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Switch } from "@kana-consultant/ui-kit";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@kana-consultant/ui-kit";
import { Calendar } from "lucide-react";
import type { Draft } from "@//services/draft";

interface ScheduleDialogProps {
    draft: Draft | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    scheduleDate: string;
    setScheduleDate: (value: string) => void;
    scheduleTime: string;
    setScheduleTime: (value: string) => void;
    dailySchedule: boolean;
    setDailySchedule: (value: boolean) => void;
    dailyTime: string;
    setDailyTime: (value: string) => void;
    onSchedule: () => void;
}

export function ScheduleDialog({
    draft,
    open,
    onOpenChange,
    scheduleDate,
    setScheduleDate,
    scheduleTime,
    setScheduleTime,
    dailySchedule,
    setDailySchedule,
    dailyTime,
    setDailyTime,
    onSchedule,
}: ScheduleDialogProps) {
    if (!draft) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Jadwalkan Posting</DialogTitle>
                    <DialogDescription>
                        Atur jadwal untuk "{draft.title}"
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div className="flex items-center justify-between">
                        <Label>Posting Berulang (Harian)</Label>
                        <Switch checked={dailySchedule} onCheckedChange={setDailySchedule} />
                    </div>

                    {dailySchedule ? (
                        <div>
                            <Label>Waktu Posting Harian</Label>
                            <Input
                                type="time"
                                value={dailyTime}
                                onChange={(e : any) => setDailyTime(e.target.value)}
                                className="mt-1"
                            />
                            <p className="mt-1 text-xs text-slate-500">
                                Konten akan diposting setiap hari jam {dailyTime}
                            </p>
                        </div>
                    ) : (
                        <div className="grid grid-cols-2 gap-3">
                            <div>
                                <Label>Tanggal</Label>
                                <Input
                                    type="date"
                                    value={scheduleDate}
                                    onChange={(e : any) => setScheduleDate(e.target.value)}
                                    className="mt-1"
                                />
                            </div>
                            <div>
                                <Label>Waktu</Label>
                                <Input
                                    type="time"
                                    value={scheduleTime}
                                    onChange={(e : any) => setScheduleTime(e.target.value)}
                                    className="mt-1"
                                />
                            </div>
                        </div>
                    )}
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Batal
                    </Button>
                    <Button onClick={onSchedule}>
                        <Calendar className="mr-2 h-4 w-4" />
                        Jadwalkan
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}