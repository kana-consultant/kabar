// src/hooks/useAuth.ts
import { useState, useEffect, useCallback } from 'react';
import { 
    login as loginService, 
    logout as logoutService, 
    getCurrentUser,
    getUserRole,
    getTeamId,
    isAdmin as checkIsAdmin,
    isSuperAdmin as checkIsSuperAdmin,
    isManagerOrAbove,
    isEditorOrAbove,
    hasRole,
    hasAnyRole,
    getRoleLevel,
    canAccessTeam,
    getToken,
    hasToken,
    clearAuthData,
    updateLocalUser
} from '@/services/authService';
import type { User } from '@/types/user';
import { toast } from 'sonner';

interface UseAuthReturn {
    // User data
    user: User | null;
    token: string | null;
    role: string | null;
    teamId: string | null;
    
    // Status
    isAuthenticated: boolean;
    isLoading: boolean;
    isAdmin: boolean;
    isSuperAdmin: boolean;
    
    // Role checks
    isManagerOrAbove: boolean;
    isEditorOrAbove: boolean;
    hasRole: (role: string) => boolean;
    hasAnyRole: (roles: string[]) => boolean;
    getRoleLevel: () => number;
    canAccessTeam: (resourceTeamId: string | null) => boolean;
    
    // Actions
    login: (email: string, password: string) => Promise<boolean>;
    logout: () => Promise<void>;
    refreshUser: () => Promise<void>;
    updateUser: (user: User) => void;
    clearAuth: () => void;
}

export function useAuth(): UseAuthReturn {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [role, setRole] = useState<string | null>(null);
    const [teamId, setTeamId] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isAdmin, setIsAdmin] = useState(false);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);

    // Load user from localStorage on mount
    const loadUser = useCallback(async () => {
        setIsLoading(true);
        try {
            const currentUser = await getCurrentUser();
            const currentToken = getToken();
            const currentRole = getUserRole();
            const currentTeamId = getTeamId();
            
            setUser(currentUser);
            setToken(currentToken);
            setRole(currentRole);
            setTeamId(currentTeamId);
            setIsAdmin(checkIsAdmin());
            setIsSuperAdmin(checkIsSuperAdmin());
        } catch (error) {
            console.error('Failed to load user:', error);
        } finally {
            setIsLoading(false);
        }
    }, []);

    useEffect(() => {
        loadUser();
    }, [loadUser]);

    const login = useCallback(async (email: string, password: string): Promise<boolean> => {
        setIsLoading(true);
        try {
            const response = await loginService(email, password);
            setToken(response.token);
            setUser(response.user);
            setRole(response.user.role);
            setTeamId(response.teamId || null);
            setIsAdmin(response.user.role === 'admin' || response.user.role === 'super_admin');
            setIsSuperAdmin(response.user.role === 'super_admin');
            
            toast.success('Login berhasil!', {
                description: `Selamat datang, ${response.user.name}!`,
            });
            return true;
        } catch (error: any) {
            const message = error.response?.data?.message || error.message || 'Login gagal';
            toast.error('Login gagal', { description: message });
            return false;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const logout = useCallback(async () => {
        setIsLoading(true);
        try {
            await logoutService();
            setUser(null);
            setToken(null);
            setRole(null);
            setTeamId(null);
            setIsAdmin(false);
            setIsSuperAdmin(false);
            toast.success('Logout berhasil');
        } catch (error) {
            console.error('Logout failed:', error);
        } finally {
            setIsLoading(false);
        }
    }, []);

    const refreshUser = useCallback(async () => {
        await loadUser();
    }, [loadUser]);

    const updateUser = useCallback((updatedUser: User) => {
        setUser(updatedUser);
        updateLocalUser(updatedUser);
        setRole(updatedUser.role);
        setIsAdmin(updatedUser.role === 'admin' || updatedUser.role === 'super_admin');
        setIsSuperAdmin(updatedUser.role === 'super_admin');
    }, []);

    const clearAuth = useCallback(() => {
        clearAuthData();
        setUser(null);
        setToken(null);
        setRole(null);
        setTeamId(null);
        setIsAdmin(false);
        setIsSuperAdmin(false);
    }, []);

    const isAuthenticated = hasToken();

    return {
        // User data
        user,
        token,
        role,
        teamId,
        
        // Status
        isAuthenticated,
        isLoading,
        isAdmin,
        isSuperAdmin,
        
        // Role checks
        isManagerOrAbove: isManagerOrAbove(),
        isEditorOrAbove: isEditorOrAbove(),
        hasRole: (roleName: string) => hasRole(roleName),
        hasAnyRole: (roles: string[]) => hasAnyRole(roles),
        getRoleLevel: () => getRoleLevel(),
        canAccessTeam: (resourceTeamId: string | null) => canAccessTeam(resourceTeamId),
        
        // Actions
        login,
        logout,
        refreshUser,
        updateUser,
        clearAuth,
    };
}