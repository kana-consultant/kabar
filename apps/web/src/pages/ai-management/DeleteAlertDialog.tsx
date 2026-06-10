// pages/admin/ai-management/DeleteAlertDialog.tsx
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { AlertTriangle } from "lucide-react";

interface DeleteAlertDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  itemName: string;
  itemType: "provider" | "model";
  onConfirm: () => void;
}

export function DeleteAlertDialog({ open, onOpenChange, itemName, itemType, onConfirm }: DeleteAlertDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <div className="flex items-center space-x-2">
            <AlertTriangle className="h-5 w-5 text-red-500" />
            <DialogTitle>Delete {itemType}</DialogTitle>
          </div>
          <DialogDescription>
            Are you sure you want to delete{" "}
            <span className="font-semibold">{itemName}</span>?
            {itemType === "provider" && (
              <p className="mt-2 text-red-600">
                Warning: This will also delete all models associated with this provider.
              </p>
            )}
            This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        
        <div className="flex justify-end space-x-2">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={onConfirm}>
            Delete {itemType}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}