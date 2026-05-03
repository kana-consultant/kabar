import { useState, useCallback } from 'react';
import { 
    getCurrentUser, 
    getToken, 
    getUserRole, 
    getTeamId,
    hasToken
} from '@/services/auth';
import type { User } from '@/types/user';

export function useAuthState() {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [role, setRole] = useState<string | null>(null);
    const [teamId, setTeamId] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isAdmin, setIsAdmin] = useState(false);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);

    const loadUser = useCallback(async () => {
        setIsLoading(true);
        try {
            const [currentUser, currentToken, currentRole, currentTeamId] = await Promise.all([
                getCurrentUser(),
                getToken(),
                getUserRole(),
                getTeamId()
            ]);
            
            setUser(currentUser);
            setToken(currentToken);
            setRole(currentRole);
            setTeamId(currentTeamId);
            setIsAdmin(isAdmin);
            setIsSuperAdmin(isSuperAdmin);
        } catch (error) {
            console.error('Failed to load user:', error);
        } finally {
            setIsLoading(false);
        }
    }, []);

    const updateUserState = useCallback((updatedUser: User) => {
        setUser(updatedUser);
        setRole(updatedUser.role);
        setIsAdmin(updatedUser.role === 'admin' || updatedUser.role === 'super_admin');
        setIsSuperAdmin(updatedUser.role === 'super_admin');
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
            isAuthenticated: hasToken()
        },
        setters: {
            setToken,
            setUser: updateUserState,
            setRole,
            setTeamId,
            setIsLoading,
            setIsAdmin,
            setIsSuperAdmin
        },
        actions: {
            loadUser,
            clearUserState,
            updateUserState
        }
    };
}