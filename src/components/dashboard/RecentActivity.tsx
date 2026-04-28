import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { CheckCircle2, Clock, Send } from "lucide-react";

const activities = [
    {
        id: 1,
        title: "Cara Memilih Sepatu Gunung",
        product: "TrekkingID",
        status: "published",
        time: "2 jam lalu",
    },
    {
        id: 2,
        title: "Tips Camping Hemat",
        product: "CampingMart",
        status: "draft",
        time: "5 jam lalu",
    },
    {
        id: 3,
        title: "Perawatan Jaket Outdoor",
        product: "OutdoorGear",
        status: "published",
        time: "kemarin",
    },
    {
        id: 4,
        title: "Rekomendasi Tenda 4 Musim",
        product: "TrekkingID",
        status: "processing",
        time: "kemarin",
    },
];

import { useHistory } from "@/hooks/useHistory";

const statusConfig = {
    published: { icon: CheckCircle2, color: "text-green-600", label: "Terbit" },
    draft: { icon: Clock, color: "text-yellow-600", label: "Draft" },
    processing: { icon: Send, color: "text-blue-600", label: "Diproses" },
};

export interface HistoryItem {
    id: string;
    title: string;
    topic: string;
    content: string;
    imageUrl?: string;
    targetProducts: string[];
    status: 'success' | 'failed' | 'pending';
    action: 'published' | 'scheduled' | 'draft_saved';
    errorMessage?: string;
    publishedAt: string;
    scheduledFor?: string;
    createdBy?: string;
    teamId?: string;
    createdAt: string;
}

export function RecentActivity() {
    const {
        history
    } = useHistory()
    return (
        <Card>
            <CardHeader>
                <CardTitle>Aktivitas Terbaru</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {history.slice(0, 5).map((activity) => {
                        const config =
                            statusConfig[activity.status as keyof typeof statusConfig];

                        const Icon = config.icon;

                        return (
                            <div
                                key={activity.id}
                                className="flex items-center justify-between border-b pb-3 last:border-0"
                            >
                                <div className="space-y-1">
                                    <p className="text-sm font-medium">
                                        {activity.title}
                                    </p>
                                    <p className="text-xs text-slate-500">
                                        {activity.status}
                                    </p>
                                </div>

                                <div className="flex items-center gap-2">
                                    <Icon className={`h-4 w-4 ${config.color}`} />
                                    <span className="text-xs text-slate-500">
                                        {activity.publishedAt}
                                    </span>
                                </div>
                            </div>
                        );
                    })}
                </div>
            </CardContent>
        </Card>
    );
}