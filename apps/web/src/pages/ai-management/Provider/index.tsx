// pages/admin/ai-management/index.tsx

import { useState, useEffect } from "react";
import { useNavigate } from "@tanstack/react-router";
import { Plus, Search, Filter } from "lucide-react";
import {
    Page,
    PageHeader,
    PageTitle,
    PageDescription,
    PageContent,
} from "@/components/page/page"
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Card } from "@kana-consultant/ui-kit";
import { useToast } from "@/hooks/use-toast";
import { ProviderTable } from "./ProviderForm/ProviderList/ProviderTable";
import { ProviderFilters } from "./ProviderForm/ProviderList/ProviderFilters";
import { deleteProvider, toggleProviderStatus } from "@/services/provider/provider-api";
import type { APIProvider } from "@/types/provider.types";
import { getProviders } from "@/services/provider/providerMutations";

export default function AIManagementPage() {
    const toast = useToast()
    const navigate = useNavigate();
    const [providers, setProviders] = useState<APIProvider[]>([]);
    const [filteredProviders, setFilteredProviders] = useState<APIProvider[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState("");
    const [filters, setFilters] = useState({
        status: "all",
        family: "all",
    });

    useEffect(() => {
        loadProviders();
    }, []);

    useEffect(() => {
        filterProviders();
    }, [searchTerm, filters, providers]);

    const loadProviders = async () => {
        try {
            const provider = await getProviders();
            setProviders(provider.data as APIProvider[]);
        } catch (error) {
            toast.error("Failed to load providers");
        } finally {
            setIsLoading(false);
        }
    };

    const filterProviders = () => {
        let filtered = [...providers];

        if (searchTerm) {
            filtered = filtered.filter(
                (p) =>
                    p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    p.display_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    p.description?.toLowerCase().includes(searchTerm.toLowerCase())
            );
        }

        if (filters.status !== "all") {
            filtered = filtered.filter((p) =>
                filters.status === "active" ? p.is_active : !p.is_active
            );
        }

        if (filters.family !== "all") {
            filtered = filtered.filter((p) =>
                p.families.some((f) => f.name === filters.family)
            );
        }

        setFilteredProviders(filtered);
    };

    const handleDelete = async (id: string) => {
        if (!confirm("Are you sure you want to delete this provider?")) return;
        
        try {
            await deleteProvider(id);
            toast.success("Provider deleted successfully");
            loadProviders();
        } catch (error) {
            toast.error("Failed to delete provider");
        }
    };

    const handleToggleStatus = async (id: string, isActive: boolean) => {
        try {
            await toggleProviderStatus(id, isActive);
            toast.success(`Provider ${isActive ? "activated" : "deactivated"}`);
            loadProviders();
        } catch (error) {
            toast.error("Failed to toggle provider status");
        }
    };

    return (
        <Page>
            <PageHeader>
                <div className="flex items-center justify-between">
                    <div>
                        <PageTitle>AI Providers</PageTitle>
                        <PageDescription>
                            Manage API providers for AI model integration
                        </PageDescription>
                    </div>
                    <Button onClick={() => navigate({ to: "/provider/add" })}>
                        <Plus className="h-4 w-4 mr-2" />
                        Add Provider
                    </Button>
                </div>
            </PageHeader>

            <PageContent>
                <Card className="p-4 mb-6">
                    <div className="flex flex-col md:flex-row gap-4">
                        <div className="flex-1 relative">
                            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder="Search providers..."
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className="pl-9"
                            />
                        </div>
                        <ProviderFilters filters={filters} onFiltersChange={setFilters} />
                    </div>
                </Card>

                <ProviderTable
                    providers={filteredProviders}
                    isLoading={isLoading}
                    onDelete={handleDelete}
                    onToggleStatus={handleToggleStatus}
                />
            </PageContent>
        </Page>
    );
}