import { useSettingsState } from "./useSettingsState";
import { useSettingsData } from "./useSettingsData";
import { useSettingsPermissions } from "./useSettingsPermissions";
import { useSettingsUserActions } from "./useSettingsUserActions";
import {
    getRoleDisplayName,
    roleOptions,
    getAvailableRoles,
} from "./useSettingsHelpers";
import { addTeamMember } from "../team";
import { getTeamId } from "../user";
import { type AddTeamMemberRequest } from "../team";
import { Toast } from  "@kana-consultant/ui-kit";
import { useAuth } from "@/hooks/auth/useAuth";

export function useSettings() {
    const {
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
        isAddingUser, setIsAddingUser
    } = useSettingsState();

    const {changePassword} = useAuth();

    const { loadData } = useSettingsData(
        setUsers, setTeams, setCurrentUser, setLoading
    );

    const { canManageUsers, canManageTeams, isAdmin } = useSettingsPermissions(currentUser);

    const { handleUpdateUser, handleDeleteUser } = useSettingsUserActions(
        loadData, currentUser, users
    );

    return {
        currentUser,
        users,
        teams,
        loading,
        canManageUsers,
        canManageTeams,
        isAdmin,
        roleOptions,
        getAvailableRoles: () => getAvailableRoles(currentUser, isAdmin),
        getRoleDisplayName,
        showAddUserDialog,
        setShowAddUserDialog,
        showEditUserDialog,
        setShowEditUserDialog,
        selectedUser,
        setSelectedUser,
        newUserName,
        setNewUserName,
        newUserEmail,
        setNewUserEmail,
        newUserRole,
        setNewUserRole,
        isAddingUser, 
        setIsAddingUser,
        handleAddUser: async () => {
            try {
                setIsAddingUser(true);
                const userTeam: AddTeamMemberRequest = {
                    Email: newUserEmail,
                    role: newUserRole
                };
                await addTeamMember(getTeamId(), userTeam);
                toast.success(`User ${newUserEmail} has been added as ${newUserRole}`);
                setNewUserEmail("");
                setNewUserRole("viewer");
                setNewUserName("");
                await loadData();
            } catch (error: any) {
                console.error("Error adding user:", error);
                toast.error(error?.message || "Failed to add user to team. Please try again.");
            } finally {
                setIsAddingUser(false);
            }
        },
        handleUpdateUser: () => handleUpdateUser(selectedUser, setShowEditUserDialog, setSelectedUser),
        handleDeleteUser,
        loadData,
        HandleChangePassword : changePassword
    };
}