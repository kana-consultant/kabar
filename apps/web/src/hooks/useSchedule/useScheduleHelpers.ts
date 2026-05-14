export function formatDate(dateString: string) {
    return new Date(dateString).toLocaleDateString("id-ID", {
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
}

export function getScheduleDisplay(scheduledFor?: string) {
    if (!scheduledFor) return "Tidak terjadwal";
    if (scheduledFor.startsWith("daily:")) {
        return `Setiap hari jam ${scheduledFor.replace("daily:", "")}`;
    }
    return formatDate(scheduledFor);
}

export function isDailySchedule(scheduledFor?: string) {
    return scheduledFor?.startsWith("daily:") || false;
}