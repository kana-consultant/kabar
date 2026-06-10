import { useNavigate } from "@tanstack/react-router";
import { DraftsHeader } from "./DraftsHeader";
import { DraftList } from "./DraftList";
import { ViewDraftDialog } from "./ViewDraftDialog";
import { ScheduleDialog } from "./ScheduleDialog";
import { DeleteAlertDialog } from "./DeleteAlertDialog";
import { LoadingDrafts } from "./LoadingDrafts";
import { useDrafts } from "@/hooks/useDrafts";
import { type Draft } from "@/services/draft";
import { SEOScoreDialog } from "./SEOScoreDialog";
import { SimilarityDialog } from "./SimilarityDialog";

export default function Drafts() {
    const navigate = useNavigate();
    const {
        filteredDrafts,
        searchQuery,
        setSearchQuery,
        loading,
        selectedDraft,
        showScheduleDialog,
        showDeleteDialog,
        scheduleDate,
        setScheduleDate,
        scheduleTime,
        setScheduleTime,
        dailySchedule,
        setDailySchedule,
        dailyTime,
        setDailyTime,
        handleDelete,
        handleSchedule,
        handlePublishNow,
        openScheduleDialog,
        openDeleteDialog,
        closeDialogs,
        formatDate,
        setSelectedDraft,
        currentPage,
        totalPages,
        totalItems,
        loadDrafts,
        checkSimilarity,
        getSeoScore,
        seoDialog, setSeoDialog,
        similarityData,
        similarityDialog, setSimilarityDialog,
        seoLoading,
        similarityLoading,
        seoData,
    } = useDrafts();

    const handleEdit = (draft: Draft) => {
        navigate({ to: `/generate?edit=${draft.id}` });
    };

    if (loading) {
        return <LoadingDrafts />;
    }

    return (
        <div className="space-y-6">

            <DraftsHeader
                searchQuery={searchQuery}
                setSearchQuery={setSearchQuery}

            />

            <DraftList
                drafts={filteredDrafts}
                onView={(draft) => setSelectedDraft(draft)}
                onEdit={handleEdit}
                onSchedule={openScheduleDialog}
                onPublishNow={handlePublishNow}
                onDelete={openDeleteDialog}
                formatDate={formatDate}
                currentPage={currentPage}
                onPageChange={(page) => loadDrafts(page)}
                totalPages={totalPages}
                totalItems={totalItems}
                checkSimilarity={checkSimilarity}
                getSeoScore={getSeoScore}
            />

            <ViewDraftDialog
                draft={selectedDraft}
                open={!!selectedDraft && !showScheduleDialog && !showDeleteDialog}
                onOpenChange={() => setSelectedDraft(null)}
                formatDate={formatDate}
            />

            <ScheduleDialog
                draft={selectedDraft}
                open={showScheduleDialog}
                onOpenChange={closeDialogs}
                scheduleDate={scheduleDate}
                setScheduleDate={setScheduleDate}
                scheduleTime={scheduleTime}
                setScheduleTime={setScheduleTime}
                dailySchedule={dailySchedule}
                setDailySchedule={setDailySchedule}
                dailyTime={dailyTime}
                setDailyTime={setDailyTime}
                onSchedule={handleSchedule}
            />

            <DeleteAlertDialog
                draft={selectedDraft}
                open={showDeleteDialog}
                onOpenChange={closeDialogs}
                onConfirm={handleDelete}
            />

            <SEOScoreDialog
                open={seoDialog.open}
                onOpenChange={(open) => setSeoDialog((prev) => ({ ...prev, open }))}
                draftTitle={seoDialog.draft?.title}
                data={seoData}
                loading={seoLoading}
            />

            <SimilarityDialog
                open={similarityDialog.open}
                onOpenChange={(open) => setSimilarityDialog((prev) => ({ ...prev, open }))}
                draftTitle={similarityDialog.draft?.title}
                data={similarityData}
                loading={similarityLoading}
            />
        </div>
    );
}