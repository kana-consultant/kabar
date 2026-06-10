// pages/admin/ai-management/components/DangerZone.tsx

import { useState } from "react";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@kana-consultant/ui-kit";
import { Trash2, AlertTriangle } from "lucide-react";

interface DangerZoneProps {
    onDelete: () => void;
    providerName: string;
}

export function DangerZone({ onDelete, providerName }: DangerZoneProps) {
    const [confirmText, setConfirmText] = useState("");
    const [isOpen, setIsOpen] = useState(false);

    const handleDelete = () => {
        if (confirmText === providerName) {
            onDelete();
            setIsOpen(false);
            setConfirmText("");
        }
    };

    return (
        <Card className="border-red-200 dark:border-red-800">
            <CardHeader>
                <CardTitle className="text-red-600 dark:text-red-400 flex items-center gap-2">
                    <AlertTriangle className="h-5 w-5" />
                    Danger Zone
                </CardTitle>
                <CardDescription>
                    Irreversible actions that cannot be undone
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="border border-red-200 dark:border-red-800 rounded-lg p-4">
                    <div className="flex items-center justify-between">
                        <div>
                            <h4 className="font-semibold text-red-600 dark:text-red-400">
                                Delete Provider
                            </h4>
                            <p className="text-sm text-muted-foreground mt-1">
                                Permanently delete "{providerName}" and all its configurations.
                                This action cannot be undone.
                            </p>
                        </div>
                        <Button variant="destructive" onClick={() => setIsOpen(true)}>
                            <Trash2 className="h-4 w-4 mr-2" />
                            Delete Provider
                        </Button>
                    </div>
                </div>
            </CardContent>

            <Dialog open={isOpen} onOpenChange={setIsOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Are you absolutely sure?</DialogTitle>
                        <DialogDescription>
                            This action cannot be undone. This will permanently delete
                            the provider "{providerName}" and remove all associated data.
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-2 py-4">
                        <Label>
                            Type <span className="font-bold">{providerName}</span> to confirm:
                        </Label>
                        <Input
                            value={confirmText}
                            onChange={(e) => setConfirmText(e.target.value)}
                            placeholder={providerName}
                            className="font-mono"
                        />
                    </div>
                    <div className="flex justify-end gap-2">
                        <Button variant="outline" onClick={() => {
                            setIsOpen(false);
                            setConfirmText("");
                        }}>
                            Cancel
                        </Button>
                        <Button
                            variant="destructive"
                            onClick={handleDelete}
                            disabled={confirmText !== providerName}
                        >
                            Delete Forever
                        </Button>
                    </div>
                </DialogContent>
            </Dialog>
        </Card>
    );
}