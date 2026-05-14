import { apiClient } from '../api';
import type { Team } from './types';

// Get all teams
export async function getTeams(): Promise<Team[]> {
    return apiClient.get('/teams');
}

// Get team by id
export async function getTeamById(id: string): Promise<Team> {
    return apiClient.get(`/teams/${id}`);
}