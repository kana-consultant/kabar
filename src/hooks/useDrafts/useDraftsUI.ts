import { useState, useCallback } from "react";
import type { Draft } from "@/types/draft";
import type { DraftsUIState, PublishResult, ScheduleConfig } from "./types";

const DEFAULT_SCHEDULE_CONFIG: ScheduleConfig = {
    date: "",
    time: "09:00",
    dailySchedule: false,
    dailyTime: "08:00",
};

export function useDraftsUI() {
    const [selectedDraft, setSelectedDraft] = useState<Draft | null>(null);
    const [showScheduleDialog, setShowScheduleDialog] = useState(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [showResultDialog, setShowResultDialog] = useState(false);
    const [publishResults, setPublishResults] = useState<PublishResult | null>(null);
    const [publishingId, setPublishingId] = useState<string | null>(null);
    const [scheduleConfig, setScheduleConfig] = useState<ScheduleConfig>(DEFAULT_SCHEDULE_CONFIG);

    const updateScheduleConfig = useCallback((updates: Partial<ScheduleConfig>) => {
        setScheduleConfig(prev => ({ ...prev, ...updates }));
    }, []);

    const resetScheduleConfig = useCallback(() => {
        setScheduleConfig(DEFAULT_SCHEDULE_CONFIG);
    }, []);

    const openScheduleDialog = useCallback((draft: Draft) => {
        setSelectedDraft(draft);
        setShowScheduleDialog(true);
        resetScheduleConfig();
    }, [resetScheduleConfig]);

    const openDeleteDialog = useCallback((draft: Draft) => {
        setSelectedDraft(draft);
        setShowDeleteDialog(true);
    }, []);

    const closeDialogs = useCallback(() => {
        setShowScheduleDialog(false);
        setShowDeleteDialog(false);
        setShowResultDialog(false);
        setSelectedDraft(null);
        setPublishResults(null);
        resetScheduleConfig();
    }, [resetScheduleConfig]);

    const formatDate = useCallback((dateString: string) => {
        return new Date(dateString).toLocaleDateString("id-ID", {
            day: "numeric",
            month: "long",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }, []);

    return {
        // State
        selectedDraft,
        showScheduleDialog,
        showDeleteDialog,
        showResultDialog,
        publishResults,
        publishingId,
        scheduleConfig,
        
        // Setters (for granular updates)
        setPublishingId,
        setPublishResults,
        setShowResultDialog,
        setScheduleConfig: updateScheduleConfig,
        
        // Actions
        openScheduleDialog,
        openDeleteDialog,
        closeDialogs,
        formatDate,
        resetScheduleConfig,
    };
}