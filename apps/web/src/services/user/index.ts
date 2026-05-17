// Types
export type { User, UserRole, Team, TeamMember } from './types';

// User Queries (GET)
export { getUsers, getUserById, getUserByEmail, getUserCount, getActiveUserCount } from './userQueries';

// User Mutations (POST, PUT, DELETE)
export { updateLastActive, register, addUser, updateUser, deleteUser } from './userMutations';

// Team Queries (GET)
export { getTeams, getTeamById, getUserTeams, getTeamCount } from './teamQueries';

// Team Mutations (POST, PUT, DELETE)
export { addTeam, updateTeam, deleteTeam } from './teamMutations';

// Team Member Mutations
export { getTeamMembers, addTeamMember, updateTeamMemberRole, removeTeamMember } from './teamMemberMutations';

// Authentication
export { getCurrentUser, login, logout, isAuthenticated } from './auth';

// Permissions
export { 
    getUserRole, 
    getUserRoleLevel, 
    getUserRoleName, 
    getUserRoleDisplayName,
    getTeamId,
    isSuperAdmin, 
    isAdmin, 
    isManagerOrAbove, 
    isEditorOrAbove, 
    isViewerOrAbove,
    hasRole, 
    hasAnyRole, 
    hasMinLevel, 
    canAccessResource
} from './permissions';

// Utilities
export { getUserName, getUserEmail, getUserAvatar } from './utils';