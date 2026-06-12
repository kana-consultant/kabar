// src/pages/products/WorkflowBuilder/WorkflowTestDialog.tsx
import { useState } from "react";
import { Button, Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";

interface WorkflowTestDialogProps {
    open: boolean;
    onClose: () => void;
    workflowId: string;
}

export function WorkflowTestDialog({ open, onClose, workflowId }: WorkflowTestDialogProps) {
    const [logs, setLogs] = useState<string[]>([]);
    const [running, setRunning] = useState(false);
    const toast = useToast();

    const handleTest = async () => {
        setRunning(true);
        setLogs([]);

        try {
            // TODO: Call execute endpoint
            setLogs((prev) => [...prev, "Starting workflow execution..."]);

            // Simulasi per node
            for (let i = 1; i <= 3; i++) {
                await new Promise((r) => setTimeout(r, 1000));
                setLogs((prev) => [
                    ...prev,
                    `Node ${i}: Sending request...`,
                    `Node ${i}: Response 200 OK`,
                    `Node ${i}: Extracted data: {...}`,
                ]);
            }

            setLogs((prev) => [...prev, "Workflow completed successfully"]);
            toast.success("Workflow executed successfully");
        } catch (error: any) {
            setLogs((prev) => [...prev, `Error: ${error.message}`]);
            toast.error("Workflow execution failed");
        } finally {
            setRunning(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-2xl max-h-[80vh]">
                <DialogHeader>
                    <DialogTitle>Test Workflow</DialogTitle>
                    <DialogDescription>
                        Execute workflow and see step-by-step logs
                    </DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                    {/* Logs */}
                    <div className="bg-black text-green-400 p-4 rounded-lg font-mono text-xs h-64 overflow-y-auto">
                        {logs.length === 0 && (
                            <p className="text-gray-500">Click "Run" to start...</p>
                        )}
                        {logs.map((log, i) => (
                            <p key={i} className="mb-1">
                                [{i + 1}] {log}
                            </p>
                        ))}
                    </div>

                    <div className="flex gap-3">
                        <Button variant="outline" onClick={onClose} className="flex-1">
                            Close
                        </Button>
                        <Button onClick={handleTest} disabled={running} className="flex-1">
                            {running ? "Running..." : "Run"}
                        </Button>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}