import { updateDraft } from '../draft';

// Reschedule draft (one-time)
export async function rescheduleDraft(id: string, newDate: string, newTime: string): Promise<void> {
    const dateTime = `${newDate}T${newTime}`;  // ← perbaiki: variabel dateTime tidak didefinisikan
    await updateDraft(id, { scheduledFor: dateTime });
}

// Reschedule to daily
export async function rescheduleToDaily(id: string, dailyTime: string): Promise<void> {
    await updateDraft(id, { scheduledFor: `daily:${dailyTime}` });
}

// Cancel schedule (back to draft)
export async function cancelSchedule(id: string): Promise<void> {
    await updateDraft(id, { 
        status: 'draft', 
        scheduledFor: undefined 
    });
}