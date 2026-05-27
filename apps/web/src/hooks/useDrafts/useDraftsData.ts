import { useState, useEffect, useCallback } from "react";
import { useToast } from "@/hooks/use-toast";
import { getDrafts } from "@/services/draft";
import type { StatusFilter } from "./types";
import type {  Draft, DraftStats ,SEOScore,SimilarDraft } from "@/services/draft";


function useDebounce<T>(value: T, delay: number): T {
    const [debounced, setDebounced] = useState(value);

    useEffect(() => {
        const timer = setTimeout(() => setDebounced(value), delay);
        return () => clearTimeout(timer);
    }, [value, delay]);

    return debounced;
}

export function useDraftsData() {
    const toast = useToast();
    const [drafts, setDrafts] = useState<Draft[]>([]);
    const [stats, setStats] = useState<DraftStats | null>(null);
    const [loading, setLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
    const [seoDialog, setSeoDialog] = useState<{ open: boolean; draft: Draft | null }>({ open: false, draft: null });
    const [similarityDialog, setSimilarityDialog] = useState<{ open: boolean; draft: Draft | null }>({ open: false, draft: null });
    const [seoData, setSeoData] = useState<SEOScore | null>(null);
    const [similarityData, setSimilarityData] = useState<SimilarDraft[] | null>(null);
    const [seoLoading, setSeoLoading] = useState(false);
    const [similarityLoading, setSimilarityLoading] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(0);
    const [totalItems, setTotalItems] = useState(0);
    const LIMIT = 5;

    const debouncedSearch = useDebounce(searchQuery, 1000);

    const loadDrafts = useCallback(async (page: number = 1) => {
        setLoading(true);
        try {
            const response = await getDrafts({
                offset: (page - 1) * LIMIT,
                limit: LIMIT,
                search: debouncedSearch || undefined, // ✅ pakai debouncedSearch
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
    }, [debouncedSearch, statusFilter]);

    // ✅ Trigger reset & fetch hanya setelah debounce selesai
    useEffect(() => {
        setCurrentPage(1);
        loadDrafts(1);
    }, [debouncedSearch, statusFilter]);

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
        setSearchQuery, // tetap update searchQuery langsung untuk value input
        statusFilter,
        setStatusFilter,
        loadDrafts: (page: number) => loadDrafts(page),
        currentPage,
        setCurrentPage,
        totalPages,
        totalItems,
        seoDialog, setSeoDialog,
        similarityData, setSimilarityData,
        similarityDialog, setSimilarityDialog,
        seoLoading, setSeoLoading,
        similarityLoading, setSimilarityLoading,
        seoData, setSeoData,
    };
}