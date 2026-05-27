import { useCallback } from "react";
import { useToast } from "@/hooks/use-toast";
import { deleteDraft, publishDraft,checkSimilarity ,getSeoScore, type Draft } from "@/services/draft";
import type { PublishResult, ScheduleConfig } from "./types";

interface UseDraftsActionsParams {
    loadDrafts: () => Promise<void>;
    setPublishingId: (id: string | null) => void;
    setPublishResults: (results: PublishResult | null) => void;
    setShowResultDialog: (show: boolean) => void;
    closeDialogs: () => void;
}

export function useDraftsActions({
    loadDrafts,
    setPublishingId,
    setPublishResults,
    setShowResultDialog,
    closeDialogs,
}: UseDraftsActionsParams) {
    const toast = useToast(); //   Gunakan hook toast

    const handleDelete = useCallback(async (draft: Draft) => {
        try {
            await deleteDraft(draft.id as string);
            toast.success("Draft dihapus", {
                description: `"${draft.title}" telah dihapus`,
            });
            await loadDrafts();
            return true;
        } catch (error) {
            toast.error("Gagal menghapus draft");
            return false;
        }
    }, [loadDrafts, toast]);

    const processPublishResponse = useCallback((
        response: any,
        draft: Draft,
        isScheduled: boolean,
        scheduleConfig?: ScheduleConfig
    ) => {
        if (response && typeof response === 'object' && response.results) {
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
                return false;
            }
        }

        // Success Toast message
        if (isScheduled && scheduleConfig) {
            const { dailySchedule, dailyTime, date, time } = scheduleConfig;
            toast.success("Draft dijadwalkan", {
                description: dailySchedule
                    ? `"${draft.title}" akan diposting setiap hari jam ${dailyTime}`
                    : `"${draft.title}" akan diposting pada ${date} jam ${time}`,
            });
        } else {
            toast.success("Draft dipublikasikan!", {
                description: `"${draft.title}" telah dipublikasikan`,
            });
        }

        return true;
    }, [setPublishResults, setShowResultDialog, toast]);

    const handlePublish = useCallback(async (
        draft: Draft,
        scheduledFor?: string,
        scheduleConfig?: ScheduleConfig
    ) => {
        setPublishingId(draft.id as string);

        try {
            const response = await publishDraft(draft.id as string, draft);
            const isScheduled = !!scheduledFor;

            const success = processPublishResponse(response, draft, isScheduled, scheduleConfig);

            if (success) {
                await loadDrafts();
                if (isScheduled) {
                    closeDialogs();
                }
            }

            return success;
        } catch (error: any) {
            console.error("Publish error:", error);
            const errorMessage = error?.response?.data?.message || error?.message ||
                (scheduledFor ? "Gagal menjadwalkan draft" : "Gagal mempublikasikan draft");

            toast.error(scheduledFor ? "Gagal menjadwalkan draft" : "Gagal mempublikasikan draft", {
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

            return false;
        } finally {
            setPublishingId(null);
        }
    }, [setPublishingId, processPublishResponse, loadDrafts, closeDialogs, setPublishResults, setShowResultDialog, toast]);

    const handlePublishNow = useCallback((draft: Draft) => {
        return handlePublish(draft);
    }, [handlePublish]);

    const handleSchedule = useCallback((draft: Draft, scheduleConfig: ScheduleConfig) => {
        const { date, time, dailySchedule, dailyTime: dailyTimeValue } = scheduleConfig;

        let scheduledFor: string | undefined;

        if (dailySchedule) {
            scheduledFor = `daily:${dailyTimeValue}`;
        } else {
            if (!date) {
                toast.error("Pilih tanggal jadwal");
                return Promise.resolve(false);
            }
            scheduledFor = `${date}T${time}:00`;
        }

        console.log(draft)

        return handlePublish(draft, scheduledFor, scheduleConfig);
    }, [handlePublish, toast]);


    const handlecheckSimilarity = useCallback(async(draft : Draft)=>{
        return await checkSimilarity(draft.id as string)
    },[checkSimilarity])

      const handlegetSeoScore = useCallback(async(draft : Draft)=>{
        return await getSeoScore(draft.id as string)
    },[getSeoScore])
    
    return {
        handleDelete,
        handlePublishNow,
        handleSchedule,
        handlecheckSimilarity,
        handlegetSeoScore
    };

}