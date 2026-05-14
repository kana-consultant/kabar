import type { User } from '@/services/user';

export interface LoginResponse {
    token: string;
    user: User;
    teamId?: string;
}