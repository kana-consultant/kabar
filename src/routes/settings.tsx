import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { SettingsHeader } from "@/pages/settings/SettingsHeader";
import { ProfileTab } from "@/pages/settings/ProfileTab";
import { UsersTab } from "@/pages/settings/UsersTab";
import { TeamsTab } from "@/pages/settings/TeamsTab";
import { ApiKeysTab } from "@/pages/settings/ApiKeysTab";
import { PreferencesTab } from "@/pages/settings/PreferencesTab";
import { AddUserDialog } from "@/pages/settings/AddUserDialog";
import { EditUserDialog } from "@/pages/settings/EditUserDialog";
import { AddTeamDialog } from "@/pages/settings/AddTeamDialog";
import { AddMemberDialog } from "@/pages/settings/AddMemberDialog";
import { EditTeamDialog } from "@/pages/settings/EditTeamDialog";
import { LoadingSettings } from "@/pages/settings/LoadingSettings";
import { useSettings } from "@/services/useSettings";
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute("/settings")({
  component: Settings,
});

export default function Settings() {
    const {
        currentUser,
        users,
        teams,
        loading,
        canManageUsers,
        canManageTeams,
        isSuperAdmin,
        isAdmin,
        roleOptions,
        getAvailableRoles,
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
        handleAddUser,
        handleUpdateUser,
        handleDeleteUser,
        handleAddTeam,
        handleUpdateTeam,
        handleDeleteTeam,
        handleAddMember,
        handleRemoveMember,
    } = useSettings();

    // Tampilkan loading
    if (loading) {
        return <LoadingSettings />;
    }

    return (
        <div className="space-y-6">
            <SettingsHeader
                title="Settings"
                description="Kelola pengaturan akun, tim, dan hak akses"
            />

            <Tabs defaultValue="profile" className="w-full">
                <TabsList className="grid w-full grid-cols-3">
                    <TabsTrigger value="profile">Profile</TabsTrigger>
                    <TabsTrigger value="users">Users</TabsTrigger>
                    <TabsTrigger value="api">API Keys</TabsTrigger>
                </TabsList>

                <TabsContent value="profile" className="space-y-4">
                    <ProfileTab currentUser={currentUser} />
                </TabsContent>

                <TabsContent value="users" className="space-y-4">
                    <UsersTab
                        users={users}
                        currentUserId={currentUser?.id}
                        canManage={canManageUsers as any}
                        isSuperAdmin={isSuperAdmin}
                        isAdmin={isAdmin}
                        roleOptions={roleOptions}
                        getAvailableRoles={getAvailableRoles}
                        getRoleDisplayName={getRoleDisplayName}
                        onAddUser={() => setShowAddUserDialog(true)}
                        onEditUser={(user) => {
                            setSelectedUser(user);
                            setShowEditUserDialog(true);
                        }}
                        onDeleteUser={handleDeleteUser}
                    />
                </TabsContent>

                {/* <TabsContent value="teams" className="space-y-4">
                    <TeamsTab
                        teams={teams}
                        currentUserId={currentUser?.id}
                        canManage={canManageTeams as any}
                        roleOptions={roleOptions}
                        getRoleDisplayName={getRoleDisplayName}
                        onAddTeam={() => setShowAddTeamDialog(true)}
                        onEditTeam={(team) => {
                            setSelectedTeam(team);
                            setShowEditTeamDialog(true);
                        }}
                        onDeleteTeam={handleDeleteTeam}
                        onAddMember={(team) => {
                            setSelectedTeam(team);
                            setShowAddMemberDialog(true);
                        }}
                        onRemoveMember={handleRemoveMember}
                    />
                </TabsContent> */}

                <TabsContent value="api" className="space-y-4">
                    <ApiKeysTab />
                </TabsContent>

                {/* <TabsContent value="preferences" className="space-y-4">
                    <PreferencesTab />
                </TabsContent> */}
            </Tabs>

            {/* Dialogs */}
            <AddUserDialog
                open={showAddUserDialog}
                onOpenChange={setShowAddUserDialog}
                name={newUserName}
                onNameChange={setNewUserName}
                email={newUserEmail}
                onEmailChange={setNewUserEmail}
                role={newUserRole}
                onRoleChange={setNewUserRole}
                roleOptions={roleOptions}
                onAdd={handleAddUser}
                isAdmin={isAdmin}
            />

            <EditUserDialog
                open={showEditUserDialog}
                onOpenChange={setShowEditUserDialog}
                user={selectedUser}
                onUserChange={setSelectedUser}
                roleOptions={roleOptions}
                getAvailableRoles={getAvailableRoles}
                onUpdate={handleUpdateUser}
                isAdmin={isAdmin}
            />

            {/* <AddTeamDialog
                open={showAddTeamDialog}
                onOpenChange={setShowAddTeamDialog}
                name={newTeamName}
                onNameChange={setNewTeamName}
                description={newTeamDesc}
                onDescriptionChange={setNewTeamDesc}
                onAdd={handleAddTeam}
            />

            <EditTeamDialog
                open={showEditTeamDialog}
                onOpenChange={setShowEditTeamDialog}
                team={selectedTeam}
                onTeamChange={setSelectedTeam}
                onUpdate={handleUpdateTeam}
            />

            <AddMemberDialog
                open={showAddMemberDialog}
                onOpenChange={setShowAddMemberDialog}
                team={selectedTeam}
                email={newMemberEmail}
                onEmailChange={setNewMemberEmail}
                role={newMemberRole}
                onRoleChange={setNewMemberRole}
                roleOptions={roleOptions}
                users={users}
                onAdd={handleAddMember}
            /> */}
        </div>
    );
}