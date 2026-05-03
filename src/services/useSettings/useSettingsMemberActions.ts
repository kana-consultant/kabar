import { toast } from "sonner";
import { addTeamMember, removeTeamMember } from "@/services/user";
import type { User, Team } from "@/types/user";
import type { UserRoleType } from "./types";
import { USER_ROLES } from "./types";

export function useSettingsMemberActions(
    loadData: () => Promise<void>,
    users: User[],
    currentUser: User | null
) {
    const handleAddMember = async (
        newMemberEmail: string,
        newMemberRole: UserRoleType,
        selectedTeam: Team | null,
        setShowAddMemberDialog: (val: boolean) => void,
        setNewMemberEmail: (val: string) => void,
        setNewMemberRole: (val: UserRoleType) => void
    ) => {
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
        const member = {
            userId: user.id,
            userEmail: user.email,
            userName: user.name,
            role: USER_ROLES[newMemberRole],
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

    return {
        handleAddMember,
        handleRemoveMember,
    };
}