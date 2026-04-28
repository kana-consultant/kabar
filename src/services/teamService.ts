// src/services/teamService.ts
import { apiClient } from './api';

export interface Team {
    id: string;
    name: string;
    description?: string;
    createdBy?: string;
    createdAt: string;
    updatedAt: string;
    members: TeamMember[];
}

export interface TeamMember {
    id: string;
    teamId: string;
    userId: string;
    userEmail: string;
    userName: string;
    role: 'manager' | 'editor' | 'viewer';
    joinedAt: string;
}

export interface CreateTeamRequest {
    name: string;
    description?: string;
}

export interface AddTeamMemberRequest {
    userId: string;
    role: 'manager' | 'editor' | 'viewer';
}

// Get all teams
export async function getTeams(): Promise<Team[]> {
    return apiClient.get('/teams');
}

// Get team by id
export async function getTeamById(id: string): Promise<Team> {
    return apiClient.get(`/teams/${id}`);
}

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
    await apiClient.post(`/teams/${teamId}/members`, req);
}

// Remove member from team
export async function removeTeamMember(teamId: string, userId: string): Promise<void> {
    await apiClient.delete(`/teams/${teamId}/members/${userId}`);
}