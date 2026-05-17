import { useNavigate } from "@tanstack/react-router";
import { ScheduleHeader } from "./ScheduleHeader";
import { ScheduleStats } from "./ScheduleStats";
import { ScheduleList } from "./ScheduleList";
import { ScheduleDetailDialog } from "./ScheduleDetailDialog";
import { RescheduleDialog } from "./RescheduleDialog";
import { DeleteScheduleDialog } from "./DeleteScheduleDialog";
import { LoadingSchedule } from "./LoadingSchedule";
import { useSchedule } from "@/hooks/useSchedule";
import { type Draft } from "@/services/draft";

export default function Schedule() {
  const navigate = useNavigate();
  const {
    schedules,
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