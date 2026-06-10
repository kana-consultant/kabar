// pages/admin/ai-management/components/ProviderForm/FormActions.tsx

import { Button } from "@kana-consultant/ui-kit";

interface FormActionsProps {
    mode: "create" | "edit";
    isSubmitting: boolean;
    onCancel: () => void;
}

export function FormActions({ mode, isSubmitting, onCancel }: FormActionsProps) {
    return (
        <div className="flex justify-end space-x-3 pt-6 border-t">
            <Button
                type="button"
                variant="outline"
                onClick={onCancel}
                disabled={isSubmitting}
            >
                Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? (
                    <>
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                        {mode === "create" ? "Creating..." : "Saving..."}
                    </>
                ) : (
                    <>{mode === "create" ? "Create Provider" : "Save Changes"}</>
                )}
            </Button>
        </div>
    );
}