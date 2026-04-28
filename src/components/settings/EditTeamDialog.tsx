import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import type { Team } from "@/types/user";

interface EditTeamDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    team: Team | null;
    onTeamChange: (team: Team) => void;
    onUpdate: () => void;
}

export function EditTeamDialog({
    open,
    onOpenChange,
    team,
    onTeamChange,
    onUpdate
}: EditTeamDialogProps) {
    if (!team) return null;

    const isValid = team.name.trim() !== "";

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Edit Tim</DialogTitle>
                    <DialogDescription>
                        Ubah informasi tim yang sudah ada
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                    <div>
                        <Label htmlFor="team-name">Nama Tim</Label>
                        <Input
                            id="team-name"
                            placeholder="Masukkan nama tim"
                            value={team.name}
                            onChange={(e) => onTeamChange({ ...team, name: e.target.value })}
                            className="mt-1"
                        />
                        {!team.name.trim() && (
                            <p className="mt-1 text-xs text-red-500">
                                Nama tim tidak boleh kosong
                            </p>
                        )}
                    </div>
                    <div>
                        <Label htmlFor="team-description">Deskripsi</Label>
                        <Textarea
                            id="team-description"
                            placeholder="Masukkan deskripsi tim (opsional)"
                            value={team.description || ""}
                            onChange={(e) => onTeamChange({ ...team, description: e.target.value })}
                            className="mt-1"
                            rows={3}
                        />
                        <p className="mt-1 text-xs text-slate-500">
                            Deskripsi singkat tentang tim ini
                        </p>
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Batal
                    </Button>
                    <Button onClick={onUpdate} disabled={!isValid}>
                        Simpan Perubahan
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}