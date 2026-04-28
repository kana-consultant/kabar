import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { Draft } from "@/types/draft";

interface ScheduleDetailDialogProps {
    schedule: Draft | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    getScheduleDisplay: (scheduledFor?: string) => string;
    formatDate: (date: string) => string;
}

export function ScheduleDetailDialog({
    schedule,
    open,
    onOpenChange,
    getScheduleDisplay,
    formatDate,
}: ScheduleDetailDialogProps) {
    if (!schedule) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-6xl max-h-[80vh] overflow-auto">
                <DialogHeader>
                    <DialogTitle>{schedule.title}</DialogTitle>
                    <DialogDescription>
                        Jadwal: {getScheduleDisplay(schedule.scheduledFor)}
                    </DialogDescription>
                </DialogHeader>
                <Tabs defaultValue="article" className="w-full">
                    <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value="article">Artikel</TabsTrigger>
                        <TabsTrigger value="info">Informasi</TabsTrigger>
                    </TabsList>
                    <TabsContent value="article" className="pt-4">
                        <div className="rounded-lg bg-slate-50 p-4 text-sm dark:bg-slate-900 max-h-[50vh] overflow-auto">
                            <pre className="whitespace-pre-wrap font-sans">
                                <div dangerouslySetInnerHTML={{ __html: schedule.article }} />
                            </pre>
                        </div>
                    </TabsContent>
                    <TabsContent value="info" className="pt-4 space-y-4">
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Topik</h4>
                            <p>{schedule.topic.replace(/<[^>]*>/g, '')}</p>
                        </div>
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Target Produk</h4>
                            <div className="flex flex-wrap gap-2 mt-1">
                                {schedule.targetProducts.map(p => (
                                    <span key={p} className="rounded-full bg-slate-100 px-2 py-1 text-xs dark:bg-slate-800">
                                        {p}
                                    </span>
                                ))}
                            </div>
                        </div>
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Dibuat</h4>
                            <p>{formatDate(schedule.createdAt)}</p>
                        </div>
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Terakhir Diupdate</h4>
                            <p>{formatDate(schedule.updatedAt)}</p>
                        </div>
                    </TabsContent>
                </Tabs>
            </DialogContent>
        </Dialog>
    );
}