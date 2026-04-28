// src/services/userService.ts
import { apiClient, getAuthToken, setAuthCookie, removeAuthCookie, getUserFromCookie } from './api';
import type { User, UserRole, Team, TeamMember } from '@/types/user';
import Cookies from 'js-cookie';

// ==================== USER CRUD ====================

// Get all users
export async function getUsers(params?: {
    page?: number;
    limit?: number;
    role?: string;
    status?: string;
}): Promise<User[]> {
    return apiClient.get<User[]>('/users', { params });
}

// Get user by id
export async function getUserById(id: string): Promise<User> {
    return apiClient.get<User>(`/users/${id}`);
}

// Get user by email
export async function getUserByEmail(email: string): Promise<User | null> {
    const users = await getUsers();
    return users.find(u => u.email === email) || null;
}

// Get current logged in user from cookie
export async function getCurrentUser(): Promise<User | null> {
    return getUserFromCookie();
}

// Update last active
export async function updateLastActive(userId: string): Promise<void> {
    await apiClient.patch(`/users/${userId}/active`);
}

// Login
export async function login(email: string, password: string): Promise<{ user: User; token: string }> {
    const response = await apiClient.post<{ token: string; user: User }>('/auth/login', { email, password });
    
    if (response.token) {
        setAuthCookie(response.token, response.user);
    }
    
    return response;
}

// Logout
export async function logout(): Promise<void> {
    removeAuthCookie();
    try {
        await apiClient.post('/auth/logout');
    } catch (error) {
        // Ignore error
    }
}

// Register
export async function register(data: {
    email: string;
    name: string;
    password: string;
    role?: string;
}): Promise<User> {
    return apiClient.post<User>('/auth/register', data);
}

// Add user (admin only)
export async function addUser(user: Omit<User, 'id' | 'createdAt' | 'updatedAt'>): Promise<User> {
    return apiClient.post<User>('/users', user);
}

// Update user
export async function updateUser(id: string, updates: Partial<Omit<User, 'id' | 'createdAt'>>): Promise<User> {
    return apiClient.put<User>(`/users/${id}`, updates);
}

// Delete user
export async function deleteUser(id: string): Promise<void> {
    await apiClient.delete(`/users/${id}`);
}

// ==================== TEAM CRUD ====================

// Get all teams
export async function getTeams(): Promise<Team[]> {
    return apiClient.get<Team[]>('/teams');
}

// Get team by id
export async function getTeamById(id: string): Promise<Team> {
    return apiClient.get<Team>(`/teams/${id}`);
}

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

// ==================== TEAM MEMBERS ====================

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

// Get user's teams
export async function getUserTeams(userId: string): Promise<Team[]> {
    const teams = await getTeams();
    return teams.filter(team => 
        team.members.some(member => member.userId === userId)
    );
}

// ==================== USER ROLES & PERMISSIONS ====================

// Get user role from cookie
export function getUserRole(): UserRole | null {
    const user = getUserFromCookie();
    return user?.role || null;
}

// Get role level
export function getUserRoleLevel(): number {
    const role = getUserRole();
    return role?.level || 0;
}

// Get role name
export function getUserRoleName(): string {
    const role = getUserRole();
    return role?.name || 'viewer';
}

// Get role display name
export function getUserRoleDisplayName(): string {
    const role = getUserRole();
    return role?.displayName || 'Viewer';
}

// Get team ID from cookie
export function getTeamId(): string | null {
    return Cookies.get('team_id') || null;
}

// Check if user is super admin
export function isSuperAdmin(): boolean {
    const roleName = getUserRoleName();
    return roleName === 'super_admin';
}

// Check if user is admin (includes super_admin)
export function isAdmin(): boolean {
    const roleName = getUserRoleName();
    return roleName === 'admin' || roleName === 'super_admin';
}

// Check if user is manager or above
export function isManagerOrAbove(): boolean {
    const roleLevel = getUserRoleLevel();
    return roleLevel >= 70; // manager level = 70
}

// Check if user is editor or above
export function isEditorOrAbove(): boolean {
    const roleLevel = getUserRoleLevel();
    return roleLevel >= 50; // editor level = 50
}

// Check if user is viewer or above (everyone)
export function isViewerOrAbove(): boolean {
    const roleLevel = getUserRoleLevel();
    return roleLevel >= 10; // viewer level = 10
}

// Check if user has specific role by name
export function hasRole(roleName: string): boolean {
    const userRoleName = getUserRoleName();
    if (userRoleName === 'super_admin') return true;
    return userRoleName === roleName;
}

// Check if user has any of the required role names
export function hasAnyRole(roleNames: string[]): boolean {
    const userRoleName = getUserRoleName();
    if (userRoleName === 'super_admin') return true;
    return roleNames.includes(userRoleName);
}

// Check if user has permission based on role level
export function hasMinLevel(requiredLevel: number): boolean {
    const userLevel = getUserRoleLevel();
    return userLevel >= requiredLevel;
}

// Check if user can access resource based on team
export async function canAccessResource(userId: string, resourceTeamId: string): Promise<boolean> {
    const user = await getUserById(userId);
    if (!user) return false;
    
    const roleName = user.role || "viewer";
    
    // Super admin and admin can access everything
    if (roleName === 'super_admin' || roleName === 'admin') return true;
    
    // Check if user is in the team
    const userTeams = await getUserTeams(userId);
    return userTeams.some(team => team.id === resourceTeamId);
}

// ==================== UTILITIES ====================

// Get user count
export async function getUserCount(): Promise<number> {
    const users = await getUsers();
    return users.length;
}

// Get active user count
export async function getActiveUserCount(): Promise<number> {
    const users = await getUsers();
    return users.filter(u => u.status === 'active').length;
}

// Get team count
export async function getTeamCount(): Promise<number> {
    const teams = await getTeams();
    return teams.length;
}

// Get user name
export function getUserName(): string {
    const user = getUserFromCookie();
    return user?.name || '';
}

// Get user email
export function getUserEmail(): string {
    const user = getUserFromCookie();
    return user?.email || '';
}

// Get user avatar
export function getUserAvatar(): string | undefined {
    const user = getUserFromCookie();
    return user?.avatar;
}

// Check if authenticated
export function isAuthenticated(): boolean {
    return !!getAuthToken();
}