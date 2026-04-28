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
        <div className="flex items-center justify-between">
            <div>
                <h2 className="text-2xl font-bold tracking-tight">Jadwal Posting</h2>
                <p className="text-slate-500">
                    Kelola konten yang dijadwalkan untuk dipublikasikan
                </p>
            </div>
            <div className="flex gap-2">
                <div className="relative">
                    <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
                    <Input
                        placeholder="Cari jadwal..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-8 w-64"
                    />
                </div>
                <Button variant="outline" onClick={onRefresh}>
                    <RefreshCw className="mr-2 h-4 w-4" />
                    Refresh
                </Button>
            </div>
        </div>
    );
}