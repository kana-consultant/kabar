import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { Team, User } from "@/types/user";

type UserRoleType = 'super_admin' | 'admin' | 'manager' | 'editor' | 'viewer';

interface AddMemberDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    team: Team | null;
    email: string;
    onEmailChange: (value: string) => void;
    role: UserRoleType;
    onRoleChange: (value: UserRoleType) => void;
    roleOptions: { value: UserRoleType; label: string }[];
    users: User[];
    onAdd: () => void;
}

export function AddMemberDialog({
    open,
    onOpenChange,
    team,
    email,
    onEmailChange,
    role,
    onRoleChange,
    roleOptions,
    users,
    onAdd,
}: AddMemberDialogProps) {
    if (!team) return null;

    // Check if user exists by email
    const userExists = users.some(u => u.email === email);
    const isValidEmail = email && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Tambah Anggota Tim</DialogTitle>
                    <DialogDescription>
                        Tambahkan anggota ke tim {team.name}
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div>
                        <Label>Email User</Label>
                        <Input
                            placeholder="email@example.com"
                            value={email}
                            onChange={(e) => onEmailChange(e.target.value)}
                        />
                        {email && !userExists && (
                            <p className="mt-1 text-xs text-red-500">
                                User dengan email ini tidak ditemukan
                            </p>
                        )}
                        {email && userExists && (
                            <p className="mt-1 text-xs text-green-500">
                                User ditemukan
                            </p>
                        )}
                        <p className="mt-1 text-xs text-slate-500">
                            Masukkan email user yang sudah terdaftar di sistem
                        </p>
                    </div>
                    <div>
                        <Label>Role dalam Tim</Label>
                        <Select value={role} onValueChange={(v) => onRoleChange(v as UserRoleType)}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                {roleOptions
                                    .filter(opt => opt.value !== 'super_admin' && opt.value !== 'admin')
                                    .map((option) => (
                                        <SelectItem key={option.value} value={option.value}>
                                            {option.label}
                                        </SelectItem>
                                    ))}
                            </SelectContent>
                        </Select>
                        <p className="mt-1 text-xs text-slate-500">
                            Manager: dapat mengelola tim<br />
                            Editor: dapat mengedit konten<br />
                            Viewer: hanya dapat melihat
                        </p>
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Batal
                    </Button>
                    <Button 
                        onClick={onAdd} 
                        disabled={!email || !userExists || !isValidEmail}
                    >
                        Tambah Anggota
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}