import { useMemo } from 'react';
import {
    isManagerOrAbove,
    isEditorOrAbove,
    hasRole,
    hasAnyRole,
    getRoleLevel,
    canAccessTeam
} from '@/services/auth';
import type { RoleChecks } from '@/types/auth';

export function useRoleCheck(): RoleChecks {
    const checks = useMemo(() => ({
        isManagerOrAbove: isManagerOrAbove(),
        isEditorOrAbove: isEditorOrAbove(),
        hasRole: (roleName: string) => hasRole(roleName),
        hasAnyRole: (roles: string[]) => hasAnyRole(roles),
        getRoleLevel: () => getRoleLevel(),
        canAccessTeam: (resourceTeamId: string | null) => canAccessTeam(resourceTeamId)
    }), []);

    return checks;
}