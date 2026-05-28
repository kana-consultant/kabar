import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { SettingsHeader } from "./SettingsHeader";
import { ProfileTab } from "./ProfileTab";
import { UsersTab } from "./UsersTab";
import { ApiKeysTab } from "./ApiKeysTab";
import { AddUserDialog } from "./AddUserDialog";
import { EditUserDialog } from "./EditUserDialog";
import { LoadingSettings } from "./LoadingSettings";
import { useSettings } from "@/services/useSettings";

export default function Settings() {
    const {
        currentUser,
        users,
        loading,
        canManageUsers,
        isAdmin,
        isSuperAdmin,
        roleOptions,
        getAvailableRoles,
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
        handleAddUser,
        handleUpdateUser,
        handleDeleteUser,
        isAddingUser, 
        HandleChangePassword
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
                    <ProfileTab 
                    currentUser={currentUser}
                    onChangePassword = {HandleChangePassword} />
                </TabsContent>

                <TabsContent value="users" className="space-y-4">
                    <UsersTab
                        users={users}
                        currentUserId={currentUser?.id}
                        canManage={canManageUsers as any}
                        isAdmin={isAdmin}
                        isSuperAdmin = {isSuperAdmin}
                        onAddUser={() => setShowAddUserDialog(true)}
                        onEditUser={(user) => {
                            setSelectedUser(user);
                            setShowEditUserDialog(true);
                        }}
                        onDeleteUser={handleDeleteUser}
                    />
                </TabsContent>

                <TabsContent value="api" className="space-y-4">
                    <ApiKeysTab />
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
                isAdmin={isAdmin}
                isLoading={isAddingUser}

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
        </div>
    );
}