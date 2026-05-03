import { apiClient } from '../api';
import type { User } from '@/types/user';

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