import { type ToastContextType } from "@/hooks/use-toast";
import { addTeam, updateTeam, deleteTeam } from "@/services/user";
import type { Team } from "@/services/user";

export function useSettingsTeamActions(loadData: () => Promise<void>) {
    const handleAddTeam = async (
        newTeamName: string,
        newTeamDesc: string,
        setShowAddTeamDialog: (val: boolean) => void,
        setNewTeamName: (val: string) => void,
        setNewTeamDesc: (val: string) => void,
        toast : ToastContextType
    ) => {
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

    const handleUpdateTeam = async (
        selectedTeam: Team | null,
        setShowEditTeamDialog: (val: boolean) => void,
        setSelectedTeam: (val: Team | null) => void,
        toast : ToastContextType
    ) => {
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

    const handleDeleteTeam = async (id: string, name: string, toast : ToastContextType) => {
        try {
            await deleteTeam(id);
            toast.success("Tim dihapus", { description: name });
            await loadData();
        } catch (error) {
            toast.error("Gagal menghapus tim");
        }
    };

    return {
        handleAddTeam,
        handleUpdateTeam,
        handleDeleteTeam,
    };
}