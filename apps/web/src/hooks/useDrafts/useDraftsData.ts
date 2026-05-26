import { useState, useEffect, useCallback } from "react";
import { useToast } from "@/hooks/use-toast";
import { getDrafts } from "@/services/draft";
import type { Draft, DraftStats } from "@/services/draft";
import type { StatusFilter } from "./types";

export function useDraftsData() {
    const toast = useToast();
    const [drafts, setDrafts] = useState<Draft[]>([]);
    const [stats, setStats] = useState<DraftStats | null>(null);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [totalItems, setTotalItems] = useState(0);
    const LIMIT = 5;

    const loadDrafts = useCallback(async (page: number = 1) => {
        setLoading(true);
        try {
            const response = await getDrafts({
                offset: (page - 1) * LIMIT,
                limit: LIMIT,
                search: searchQuery || undefined,
                status: statusFilter !== "all" ? statusFilter : undefined,
            });

            setDrafts(response.drafts.data ?? []);
            setTotalPages(response.drafts.total_pages);
            setTotalItems(response.drafts.total_items);
            setCurrentPage(response.drafts.current_page);
            setStats(response.stats ?? null);
        } catch (error) {
            console.error("Failed to load drafts:", error);
            toast.error("Gagal memuat draft");
        } finally {
            setLoading(false);
        }
    }, [searchQuery, statusFilter]);

    useEffect(() => {
        setCurrentPage(1);
        loadDrafts(1);
    }, [searchQuery, statusFilter]);

    useEffect(() => {
        if (currentPage === 1) return;
        loadDrafts(currentPage);
    }, [currentPage]);

    return {
        drafts,
        filteredDrafts: drafts,
        stats,
        loading,
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,
        loadDrafts: (page : number) => loadDrafts(page),
        currentPage,
        setCurrentPage,
        totalPages,
        totalItems,
    };
}