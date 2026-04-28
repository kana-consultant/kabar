import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { HistoryItem } from "./HistoryItem";
import { HistoryIcon } from "lucide-react";
import type { HistoryItem as HistoryItemType } from "@/types/history";

interface HistoryListProps {
    items: HistoryItemType[];
    onView: (item: HistoryItemType) => void;
    onRepost: (item: HistoryItemType) => void;
    onDelete: (item: HistoryItemType) => void;
    formatDate: (date: string) => string;
    getStatusData: (status: string) => { label: string; icon: string; color: string };
    getActionData: (action: string) => { label: string; icon: string };
}

export function HistoryList({
    items,
    onView,
    onRepost,
    onDelete,
    formatDate,
    getStatusData,
    getActionData,
}: HistoryListProps) {
    if (items.length === 0) {
        return (
            <Card>
                <CardContent className="py-12 text-center">
                    <HistoryIcon className="mx-auto h-12 w-12 text-slate-300" />
                    <p className="mt-2 text-slate-500">Belum ada riwayat</p>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Riwayat Konten</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {items.map((item) => (
                        <HistoryItem
                            key={item.id}
                            item={item}
                            onView={onView}
                            onRepost={onRepost}
                            onDelete={onDelete}
                            formatDate={formatDate}
                            getStatusData={getStatusData}
                            getActionData={getActionData}
                        />
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}