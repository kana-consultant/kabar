import { apiClient } from '../api';
import type { Team } from '@/types/user';

// Add team
export async function addTeam(team: Omit<Team, 'id' | 'createdAt' | 'updatedAt'>): Promise<Team> {
    return apiClient.post<Team>('/teams', team);
}

// Update team
export async function updateTeam(id: string, updates: Partial<Omit<Team, 'id' | 'createdAt'>>): Promise<Team> {
    return apiClient.put<Team>(`/teams/${id}`, updates);
}

// Delete team
export async function deleteTeam(id: string): Promise<void> {
    await apiClient.delete(`/teams/${id}`);
}