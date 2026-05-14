import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
    Select, SelectContent, SelectItem,
    SelectTrigger, SelectValue,
} from "@/components/ui/select";
import { Search, Trash2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface HistoryHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
    statusFilter: string;
    setStatusFilter: (value: "all" | "success" | "failed" | "pending") => void;
    actionFilter: string;
    setActionFilter: (value: "all" | "published" | "scheduled" | "draft_saved") => void;
    onClearAll: () => void;
}

export default function HistoryHeader({
    searchQuery, setSearchQuery,
    statusFilter, setStatusFilter,
    actionFilter, setActionFilter,
    onClearAll,
}: HistoryHeaderProps) {
    return (
        <div className="space-y-4">
            <div>
                <h2 className="text-xl font-semibold tracking-tight text-slate-900 dark:text-white">
                    History
                </h2>
                <p className="mt-0.5 text-sm text-slate-400 dark:text-slate-500">
                    Riwayat generate dan publish konten
                </p>
            </div>

            <div className="flex flex-wrap items-center gap-2">
                {/* Search */}
                <div className="relative flex-1 min-w-[180px]">
                    <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-slate-400 pointer-events-none" />
                    <Input
                        placeholder="Cari riwayat..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className={cn(
                            "h-8 pl-8 text-sm",
                            "border-slate-200/80 bg-white placeholder:text-slate-400",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
                            "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
                            "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40"
                        )}
                    />
                </div>

                {/* Status filter */}
                <Select value={statusFilter} onValueChange={(v) => setStatusFilter(v as any)}>
                    <SelectTrigger className={cn(
                        "h-8 w-[148px] text-xs",
                        "border-slate-200/80 bg-white text-slate-600",
                        "dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-slate-400",
                        "focus:ring-1 focus:ring-green-500/40 dark:focus:ring-purple-500/40"
                    )}>
                        <SelectValue placeholder="Semua Status" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">Semua Status</SelectItem>
                        <SelectItem value="published">Berhasil</SelectItem>
                        <SelectItem value="failed">Gagal</SelectItem>
                    </SelectContent>
                </Select>

                {/* Action filter */}
                <Select value={actionFilter} onValueChange={(v) => setActionFilter(v as any)}>
                    <SelectTrigger className={cn(
                        "h-8 w-[148px] text-xs",
                        "border-slate-200/80 bg-white text-slate-600",
                        "dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-slate-400",
                        "focus:ring-1 focus:ring-green-500/40 dark:focus:ring-purple-500/40"
                    )}>
                        <SelectValue placeholder="Semua Aksi" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">Semua Aksi</SelectItem>
                        <SelectItem value="published">Published</SelectItem>
                        <SelectItem value="scheduled">Scheduled</SelectItem>
                        <SelectItem value="draft_saved">Draft Saved</SelectItem>
                    </SelectContent>
                </Select>

                {/* Clear all — pushed to end */}
                <Button
                    size="sm"
                    onClick={onClearAll}
                    className={cn(
                        "h-8 gap-1.5 px-3 text-xs font-medium ml-auto",
                        "bg-red-50 text-red-600 border border-red-200/80 shadow-none",
                        "hover:bg-red-100 hover:border-red-300",
                        "dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20",
                        "dark:hover:bg-red-500/20"
                    )}
                    variant="ghost"
                >
                    <Trash2 className="h-3.5 w-3.5" />
                    Hapus Semua
                </Button>
            </div>
        </div>
    );
}