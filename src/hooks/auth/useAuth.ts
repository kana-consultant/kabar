import { useEffect } from 'react';
import { useAuthState } from './useAuthState';
import { useAuthActions } from './useAuthActions';
import { useRoleCheck } from './useRoleCheck';
import type { UseAuthReturn } from '@/types/auth';

export function useAuth(): UseAuthReturn {
    const { state, setters, actions } = useAuthState();
    const authActions = useAuthActions({
        setToken: setters.setToken,
        setUser: setters.setUser,
        setRole: setters.setRole,
        setTeamId: setters.setTeamId,
        setIsAdmin: setters.setIsAdmin,
        setIsSuperAdmin: setters.setIsSuperAdmin,
        setIsLoading: setters.setIsLoading,
        loadUser: actions.loadUser,
        clearUserState: actions.clearUserState
    });
    const roleChecks = useRoleCheck();

    // Load user on mount
    useEffect(() => {
        actions.loadUser();
    }, []);

    return {
        // User data
        user: state.user,
        token: state.token,
        role: state.role,
        teamId: state.teamId,
        
        // Status
        isAuthenticated: state.isAuthenticated,
        isLoading: state.isLoading,
        isAdmin: state.isAdmin,
        isSuperAdmin: state.isSuperAdmin,
        
        // Role checks
        ...roleChecks,
        
        // Actions
        ...authActions
    };
}