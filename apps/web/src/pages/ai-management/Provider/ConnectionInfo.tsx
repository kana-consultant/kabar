// pages/admin/ai-management/components/ConnectionInfo.tsx

import { Card, CardContent, CardHeader, CardTitle } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import type { APIProvider } from "@/types/provider.types";

interface ConnectionInfoProps {
    provider: APIProvider;
}

export function ConnectionInfo({ provider }: ConnectionInfoProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Connection Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                <div>
                    <h4 className="text-sm font-medium mb-1">Base URL</h4>
                    <code className="text-sm bg-muted px-2 py-1 rounded block">
                        {provider.base_url}
                    </code>
                </div>

                <div>
                    <h4 className="text-sm font-medium mb-2">Authentication</h4>
                    <div className="space-y-2">
                        <div className="flex items-center gap-2">
                            <span className="text-sm text-muted-foreground">Type:</span>
                            <Badge tone="outline">{provider.auth_type}</Badge>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="text-sm text-muted-foreground">Header:</span>
                            <code className="text-xs bg-muted px-2 py-0.5 rounded">
                                {provider.auth_header}
                            </code>
                        </div>
                        {provider.auth_prefix && (
                            <div className="flex items-center gap-2">
                                <span className="text-sm text-muted-foreground">Prefix:</span>
                                <code className="text-xs bg-muted px-2 py-0.5 rounded">
                                    {provider.auth_prefix}
                                </code>
                            </div>
                        )}
                    </div>
                </div>

                {Object.keys(provider.default_headers).length > 0 && (
                    <div>
                        <h4 className="text-sm font-medium mb-2">Default Headers</h4>
                        <pre className="text-xs bg-muted p-3 rounded-lg overflow-x-auto">
                            {JSON.stringify(provider.default_headers, null, 2)}
                        </pre>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}