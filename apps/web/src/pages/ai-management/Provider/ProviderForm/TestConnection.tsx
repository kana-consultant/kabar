// pages/admin/ai-management/components/TestConnection.tsx

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { AlertCircle, CheckCircle, Loader2 } from "lucide-react";
import type { ProviderFormData } from "@/types/provider.types";
import { testConnection } from "@/services/provider/provider-api";
interface TestConnectionProps {
    config: Partial<ProviderFormData>;
}

export function TestConnection({ config }: TestConnectionProps) {
    const [isTesting, setIsTesting] = useState(false);
    const [result, setResult] = useState<{ success: boolean; message?: string; latency?: number } | null>(null);

    const handleTest = async () => {
        if (!config.base_url) {
            setResult({ success: false, message: "Base URL is required" });
            return;
        }

        setIsTesting(true);
        try {
            const response = await testConnection({
                base_url: config.base_url,
                auth_type: config.auth_type || "bearer",
                auth_header: config.auth_header || "Authorization",
                auth_prefix: config.auth_prefix || "",
            });
            setResult({
                success: response.success,
                message: response.message,
                latency: response.latency,
            });
        } catch (error) {
            setResult({
                success: false,
                message: error instanceof Error ? error.message : "Connection failed",
            });
        } finally {
            setIsTesting(false);
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-sm">Test Connection</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
                <Button 
                    variant="outline" 
                    size="sm" 
                    onClick={handleTest}
                    disabled={isTesting || !config.base_url}
                    className="w-full"
                >
                    {isTesting ? (
                        <>
                            <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                            Testing...
                        </>
                    ) : (
                        "Run Connection Test"
                    )}
                </Button>

                {result && (
                    <div className={`p-3 rounded-lg ${result.success ? 'bg-green-50 dark:bg-green-950' : 'bg-red-50 dark:bg-red-950'}`}>
                        <div className="flex items-start space-x-2">
                            {result.success ? (
                                <CheckCircle className="h-4 w-4 text-green-600 mt-0.5" />
                            ) : (
                                <AlertCircle className="h-4 w-4 text-red-600 mt-0.5" />
                            )}
                            <div className="flex-1">
                                <p className={`text-sm ${result.success ? 'text-green-700 dark:text-green-300' : 'text-red-700 dark:text-red-300'}`}>
                                    {result.message || (result.success ? "Connection successful!" : "Connection failed")}
                                </p>
                                {result.latency && (
                                    <p className="text-xs mt-1 text-muted-foreground">
                                        Latency: {result.latency}ms
                                    </p>
                                )}
                            </div>
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}