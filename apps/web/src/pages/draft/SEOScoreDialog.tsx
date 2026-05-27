// SEOScoreDialog.tsx
import { cn } from "@/lib/utils";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@kana-consultant/ui-kit"
import { BarChart2, CheckCircle, AlertCircle, TrendingUp } from "lucide-react";
import type { SEOScore } from "@/services/draft";

interface SEOScoreDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    draftTitle?: string;
    data: SEOScore | null;
    loading?: boolean;
}

export function SEOScoreDialog({ open, onOpenChange, draftTitle, data, loading }: SEOScoreDialogProps) {
    const scoreColor = (score: number) => {
        if (score >= 80) return "text-green-600 dark:text-green-400";
        if (score >= 60) return "text-yellow-600 dark:text-yellow-400";
        return "text-red-500 dark:text-red-400";
    };

    const scoreBg = (score: number) => {
        if (score >= 80) return "bg-green-50 dark:bg-green-500/10 ring-green-200 dark:ring-green-500/20";
        if (score >= 60) return "bg-yellow-50 dark:bg-yellow-500/10 ring-yellow-200 dark:ring-yellow-500/20";
        return "bg-red-50 dark:bg-red-500/10 ring-red-200 dark:ring-red-500/20";
    };

    const barColor = (score: number) => {
        if (score >= 80) return "bg-green-500";
        if (score >= 60) return "bg-yellow-500";
        return "bg-red-500";
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-md bg-white dark:bg-[#0f0d1a] border-slate-200/80 dark:border-white/[0.06]">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2 text-slate-900 dark:text-white">
                        <BarChart2 className="h-4 w-4 text-orange-500" />
                        SEO Score
                    </DialogTitle>
                    {draftTitle && (
                        <p className="text-xs text-slate-400 dark:text-slate-500 truncate mt-1">
                            {draftTitle}
                        </p>
                    )}
                </DialogHeader>

                {loading ? (
                    <div className="space-y-3 py-4">
                        <div className="h-20 w-20 rounded-full bg-slate-100 dark:bg-white/5 animate-pulse mx-auto" />
                        {[1, 2, 3].map((i) => (
                            <div key={i} className="h-8 rounded-lg bg-slate-100 dark:bg-white/5 animate-pulse" />
                        ))}
                    </div>
                ) : data ? (
                    <div className="space-y-5">
                        {/* Total score */}
                        <div className="flex justify-center">
                            <div className={cn(
                                "flex h-24 w-24 flex-col items-center justify-center rounded-full ring-2",
                                scoreBg(data.total)
                            )}>
                                <span className={cn("text-3xl font-bold tabular-nums", scoreColor(data.total))}>
                                    {data.total}
                                </span>
                                <span className="text-[10px] text-slate-400 dark:text-slate-500">/ 100</span>
                            </div>
                        </div>

                        {/* Details */}
                        {Object.keys(data.details).length > 0 && (
                            <div className="space-y-2.5">
                                <p className="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wide">
                                    Detail
                                </p>
                                {Object.entries(data.details).map(([key, value]) => (
                                    <div key={key} className="space-y-1">
                                        <div className="flex justify-between text-xs">
                                            <span className="text-slate-600 dark:text-slate-400 capitalize">
                                                {key.replace(/_/g, " ")}
                                            </span>
                                            <span className={cn("font-medium tabular-nums", scoreColor(value))}>
                                                {value}
                                            </span>
                                        </div>
                                        <div className="h-1.5 w-full rounded-full bg-slate-100 dark:bg-white/5">
                                            <div
                                                className={cn("h-1.5 rounded-full transition-all duration-500", barColor(value))}
                                                style={{ width: `${Math.min(value, 100)}%` }}
                                            />
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}

                        {/* Suggestions */}
                        {data.suggestions.length > 0 && (
                            <div className="space-y-2">
                                <p className="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wide">
                                    Saran
                                </p>
                                <ul className="space-y-1.5">
                                    {data.suggestions.map((s, i) => (
                                        <li key={i} className="flex items-start gap-2 text-xs text-slate-600 dark:text-slate-400">
                                            <AlertCircle className="h-3.5 w-3.5 shrink-0 mt-0.5 text-yellow-500" />
                                            {s}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        )}

                        {data.suggestions.length === 0 && (
                            <div className="flex items-center gap-2 rounded-lg bg-green-50 dark:bg-green-500/10 px-3 py-2.5">
                                <CheckCircle className="h-4 w-4 shrink-0 text-green-500" />
                                <p className="text-xs text-green-700 dark:text-green-400">
                                    Tidak ada saran — konten sudah optimal!
                                </p>
                            </div>
                        )}
                    </div>
                ) : (
                    <div className="flex flex-col items-center gap-2 py-8 text-slate-400">
                        <TrendingUp className="h-8 w-8 opacity-30" />
                        <p className="text-sm">Gagal memuat SEO score</p>
                    </div>
                )}
            </DialogContent>
        </Dialog>
    );
}