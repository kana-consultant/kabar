import { Input } from "@kana-consultant/ui-kit";
import { Search } from "lucide-react";

interface DraftsHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
}

export function DraftsHeader({
    searchQuery,
    setSearchQuery,
}: DraftsHeaderProps) {
    return (
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="space-y-1">
                <h2 className="text-2xl font-bold tracking-tight">Draft & Terjadwal</h2>
                <p className="text-sm text-muted-foreground">
                    Kelola draft, jadwalkan posting, atau publikasikan langsung
                </p>
            </div>

            <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
                <div className="relative">
                    <Search className="absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                        placeholder="Cari draft..."
                        value={searchQuery}
                        onChange={(e : any) => setSearchQuery(e.target.value)}
                        className="pl-8 w-full sm:w-64"
                    />
                </div>
            </div>
        </div>
    );
}