// pages/admin/ai-management/components/ProviderList/ProviderTable.tsx

import { useNavigate } from "@tanstack/react-router";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { MoreVertical, Edit, Power, PowerOff, Trash2, Eye } from "lucide-react";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@kana-consultant/ui-kit";
import { type APIProvider } from "@/types/provider.types";

interface ProviderTableProps {
    providers: APIProvider[];
    isLoading: boolean;
    onDelete?: (id: string) => void;
    onToggleStatus?: (id: string, isActive: boolean) => void;
}

export function ProviderTable({ 
    providers, 
    isLoading,  
    onDelete, 
    onToggleStatus 
}: ProviderTableProps) {
    const navigate = useNavigate();

    const handleNavigateToDetail = (id: string) => {
        navigate({
            to: "/protected/provider/$id",
            params: { id }
        });
    };

    const handleNavigateToEdit = (id: string) => {
        navigate({
            to: "/protected/provider/$id/edit",
            params: { id }
        });
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
        );
    }

    if (providers.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center h-64 text-center">
                <div className="text-muted-foreground mb-2">No providers found</div>
                <Button 
                    variant="outline" 
                    onClick={() => navigate({ to: "/provider/add" })}
                >
                    Create your first provider
                </Button>
            </div>
        );
    }

    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>Provider</TableHead>
                    <TableHead>Base URL</TableHead>
                    <TableHead>Families</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Last Updated</TableHead>
                    <TableHead className="w-20"></TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {providers.map((provider) => (
                    <TableRow 
                        key={provider.id}
                        className="cursor-pointer hover:bg-muted/50"
                        onClick={() => handleNavigateToDetail(provider.id!)}
                    >
                        <TableCell>
                            <div>
                                <div className="font-medium">{provider.display_name}</div>
                                <div className="text-sm text-muted-foreground">{provider.name}</div>
                                {provider.description && (
                                    <div className="text-xs text-muted-foreground truncate max-w-md">
                                        {provider.description}
                                    </div>
                                )}
                            </div>
                        </TableCell>
                        <TableCell>
                            <code className="text-xs bg-muted px-1 py-0.5 rounded">
                                {provider.base_url}
                            </code>
                        </TableCell>
                        <TableCell>
                            <div className="flex flex-wrap gap-1">
                                {provider.families.map((family: any) => (
                                    <Badge key={family.name} tone="outline" className="text-xs">
                                        {family.display_name}
                                    </Badge>
                                ))}
                            </div>
                        </TableCell>
                        <TableCell>
                            <Badge tone={provider.is_active ? "success" : "outline"}>
                                {provider.is_active ? "Active" : "Inactive"}
                            </Badge>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground">
                            {new Date(provider.updated_at as string).toLocaleDateString()}
                        </TableCell>
                        <TableCell onClick={(e) => e.stopPropagation()}>
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="sm">
                                        <MoreVertical className="h-4 w-4" />
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end">
                                    <DropdownMenuItem onClick={() => handleNavigateToDetail(provider.id!)}>
                                        <Eye className="h-4 w-4 mr-2" />
                                        View Details
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => handleNavigateToEdit(provider.id!)}>
                                        <Edit className="h-4 w-4 mr-2" />
                                        Edit
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onClick={() => onToggleStatus?.(provider.id!, !provider.is_active)}>
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
                                    </DropdownMenuItem>
                                    <DropdownMenuItem 
                                        onClick={() => onDelete?.(provider.id!)}
                                        className="text-red-600"
                                    >
                                        <Trash2 className="h-4 w-4 mr-2" />
                                        Delete
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}