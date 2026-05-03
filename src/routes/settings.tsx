import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { SettingsHeader } from "@/components/settings/SettingsHeader";
import { ProfileTab } from "@/components/settings/ProfileTab";
import { UsersTab } from "@/components/settings/UsersTab";
import { TeamsTab } from "@/components/settings/TeamsTab";
import { ApiKeysTab } from "@/components/settings/ApiKeysTab";
import { PreferencesTab } from "@/components/settings/PreferencesTab";
import { AddUserDialog } from "@/components/settings/AddUserDialog";
import { EditUserDialog } from "@/components/settings/EditUserDialog";
import { AddTeamDialog } from "@/components/settings/AddTeamDialog";
import { AddMemberDialog } from "@/components/settings/AddMemberDialog";
import { EditTeamDialog } from "@/components/settings/EditTeamDialog";
import { LoadingSettings } from "@/components/settings/LoadingSettings";
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
                <TabsList className="grid w-full grid-cols-5">
                    <TabsTrigger value="profile">Profile</TabsTrigger>
                    <TabsTrigger value="users">Users</TabsTrigger>
                    <TabsTrigger value="teams">Teams</TabsTrigger>
                    <TabsTrigger value="api">API Keys</TabsTrigger>
                    <TabsTrigger value="preferences">Preferences</TabsTrigger>
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

                <TabsContent value="teams" className="space-y-4">
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
                </TabsContent>

                <TabsContent value="api" className="space-y-4">
                    <ApiKeysTab />
                </TabsContent>

                <TabsContent value="preferences" className="space-y-4">
                    <PreferencesTab />
                </TabsContent>
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
            />

            <EditUserDialog
                open={showEditUserDialog}
                onOpenChange={setShowEditUserDialog}
                user={selectedUser}
                onUserChange={setSelectedUser}
                roleOptions={roleOptions}
                getAvailableRoles={getAvailableRoles}
                onUpdate={handleUpdateUser}
            />

            <AddTeamDialog
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
            />
        </div>
    );
}