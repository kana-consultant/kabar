import { Card, CardContent } from "@/components/ui/card";
import { FileText, AlertCircle, Calendar, CheckCircle } from "lucide-react";
import type { Draft } from "@/services/draft";

interface DraftStatsProps {
    drafts: Draft[];
}

export function DraftStats({ drafts }: DraftStatsProps) {
    const stats = [
        {
            title: "Total Draft",
            value: drafts.length,
            icon: FileText,
            color: "text-blue-500",
        },
        {
            title: "Belum Terbit",
            value: drafts.filter(d => d.status === "draft").length,
            icon: AlertCircle,
            color: "text-yellow-500",
        },
        {
            title: "Sudah Terbit",
            value: drafts.filter(d => d.status === "published").length,
            icon: CheckCircle,
            color: "text-green-500",
        },
    ];

    return (
        <div className="grid gap-4 md:grid-cols-3">
            {stats.map((stat) => (
                <Card key={stat.title}>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-slate-500">{stat.title}</p>
                                <p className="text-2xl font-bold">{stat.value}</p>
                            </div>
                            <stat.icon className={`h-8 w-8 ${stat.color}`} />
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}