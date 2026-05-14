
export interface Team {
    id: string;
    name: string;
    description?: string;
    createdBy?: string;
    createdAt: string;
    updatedAt: string;
    members: TeamMember[];
}

export interface TeamMember {
    id: string;
    teamId: string;
    userId: string;
    userEmail: string;
    userName: string;
    role: 'manager' | 'editor' | 'viewer';
    joinedAt: string;
}

export interface CreateTeamRequest {
    name: string;
    description?: string;
}

export interface AddTeamMemberRequest {
    userId: string;
    role: 'manager' | 'editor' | 'viewer';
}