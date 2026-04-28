import { Card, CardContent } from "@/components/ui/card";
import { Calendar, Clock, RefreshCw } from "lucide-react";
import type { Draft } from "@/types/draft";

interface ScheduleStatsProps {
    schedules: Draft[];
    isDailySchedule: (scheduledFor?: string) => boolean;
}

export function ScheduleStats({ schedules, isDailySchedule }: ScheduleStatsProps) {
    const total = schedules.length;
    const oneTime = schedules.filter(s => !isDailySchedule(s.scheduledFor)).length;
    const daily = schedules.filter(s => isDailySchedule(s.scheduledFor)).length;

    return (
        <div className="grid gap-4 md:grid-cols-3">
            <Card>
                <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-slate-500">Total Jadwal</p>
                            <p className="text-2xl font-bold">{total}</p>
                        </div>
                        <Calendar className="h-8 w-8 text-blue-500" />
                    </div>
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-slate-500">One Time</p>
                            <p className="text-2xl font-bold">{oneTime}</p>
                        </div>
                        <Clock className="h-8 w-8 text-green-500" />
                    </div>
                </CardContent>
            </Card>
            <Card>
                <CardContent className="p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm text-slate-500">Daily Schedule</p>
                            <p className="text-2xl font-bold">{daily}</p>
                        </div>
                        <RefreshCw className="h-8 w-8 text-purple-500" />
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}