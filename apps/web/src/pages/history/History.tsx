import { useNavigate } from "@tanstack/react-router";
import HistoryHeader from "./HistoryHeader";
import { HistoryStats } from "./HistoryStats";
import { HistoryList } from "./HistoryList";
import { ViewHistoryDialog } from "./ViewHistoryDialog";
import { DeleteHistoryDialog } from "./DeleteHistoryDialog";
import { LoadingHistory } from "./LoadingHistory";
import { useHistory } from "@/hooks/useHistory";
import { useToast } from "@/hooks/use-toast";

export default function History() {
    const navigate = useNavigate();
    const toast = useToast();
    const {
        filteredHistory,
        searchQuery,
        setSearchQuery,
        statusFilter,
        setStatusFilter,
        actionFilter,
        setActionFilter,
        loading,
        selectedHistory,
        setSelectedHistory,
        showDetailDialog,
        setShowDetailDialog,
        showDeleteDialog,
        setShowDeleteDialog,
        handleDelete,
        handleClearAll,
        formatDate,
        getStatusData,
        getActionData,
        currentPage,
        totalPages,
        totalItems,
        handlePageChange,
        totalSuccess,
        totalFailed,
    } = useHistory();

    const handleRepost = (item: any) => {
        navigate({ to: `/generate?topic=${encodeURIComponent(item.id)}` });
        toast.info("Memuat ulang konten", {
            description: `"${item.title}" akan dimuat ulang`,
        });
    };

    if (loading) {
        return <LoadingHistory />;
    }

    return (
        <div className="space-y-6">
            <HistoryHeader
                searchQuery={searchQuery}
                setSearchQuery={setSearchQuery}
                statusFilter={statusFilter}
                setStatusFilter={setStatusFilter}
                actionFilter={actionFilter}
                setActionFilter={setActionFilter}
                onClearAll={handleClearAll}
            />

            <HistoryStats
                totalItems={totalItems}
                totalSuccess={totalSuccess}
                totalFailed={totalFailed}
            />

            <HistoryList
                items={filteredHistory}
                onView={(item) => {
                    setSelectedHistory(item);
                    setShowDetailDialog(true);
                }}
                onRepost={handleRepost}
                onDelete={(item) => {
                    setSelectedHistory(item);
                    setShowDeleteDialog(true);
                }}
                formatDate={formatDate}
                getStatusData={getStatusData}
                getActionData={getActionData}
                currentPage={currentPage}
                totalPages={totalPages}
                totalItems={totalItems}
                onPageChange={handlePageChange}  // ✅ pass ke HistoryList
            />

            <ViewHistoryDialog
                item={selectedHistory}
                open={showDetailDialog}
                onOpenChange={setShowDetailDialog}
                formatDate={formatDate}
                getStatusData={getStatusData}
                getActionData={getActionData}
            />

            <DeleteHistoryDialog
                item={selectedHistory}
                open={showDeleteDialog}
                onOpenChange={setShowDeleteDialog}
                onConfirm={handleDelete}
            />
        </div>
    );
}