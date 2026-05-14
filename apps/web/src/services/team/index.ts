// Types
export type { 
    Team, 
    TeamMember, 
    CreateTeamRequest, 
    AddTeamMemberRequest 
} from './types';

// Queries (GET)
export { getTeams, getTeamById } from './teamQueries';

// Mutations (POST, PUT, DELETE)
export { 
    createTeam, 
    updateTeam, 
    deleteTeam, 
    addTeamMember, 
    removeTeamMember 
} from './teamMutations';