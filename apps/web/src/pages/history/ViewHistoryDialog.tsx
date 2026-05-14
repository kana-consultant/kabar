import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { HistoryItem } from "@/services/history";

interface ViewHistoryDialogProps {
    item: HistoryItem | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    formatDate: (date: string) => string;
    getStatusData: (status: string) => { label: string; icon: string; color: string };
    getActionData: (action: string) => { label: string; icon: string };
}

export function ViewHistoryDialog({
    item,
    open,
    onOpenChange,
    formatDate,
    getStatusData,
    getActionData,
}: ViewHistoryDialogProps) {
    if (!item) return null;

    const status = getStatusData(item.status);
    const action = getActionData(item.action);

    const statusColorClass = {
        green: "text-green-700 dark:text-green-400",
        red: "text-red-700 dark:text-red-400",
        yellow: "text-yellow-700 dark:text-yellow-400",
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-3xl max-h-[80vh] overflow-auto">
                <DialogHeader>
                    <DialogTitle>{item.title}</DialogTitle>
                    <DialogDescription>
                        Dipublikasikan: {formatDate(item.publishedAt)}
                    </DialogDescription>
                </DialogHeader>
                <Tabs defaultValue="content" className="w-full">
                    <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value="content">Konten</TabsTrigger>
                        <TabsTrigger value="info">Informasi</TabsTrigger>
                    </TabsList>
                    <TabsContent value="content" className="pt-4">
                        <div className="rounded-lg bg-slate-50 p-4 text-sm dark:bg-slate-900 max-h-[50vh] overflow-auto"
                        dangerouslySetInnerHTML={{ __html: item.content }}>
                            
                        </div>
                    </TabsContent>
                    <TabsContent value="info" className="pt-4 space-y-4 ">
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Topik</h4>
                            <p>{item.topic.replace(/<[^>]*>/g, '')}</p>
                        </div>
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Target Produk</h4>
                            <div className="flex flex-wrap gap-2 mt-1">
                                {item.targetProducts.map(p => (
                                    <span key={p} className="rounded-full bg-slate-100 px-2 py-1 text-xs dark:bg-slate-800">
                                        {p}
                                    </span>
                                ))}
                            </div>
                        </div>
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Status</h4>
                            <p className={`${statusColorClass[status.color as keyof typeof statusColorClass]}`}>
                                {status.icon} {status.label}
                            </p>
                        </div>
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Aksi</h4>
                            <p>{action.icon} {action.label}</p>
                        </div>
                    </TabsContent>
                </Tabs>
            </DialogContent>
        </Dialog>
    );
}