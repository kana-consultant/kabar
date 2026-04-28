import { useState, useEffect } from "react";
import { toast } from "sonner";
import { 
    getUsers, 
    addUser, 
    updateUser, 
    deleteUser, 
    getCurrentUser,
    getTeams,
    addTeam,
    updateTeam,
    deleteTeam,
    addTeamMember,
    removeTeamMember
} from "@/services/userService";
import type { User, Team, TeamMember, UserRole } from "@/types/user";

// Role options yang sesuai dengan User.role (string untuk User)
export type UserRoleType = 'super_admin' | 'admin' | 'manager' | 'editor' | 'viewer';

// Predefined UserRole objects untuk TeamMember.role
const USER_ROLES: Record<UserRoleType, UserRole> = {
    super_admin: {
        id: "role_super_admin",
        name: "super_admin",
        displayName: "Super Admin",
        description: "Full access to everything",
        scope: "system",
        level: 100
    },
    admin: {
        id: "role_admin",
        name: "admin",
        displayName: "Admin",
        description: "Admin access",
        scope: "global",
        level: 80
    },
    manager: {
        id: "role_manager",
        name: "manager",
        displayName: "Manager",
        description: "Can manage team",
        scope: "team",
        level: 60
    },
    editor: {
        id: "role_editor",
        name: "editor",
        displayName: "Editor",
        description: "Can edit content",
        scope: "team",
        level: 40
    },
    viewer: {
        id: "role_viewer",
        name: "viewer",
        displayName: "Viewer",
        description: "View only",
        scope: "team",
        level: 20
    }
};

export function useSettings() {
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

    const loadData = async () => {
        setLoading(true);
        try {
            const [usersData, teamsData, currentUserData] = await Promise.all([
                getUsers(),
                getTeams(),
                getCurrentUser(),
            ]);
            setUsers(usersData || []);
            setTeams(teamsData || []);
            setCurrentUser(currentUserData);
        } catch (error) {
            console.error("Failed to load settings data:", error);
            toast.error("Gagal memuat data pengaturan");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadData();
    }, []);

    // Permission checks berdasarkan role
    const canManageUsers = currentUser && ['super_admin', 'admin', 'manager'].includes(currentUser.role);
    const canManageTeams = currentUser && ['super_admin', 'admin', 'manager'].includes(currentUser.role);
    const isSuperAdmin = currentUser?.role === 'super_admin';
    const isAdmin = currentUser?.role === 'admin' || currentUser?.role === 'super_admin';

    const handleAddUser = async () => {
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

    const handleUpdateUser = async () => {
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

    const handleAddTeam = async () => {
        if (!newTeamName) {
            toast.error("Isi nama tim");
            return;
        }
        try {
            await addTeam({
                name: newTeamName,
                description: newTeamDesc,
                members: [],
            });
            toast.success("Tim ditambahkan", { description: newTeamName });
            await loadData();
            setShowAddTeamDialog(false);
            setNewTeamName("");
            setNewTeamDesc("");
        } catch (error) {
            toast.error("Gagal menambahkan tim");
        }
    };

    const handleUpdateTeam = async () => {
        if (selectedTeam) {
            try {
                await updateTeam(selectedTeam.id, {
                    name: selectedTeam.name,
                    description: selectedTeam.description,
                });
                toast.success("Tim diperbarui");
                await loadData();
                setShowEditTeamDialog(false);
                setSelectedTeam(null);
            } catch (error) {
                toast.error("Gagal memperbarui tim");
            }
        }
    };

    const handleDeleteTeam = async (id: string, name: string) => {
        try {
            await deleteTeam(id);
            toast.success("Tim dihapus", { description: name });
            await loadData();
        } catch (error) {
            toast.error("Gagal menghapus tim");
        }
    };

    const handleAddMember = async () => {
        if (!newMemberEmail) {
            toast.error("Isi email anggota");
            return;
        }
        
        const user = users.find(u => u.email === newMemberEmail);
        if (!user) {
            toast.error("User tidak ditemukan");
            return;
        }
        
        if (!selectedTeam) {
            toast.error("Pilih tim terlebih dahulu");
            return;
        }
        
        // Cek apakah user sudah menjadi member
        const isAlreadyMember = selectedTeam.members?.some(m => m.userId === user.id);
        if (isAlreadyMember) {
            toast.error("User sudah menjadi anggota tim ini");
            return;
        }
        
        // Buat TeamMember dengan role sebagai UserRole OBJECT
        const member: TeamMember = {
            userId: user.id,
            userEmail: user.email,
            userName: user.name,
            role: USER_ROLES[newMemberRole], // ← UserRole OBJECT, bukan string
            joinedAt: new Date().toISOString(),
        };
        
        try {
            await addTeamMember(selectedTeam.id, member);
            toast.success("Anggota ditambahkan", { description: user.name });
            await loadData();
            setShowAddMemberDialog(false);
            setNewMemberEmail("");
            setNewMemberRole("viewer");
        } catch (error) {
            toast.error("Gagal menambahkan anggota");
        }
    };

    const handleRemoveMember = async (teamId: string, userId: string, userName: string) => {
        if (userId === currentUser?.id) {
            toast.error("Tidak bisa menghapus diri sendiri dari tim");
            return;
        }
        
        try {
            await removeTeamMember(teamId, userId);
            toast.success("Anggota dihapus", { description: userName });
            await loadData();
        } catch (error) {
            toast.error("Gagal menghapus anggota");
        }
    };

    // Helper function untuk mendapatkan UserRole display name
    const getRoleDisplayName = (roleType: UserRoleType): string => {
        return USER_ROLES[roleType]?.displayName || roleType;
    };

    // Role options untuk dropdown
    const roleOptions: { value: UserRoleType; label: string }[] = [
        { value: 'viewer', label: 'Viewer' },
        { value: 'editor', label: 'Editor' },
        { value: 'manager', label: 'Manager' },
        { value: 'admin', label: 'Admin' },
        { value: 'super_admin', label: 'Super Admin' },
    ];

    // Filter role options berdasarkan permission current user
    const getAvailableRoles = (): UserRoleType[] => {
        if (isSuperAdmin) {
            return ['super_admin', 'admin', 'manager', 'editor', 'viewer'];
        }
        if (isAdmin) {
            return ['admin', 'manager', 'editor', 'viewer'];
        }
        if (currentUser?.role === 'manager') {
            return ['manager', 'editor', 'viewer'];
        }
        return ['viewer'];
    };

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
        getAvailableRoles,
        getRoleDisplayName,
        USER_ROLES,
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
        loadData,
    };
}