import { apiClient } from '../api';
import type { Team, CreateTeamRequest, AddTeamMemberRequest,AcceptInviteResponse,VerifyInviteResponse } from './types';

// Create team
export async function createTeam(req: CreateTeamRequest): Promise<Team> {
    return apiClient.post('/teams', req);
}

// Update team
export async function updateTeam(id: string, updates: Partial<Team>): Promise<void> {
    await apiClient.put(`/teams/${id}`, updates);
}

// Delete team
export async function deleteTeam(id: string): Promise<void> {
    await apiClient.delete(`/teams/${id}`);
}

// Add member to team
export async function addTeamMember(teamId: string, req: AddTeamMemberRequest): Promise<void> {
    await apiClient.post(`/teams/${teamId}/invites`, req);
}

// Remove member from team
export async function removeTeamMember(teamId: string, userId: string): Promise<void> {
    await apiClient.delete(`/teams/${teamId}/members/${userId}`);
}





// Accept invite (with optional registration data)
export const acceptInvite = async (
    token: string, 
    userData?: { password?: string; fullName?: string }
): Promise<AcceptInviteResponse> => {
    try {
        return await apiClient.post(`/teams/accept-invite`, {
            token,
            ...userData
        });
    } catch (error: any) {
        throw new Error(error.response?.data?.message || "Failed to accept invitation");
    }
};

