import type { User } from '@/services/user';

export interface LoginResponse {
    token: string;
    user: User;
    teamId: string | null;
    role : string;
    permissions? : string[];
}