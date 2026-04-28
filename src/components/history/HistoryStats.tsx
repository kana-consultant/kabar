import { Card, CardContent } from "@/components/ui/card";
import { History, CheckCircle, XCircle, Clock } from "lucide-react";
import type { HistoryItem } from "@/types/history";

interface HistoryStatsProps {
    history: HistoryItem[];
}

export function HistoryStats({ history }: HistoryStatsProps) {
    const stats = [
        {
            title: "Total",
            value: history.length,
            icon: History,
            color: "text-blue-500",
            bg: "bg-blue-50 dark:bg-blue-950/30",
        },
        {
            title: "Berhasil",
            value: history.filter(h => h.status === "success").length,
            icon: CheckCircle,
            color: "text-green-500",
            bg: "bg-green-50 dark:bg-green-950/30",
        },
        {
            title: "Gagal",
            value: history.filter(h => h.status === "failed").length,
            icon: XCircle,
            color: "text-red-500",
            bg: "bg-red-50 dark:bg-red-950/30",
        },
        {
            title: "Pending",
            value: history.filter(h => h.status === "pending").length,
            icon: Clock,
            color: "text-yellow-500",
            bg: "bg-yellow-50 dark:bg-yellow-950/30",
        },
    ];

    return (
        <div className="grid gap-4 md:grid-cols-4">
            {stats.map((stat) => (
                <Card key={stat.title}>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-slate-500">{stat.title}</p>
                                <p className="text-2xl font-bold mt-1">{stat.value}</p>
                            </div>
                            <div className={`rounded-full p-3 ${stat.bg}`}>
                                <stat.icon className={`h-5 w-5 ${stat.color}`} />
                            </div>
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}