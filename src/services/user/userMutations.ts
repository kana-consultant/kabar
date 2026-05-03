import { apiClient } from '../api';
import type { User } from '@/types/user';

// Update last active
export async function updateLastActive(userId: string): Promise<void> {
    await apiClient.patch(`/users/${userId}/active`);
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