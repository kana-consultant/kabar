import { useState } from "react";
import type { User, Team } from "@/types/user";
import type { UserRoleType } from "./types";

export function useSettingsState() {
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [users, setUsers] = useState<User[]>([]);
    const [teams, setTeams] = useState<Team[]>([]);
    const [loading, setLoading] = useState(true);
    
    // Dialog states
    const [showAddUserDialog, setShowAddUserDialog] = useState(false);
    const [showEditUserDialog, setShowEditUserDialog] = useState(false);
    const [showAddTeamDialog, setShowAddTeamDialog] = useState(false);
    const [showEditTeamDialog, setShowEditTeamDialog] = useState(false);
    const [showAddMemberDialog, setShowAddMemberDialog] = useState(false);
    const [selectedUser, setSelectedUser] = useState<User | null>(null);
    const [selectedTeam, setSelectedTeam] = useState<Team | null>(null);
    
    // Form states
    const [newUserName, setNewUserName] = useState("");
    const [newUserEmail, setNewUserEmail] = useState("");
    const [newUserRole, setNewUserRole] = useState<UserRoleType>("viewer");
    const [newTeamName, setNewTeamName] = useState("");
    const [newTeamDesc, setNewTeamDesc] = useState("");
    const [newMemberEmail, setNewMemberEmail] = useState("");
    const [newMemberRole, setNewMemberRole] = useState<UserRoleType>("viewer");

    return {
        currentUser, setCurrentUser,
        users, setUsers,
        teams, setTeams,
        loading, setLoading,
        showAddUserDialog, setShowAddUserDialog,
        showEditUserDialog, setShowEditUserDialog,
        showAddTeamDialog, setShowAddTeamDialog,
        showEditTeamDialog, setShowEditTeamDialog,
        showAddMemberDialog, setShowAddMemberDialog,
        selectedUser, setSelectedUser,
        selectedTeam, setSelectedTeam,
        newUserName, setNewUserName,
        newUserEmail, setNewUserEmail,
        newUserRole, setNewUserRole,
        newTeamName, setNewTeamName,
        newTeamDesc, setNewTeamDesc,
        newMemberEmail, setNewMemberEmail,
        newMemberRole, setNewMemberRole,
    };
}