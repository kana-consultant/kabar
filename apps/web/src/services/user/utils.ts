import { getUserFromCookie } from '../api';

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