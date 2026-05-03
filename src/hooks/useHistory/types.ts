import type { HistoryItem } from "@/types/history";

export type StatusFilter = "all" | "success" | "failed" | "pending";
export type ActionFilter = "all" | "published" | "scheduled" | "draft_saved";

export interface StatusData {
    label: string;
    icon: string;
    color: string;
}

export interface ActionData {
    label: string;
    icon: string;
}