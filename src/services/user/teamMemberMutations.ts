import { apiClient } from '../api';
import type { Team, TeamMember } from '@/types/user';
import { getTeamById } from './teamQueries';

// Get team members
export async function getTeamMembers(teamId: string): Promise<TeamMember[]> {
    const team = await getTeamById(teamId);
    return team?.members || [];
}

// Add member to team
export async function addTeamMember(
    teamId: string, 
    member: Omit<TeamMember, 'id' | 'joinedAt'>
): Promise<Team> {
    return apiClient.post<Team>(`/teams/${teamId}/members`, member);
}

// Update team member role
export async function updateTeamMemberRole(
    teamId: string, 
    userId: string, 
    role: TeamMember['role']
): Promise<Team> {
    return apiClient.put<Team>(`/teams/${teamId}/members/${userId}`, { role });
}

// Remove member from team
export async function removeTeamMember(teamId: string, userId: string): Promise<Team> {
    return apiClient.delete(`/teams/${teamId}/members/${userId}`);
}