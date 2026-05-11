import { getDrafts } from '../draft';
import type { Draft } from '@/services/draft';

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
            if (s.scheduled_for?.startsWith('daily:')) return true;
            // One-time schedule: cek apakah masih di masa depan
            if (s.scheduled_for) {
                const scheduleDate = new Date(s.scheduled_for);
                return scheduleDate > now;
            }
            return false;
        })
        .slice(0, limit);
}

// Get draft schedule info
export async function getScheduleInfo(id: string): Promise<{
    isScheduled: boolean;
    isDaily: boolean;
    scheduled_for?: string;
    dailyTime?: string;
} | null> {
    const drafts = await getDrafts();
    const draft = drafts.find(d => d.id === id);
    
    if (!draft || draft.status !== 'scheduled') {
        return null;
    }

    const isDaily = draft.scheduled_for?.startsWith('daily:') || false;
    
    return {
        isScheduled: true,
        isDaily,
        scheduled_for: isDaily ? undefined : draft.scheduled_for,
        dailyTime: isDaily ? draft.scheduled_for?.replace('daily:', '') : undefined,
    };
}