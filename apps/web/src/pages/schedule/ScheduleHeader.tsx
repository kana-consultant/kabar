import { Input } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Search, RefreshCw } from "lucide-react";
import { cn } from "@/lib/utils";

interface ScheduleHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
    onRefresh: () => void;
}

export function ScheduleHeader({ searchQuery, setSearchQuery, onRefresh }: ScheduleHeaderProps) {
    return (
        <div className="flex items-start justify-between gap-4">
            <div>
                <h2 className="text-xl font-semibold tracking-tight text-slate-900 dark:text-white">
                    Jadwal Posting
                </h2>
                <p className="mt-0.5 text-sm text-slate-400 dark:text-slate-500">
                    Kelola konten yang dijadwalkan untuk dipublikasikan
                </p>
            </div>

            <div className="flex items-center gap-2 shrink-0">
                {/* Search */}
                <div className="relative">
                    <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-slate-400 pointer-events-none" />
                    <Input
                        placeholder="Cari jadwal..."
                        value={searchQuery}
                        onChange={(e : any) => setSearchQuery(e.target.value)}
                        className={cn(
                            "h-8 pl-8 pr-3 w-52 text-sm",
                            "border-slate-200/80 bg-white placeholder:text-slate-400",
                            "dark:border-white/[0.08] dark:bg-white/[0.03] dark:placeholder:text-slate-600",
                            "focus-visible:ring-1 focus-visible:ring-green-500/40 focus-visible:border-green-400/60",
                            "dark:focus-visible:ring-purple-500/40 dark:focus-visible:border-purple-400/40"
                        )}
                    />
                </div>

                {/* Refresh */}
                <Button
                    variant="outline"
                    size="sm"
                    onClick={onRefresh}
                    className={cn(
                        "h-8 gap-1.5 px-3 text-xs font-medium",
                        "border-slate-200/80 text-slate-500 bg-white",
                        "hover:text-green-600 hover:border-green-300/60 hover:bg-green-50/50",
                        "dark:border-white/[0.08] dark:bg-white/[0.02] dark:text-slate-400",
                        "dark:hover:text-purple-400 dark:hover:border-purple-500/30 dark:hover:bg-purple-500/5"
                    )}
                >
                    <RefreshCw className="h-3.5 w-3.5" />
                    Refresh
                </Button>
            </div>
        </div>
    );
}