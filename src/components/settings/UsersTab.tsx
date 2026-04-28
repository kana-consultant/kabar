import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus, Edit2, Trash2, Shield, UserCog, Eye, Pencil } from "lucide-react";
import type { User } from "@/types/user";

type UserRoleType = 'super_admin' | 'admin' | 'manager' | 'editor' | 'viewer';

const roleLabels: Record<UserRoleType, string> = {
    super_admin: "Super Admin",
    admin: "Administrator",
    manager: "Manager",
    editor: "Editor",
    viewer: "Viewer",
};

const roleColors: Record<UserRoleType, string> = {
    super_admin: "text-purple-600 bg-purple-50 dark:bg-purple-950",
    admin: "text-red-600 bg-red-50 dark:bg-red-950",
    manager: "text-blue-600 bg-blue-50 dark:bg-blue-950",
    editor: "text-green-600 bg-green-50 dark:bg-green-950",
    viewer: "text-gray-600 bg-gray-50 dark:bg-gray-950",
};

const roleIcons: Record<UserRoleType, React.ElementType> = {
    super_admin: Shield,
    admin: Shield,
    manager: UserCog,
    editor: Pencil,
    viewer: Eye,
};

interface UsersTabProps {
    users: User[];
    currentUserId?: string;
    canManage: boolean;
    isSuperAdmin: boolean;
    isAdmin: boolean;
    roleOptions: { value: UserRoleType; label: string }[];
    getAvailableRoles: () => UserRoleType[];
    getRoleDisplayName: (role: UserRoleType) => string;
    onAddUser: () => void;
    onEditUser: (user: User) => void;
    onDeleteUser: (id: string, name: string) => void;
}

export function UsersTab({ 
    users, 
    currentUserId, 
    canManage,
    isSuperAdmin,
    isAdmin,
    onAddUser, 
    onEditUser, 
    onDeleteUser 
}: UsersTabProps) {
    
    // Filter users based on permission
    const getVisibleUsers = () => {
        if (isSuperAdmin) {
            return users;
        }
        if (isAdmin) {
            // Admin can see all except super_admin
            return users.filter(u => u.role !== 'super_admin');
        }
        return users;
    };

    const visibleUsers = getVisibleUsers();

    // Check if user can be edited
    const canEditUser = (user: User) => {
        if (user.id === currentUserId) return true;
        if (user.role === 'super_admin' && !isSuperAdmin) return false;
        if (user.role === 'admin' && !isSuperAdmin && !isAdmin) return false;
        return canManage;
    };

    // Check if user can be deleted
    const canDeleteUser = (user: User) => {
        if (user.id === currentUserId) return false;
        if (user.role === 'super_admin' && !isSuperAdmin) return false;
        if (user.role === 'admin' && !isSuperAdmin && !isAdmin) return false;
        return canManage;
    };

    // Get status badge color
    const getStatusColor = (status: string) => {
        switch (status) {
            case 'active':
                return 'text-green-600 bg-green-50 dark:bg-green-950';
            case 'inactive':
                return 'text-yellow-600 bg-yellow-50 dark:bg-yellow-950';
            case 'suspended':
                return 'text-red-600 bg-red-50 dark:bg-red-950';
            default:
                return 'text-gray-600 bg-gray-50 dark:bg-gray-950';
        }
    };

    const getStatusLabel = (status: string) => {
        switch (status) {
            case 'active':
                return 'Active';
            case 'inactive':
                return 'Inactive';
            case 'suspended':
                return 'Suspended';
            default:
                return status;
        }
    };

    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <div>
                    <CardTitle>Manajemen Pengguna</CardTitle>
                    <CardDescription>Kelola pengguna dan hak akses mereka</CardDescription>
                </div>
                {canManage && (
                    <Button onClick={onAddUser}>
                        <Plus className="mr-2 h-4 w-4" />
                        Tambah Pengguna
                    </Button>
                )}
            </CardHeader>
            <CardContent>
                {visibleUsers.length === 0 ? (
                    <div className="text-center py-8 text-slate-500">
                        Belum ada pengguna yang terdaftar
                    </div>
                ) : (
                    <div className="space-y-3">
                        {visibleUsers.map((user) => {
                            const RoleIcon = roleIcons[user.role as UserRoleType] || UserCog;
                            const canEdit = canEditUser(user);
                            const canDelete = canDeleteUser(user);
                            
                            return (
                                <div key={user.id} className="flex items-center justify-between rounded-lg border p-4">
                                    <div className="flex items-center gap-3">
                                        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-r from-blue-500 to-blue-600 text-sm font-semibold text-white">
                                            {user.name.charAt(0).toUpperCase()}
                                        </div>
                                        <div>
                                            <p className="font-medium">{user.name}</p>
                                            <p className="text-sm text-slate-500">{user.email}</p>
                                            <div className="flex items-center gap-2 mt-1">
                                                <span className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs ${roleColors[user.role as UserRoleType]}`}>
                                                    <RoleIcon className="h-3 w-3" />
                                                    {roleLabels[user.role as UserRoleType]}
                                                </span>
                                                <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs ${getStatusColor(user.status)}`}>
                                                    {getStatusLabel(user.status)}
                                                </span>
                                                {user.id === currentUserId && (
                                                    <span className="text-xs text-blue-600">(Anda)</span>
                                                )}
                                                {user.emailVerified === false && (
                                                    <span className="text-xs text-yellow-600">(Email belum diverifikasi)</span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                    {canManage && (
                                        <div className="flex gap-2">
                                            <Button 
                                                variant="outline" 
                                                size="sm" 
                                                onClick={() => onEditUser(user)}
                                                disabled={!canEdit}
                                                title={!canEdit ? "Tidak memiliki izin untuk mengedit user ini" : "Edit user"}
                                            >
                                                <Edit2 className="h-4 w-4" />
                                            </Button>
                                            <Button 
                                                variant="destructive" 
                                                size="sm" 
                                                onClick={() => onDeleteUser(user.id, user.name)}
                                                disabled={!canDelete}
                                                title={!canDelete ? "Tidak memiliki izin untuk menghapus user ini" : "Hapus user"}
                                            >
                                                <Trash2 className="h-4 w-4" />
                                            </Button>
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                )}
                
                {/* Stats summary */}
                {users.length > 0 && (
                    <div className="mt-6 pt-4 border-t text-xs text-slate-500">
                        <div className="flex gap-4">
                            <span>Total: {users.length} users</span>
                            <span>Active: {users.filter(u => u.status === 'active').length}</span>
                            <span>Inactive: {users.filter(u => u.status === 'inactive').length}</span>
                            <span>Suspended: {users.filter(u => u.status === 'suspended').length}</span>
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}