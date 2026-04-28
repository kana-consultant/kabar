// src/services/scheduleService.ts
import { getDrafts, updateDraft } from './draftService';
import type { Draft } from '@/types/draft';

// Get all scheduled drafts
export async function getScheduledDrafts(): Promise<Draft[]> {
    const drafts = await getDrafts();
    return drafts.filter(d => d.status === 'scheduled');
}

// Get upcoming schedules
export async function getUpcomingSchedules(limit: number = 10): Promise<Draft[]> {
    const scheduled = await getScheduledDrafts();
    const now = new Date();

    return scheduled
        .filter(s => {
            // Daily schedule selalu dianggap upcoming
            if (s.scheduledFor?.startsWith('daily:')) return true;
            // One-time schedule: cek apakah masih di masa depan
            if (s.scheduledFor) {
                const scheduleDate = new Date(s.scheduledFor);
                return scheduleDate > now;
            }
            return false;
        })
        .slice(0, limit);
}

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

// Process daily schedules (untuk cron job nanti)
export async function processDailySchedules(): Promise<Draft[]> {
    const scheduled = await getScheduledDrafts();
    const now = new Date();
    const currentHour = now.getHours().toString().padStart(2, '0');
    const currentMinute = now.getMinutes().toString().padStart(2, '0');
    const currentTime = `${currentHour}:${currentMinute}`;

    // Cari draft dengan daily schedule yang waktunya sekarang
    const toPublish = scheduled.filter(s => 
        s.scheduledFor === `daily:${currentTime}`
    );

    // Update status to published
    for (const draft of toPublish) {
        await updateDraft(draft.id, { 
            status: 'published',
            scheduledFor: undefined, // clear schedule setelah publish
        });
    }

    return toPublish;
}

// Process one-time schedules (untuk cron job)
export async function processOneTimeSchedules(): Promise<Draft[]> {
    const scheduled = await getScheduledDrafts();
    const now = new Date();

    // Cari draft dengan one-time schedule yang waktunya sudah lewat
    const toPublish = scheduled.filter(s => {
        if (s.scheduledFor && !s.scheduledFor.startsWith('daily:')) {
            const scheduleDate = new Date(s.scheduledFor);
            return scheduleDate <= now && s.status === 'scheduled';
        }
        return false;
    });

    // Update status to published
    for (const draft of toPublish) {
        await updateDraft(draft.id, { 
            status: 'published',
            scheduledFor: undefined,
        });
    }

    return toPublish;
}

// Get draft schedule info
export async function getScheduleInfo(id: string): Promise<{
    isScheduled: boolean;
    isDaily: boolean;
    scheduledFor?: string;
    dailyTime?: string;
} | null> {
    const drafts = await getDrafts();
    const draft = drafts.find(d => d.id === id);
    
    if (!draft || draft.status !== 'scheduled') {
        return null;
    }

    const isDaily = draft.scheduledFor?.startsWith('daily:') || false;
    
    return {
        isScheduled: true,
        isDaily,
        scheduledFor: isDaily ? undefined : draft.scheduledFor,
        dailyTime: isDaily ? draft.scheduledFor?.replace('daily:', '') : undefined,
    };
}