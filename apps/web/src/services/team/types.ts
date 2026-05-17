// team/models.ts
import { type UserRoleType } from "../useSettings/types";

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
    role: UserRoleType
    joinedAt: string;
}

export interface CreateTeamRequest {
    name: string;
    description?: string;
}

export interface AddTeamMemberRequest {
    Email: string;
    role: UserRoleType;
}

// TeamInvite - dari Go struct
export interface TeamInvite {
    id: string;
    email: string;
    teamId: string;
    TeamName : string;
    role: string;
    token: string;
    status: string;
    invitedBy: string;
    expiresAt: string; 
    createdAt: string;
    updatedAt: string;
}

// InviteTeamMemberRequest - dari Go struct dengan validation
export interface InviteTeamMemberRequest {
    email: string;
    role: UserRoleType
}

// Optional: Enum untuk status invite
export type InviteStatus = 'pending' | 'accepted' | 'expired' | 'cancelled';

// Optional: Extended TeamInvite dengan status type
export interface TeamInviteWithStatus extends TeamInvite {
    status: InviteStatus;
}

export interface VerifyInviteResponse {
     id: string;
    email: string;
    teamId: string;
    role: string;
    token: string;
    status: string;
    invitedBy: string;
    expiresAt: string; 
    createdAt: string;
    updatedAt: string;
}

export interface AcceptInviteResponse {
    success: boolean;
    teamId?: string;
    message?: string;
}