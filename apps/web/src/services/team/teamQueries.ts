import { apiClient } from '../api';
import type { Team, TeamInvite, VerifyInviteResponse } from './types';

// Get all teams
export async function getTeams(): Promise<Team[]> {
    return apiClient.get('/teams');
}

// Get team by id
export async function getTeamById(id: string): Promise<Team> {
    return apiClient.get(`/teams/${id}`);
}

export async function acceptInvite(Name :string, token: string, Password: string): Promise<TeamInvite> {
    const payload = { token, Password}
    const response = await apiClient.post(`/accept-invite`,payload);
    return response as TeamInvite;
}
// Verify invite token before showing form
export const verifyInviteToken = async (token: string): Promise<VerifyInviteResponse> => {
    const response = await apiClient.get(`/verify-invite/${token}`);
    return response as VerifyInviteResponse
}