import { useNavigate } from "@tanstack/react-router";
import { ScheduleHeader } from "./ScheduleHeader";
import { ScheduleList } from "./ScheduleList";
import { ScheduleDetailDialog } from "./ScheduleDetailDialog";
import { RescheduleDialog } from "./RescheduleDialog";
import { DeleteScheduleDialog } from "./DeleteScheduleDialog";
import { LoadingSchedule } from "./LoadingSchedule";
import { useSchedule } from "@/hooks/useSchedule";
import { type Draft } from "@/services/draft";
import { useAuth } from "@/contexts/AuthContext";

export default function Schedule() {
  const navigate = useNavigate();
  const {
    filteredSchedules,
    searchQuery,
    setSearchQuery,
    loading,
    selectedSchedule,
    setSelectedSchedule,
    showDetailDialog,
    setShowDetailDialog,
    showDeleteDialog,
    setShowDeleteDialog,
    showRescheduleDialog,
    setShowRescheduleDialog,
    newScheduleDate,
    setNewScheduleDate,
    newScheduleTime,
    setNewScheduleTime,
    loadSchedules,
    handlePublishNow,
    handleDelete,
    handleReschedule,
    formatDate,
    getScheduleDisplay,
    isDailySchedule,
    currentPage,
    setCurrentPage,
    totalPages,
    totalItems,
  } = useSchedule();

  const { can } = useAuth();

  const handleEdit = (schedule: Draft) => {
    navigate({ to: `/generate?edit=${schedule.id}` });
  };

  if (loading) {
    return <LoadingSchedule />;
  }

  return (
    <div className="space-y-6">
      <ScheduleHeader
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        onRefresh={loadSchedules}
      />

      <ScheduleList
        schedules={filteredSchedules}
        isDailySchedule={isDailySchedule}
        getScheduleDisplay={getScheduleDisplay}

        onView={(schedule) => {
          setSelectedSchedule(schedule);
          setShowDetailDialog(true);
        }}

        onEdit={can("schedule:edit:team") ? handleEdit : ()=>{}}
        onReschedule={can("schedule:edit:team") ? (schedule) => {
          setSelectedSchedule(schedule);
          setShowRescheduleDialog(true);
        } : ()=>{}}

        onPublishNow={can("schedule:publish:team") ? handlePublishNow : ()=>{}}

        onDelete={can("schedule:delete:team") ? (schedule) => {
          setSelectedSchedule(schedule);
          setShowDeleteDialog(true);
        } : ()=>{}}

        currentPage={currentPage}
        totalPages={totalPages}
        totalItems={totalItems}
        onPageChange={setCurrentPage}
      />

      {/* Detail — schedule:view:team */}
      {can("schedule:view:team") && (
        <ScheduleDetailDialog
          schedule={selectedSchedule}
          open={showDetailDialog}
          onOpenChange={setShowDetailDialog}
          getScheduleDisplay={getScheduleDisplay}
          formatDate={formatDate}
        />
      )}

      {/* Reschedule — schedule:edit:team */}
      {can("schedule:edit:team") && (
        <RescheduleDialog
          schedule={selectedSchedule}
          open={showRescheduleDialog}
          onOpenChange={setShowRescheduleDialog}
          newDate={newScheduleDate}
          onDateChange={setNewScheduleDate}
          newTime={newScheduleTime}
          onTimeChange={setNewScheduleTime}
          onReschedule={handleReschedule}
        />
      )}

      {/* Delete — schedule:delete:team */}
      {can("schedule:delete:team") && (
        <DeleteScheduleDialog
          schedule={selectedSchedule}
          open={showDeleteDialog}
          onOpenChange={setShowDeleteDialog}
          onConfirm={handleDelete}
        />
      )}
    </div>
  );
}