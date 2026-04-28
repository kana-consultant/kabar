import { useState, useEffect } from "react";
import { toast } from "sonner";
import { getDrafts, publishDraft, deleteDraft, updateDraft } from "@/services/draftService";
import type { Draft } from "@/types/draft";

export function useSchedule() {
    const [schedules, setSchedules] = useState<Draft[]>([]);
    const [filteredSchedules, setFilteredSchedules] = useState<Draft[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [loading, setLoading] = useState(true); // ← tambah loading state
    const [selectedSchedule, setSelectedSchedule] = useState<Draft | null>(null);
    const [showDetailDialog, setShowDetailDialog] = useState(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [showRescheduleDialog, setShowRescheduleDialog] = useState(false);
    const [newScheduleDate, setNewScheduleDate] = useState("");
    const [newScheduleTime, setNewScheduleTime] = useState("09:00");

    const loadSchedules = async () => {
        setLoading(true);
        try {
            const allDrafts = await getDrafts();
            const scheduled = allDrafts ? allDrafts.filter(d => d.status === "scheduled") : [];
            setSchedules(scheduled || []);
        } catch (error) {
            console.error("Failed to load schedules:", error);
            toast.error("Gagal memuat jadwal");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadSchedules();
    }, []);

    useEffect(() => {
        let filtered = [...schedules];
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(s =>
                s.title.toLowerCase().includes(query) ||
                s.topic.toLowerCase().includes(query)
            );
        }
        setFilteredSchedules(filtered);
    }, [schedules, searchQuery]);

    const handlePublishNow = async (schedule: Draft) => {
        try {
            await publishDraft(schedule.id);
            toast.success("Konten dipublikasikan!", {
                description: `"${schedule.title}" telah dipublikasikan`,
            });
            await loadSchedules();
        } catch (error) {
            toast.error("Gagal mempublikasikan");
        }
    };

    const handleDelete = async () => {
        if (selectedSchedule) {
            try {
                await deleteDraft(selectedSchedule.id);
                toast.success("Jadwal dihapus", {
                    description: `"${selectedSchedule.title}" telah dihapus dari jadwal`,
                });
                await loadSchedules();
                setShowDeleteDialog(false);
                setSelectedSchedule(null);
            } catch (error) {
                toast.error("Gagal menghapus jadwal");
            }
        }
    };

    const handleReschedule = async () => {
        if (selectedSchedule && newScheduleDate && newScheduleTime) {
            const dateTime = `${newScheduleDate}T${newScheduleTime}`;
            try {
                await updateDraft(selectedSchedule.id, { scheduledFor: dateTime });
                toast.success("Jadwal diperbarui", {
                    description: `"${selectedSchedule.title}" dijadwalkan ulang pada ${newScheduleDate} jam ${newScheduleTime}`,
                });
                await loadSchedules();
                setShowRescheduleDialog(false);
                setSelectedSchedule(null);
                setNewScheduleDate("");
                setNewScheduleTime("09:00");
            } catch (error) {
                toast.error("Gagal memperbarui jadwal");
            }
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    const getScheduleDisplay = (scheduledFor?: string) => {
        if (!scheduledFor) return "Tidak terjadwal";
        if (scheduledFor.startsWith("daily:")) {
            return `Setiap hari jam ${scheduledFor.replace("daily:", "")}`;
        }
        return formatDate(scheduledFor);
    };

    const isDailySchedule = (scheduledFor?: string) => {
        return scheduledFor?.startsWith("daily:") || false;
    };

    return {
        schedules,
        filteredSchedules,
        searchQuery,
        setSearchQuery,
        loading, // ← export loading
        selectedSchedule,
        setSelectedSchedule,
        showDetailDialog,
        setShowDetailDialog,
        showDeleteDialog,
        setShowDeleteDialog,
        showRescheduleDialog,
        setShowRescheduleDialog,
        newScheduleDate,
        setNewScheduleDate,
        newScheduleTime,
        setNewScheduleTime,
        loadSchedules,
        handlePublishNow,
        handleDelete,
        handleReschedule,
        formatDate,
        getScheduleDisplay,
        isDailySchedule,
    };
}