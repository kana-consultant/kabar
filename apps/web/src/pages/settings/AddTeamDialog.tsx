import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, Button, Input, Label } from  "@kana-consultant/ui-kit";

interface AddTeamDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    name: string;
    onNameChange: (value: string) => void;
    description: string;
    onDescriptionChange: (value: string) => void;
    onAdd: () => void;
}

export function AddTeamDialog({
    open,
    onOpenChange,
    name,
    onNameChange,
    description,
    onDescriptionChange,
    onAdd,
}: AddTeamDialogProps) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Add New Team</DialogTitle>
                    <DialogDescription>Buat tim baru untuk kolaborasi</DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <div>
                        <Label>Team Name</Label>
                        <Input value={name} onChange={(e : any) => onNameChange(e.target.value)} />
                    </div>
                    <div>
                        <Label>Description</Label>
                        <Input value={description} onChange={(e : any) => onDescriptionChange(e.target.value)} />
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>Cancel</Button>
                    <Button onClick={onAdd}>Add Team</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}