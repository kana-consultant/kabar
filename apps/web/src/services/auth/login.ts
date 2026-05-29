import { apiClient } from '../api';
import Cookies from 'js-cookie';
import { COOKIE_OPTIONS } from './config';
import type { LoginResponse } from './types';
import {
    getUserRole,
    getPermissions,
    getTeamId,
} from '../user/permissions';

export async function login(email: string, password: string): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>('/auth/login', { email, password });

    const token = response.token;
    const user = response.user;

    // Save ke cookie
    Cookies.set('auth_token', token, COOKIE_OPTIONS);
    Cookies.set('user', JSON.stringify(user), COOKIE_OPTIONS);

    // Tunggu sebentar untuk memastikan cookie tersimpan
    await new Promise(resolve => setTimeout(resolve, 50));

    return {
        token,
        user,
        teamId: getTeamId() || null,
        role: getUserRole(),
        permissions: getPermissions(),
    };
}