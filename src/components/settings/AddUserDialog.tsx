import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

type UserRoleType = 'super_admin' | 'admin' | 'manager' | 'editor' | 'viewer';

interface AddUserDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    name: string;
    onNameChange: (value: string) => void;
    email: string;
    onEmailChange: (value: string) => void;
    role: UserRoleType;
    onRoleChange: (value: UserRoleType) => void;
    roleOptions: { value: UserRoleType; label: string }[];
    onAdd: () => void;
}

export function AddUserDialog({
    open,
    onOpenChange,
    name,
    onNameChange,
    email,
    onEmailChange,
    role,
    onRoleChange,
    roleOptions,
    onAdd,
}: AddUserDialogProps) {
    const isValidEmail = email && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
    const isValid = name.trim() !== "" && isValidEmail;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Tambah Pengguna Baru</DialogTitle>
                    <DialogDescription>Tambahkan pengguna baru ke sistem</DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div>
                        <Label>Nama Lengkap</Label>
                        <Input 
                            placeholder="John Doe" 
                            value={name} 
                            onChange={(e) => onNameChange(e.target.value)} 
                        />
                    </div>
                    <div>
                        <Label>Email</Label>
                        <Input 
                            type="email" 
                            placeholder="email@example.com"
                            value={email} 
                            onChange={(e) => onEmailChange(e.target.value)} 
                        />
                        {email && !isValidEmail && (
                            <p className="mt-1 text-xs text-red-500">
                                Format email tidak valid
                            </p>
                        )}
                    </div>
                    <div>
                        <Label>Role</Label>
                        <Select value={role} onValueChange={(v) => onRoleChange(v as UserRoleType)}>
                            <SelectTrigger>
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                {roleOptions.map((option) => (
                                    <SelectItem key={option.value} value={option.value}>
                                        {option.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        <p className="mt-1 text-xs text-slate-500">
                            Admin: akses penuh<br />
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
                    <Button onClick={onAdd} disabled={!isValid}>
                        Tambah Pengguna
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}