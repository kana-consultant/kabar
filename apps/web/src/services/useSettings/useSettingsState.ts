import { useState } from "react";
import type { User } from "@/types/user";
import type { Team } from "@/services/user";
import { type UserRoleType } from "@/services/useSettings/types";

export function useSettingsState() {
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [users, setUsers] = useState<User[]>([]);
    const [teams, setTeams] = useState<Team[]>([]);
    const [loading, setLoading] = useState(true);
    
    // Dialog states (HANYA untuk User)
    const [showAddUserDialog, setShowAddUserDialog] = useState(false);
    const [showEditUserDialog, setShowEditUserDialog] = useState(false);
    
    // Selected items (HANYA untuk User)
    const [selectedUser, setSelectedUser] = useState<User | null>(null);
    
    // User form states
    const [newUserName, setNewUserName] = useState("");
    const [newUserEmail, setNewUserEmail] = useState("");
    const [newUserRole, setNewUserRole] = useState<UserRoleType>("viewer");
    const [isAddingUser, setIsAddingUser] = useState(false);
    
    return {
        currentUser, setCurrentUser,
        users, setUsers,
        teams, setTeams,
        loading, setLoading,
        showAddUserDialog, setShowAddUserDialog,
        showEditUserDialog, setShowEditUserDialog,
        selectedUser, setSelectedUser,
        newUserName, setNewUserName,
        newUserEmail, setNewUserEmail,
        newUserRole, setNewUserRole,
        isAddingUser, setIsAddingUser,
    };
}