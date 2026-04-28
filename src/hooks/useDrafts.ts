import { useState, useEffect } from "react";
import { toast } from "sonner";
import { getDrafts, deleteDraft, publishDraft } from "@/services/draftService";
import type { Draft } from "@/types/draft";

export function useDrafts() {
    const [drafts, setDrafts] = useState<Draft[]>([]);
    const [filteredDrafts, setFilteredDrafts] = useState<Draft[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [loading, setLoading] = useState(true);
    const [publishingId, setPublishingId] = useState<string | null>(null); // track which draft is publishing
    const [statusFilter, setStatusFilter] = useState<"all" | "draft" | "scheduled" | "published">("all");
    const [selectedDraft, setSelectedDraft] = useState<Draft | null>(null);
    const [showScheduleDialog, setShowScheduleDialog] = useState(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [showResultDialog, setShowResultDialog] = useState(false);
    const [publishResults, setPublishResults] = useState<{ success: boolean; message: string; results?: any[] } | null>(null);
    const [scheduleDate, setScheduleDate] = useState("");
    const [scheduleTime, setScheduleTime] = useState("09:00");
    const [dailySchedule, setDailySchedule] = useState(false);
    const [dailyTime, setDailyTime] = useState("08:00");

    const loadDrafts = async () => {
        setLoading(true);
        try {
            const allDrafts = await getDrafts();
            setDrafts(allDrafts ?? []);
        } catch (error) {
            console.error("Failed to load drafts:", error);
            toast.error("Gagal memuat draft");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadDrafts();
    }, []);

    useEffect(() => {
        let filtered = [...drafts];
        
        if (statusFilter !== "all") {
            filtered = filtered.filter(d => d.status === statusFilter);
        }
        
        if (searchQuery) {
            const query = searchQuery.toLowerCase();
            filtered = filtered.filter(d =>
                d.title.toLowerCase().includes(query) ||
                d.topic.toLowerCase().includes(query)
            );
        }
        
        setFilteredDrafts(filtered || []);
    }, [drafts, searchQuery, statusFilter]);

    const handleDelete = async () => {
        if (selectedDraft) {
            try {
                await deleteDraft(selectedDraft.id);
                toast.success("Draft dihapus", {
                    description: `"${selectedDraft.title}" telah dihapus`,
                });
                await loadDrafts();
                setShowDeleteDialog(false);
                setSelectedDraft(null);
            } catch (error) {
                toast.error("Gagal menghapus draft");
            }
        }
    };

    const handleSchedule = async () => {
        if (selectedDraft) {
            let scheduledFor: string | undefined;
            
            if (dailySchedule) {
                scheduledFor = `daily:${dailyTime}`;
            } else {
                if (!scheduleDate) {
                    toast.error("Pilih tanggal jadwal");
                    return;
                }
                scheduledFor = `${scheduleDate}T${scheduleTime}:00`;
            }
            
            setPublishingId(selectedDraft.id);
            
            try {
                const response = await publishDraft(selectedDraft.id, scheduledFor);
                
                // Check if response contains results (multi-product publish)
                if (response && typeof response === 'object') {
                    const hasErrors = response.results?.some((r: any) => !r.success);
                    
                    if (hasErrors) {
                        // Show detailed error dialog
                        setPublishResults({
                            success: false,
                            message: response.message || "Sebagian produk gagal dipublikasikan",
                            results: response.results
                        });
                        setShowResultDialog(true);
                        toast.error("Publikasi sebagian gagal", {
                            description: "Beberapa produk tidak dapat dijangkau. Lihat detail untuk info lebih lanjut."
                        });
                    } else {
                        // All success
                        toast.success("Draft dijadwalkan", {
                            description: dailySchedule
                                ? `"${selectedDraft.title}" akan diposting setiap hari jam ${dailyTime}`
                                : `"${selectedDraft.title}" akan diposting pada ${scheduleDate} jam ${scheduleTime}`,
                        });
                        await loadDrafts();
                        setShowScheduleDialog(false);
                        setSelectedDraft(null);
                        setScheduleDate("");
                        setScheduleTime("09:00");
                        setDailySchedule(false);
                        setDailyTime("08:00");
                    }
                } else {
                    // Old response format (just success)
                    toast.success("Draft dijadwalkan", {
                        description: dailySchedule
                            ? `"${selectedDraft.title}" akan diposting setiap hari jam ${dailyTime}`
                            : `"${selectedDraft.title}" akan diposting pada ${scheduleDate} jam ${scheduleTime}`,
                    });
                    await loadDrafts();
                    setShowScheduleDialog(false);
                    setSelectedDraft(null);
                    setScheduleDate("");
                    setScheduleTime("09:00");
                    setDailySchedule(false);
                    setDailyTime("08:00");
                }
            } catch (error: any) {
                console.error("Schedule error:", error);
                const errorMessage = error?.response?.data?.message || error?.message || "Gagal menjadwalkan draft";
                toast.error("Gagal menjadwalkan draft", {
                    description: errorMessage,
                });
                
                // Try to extract results from error if available
                if (error?.response?.data?.results) {
                    setPublishResults({
                        success: false,
                        message: error.response.data.message || "Publikasi gagal",
                        results: error.response.data.results
                    });
                    setShowResultDialog(true);
                }
            } finally {
                setPublishingId(null);
            }
        }
    };

    const handlePublishNow = async (draft: Draft) => {
        setPublishingId(draft.id);
        
        try {
            const response = await publishDraft(draft.id);
            
            // Check if response contains results (multi-product publish)
            if (response && typeof response === 'object') {
                const hasErrors = response.results?.some((r: any) => !r.success);
                
                if (hasErrors) {
                    setPublishResults({
                        success: false,
                        message: response.message || "Sebagian produk gagal dipublikasikan",
                        results: response.results
                    });
                    setShowResultDialog(true);
                    toast.error("Publikasi sebagian gagal", {
                        description: "Beberapa produk tidak dapat dijangkau. Lihat detail untuk info lebih lanjut."
                    });
                } else {
                    toast.success("Draft dipublikasikan!", {
                        description: `"${draft.title}" telah dipublikasikan`,
                    });
                    await loadDrafts();
                }
            } else {
                toast.success("Draft dipublikasikan!", {
                    description: `"${draft.title}" telah dipublikasikan`,
                });
                await loadDrafts();
            }
        } catch (error: any) {
            console.error("Publish error:", error);
            const errorMessage = error?.response?.data?.message || error?.message || "Gagal mempublikasikan draft";
            toast.error("Gagal mempublikasikan draft", {
                description: errorMessage,
            });
            
            if (error?.response?.data?.results) {
                setPublishResults({
                    success: false,
                    message: error.response.data.message || "Publikasi gagal",
                    results: error.response.data.results
                });
                setShowResultDialog(true);
            }
        } finally {
            setPublishingId(null);
        }
    };

    const openScheduleDialog = (draft: Draft) => {
        setSelectedDraft(draft);
        setShowScheduleDialog(true);
    };

    const openDeleteDialog = (draft: Draft) => {
        setSelectedDraft(draft);
        setShowDeleteDialog(true);
    };

    const closeDialogs = () => {
        setShowScheduleDialog(false);
        setShowDeleteDialog(false);
        setShowResultDialog(false);
        setSelectedDraft(null);
        setPublishResults(null);
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

    return {
        drafts,
        filteredDrafts,
        searchQuery,
        setSearchQuery,
        loading,
        publishingId,
        statusFilter,
        setStatusFilter,
        selectedDraft,
        setSelectedDraft,
        showScheduleDialog,
        showDeleteDialog,
        showResultDialog,
        publishResults,
        scheduleDate,
        setScheduleDate,
        scheduleTime,
        setScheduleTime,
        dailySchedule,
        setDailySchedule,
        dailyTime,
        setDailyTime,
        loadDrafts,
        handleDelete,
        handleSchedule,
        handlePublishNow,
        openScheduleDialog,
        openDeleteDialog,
        closeDialogs,
        formatDate,
    };
}