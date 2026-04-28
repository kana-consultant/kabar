import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus, Edit2, Trash2, UserPlus, Trash2 as TrashIcon } from "lucide-react";
import type { Team, UserRole } from "@/types/user";

type UserRoleType = 'super_admin' | 'admin' | 'manager' | 'editor' | 'viewer';


const roleColors: Record<UserRoleType, string> = {
    super_admin: "text-purple-600 bg-purple-50 dark:bg-purple-950",
    admin: "text-red-600 bg-red-50 dark:bg-red-950",
    manager: "text-blue-600 bg-blue-50 dark:bg-blue-950",
    editor: "text-green-600 bg-green-50 dark:bg-green-950",
    viewer: "text-gray-600 bg-gray-50 dark:bg-gray-950",
};

interface TeamsTabProps {
    teams: Team[];
    currentUserId?: string;
    canManage: boolean;
    roleOptions: { value: UserRoleType; label: string }[];
    getRoleDisplayName: (role: UserRoleType) => string;
    onAddTeam: () => void;
    onEditTeam: (team: Team) => void;
    onDeleteTeam: (id: string, name: string) => void;
    onAddMember: (team: Team) => void;
    onRemoveMember: (teamId: string, userId: string, userName: string) => void;
}

export function TeamsTab({ 
    teams, 
    currentUserId,
    canManage, 
    onAddTeam, 
    onEditTeam, 
    onDeleteTeam, 
    onAddMember, 
    onRemoveMember 
}: TeamsTabProps) {
    

    // Helper untuk mendapatkan role display name dari object UserRole
    const getRoleDisplayNameFromObject = (role: UserRole): string => {
        return role.displayName || role.name;
    };

    // Helper untuk mendapatkan role color dari object UserRole
    const getRoleColor = (role: UserRole): string => {
        const roleName = role.name as UserRoleType;
        return roleColors[roleName] || roleColors.viewer;
    };

    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between">
                <div>
                    <CardTitle>Manajemen Tim</CardTitle>
                    <CardDescription>Kelola tim dan anggota tim</CardDescription>
                </div>
                {canManage && (
                    <Button onClick={onAddTeam}>
                        <Plus className="mr-2 h-4 w-4" />
                        Tambah Tim
                    </Button>
                )}
            </CardHeader>
            <CardContent>
                {teams.length === 0 ? (
                    <div className="text-center py-8 text-slate-500">
                        Belum ada tim yang dibuat
                    </div>
                ) : (
                    <div className="space-y-6">
                        {teams.map((team) => (
                            <div key={team.id} className="rounded-lg border p-4">
                                <div className="flex items-center justify-between">
                                    <div>
                                        <h3 className="font-semibold">{team.name}</h3>
                                        <p className="text-sm text-slate-500">{team.description || "Tidak ada deskripsi"}</p>
                                        <p className="text-xs text-slate-400 mt-1">
                                            Dibuat: {new Date(team.createdAt).toLocaleDateString('id-ID')}
                                        </p>
                                    </div>
                                    {canManage && (
                                        <div className="flex gap-2">
                                            <Button variant="outline" size="sm" onClick={() => onAddMember(team)}>
                                                <UserPlus className="mr-1 h-3 w-3" />
                                                Tambah Anggota
                                            </Button>
                                            <Button variant="outline" size="sm" onClick={() => onEditTeam(team)}>
                                                <Edit2 className="h-4 w-4" />
                                            </Button>
                                            <Button variant="destructive" size="sm" onClick={() => onDeleteTeam(team.id, team.name)}>
                                                <Trash2 className="h-4 w-4" />
                                            </Button>
                                        </div>
                                    )}
                                </div>
                                
                                <div className="mt-4">
                                    <h4 className="text-sm font-medium">
                                        Anggota ({team.members?.length || 0})
                                    </h4>
                                    <div className="mt-2 space-y-2">
                                        {team.members && team.members.length > 0 ? (
                                            team.members.map((member) => (
                                                <div key={member.userId} className="flex items-center justify-between rounded-lg bg-slate-50 p-2 dark:bg-slate-900">
                                                    <div>
                                                        <p className="text-sm font-medium">{member.userName}</p>
                                                        <p className="text-xs text-slate-500">{member.userEmail}</p>
                                                        <p className="text-xs text-slate-400">
                                                            Bergabung: {new Date(member.joinedAt).toLocaleDateString('id-ID')}
                                                        </p>
                                                    </div>
                                                    <div className="flex items-center gap-2">
                                                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs ${getRoleColor(member.role)}`}>
                                                            {getRoleDisplayNameFromObject(member.role)}
                                                        </span>
                                                        {canManage && member.userId !== currentUserId && (
                                                            <Button 
                                                                variant="ghost" 
                                                                size="sm" 
                                                                onClick={() => onRemoveMember(team.id, member.userId, member.userName)}
                                                                className="text-red-500 hover:text-red-700 hover:bg-red-50"
                                                            >
                                                                <TrashIcon className="h-3 w-3" />
                                                            </Button>
                                                        )}
                                                        {member.userId === currentUserId && (
                                                            <span className="text-xs text-blue-600">(Anda)</span>
                                                        )}
                                                    </div>
                                                </div>
                                            ))
                                        ) : (
                                            <p className="text-sm text-slate-500">Belum ada anggota dalam tim ini</p>
                                        )}
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </CardContent>
        </Card>
    );
}