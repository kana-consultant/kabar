// pages/admin/ai-management/[id]/index.tsx

import { useState, useEffect } from "react";
import { useNavigate, useParams } from "@tanstack/react-router";
import { Edit, Play, Power, PowerOff, ArrowLeft } from "lucide-react";
import { Page, PageHeader, PageTitle, PageDescription, PageContent } from "@/components/page/page";
import { Button } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Card, CardContent, CardHeader, CardTitle } from "@kana-consultant/ui-kit";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";
import { toggleProviderStatus, testProviderConnection, getProviderById  } from "@/services/provider/provider-api";
import { FamiliesList } from "../FamiliesList";
import { ConnectionInfo } from "../ConnectionInfo";
import { ProviderLogs } from "../ProviderLogs";
import type { APIProvider } from "@/types/provider.types";

export default function ProviderDetailPage() {
    const navigate = useNavigate();
    const toast = useToast()
    const { id } = useParams({ from: '/admin/ai-management/$id' });
    const [provider, setProvider] = useState<APIProvider | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isTesting, setIsTesting] = useState(false);

    useEffect(() => {
        if (id) {
            loadProvider();
        }
    }, [id]);

    const loadProvider = async () => {
        try {
            const provider = await getProviderById(id as string);
            setProvider(provider as APIProvider);
        } catch (error) {
            toast.error("Failed to load provider");
            navigate({ to: "/admin/ai-management" });
        } finally {
            setIsLoading(false);
        }
    };

    const handleToggleStatus = async () => {
        if (!provider) return;

        try {
            const updated = await toggleProviderStatus(id as string, !provider.is_active);
            setProvider(updated);
            toast.success(`Provider ${updated.is_active ? "activated" : "deactivated"}`);
        } catch (error) {
            toast.error("Failed to toggle provider status");
        }
    };

    const handleTestConnection = async () => {
        setIsTesting(true);
        try {
            const result = await testProviderConnection(id as string);
            if (result.success) {
                toast.success("Connection test successful!");
            } else {
                toast.error(`Connection failed: ${result.error}`);
            }
        } catch (error) {
            toast.error("Failed to test connection");
        } finally {
            setIsTesting(false);
        }
    };

    const handleEdit = () => {
        navigate({
            to: "/protected/provider/$id/edit",
            params: { id: id as string }
        });
    };

    const handleBack = () => {
        navigate({ to: "/admin/ai-management" });
    };

    if (isLoading) {
        return (
            <Page>
                <PageContent>
                    <div className="flex items-center justify-center h-64">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                    </div>
                </PageContent>
            </Page>
        );
    }

    if (!provider) {
        return null;
    }

    return (
        <Page>
            <PageHeader>
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={handleBack}
                        >
                            <ArrowLeft className="h-4 w-4 mr-2" />
                            Back
                        </Button>
                        <div>
                            <div className="flex items-center space-x-3">
                                <PageTitle>{provider.display_name}</PageTitle>
                                <Badge tone={provider.is_active ? "success" : "outline"}>
                                    {provider.is_active ? "Active" : "Inactive"}
                                </Badge>
                            </div>
                            <PageDescription>{provider.description || "No description"}</PageDescription>
                        </div>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Button
                            variant="outline"
                            onClick={handleTestConnection}
                            disabled={isTesting}
                        >
                            <Play className="h-4 w-4 mr-2" />
                            {isTesting ? "Testing..." : "Test Connection"}
                        </Button>
                        <Button
                            variant={provider.is_active ? "destructive" : "primary"}
                            onClick={handleToggleStatus}
                        >
                            {provider.is_active ? (
                                <>
                                    <PowerOff className="h-4 w-4 mr-2" />
                                    Deactivate
                                </>
                            ) : (
                                <>
                                    <Power className="h-4 w-4 mr-2" />
                                    Activate
                                </>
                            )}
                        </Button>
                        <Button onClick={handleEdit}>
                            <Edit className="h-4 w-4 mr-2" />
                            Edit
                        </Button>
                    </div>
                </div>
            </PageHeader>

            <PageContent>
                <Tabs defaultValue="families" className="space-y-6">
                    <TabsList>
                        <TabsTrigger value="families">API Families</TabsTrigger>
                        <TabsTrigger value="connection">Connection</TabsTrigger>
                        <TabsTrigger value="logs">Logs</TabsTrigger>
                        <TabsTrigger value="settings">Settings</TabsTrigger>
                    </TabsList>

                    <TabsContent value="families">
                        <FamiliesList families={provider.families} />
                    </TabsContent>

                    <TabsContent value="connection">
                        <ConnectionInfo provider={provider} />
                    </TabsContent>

                    <TabsContent value="logs">
                        <ProviderLogs providerId={id as string} />
                    </TabsContent>

                    <TabsContent value="settings">
                        <Card>
                            <CardHeader>
                                <CardTitle>Provider Settings</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    <div className="grid grid-cols-2 gap-4">
                                        <div>
                                            <label className="text-sm font-medium">Provider Name</label>
                                            <p className="text-sm text-muted-foreground mt-1">{provider.name}</p>
                                        </div>
                                        <div>
                                            <label className="text-sm font-medium">Base URL</label>
                                            <p className="text-sm text-muted-foreground mt-1 font-mono">{provider.base_url}</p>
                                        </div>
                                    </div>
                                    <div>
                                        <label className="text-sm font-medium">Default Headers</label>
                                        <pre className="mt-2 p-2 bg-muted rounded-md text-xs">
                                            {JSON.stringify(provider.default_headers, null, 2)}
                                        </pre>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </TabsContent>
                </Tabs>
            </PageContent>
        </Page>
    );
}