// src/components/drafts/PublishResultDialog.tsx
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { AlertCircle, CheckCircle, XCircle } from "lucide-react";

interface PublishResultDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    results: { success: boolean; message: string; results?: any[] } | null;
}

export function PublishResultDialog({ open, onOpenChange, results }: PublishResultDialogProps) {
    if (!results) return null;

    const successCount = results.results?.filter(r => r.success).length || 0;
    const failCount = results.results?.filter(r => !r.success).length || 0;

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-3xl">
                <DialogHeader>
                    <DialogTitle className="flex items-center gap-2">
                        {results.success ? (
                            <CheckCircle className="h-5 w-5 text-green-500" />
                        ) : (
                            <AlertCircle className="h-5 w-5 text-yellow-500" />
                        )}
                        Hasil Publikasi
                    </DialogTitle>
                    <DialogDescription>
                        {results.message}
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-3 max-h-96 overflow-y-auto">
                    <div className="flex gap-4 text-sm border-b pb-2">
                        <span className="text-green-600">✅ Berhasil: {successCount}</span>
                        <span className="text-red-600">❌ Gagal: {failCount}</span>
                    </div>

                    {results.results?.map((result, idx) => (
                        <div key={idx} className="text-sm border rounded-lg p-3">
                            <div className="flex items-center gap-2">
                                {result.success ? (
                                    <CheckCircle className="h-4 w-4 text-green-500" />
                                ) : (
                                    <XCircle className="h-4 w-4 text-red-500" />
                                )}
                                <span className="font-medium">{result.product}</span>
                            </div>
                            {!result.success && result.error && (
                                <p className="text-xs text-red-500 mt-1 ml-6">{result.error}</p>
                            )}
                            {result.success && result.response && (
                                <p className="text-xs text-green-500 mt-1 ml-6 truncate">{result.response}</p>
                            )}
                        </div>
                    ))}
                </div>

                <DialogFooter>
                    <Button onClick={() => onOpenChange(false)}>Tutup</Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}