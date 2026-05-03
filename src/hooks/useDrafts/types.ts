import type { Draft } from "@/types/draft";

export type StatusFilter = "all" | "draft" | "scheduled" | "published";

export interface PublishResult {
    success: boolean;
    message: string;
    results?: any[];
}

export interface ScheduleConfig {
    date: string;
    time: string;
    dailySchedule: boolean;
    dailyTime: string;
}

export interface DraftsUIState {
    selectedDraft: Draft | null;
    showScheduleDialog: boolean;
    showDeleteDialog: boolean;
    showResultDialog: boolean;
    publishResults: PublishResult | null;
    publishingId: string | null;
}

export interface DraftsFilterState {
    searchQuery: string;
    statusFilter: StatusFilter;
}