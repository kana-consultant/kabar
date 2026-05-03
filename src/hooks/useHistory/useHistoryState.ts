import { useState } from "react";
import type { HistoryItem } from "@/types/history";
import type { StatusFilter, ActionFilter } from "./types";

export function useHistoryState() {
    const [history, setHistory] = useState<HistoryItem[]>([]);
    const [filteredHistory, setFilteredHistory] = useState<HistoryItem[]>([]);
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
    const [actionFilter, setActionFilter] = useState<ActionFilter>("all");
    const [loading, setLoading] = useState<boolean>(true);
    const [selectedHistory, setSelectedHistory] = useState<HistoryItem | null>(null);
    const [showDetailDialog, setShowDetailDialog] = useState<boolean>(false);
    const [showDeleteDialog, setShowDeleteDialog] = useState<boolean>(false);

    return {
        history, setHistory,
        filteredHistory, setFilteredHistory,
        searchQuery, setSearchQuery,
        statusFilter, setStatusFilter,
        actionFilter, setActionFilter,
        loading, setLoading,
        selectedHistory, setSelectedHistory,
        showDetailDialog, setShowDetailDialog,
        showDeleteDialog, setShowDeleteDialog,
    };
}