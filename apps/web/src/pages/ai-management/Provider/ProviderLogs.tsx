// pages/admin/ai-management/components/ProviderLogs.tsx

import { useState, useEffect } from "react";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@kana-consultant/ui-kit";
import { Search, RefreshCw, Clock } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

interface LogEntry {
    id: string;
    timestamp: string;
    level: "info" | "error" | "warning";
    message: string;
    status_code?: number;
    latency?: number;
    endpoint?: string;
}

interface ProviderLogsProps {
    providerId: string;
}

export function ProviderLogs({ providerId }: ProviderLogsProps) {
    const [logs, setLogs] = useState<LogEntry[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState("");

    const toast = useToast()
    useEffect(() => {
        loadLogs();
    }, [providerId]);

    const loadLogs = async () => {
        setIsLoading(true);
        try {
            // TODO: Replace with actual API call
            await new Promise(resolve => setTimeout(resolve, 1000));
            
            // Mock data
            const mockLogs: LogEntry[] = [
                {
                    id: "1",
                    timestamp: new Date().toISOString(),
                    level: "info",
                    message: "API request successful",
                    status_code: 200,
                    latency: 245,
                    endpoint: "/chat/completions",
                },
                {
                    id: "2",
                    timestamp: new Date(Date.now() - 5 * 60000).toISOString(),
                    level: "error",
                    message: "Rate limit exceeded",
                    status_code: 429,
                    latency: 120,
                    endpoint: "/chat/completions",
                },
                {
                    id: "3",
                    timestamp: new Date(Date.now() - 10 * 60000).toISOString(),
                    level: "warning",
                    message: "Slow response detected",
                    status_code: 200,
                    latency: 3150,
                    endpoint: "/embeddings",
                },
                {
                    id: "4",
                    timestamp: new Date(Date.now() - 30 * 60000).toISOString(),
                    level: "info",
                    message: "Connection established",
                    status_code: 200,
                    latency: 89,
                    endpoint: "/chat/completions",
                },
                {
                    id: "5",
                    timestamp: new Date(Date.now() - 60 * 60000).toISOString(),
                    level: "error",
                    message: "Authentication failed",
                    status_code: 401,
                    latency: 45,
                    endpoint: "/chat/completions",
                },
            ];
            setLogs(mockLogs);
        } catch (error) {
            toast.error("Failed to load logs");
        } finally {
            setIsLoading(false);
        }
    };

    const getLevelBadge = (level: LogEntry["level"]) => {
        switch (level) {
            case "error":
                return <Badge tone="danger" className="text-xs">Error</Badge>;
            case "warning":
                return <Badge tone="warning" className="text-xs">Warning</Badge>;
            default:
                return <Badge tone="success" className="text-xs">Info</Badge>;
        }
    };

    const getStatusBadge = (statusCode?: number) => {
        if (!statusCode) return null;
        if (statusCode >= 200 && statusCode < 300) {
            return <Badge tone="success" className="text-xs">{statusCode}</Badge>;
        }
        if (statusCode >= 400 && statusCode < 500) {
            return <Badge tone="danger" className="text-xs">{statusCode}</Badge>;
        }
        if (statusCode >= 500) {
            return <Badge tone="warning" className="text-xs">{statusCode}</Badge>;
        }
        return <Badge tone="outline" className="text-xs">{statusCode}</Badge>;
    };

    const filteredLogs = logs.filter(log =>
        log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
        log.endpoint?.toLowerCase().includes(searchTerm.toLowerCase())
    );

    if (isLoading) {
        return (
            <Card>
                <CardContent className="flex items-center justify-center py-12">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card>
            <CardHeader>
                <div className="flex items-center justify-between">
                    <CardTitle>Request Logs</CardTitle>
                    <Button variant="outline" size="sm" onClick={loadLogs}>
                        <RefreshCw className="h-4 w-4 mr-2" />
                        Refresh
                    </Button>
                </div>
                <div className="relative mt-4">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder="Search logs..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="pl-9"
                    />
                </div>
            </CardHeader>
            <CardContent>
                {filteredLogs.length === 0 ? (
                    <div className="text-center py-8 text-muted-foreground">
                        No logs found
                    </div>
                ) : (
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Time</TableHead>
                                <TableHead>Level</TableHead>
                                <TableHead>Message</TableHead>
                                <TableHead>Endpoint</TableHead>
                                <TableHead>Status</TableHead>
                                <TableHead>Latency</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {filteredLogs.map((log) => (
                                <TableRow key={log.id}>
                                    <TableCell className="text-sm text-muted-foreground whitespace-nowrap">
                                        <div className="flex items-center gap-1">
                                            <Clock className="h-3 w-3" />
                                            {new Date(log.timestamp).toLocaleTimeString()}
                                        </div>
                                    </TableCell>
                                    <TableCell>{getLevelBadge(log.level)}</TableCell>
                                    <TableCell className="font-mono text-sm">{log.message}</TableCell>
                                    <TableCell>
                                        <code className="text-xs bg-muted px-1 py-0.5 rounded">
                                            {log.endpoint || "-"}
                                        </code>
                                    </TableCell>
                                    <TableCell>{getStatusBadge(log.status_code)}</TableCell>
                                    <TableCell className="text-sm">
                                        {log.latency ? `${log.latency}ms` : "-"}
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                )}
            </CardContent>
        </Card>
    );
}