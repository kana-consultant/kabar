import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search, RefreshCw } from "lucide-react";

interface ScheduleHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
    onRefresh: () => void;
}

export function ScheduleHeader({ searchQuery, setSearchQuery, onRefresh }: ScheduleHeaderProps) {
    return (
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
                <h2 className="text-xl font-bold tracking-tight sm:text-2xl">Jadwal Posting</h2>
                <p className="text-sm text-slate-500 sm:text-base">
                    Kelola konten yang dijadwalkan untuk dipublikasikan
                </p>
            </div>
            
            <div className="flex flex-col gap-2 sm:flex-row sm:gap-2">
                <div className="relative w-full sm:w-auto">
                    <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
                    <Input
                        placeholder="Cari jadwal..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-8 w-full sm:w-64"
                    />
                </div>
                <Button variant="outline" onClick={onRefresh} className="w-full sm:w-auto">
                    <RefreshCw className="mr-2 h-4 w-4" />
                    <span>Refresh</span>
                </Button>
            </div>
        </div>
    );
}