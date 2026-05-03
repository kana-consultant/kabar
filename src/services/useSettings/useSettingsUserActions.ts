import { toast } from "sonner";
import { addUser, updateUser, deleteUser } from "@/services/user";
import type { User } from "@/types/user";
import type { UserRoleType } from "./types";

export function useSettingsUserActions(
    loadData: () => Promise<void>,
    currentUser: User | null,
    users: User[],
    isSuperAdmin: boolean
) {
    const handleAddUser = async (
        newUserName: string,
        newUserEmail: string,
        newUserRole: UserRoleType,
        setShowAddUserDialog: (val: boolean) => void,
        setNewUserName: (val: string) => void,
        setNewUserEmail: (val: string) => void,
        setNewUserRole: (val: UserRoleType) => void
    ) => {
        if (!newUserName || !newUserEmail) {
            toast.error("Isi nama dan email");
            return;
        }
        try {
            await addUser({
                name: newUserName,
                email: newUserEmail,
                role: newUserRole,
                status: "active",
                emailVerified: false,
                twoFactorEnabled: false,
            });
            toast.success("User ditambahkan", { description: `${newUserName} (${newUserRole})` });
            await loadData();
            setShowAddUserDialog(false);
            setNewUserName("");
            setNewUserEmail("");
            setNewUserRole("viewer");
        } catch (error) {
            toast.error("Gagal menambahkan user");
        }
    };

    const handleUpdateUser = async (
        selectedUser: User | null,
        setShowEditUserDialog: (val: boolean) => void,
        setSelectedUser: (val: User | null) => void
    ) => {
        if (selectedUser) {
            try {
                await updateUser(selectedUser.id, {
                    name: selectedUser.name,
                    email: selectedUser.email,
                    role: selectedUser.role,
                    status: selectedUser.status,
                });
                toast.success("User diperbarui");
                await loadData();
                setShowEditUserDialog(false);
                setSelectedUser(null);
            } catch (error) {
                toast.error("Gagal memperbarui user");
            }
        }
    };

    const handleDeleteUser = async (id: string, name: string) => {
        if (id === currentUser?.id) {
            toast.error("Tidak bisa menghapus diri sendiri");
            return;
        }
        
        const userToDelete = users.find(u => u.id === id);
        if (userToDelete?.role === 'super_admin' && !isSuperAdmin) {
            toast.error("Tidak bisa menghapus Super Admin");
            return;
        }
        
        try {
            await deleteUser(id);
            toast.success("User dihapus", { description: name });
            await loadData();
        } catch (error) {
            toast.error("Gagal menghapus user");
        }
    };

    return {
        handleAddUser,
        handleUpdateUser,
        handleDeleteUser,
    };
}