import { useSettingsState } from "./useSettingsState";
import { useSettingsData } from "./useSettingsData";
import { useSettingsPermissions } from "./useSettingsPermissions";
import { useSettingsUserActions } from "./useSettingsUserActions";
import { useSettingsTeamActions } from "./useSettingsTeamActions";
import { useSettingsMemberActions } from "./useSettingsMemberActions";
import {
    getRoleDisplayName,
    roleOptions,
    getAvailableRoles,
} from "./useSettingsHelpers";
import { addTeamMember } from "../team";
import { getTeamId } from "../user";
import { type AddTeamMemberRequest } from "../team";
import { toast } from "sonner";

export function useSettings() {
    const {
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
        isAddingUser, setIsAddingUser
    } = useSettingsState();
  

    const { loadData } = useSettingsData(
        setUsers, setTeams, setCurrentUser, setLoading
    );


    const { canManageUsers, canManageTeams, isSuperAdmin, isAdmin } = useSettingsPermissions(currentUser);

    const { handleUpdateUser, handleDeleteUser } = useSettingsUserActions(
        loadData, currentUser, users, isSuperAdmin
    );

    const { handleAddTeam, handleUpdateTeam, handleDeleteTeam } = useSettingsTeamActions(loadData);

    const { handleAddMember, handleRemoveMember } = useSettingsMemberActions(
        loadData, users, currentUser
    );

    return {
        currentUser,
        users,
        teams,
        loading,
        canManageUsers,
        canManageTeams,
        isSuperAdmin,
        isAdmin,
        roleOptions,
        getAvailableRoles: () => getAvailableRoles(currentUser, isAdmin),
        getRoleDisplayName,
        showAddUserDialog,
        setShowAddUserDialog,
        showEditUserDialog,
        setShowEditUserDialog,
        showAddTeamDialog,
        setShowAddTeamDialog,
        showEditTeamDialog,
        setShowEditTeamDialog,
        showAddMemberDialog,
        setShowAddMemberDialog,
        selectedUser,
        setSelectedUser,
        selectedTeam,
        setSelectedTeam,
        newUserName,
        setNewUserName,
        newUserEmail,
        setNewUserEmail,
        newUserRole,
        setNewUserRole,
        newTeamName,
        setNewTeamName,
        newTeamDesc,
        setNewTeamDesc,
        newMemberEmail,
        setNewMemberEmail,
        newMemberRole,
        setNewMemberRole,
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

                // Success toast
                toast.success(`User ${newUserEmail} has been added as ${newUserRole}`);
                setNewUserEmail("");
                setNewUserRole("viewer");
                setNewUserName("")
                setIsAddingUser(false);

            } catch (error: any) {
                console.error("Error adding user:", error);

                // Error toast
                toast.error(error?.message || "Failed to add user to team. Please try again.")
            }finally{
                 setIsAddingUser(false);
            }
        },
        handleUpdateUser: () => handleUpdateUser(selectedUser, setShowEditUserDialog, setSelectedUser),
        handleDeleteUser,
        handleAddTeam: () => handleAddTeam(
            newTeamName, newTeamDesc,
            setShowAddTeamDialog, setNewTeamName, setNewTeamDesc
        ),
        handleUpdateTeam: () => handleUpdateTeam(selectedTeam, setShowEditTeamDialog, setSelectedTeam),
        handleDeleteTeam,
        handleAddMember: () => handleAddMember(
            newMemberEmail, newMemberRole, selectedTeam,
            setShowAddMemberDialog, setNewMemberEmail, setNewMemberRole
        ),
        handleRemoveMember,
        loadData,
    };
}