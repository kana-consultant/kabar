import { updateDraft } from '../draft';
import { getScheduledDrafts } from './scheduleQueries';
import type { Draft } from '@/types/draft';

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