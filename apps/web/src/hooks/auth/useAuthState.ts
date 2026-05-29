import { useState, useCallback } from 'react';
import {
    getCurrentUser,
    getToken,
    getUserRole,
    getTeamIdUser,
    hasToken,

} from '@/services/auth';
import type { User } from '@/services/user';
import { isAdmin as AdminCheck, isSuperAdmin as superAdminCheck, getPermissions } from '@/services/user/permissions';

export function useAuthState() {
    const [user, setUser] = useState<User | null>(null);
    const [permissions, setPermissions] = useState<String[]>([])
    const [token, setToken] = useState<string | null>(null);
    const [role, setRole] = useState<string | null>(null);
    const [teamId, setTeamId] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isAdmin, setIsAdmin] = useState(false);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);

    const loadUser = useCallback(async () => {
        setIsLoading(true);
        try {
            const [currentUser, currentToken, currentRole, currentTeamId, permissions, isAdmin, isSuperAdmin] = await Promise.all([
                getCurrentUser(),
                getToken(),
                getUserRole(),
                getTeamIdUser(),
                getPermissions(),
                AdminCheck(),
                superAdminCheck(),
            ]);

            setUser(currentUser);
            setToken(currentToken as string);
            setRole(currentRole);
            setTeamId(currentTeamId);
            setPermissions(permissions);
            setIsAdmin(isAdmin);
            setIsSuperAdmin(isSuperAdmin);
        } catch (error) {
            console.error('Failed to load user:', error);
        } finally {
            setIsLoading(false);
        }
    }, [setToken, setUser, setRole, setTeamId, setPermissions, setIsAdmin, setIsSuperAdmin, setIsLoading]);

    const updateUserState = useCallback((updatedUser: User) => {
        setUser(updatedUser);
        setRole(updatedUser.role);
        setIsAdmin(updatedUser.role === 'admin')
    }, []);

    const clearUserState = useCallback(() => {
        setUser(null);
        setToken(null);
        setRole(null);
        setTeamId(null);
        setIsAdmin(false);
        setIsSuperAdmin(false);
    }, []);

    return {
        state: {
            user,
            token,
            role,
            teamId,
            isLoading,
            isAdmin,
            isSuperAdmin,
            isAuthenticated: hasToken(),
            permissions
        },
        setters: {
            setToken,
            setUser: updateUserState,
            setRole,
            setTeamId,
            setIsLoading,
            setIsAdmin,
            setIsSuperAdmin,
            setPermissions
        },
        actions: {
            loadUser,
            clearUserState,
            updateUserState
        }
    };
}