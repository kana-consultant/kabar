import { useNavigate } from "@tanstack/react-router";
import { ScheduleHeader } from "@/components/schedule/ScheduleHeader";
import { ScheduleStats } from "@/components/schedule/ScheduleStats";
import { ScheduleList } from "@/components/schedule/ScheduleList";
import { ScheduleDetailDialog } from "@/components/schedule/ScheduleDetailDialog";
import { RescheduleDialog } from "@/components/schedule/RescheduleDialog";
import { DeleteScheduleDialog } from "@/components/schedule/DeleteScheduleDialog";
import { LoadingSchedule } from "@/components/schedule/LoadingSchedule";
import { useSchedule } from "@/hooks/useSchedule";
import { type Draft } from "@/types/draft";
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute("/schedule")({
  component: Schedule,
});

export default function Schedule() {
  const navigate = useNavigate();
  const {
    schedules,
    filteredSchedules,
    searchQuery,
    setSearchQuery,
    loading, // ← ambil loading
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
  } = useSchedule();

  const handleEdit = (schedule: Draft) => {
    navigate({ to: `/generate?edit=${schedule.id}` });
  };

  // Tampilkan loading
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

      <ScheduleStats
        schedules={schedules}
        isDailySchedule={isDailySchedule}
      />

      <ScheduleList
        schedules={filteredSchedules}
        isDailySchedule={isDailySchedule}
        getScheduleDisplay={getScheduleDisplay}
        onView={(schedule) => {
          setSelectedSchedule(schedule);
          setShowDetailDialog(true);
        }}
        onEdit={handleEdit}
        onReschedule={(schedule) => {
          setSelectedSchedule(schedule);
          setShowRescheduleDialog(true);
        }}
        onPublishNow={handlePublishNow}
        onDelete={(schedule) => {
          setSelectedSchedule(schedule);
          setShowDeleteDialog(true);
        }}
      />

      <ScheduleDetailDialog
        schedule={selectedSchedule}
        open={showDetailDialog}
        onOpenChange={setShowDetailDialog}
        getScheduleDisplay={getScheduleDisplay}
        formatDate={formatDate}
      />

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

      <DeleteScheduleDialog
        schedule={selectedSchedule}
        open={showDeleteDialog}
        onOpenChange={setShowDeleteDialog}
        onConfirm={handleDelete}
      />
    </div>
  );
}