import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Search, FileText, Calendar, CheckCircle } from "lucide-react";

interface DraftsHeaderProps {
    searchQuery: string;
    setSearchQuery: (value: string) => void;
    statusFilter: "all" | "draft" | "scheduled" | "published";
    setStatusFilter: (value: "all" | "draft" | "scheduled" | "published") => void;
}

export function DraftsHeader({
    searchQuery,
    setSearchQuery,
    statusFilter,
    setStatusFilter,
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
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="pl-8 w-full sm:w-64"
                    />
                </div>

                <Select
                    value={statusFilter}
                    onValueChange={(value) =>
                        setStatusFilter(value as "all" | "draft" | "scheduled" | "published")
                    }
                >
                    <SelectTrigger className="w-full sm:w-36">
                        <SelectValue placeholder="Semua" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="all">
                            <div className="flex items-center gap-2">
                                <span>Semua</span>
                            </div>
                        </SelectItem>
                        <SelectItem value="draft">
                            <div className="flex items-center gap-2">
                                <FileText className="h-4 w-4" />
                                <span>Draft</span>
                            </div>
                        </SelectItem>
                        <SelectItem value="published">
                            <div className="flex items-center gap-2">
                                <CheckCircle className="h-4 w-4" />
                                <span>Terbit</span>
                            </div>
                        </SelectItem>
                    </SelectContent>
                </Select>
            </div>
        </div>
    );
}