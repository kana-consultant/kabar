import { useCallback } from "react";
import { useDraftsData } from "./useDraftsData";
import { useDraftsUI } from "./useDraftsUI";
import { useDraftsActions } from "./useDraftsActions";
import { type Draft } from "@//services/draft";

export function useDrafts() {
    // Data & filtering
    const {
        drafts,
        filteredDrafts,
        loading,
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,
        loadDrafts,
        currentPage, setCurrentPage,
        totalPages,
        totalItems,
        stats,
        seoDialog, setSeoDialog,
        similarityData, setSimilarityData,
        similarityDialog, setSimilarityDialog,
        seoLoading, setSeoLoading,
        similarityLoading, setSimilarityLoading,
        seoData, setSeoData,

    } = useDraftsData();

    // UI state
    const {
        selectedDraft,
        setSelectedDraft,
        showScheduleDialog,
        showDeleteDialog,
        showResultDialog,
        publishResults,
        publishingId,
        scheduleConfig,
        setPublishingId,
        setPublishResults,
        setShowResultDialog,
        openScheduleDialog,
        openDeleteDialog,
        closeDialogs,
        formatDate,
        setScheduleConfig,
    } = useDraftsUI();

    // Actions
    const {
        handleDelete,
        handlePublishNow,
        handleSchedule,
        handlecheckSimilarity,
        handlegetSeoScore,
    } = useDraftsActions({
        loadDrafts: () => loadDrafts(currentPage),
        setPublishingId,
        setPublishResults,
        setShowResultDialog,
        closeDialogs,

    });

    // Wrapper untuk delete dengan draft parameter
    const onDelete = useCallback(async () => {
        if (selectedDraft) {
            const success = await handleDelete(selectedDraft);
            if (success) {
                closeDialogs();
            }
        }
    }, [selectedDraft, handleDelete, closeDialogs]);

    const checkSimilarity = useCallback(async (draft: Draft) => {
        setSimilarityDialog({ open: true, draft });
        setSimilarityLoading(true);
        setSimilarityData(null);
        try {
            const result = await handlecheckSimilarity(draft);
            setSimilarityData(result?.similar_drafts || []);
        } finally {
            setSimilarityLoading(false);
        }
    }, [handlecheckSimilarity, setSimilarityDialog, setSimilarityLoading, setSimilarityData]);

    const getSeoScore = useCallback(async (draft: Draft) => {
        setSeoDialog({ open: true, draft });
        setSeoLoading(true);
        setSeoData(null);
        try {
            const result = await handlegetSeoScore(draft);
            setSeoData(result);
        } finally {
            setSeoLoading(false);
        }
    }, [handlegetSeoScore, setSeoDialog, setSeoLoading, setSeoData]);
    // Wrapper untuk schedule
    const onSchedule = useCallback(async () => {
        if (selectedDraft) {
            await handleSchedule(selectedDraft, scheduleConfig);
        }
    }, [selectedDraft, handleSchedule, scheduleConfig]);

    // Wrapper untuk publish now
    const onPublishNow = useCallback((draft: Draft) => {
        return handlePublishNow(draft);
    }, [handlePublishNow]);

    return {
        // Data
        drafts,
        filteredDrafts,
        loading,

        // Filters
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,

        // UI State
        selectedDraft,
        showScheduleDialog,
        showDeleteDialog,
        showResultDialog,
        publishResults,
        publishingId,
        scheduleConfig,
        seoDialog, setSeoDialog,
        similarityData, setSimilarityData,
        similarityDialog, setSimilarityDialog,
        seoLoading, setSeoLoading,
        similarityLoading, setSimilarityLoading,
        seoData, setSeoData,


        // Schedule config getters/setters (convenience)
        scheduleDate: scheduleConfig.date,
        setScheduleDate: (date: string) => setScheduleConfig({ date }),
        scheduleTime: scheduleConfig.time,
        setScheduleTime: (time: string) => setScheduleConfig({ time }),
        dailySchedule: scheduleConfig.dailySchedule,
        setDailySchedule: (dailySchedule: boolean) => setScheduleConfig({ dailySchedule }),
        dailyTime: scheduleConfig.dailyTime,
        setDailyTime: (dailyTime: string) => setScheduleConfig({ dailyTime }),

        // Actions
        loadDrafts,
        handleDelete: onDelete,
        handleSchedule: onSchedule,
        handlePublishNow: onPublishNow,
        openScheduleDialog,
        openDeleteDialog,
        closeDialogs,
        formatDate,
        setSelectedDraft,
        currentPage, setCurrentPage,
        totalPages,
        totalItems,
        stats,
        checkSimilarity,
        getSeoScore
    };
}

// Re-export types for convenience
export type { StatusFilter, PublishResult, ScheduleConfig } from "./types";