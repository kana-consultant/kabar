import { useCallback } from "react";
import { useDraftsData } from "./useDraftsData";
import { useDraftsUI } from "./useDraftsUI";
import { useDraftsActions } from "./useDraftsActions";
import type { Draft } from "@/types/draft";

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
    } = useDraftsData();

    // UI state
    const {
        selectedDraft,
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
    } = useDraftsActions({
        loadDrafts,
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
    };
}

// Re-export types for convenience
export type { StatusFilter, PublishResult, ScheduleConfig } from "./types";