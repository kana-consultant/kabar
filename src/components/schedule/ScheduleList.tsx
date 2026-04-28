import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar } from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import { ScheduleItem } from "./ScheduleItem";
import type { Draft } from "@/types/draft";

interface ScheduleListProps {
    schedules: Draft[];
    isDailySchedule: (scheduledFor?: string) => boolean;
    getScheduleDisplay: (scheduledFor?: string) => string;
    onView: (schedule: Draft) => void;
    onEdit: (schedule: Draft) => void;
    onReschedule: (schedule: Draft) => void;
    onPublishNow: (schedule: Draft) => void;
    onDelete: (schedule: Draft) => void;
}

export function ScheduleList({
    schedules,
    isDailySchedule,
    getScheduleDisplay,
    onView,
    onEdit,
    onReschedule,
    onPublishNow,
    onDelete,
}: ScheduleListProps) {
    const navigate = useNavigate();

    if (schedules.length === 0) {
        return (
            <Card>
                <CardContent className="py-12 text-center">
                    <Calendar className="mx-auto h-12 w-12 text-slate-300" />
                    <p className="mt-2 text-slate-500">Belum ada jadwal posting</p>
                    <Button
                        variant="outline"
                        className="mt-4"
                        onClick={() => navigate({ to: "/generate" })}
                    >
                        Buat Jadwal Baru
                    </Button>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Daftar Jadwal</CardTitle>
                <CardDescription>{schedules.length} jadwal ditemukan</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {schedules.map((schedule) => (
                        <ScheduleItem
                            key={schedule.id}
                            schedule={schedule}
                            isDailySchedule={isDailySchedule}
                            getScheduleDisplay={getScheduleDisplay}
                            onView={onView}
                            onEdit={onEdit}
                            onReschedule={onReschedule}
                            onPublishNow={onPublishNow}
                            onDelete={onDelete}
                        />
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}