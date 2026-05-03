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
    } = useSettingsState();

    const { loadData } = useSettingsData(
        setUsers, setTeams, setCurrentUser, setLoading
    );

    const { canManageUsers, canManageTeams, isSuperAdmin, isAdmin } = useSettingsPermissions(currentUser);

    const { handleAddUser, handleUpdateUser, handleDeleteUser } = useSettingsUserActions(
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
        getAvailableRoles: () => getAvailableRoles(currentUser, isSuperAdmin, isAdmin),
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
        handleAddUser: () => handleAddUser(
            newUserName, newUserEmail, newUserRole,
            setShowAddUserDialog, setNewUserName, setNewUserEmail, setNewUserRole
        ),
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