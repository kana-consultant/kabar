import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { User } from "@/services/user";
import { type UserRoleType } from "@/services/useSettings/types";

interface EditUserDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    user: User | null;
    onUserChange: (user: User) => void;
    roleOptions: { value: UserRoleType; label: string }[];
    getAvailableRoles: () => UserRoleType[];
    onUpdate: () => void;
    isAdmin: boolean
}

export function EditUserDialog({
    open,
    onOpenChange,
    user,
    onUserChange,
    roleOptions,
    getAvailableRoles,
    onUpdate,
    isAdmin
}: EditUserDialogProps) {
    if (!user) return null;

    const availableRoles = getAvailableRoles();
    const filteredRoleOptions = roleOptions.filter(opt => availableRoles.includes(opt.value));

    const isValidEmail = user.email && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(user.email);
    const isValid = user.name.trim() !== "" && isValidEmail;

    // Cek apakah user bisa diubah role nya (super_admin hanya bisa diubah oleh super_admin lain)
    const canChangeRole = user.role !== 'admin' || availableRoles.includes('admin');

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Edit Pengguna</DialogTitle>
                    <DialogDescription>Ubah informasi pengguna</DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div>
                        <Label>Nama Lengkap</Label>
                        <Input
                            value={user.name}
                            onChange={(e) => onUserChange({ ...user, name: e.target.value })}
                            placeholder="Nama lengkap"
                        />
                    </div>
                    <div>
                        <Label>Email</Label>
                        <Input
                            type="email"
                            value={user.email}
                            onChange={(e) => onUserChange({ ...user, email: e.target.value })}
                            placeholder="email@example.com"
                        />
                        {user.email && !isValidEmail && (
                            <p className="mt-1 text-xs text-red-500">
                                Format email tidak valid
                            </p>
                        )}
                    </div>
                    <div>
                        <Label>Role</Label>
                        <Select
                            value={user.role}
                            onValueChange={(v) => onUserChange({ ...user, role: v as UserRoleType })}
                            disabled={!canChangeRole}
                        >
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                {filteredRoleOptions.filter(option => !(option.value === 'admin' && !isAdmin))
                                    .map((option) => (
                                        <SelectItem key={option.value} value={option.value}>
                                            {option.label}
                                        </SelectItem>
                                    ))
                                }
                            </SelectContent>
                        </Select>
                        {!canChangeRole && (
                            <p className="mt-1 text-xs text-amber-500">
                                Role Super Admin tidak dapat diubah
                            </p>
                        )}
                        <p className="mt-1 text-xs text-slate-500">
                            Admin: akses penuh<br />
                            Manager: dapat mengelola tim<br />
                            Editor: dapat mengedit konten<br />
                            Viewer: hanya dapat melihat
                        </p>
                    </div>
                    <div>
                        <Label>Status</Label>
                        <Select
                            value={user.status}
                            onValueChange={(v) => onUserChange({ ...user, status: v as 'active' | 'inactive' | 'suspended' })}
                        >
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="active">Active</SelectItem>
                                <SelectItem value="inactive">Inactive</SelectItem>
                                <SelectItem value="suspended">Suspended</SelectItem>
                            </SelectContent>
                        </Select>
                        <p className="mt-1 text-xs text-slate-500">
                            Active: dapat mengakses sistem<br />
                            Inactive: tidak dapat mengakses<br />
                            Suspended: akun ditangguhkan
                        </p>
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Batal
                    </Button>
                    <Button onClick={onUpdate} disabled={!isValid}>
                        Update
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}