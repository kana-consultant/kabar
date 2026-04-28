import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ImageIcon } from "lucide-react";
import type { Draft } from "@/types/draft";

interface ViewDraftDialogProps {
    draft: Draft | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
    formatDate: (date: string) => string;
}

export function ViewDraftDialog({ draft, open, onOpenChange, formatDate }: ViewDraftDialogProps) {
    if (!draft) return null;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-6xl max-h-[85vh] overflow-auto">
                <DialogHeader>
                    <DialogTitle className="text-xl">{draft.title}</DialogTitle>
                    <DialogDescription>
                        Dibuat: {formatDate(draft.createdAt)}
                        {draft.updatedAt !== draft.createdAt && (
                            <span className="ml-2">
                                | Diperbarui: {formatDate(draft.updatedAt)}
                            </span>
                        )}
                    </DialogDescription>
                </DialogHeader>

                <Tabs defaultValue="article" className="w-full">
                    <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value="article">Artikel</TabsTrigger>
                        <TabsTrigger value="image">Gambar</TabsTrigger>
                    </TabsList>

                    <TabsContent value="article" className="pt-4">
                        <div>
                            <h4 className="font-medium text-sm text-slate-500">Topik</h4>
                            <p className="text-slate-700 dark:text-slate-300">{draft.topic}</p>
                        </div>
                        <div className="mt-4">
                            <h4 className="font-medium text-sm text-slate-500">Target Produk</h4>
                            <div className="flex flex-wrap gap-2 mt-1">
                                {draft.targetProducts.map(p => (
                                    <span key={p} className="rounded-full bg-slate-100 px-2 py-1 text-xs dark:bg-slate-800">
                                        {p}
                                    </span>
                                ))}
                            </div>
                        </div>
                        <div className="mt-4">
                            <h4 className="font-medium text-sm text-slate-500">Artikel</h4>
                            <div className="mt-2 rounded-lg bg-slate-50 p-4 text-sm dark:bg-slate-900 max-h-[50vh] overflow-auto">
                                <pre className="whitespace-pre-wrap font-sans">
                                    {draft.article}
                                </pre>
                            </div>
                        </div>
                    </TabsContent>

                    <TabsContent value="image" className="pt-4">
                        {draft.imageUrl ? (
                            <div className="space-y-4">
                                <img
                                    src={draft.imageUrl}
                                    alt={draft.title}
                                    className="w-full rounded-lg border object-cover"
                                />
                                <p className="text-center text-sm text-slate-500">
                                    Klik kanan pada gambar → Save Image As untuk menyimpan
                                </p>
                            </div>
                        ) : (
                            <div className="flex min-h-[200px] flex-col items-center justify-center rounded-lg border-2 border-dashed">
                                <ImageIcon className="h-8 w-8 text-slate-400" />
                                <p className="mt-2 text-slate-500">Tidak ada gambar</p>
                            </div>
                        )}
                    </TabsContent>
                </Tabs>
            </DialogContent>
        </Dialog>
    );
}