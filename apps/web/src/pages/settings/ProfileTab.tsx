import { useState } from "react";
import { Input, Button, Label, Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from  "@kana-consultant/ui-kit";
import { cn } from "@/lib/utils";
import type { User } from "@/types/user";
import { Eye, EyeOff, Lock, KeyRound } from "lucide-react";
import { Alert } from  "@kana-consultant/ui-kit";

const roleConfig: Record<string, { label: string; class: string }> = {
    admin: {
        label: "Administrator",
        class: "bg-red-50 text-red-700 border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20",
    },
    manager: {
        label: "Manager",
        class: "bg-blue-50 text-blue-700 border-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20",
    },
    editor: {
        label: "Editor",
        class: "bg-green-50 text-green-700 border-green-200/60 dark:bg-green-500/10 dark:text-green-400 dark:border-green-500/20",
    },
    viewer: {
        label: "Viewer",
        class: "bg-slate-100 text-slate-600 border-slate-200/60 dark:bg-white/[0.05] dark:text-slate-400 dark:border-white/[0.08]",
    },
};

interface ProfileTabProps {
    currentUser: User | null;
    onChangePassword?: (oldPassword: string, newPassword: string) => Promise<void>;
}

export function ProfileTab({ currentUser, onChangePassword }: ProfileTabProps) {
    const role = currentUser?.role ?? "viewer";
    const cfg = roleConfig[role] ?? roleConfig.viewer;
    const initials = currentUser?.name?.charAt(0)?.toUpperCase() ?? "?";

    // State untuk change password
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [showOldPassword, setShowOldPassword] = useState(false);
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);
    const [oldPassword, setOldPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // Validasi password
    const validatePassword = (password: string): string | null => {
        if (password.length < 8) {
            return "Password harus minimal 8 karakter";
        }
        if (!/[A-Z]/.test(password)) {
            return "Password harus mengandung minimal 1 huruf kapital";
        }
        if (!/[a-z]/.test(password)) {
            return "Password harus mengandung minimal 1 huruf kecil";
        }
        if (!/[0-9]/.test(password)) {
            return "Password harus mengandung minimal 1 angka";
        }
        return null;
    };

    const handleSubmit = async () => {
        // Reset messages
        setError(null);
        setSuccess(null);

        // Validasi input
        if (!oldPassword) {
            setError("Password lama harus diisi");
            return;
        }

        if (!newPassword) {
            setError("Password baru harus diisi");
            return;
        }

        const passwordError = validatePassword(newPassword);
        if (passwordError) {
            setError(passwordError);
            return;
        }

        if (newPassword !== confirmPassword) {
            setError("Konfirmasi password tidak sesuai");
            return;
        }

        if (oldPassword === newPassword) {
            setError("Password baru harus berbeda dengan password lama");
            return;
        }

        setIsLoading(true);
        try {
            if (onChangePassword) {
                await onChangePassword(oldPassword, newPassword);
                setSuccess("Password berhasil diubah");
                // Reset form
                setOldPassword("");
                setNewPassword("");
                setConfirmPassword("");
                // Close dialog after 2 seconds
                setTimeout(() => {
                    setIsDialogOpen(false);
                    setSuccess(null);
                }, 2000);
            } else {
                setError("Fungsi ubah password tidak tersedia");
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : "Gagal mengubah password");
        } finally {
            setIsLoading(false);
        }
    };

    const handleCloseDialog = () => {
        setIsDialogOpen(false);
        setOldPassword("");
        setNewPassword("");
        setConfirmPassword("");
        setError(null);
        setSuccess(null);
        setShowOldPassword(false);
        setShowNewPassword(false);
        setShowConfirmPassword(false);
    };

    return (
        <>
            <div className={cn(
                "overflow-hidden rounded-2xl border",
                "bg-white border-slate-200/80",
                "dark:bg-[#0f0d1a] dark:border-white/[0.06]"
            )}>
                {/* Card header */}
                <div className={cn(
                    "px-6 py-4 border-b",
                    "border-slate-100 bg-slate-50/60",
                    "dark:border-white/[0.05] dark:bg-white/[0.02]"
                )}>
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium text-slate-800 dark:text-slate-100">
                                Profile Information
                            </p>
                            <p className="text-xs text-slate-400 dark:text-slate-600 mt-0.5">
                                Informasi akun Anda
                            </p>
                        </div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setIsDialogOpen(true)}
                            className="gap-2"
                        >
                            <KeyRound className="h-3.5 w-3.5" />
                            Ubah Password
                        </Button>
                    </div>
                </div>

                {/* Card body */}
                <div className="px-6 py-5 space-y-5">
                    {/* Avatar row */}
                    <div className={cn(
                        "flex items-center gap-4 pb-5 border-b",
                        "border-slate-100 dark:border-white/[0.05]"
                    )}>
                        <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl bg-gradient-to-br from-green-500 to-emerald-600 dark:from-purple-600 dark:to-violet-700 text-white text-2xl font-semibold shadow-sm">
                            {initials}
                        </div>
                        <div>
                            <p className="text-sm font-semibold text-slate-900 dark:text-white">
                                {currentUser?.name}
                            </p>
                            <p className="text-xs text-slate-400 dark:text-slate-500 mt-0.5">
                                {currentUser?.email}
                            </p>
                            <span className={cn(
                                "mt-2 inline-flex items-center rounded-md border px-2 py-0.5 text-[10px] font-medium",
                                cfg.class
                            )}>
                                {cfg.label}
                            </span>
                        </div>
                    </div>

                    {/* Fields */}
                    <div className="grid gap-4 sm:grid-cols-2">
                        {[
                            { label: "Name", value: currentUser?.name },
                            { label: "Email", value: currentUser?.email },
                        ].map(({ label, value }) => (
                            <div key={label} className="space-y-1.5">
                                <Label className="text-[11px] font-medium uppercase tracking-wide text-slate-400 dark:text-slate-600">
                                    {label}
                                </Label>
                                <Input
                                    defaultValue={value}
                                    disabled
                                    className={cn(
                                        "h-8 text-sm",
                                        "border-slate-200/80 bg-slate-50/80 text-slate-600",
                                        "dark:border-white/[0.06] dark:bg-white/[0.02] dark:text-slate-400",
                                        "disabled:opacity-100 disabled:cursor-default"
                                    )}
                                />
                            </div>
                        ))}
                    </div>
                </div>
            </div>

            {/* Change Password Dialog */}
            <Dialog open={isDialogOpen} onOpenChange={handleCloseDialog}>
                <DialogContent className="sm:max-w-[425px]">
                    <DialogHeader>
                        <DialogTitle className="flex items-center gap-2">
                            <Lock className="h-4 w-4" />
                            Ubah Password
                        </DialogTitle>
                        <DialogDescription>
                            Masukkan password lama Anda dan buat password baru yang aman.
                        </DialogDescription>
                    </DialogHeader>

                    <div className="space-y-4 py-4">
                        {error && (
                            <Alert tone="danger" title={error} className="py-2" />
                        )}

                        {success && (
                            <Alert tone="success" title={success} className="py-2" />
                        )}

                        {/* Password Lama */}
                        <div className="space-y-2">
                            <Label htmlFor="old-password" className="text-sm font-medium">
                                Password Lama
                            </Label>
                            <div className="relative">
                                <Input
                                    id="old-password"
                                    type={showOldPassword ? "text" : "password"}
                                    value={oldPassword}
                                    onChange={(e: any) => setOldPassword(e.target.value)}
                                    placeholder="Masukkan password lama"
                                    className="pr-10"
                                    disabled={isLoading}
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => setShowOldPassword(!showOldPassword)}
                                    className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
                                >
                                    {showOldPassword ? (
                                        <EyeOff className="h-4 w-4" />
                                    ) : (
                                        <Eye className="h-4 w-4" />
                                    )}
                                </Button>
                            </div>
                        </div>

                        {/* Password Baru */}
                        <div className="space-y-2">
                            <Label htmlFor="new-password" className="text-sm font-medium">
                                Password Baru
                            </Label>
                            <div className="relative">
                                <Input
                                    id="new-password"
                                    type={showNewPassword ? "text" : "password"}
                                    value={newPassword}
                                    onChange={(e: any) => setNewPassword(e.target.value)}
                                    placeholder="Masukkan password baru"
                                    className="pr-10"
                                    disabled={isLoading}
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => setShowNewPassword(!showNewPassword)}
                                    className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
                                >
                                    {showNewPassword ? (
                                        <EyeOff className="h-4 w-4" />
                                    ) : (
                                        <Eye className="h-4 w-4" />
                                    )}
                                </Button>
                            </div>
                            <div className="text-xs text-slate-400 space-y-1">
                                <p>Password harus memenuhi kriteria:</p>
                                <ul className="list-disc list-inside space-y-0.5 ml-2">
                                    <li>Minimal 8 karakter</li>
                                    <li>Mengandung huruf kapital (A-Z)</li>
                                    <li>Mengandung huruf kecil (a-z)</li>
                                    <li>Mengandung angka (0-9)</li>
                                </ul>
                            </div>
                        </div>

                        {/* Konfirmasi Password Baru */}
                        <div className="space-y-2">
                            <Label htmlFor="confirm-password" className="text-sm font-medium">
                                Konfirmasi Password Baru
                            </Label>
                            <div className="relative">
                                <Input
                                    id="confirm-password"
                                    type={showConfirmPassword ? "text" : "password"}
                                    value={confirmPassword}
                                    onChange={(e: any) => setConfirmPassword(e.target.value)}
                                    placeholder="Konfirmasi password baru"
                                    className="pr-10"
                                    disabled={isLoading}
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                                    className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 p-0"
                                >
                                    {showConfirmPassword ? (
                                        <EyeOff className="h-4 w-4" />
                                    ) : (
                                        <Eye className="h-4 w-4" />
                                    )}
                                </Button>
                            </div>
                        </div>
                    </div>

                    <DialogFooter>
                        <Button
                            variant="outline"
                            onClick={handleCloseDialog}
                            disabled={isLoading}
                        >
                            Batal
                        </Button>
                        <Button
                            onClick={handleSubmit}
                            disabled={isLoading}
                            className="gap-2"
                        >
                            {isLoading ? (
                                <>
                                    <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                                    Memproses...
                                </>
                            ) : (
                                <>
                                    <KeyRound className="h-4 w-4" />
                                    Ubah Password
                                </>
                            )}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}