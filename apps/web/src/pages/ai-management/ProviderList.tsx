// pages/admin/ai-management/ProviderList.tsx
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@kana-consultant/ui-kit";
import { Button } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@kana-consultant/ui-kit";
import { MoreHorizontal, Pencil, Trash2, Globe, Key } from "lucide-react";
import { type APIProvider } from "@/types/provider.types";

interface ProviderListProps {
  providers: APIProvider[];
  loading: boolean;
  onEdit: (provider: APIProvider) => void;
  onDelete: (provider: APIProvider) => void;
}

export function ProviderList({ providers, loading, onEdit, onDelete }: ProviderListProps) {
  const getAuthTypeLabel = (type: string) => {
    switch (type) {
      case "bearer": return "Bearer Token";
      case "api_key": return "API Key";
      case "custom": return "Custom";
      default: return type;
    }
  };
  if (loading) {
    return <div className="text-center py-8">Loading providers...</div>;
  }

  if (providers.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        No providers found. Add your first provider to get started.
      </div>
    );
  }

  return (
    <div className="border rounded-lg">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Provider</TableHead>
            <TableHead>Base URL</TableHead>
            <TableHead>Auth Type</TableHead>
           
            <TableHead className="w-[80px]">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {providers.map((provider) => (
            <TableRow key={provider.id}>
              <TableCell>
                <div className="space-y-1">
                  <div className="font-medium">{provider.display_name}</div>
                  <div className="text-xs text-muted-foreground font-mono">
                    {provider.name}
                  </div>
                  {provider.description && (
                    <div className="text-xs text-muted-foreground line-clamp-1">
                      {provider.description}
                    </div>
                  )}
                </div>
              </TableCell>
              
              <TableCell>
                <div className="flex items-center space-x-1">
                  <Globe className="h-3 w-3 text-muted-foreground shrink-0" />
                  <span className="text-sm font-mono break-all">{provider.base_url}</span>
                </div>
              </TableCell>
              
              <TableCell>
                <div className="flex items-start space-x-2">
                  <Key className="h-3 w-3 text-muted-foreground shrink-0 mt-0.5" />
                  <div>
                    <Badge tone="outline">
                      {getAuthTypeLabel(provider.auth_type)}
                    </Badge>
                    <div className="text-xs text-muted-foreground mt-1">
                      {provider.auth_header}
                      {provider.auth_prefix && ` (${provider.auth_prefix})`}
                    </div>
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => onEdit(provider)}>
                      <Pencil className="h-4 w-4 mr-2" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      onClick={() => onDelete(provider)}
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
    </div>
  );
}