// pages/admin/ai-management/TestConnectionDialog.tsx
import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { CheckCircle, XCircle, Loader2 } from "lucide-react";

interface TestConnectionDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  providerConfig: {
    base_url: string;
    auth_type: string;
    auth_header: string;
    auth_prefix: string;
  };
}

export function TestConnectionDialog({ open, onOpenChange, providerConfig }: TestConnectionDialogProps) {
  const [apiKey, setApiKey] = useState("");
  const [testing, setTesting] = useState(false);
  const [result, setResult] = useState<{
    success: boolean;
    status?: number;
    message?: string;
  } | null>(null);

  const handleTest = async () => {
    setTesting(true);
    setResult(null);

    // Simulasi test connection
    setTimeout(() => {
      setResult({
        success: true,
        status: 200,
        message: "Connection successful! Provider is reachable.",
      });
      setTesting(false);
    }, 1500);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Test Connection</DialogTitle>
          <DialogDescription>
            Test if the provider is reachable with your configuration
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label>API Key (for testing only)</Label>
            <Input
              type="password"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder="Enter temporary API key"
            />
            <p className="text-xs text-muted-foreground">
              This key will not be saved
            </p>
          </div>

          <div className="bg-muted p-3 rounded-md text-sm space-y-1">
            <div><strong>URL:</strong> {providerConfig.base_url}</div>
            <div><strong>Auth:</strong> {providerConfig.auth_header}</div>
            {providerConfig.auth_prefix && (
              <div><strong>Prefix:</strong> {providerConfig.auth_prefix}</div>
            )}
          </div>

          {result && (
            <div className={`p-3 rounded-md flex items-start space-x-2 ${
              result.success ? "bg-green-50" : "bg-red-50"
            }`}>
              {result.success ? (
                <CheckCircle className="h-5 w-5 text-green-500 mt-0.5" />
              ) : (
                <XCircle className="h-5 w-5 text-red-500 mt-0.5" />
              )}
              <div>
                <div className="font-medium text-sm">
                  {result.success ? "Success" : "Failed"}
                </div>
                {result.status && (
                  <Badge tone="outline" className="mt-1">
                    HTTP {result.status}
                  </Badge>
                )}
                <p className="text-sm mt-1">{result.message}</p>
              </div>
            </div>
          )}

          <div className="flex justify-end space-x-2">
            <Button variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button onClick={handleTest} disabled={!apiKey || testing}>
              {testing ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Testing...
                </>
              ) : (
                "Test Connection"
              )}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}