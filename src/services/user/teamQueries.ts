import { apiClient } from '../api';
import type { Team } from '@/types/user';

// Get all teams
export async function getTeams(): Promise<Team[]> {
    return apiClient.get<Team[]>('/teams');
}

// Get team by id
export async function getTeamById(id: string): Promise<Team> {
    return apiClient.get<Team>(`/teams/${id}`);
}

// Get user's teams
export async function getUserTeams(userId: string): Promise<Team[]> {
    const teams = await getTeams();
    return teams.filter(team => 
        team.members.some(member => member.userId === userId)
    );
}

// Get team count
export async function getTeamCount(): Promise<number> {
    const teams = await getTeams();
    return teams.length;
}