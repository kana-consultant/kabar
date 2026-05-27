// SimilarityDialog.tsx
import { cn } from "@/lib/utils";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@kana-consultant/ui-kit"
import { GitCompare, AlertTriangle, CheckCircle } from "lucide-react";
import type { SimilarDraft } from "@/services/draft";

interface SimilarityDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    draftTitle?: string;
    data: SimilarDraft[] | null;
    loading?: boolean;
}

export function SimilarityDialog({ open, onOpenChange, draftTitle, data, loading }: SimilarityDialogProps) {
    const similarityColor = (score: number) => {
        if (score >= 0.8) return "text-red-500 dark:text-red-400";
        if (score >= 0.5) return "text-yellow-600 dark:text-yellow-400";
        return "text-green-600 dark:text-green-400";
    };

    const similarityBg = (score: number) => {
        if (score >= 0.8) return "bg-red-50 dark:bg-red-500/10 border-red-200/80 dark:border-red-500/20";
        if (score >= 0.5) return "bg-yellow-50 dark:bg-yellow-500/10 border-yellow-200/80 dark:border-yellow-500/20";
        return "bg-green-50 dark:bg-green-500/10 border-green-200/80 dark:border-green-500/20";
    };

    const similarityLabel = (score: number) => {
        if (score >= 0.8) return "Sangat mirip";
        if (score >= 0.5) return "Cukup mirip";
        return "Berbeda";
    };

    const barColor = (score: number) => {
        if (score >= 0.8) return "bg-red-500";
        if (score >= 0.5) return "bg-yellow-500";
        return "bg-green-500";
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-md bg-white dark:bg-[#0f0d1a] border-slate-200/80 dark:border-white/[0.06]">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2 text-slate-900 dark:text-white">
                        <GitCompare className="h-4 w-4 text-blue-500" />
                        Cek Kemiripan
                    </DialogTitle>
                    {draftTitle && (
                        <p className="text-xs text-slate-400 dark:text-slate-500 truncate mt-1">
                            {draftTitle}
                        </p>
                    )}
                </DialogHeader>

                {loading ? (
                    <div className="space-y-3 py-4">
                        {[1, 2, 3].map((i) => (
                            <div key={i} className="h-16 rounded-xl bg-slate-100 dark:bg-white/5 animate-pulse" />
                        ))}
                    </div>
                ) : data && data.length > 0 ? (
                    <div className="space-y-2.5 max-h-[60vh] overflow-y-auto pr-1">
                        {data.map((item) => (
                            <div key={item.draft_id} className={cn(
                                "rounded-xl border p-3 space-y-2",
                                similarityBg(item.similarity)
                            )}>
                                <div className="flex items-start justify-between gap-2">
                                    <p className="text-xs font-medium text-slate-700 dark:text-slate-300 line-clamp-2 flex-1">
                                        {item.title}
                                    </p>
                                    <span className={cn(
                                        "shrink-0 text-xs font-bold tabular-nums",
                                        similarityColor(item.similarity)
                                    )}>
                                        {Math.round(item.similarity * 100)}%
                                    </span>
                                </div>

                                <div className="space-y-1">
                                    <div className="h-1.5 w-full rounded-full bg-white/60 dark:bg-black/20">
                                        <div
                                            className={cn("h-1.5 rounded-full transition-all duration-500", barColor(item.similarity))}
                                            style={{ width: `${item.similarity * 100}%` }}
                                        />
                                    </div>
                                    <p className={cn("text-[10px]", similarityColor(item.similarity))}>
                                        {similarityLabel(item.similarity)}
                                    </p>
                                </div>
                            </div>
                        ))}
                    </div>
                ) : data && data.length === 0 ? (
                    <div className="flex flex-col items-center gap-2 py-8 text-slate-400">
                        <CheckCircle className="h-8 w-8 text-green-500 opacity-70" />
                        <p className="text-sm text-green-600 dark:text-green-400">
                            Tidak ditemukan konten yang mirip
                        </p>
                    </div>
                ) : (
                    <div className="flex flex-col items-center gap-2 py-8 text-slate-400">
                        <AlertTriangle className="h-8 w-8 opacity-30" />
                        <p className="text-sm">Gagal memuat data kemiripan</p>
                    </div>
                )}
            </DialogContent>
        </Dialog>
    );
}