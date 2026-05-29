import { useEffect, useCallback } from 'react';
import { useAuthState } from './useAuthState';
import { useAuthActions } from './useAuthActions';
import { useRoleCheck } from './useRoleCheck';
import type { UseAuthReturn } from '@/types/auth';
import { useToast } from '../use-toast';
export function useAuth(): UseAuthReturn {
    const { state, setters, actions } = useAuthState();
    const toast = useToast();
    
    const authActions = useAuthActions({
        setToken: setters.setToken,
        setUser: setters.setUser,
        setRole: setters.setRole,
        setPermissions: setters.setPermissions,
        setTeamId: setters.setTeamId,
        setIsAdmin: setters.setIsAdmin,
        setIsSuperAdmin: setters.setIsSuperAdmin,
        setIsLoading: setters.setIsLoading,
        loadUser: actions.loadUser,
        clearUserState: actions.clearUserState,
        toast,
    });
    
    const roleChecks = useRoleCheck();

    // Load user on mount
    useEffect(() => {
        actions.loadUser();
    }, []); // Hanya jalan sekali saat mount

    return {
        // User data
        user: state.user,
        token: state.token,
        role: state.role,
        permissions: state.permissions as string[],
        teamId: state.teamId,

        // Setters
        setUser: setters.setUser || null,
        setToken: setters.setToken,
        setRole: setters.setRole,
        setPermissions: setters.setPermissions,
        setTeamId: setters.setTeamId,

        // Status
        isAuthenticated: state.isAuthenticated,
        isLoading: state.isLoading,
        isAdmin: state.isAdmin,
        isSuperAdmin: state.isSuperAdmin,

        // Role checks
        ...roleChecks,

        // Actions
        ...authActions,
    };
}