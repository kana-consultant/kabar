import { useNavigate } from "@tanstack/react-router";
import { DraftsHeader } from "./DraftsHeader";
import { DraftStats } from "./DraftStats";
import { DraftList } from "./DraftList";
import { ViewDraftDialog } from "./ViewDraftDialog";
import { ScheduleDialog } from "./ScheduleDialog";
import { DeleteAlertDialog } from "./DeleteAlertDialog";
import { LoadingDrafts } from "./LoadingDrafts";
import { useDrafts } from "@/hooks/useDrafts";
import { type Draft } from "@/services/draft";

export default function Drafts() {
    const navigate = useNavigate();
    const {
        drafts,
        filteredDrafts,
        searchQuery,
        setSearchQuery,
        loading,
        statusFilter,
        setStatusFilter,
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
        setSelectedDraft
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
                statusFilter={statusFilter}
                setStatusFilter={setStatusFilter}
            />

            <DraftStats drafts={drafts} />

            <DraftList
                drafts={filteredDrafts}
                onView={(draft) => setSelectedDraft(draft)}
                onEdit={handleEdit}
                onSchedule={openScheduleDialog}
                onPublishNow={handlePublishNow}
                onDelete={openDeleteDialog}
                formatDate={formatDate}
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
        </div>
    );
}