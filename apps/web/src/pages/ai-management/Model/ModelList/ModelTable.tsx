import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@kana-consultant/ui-kit";
import { Badge, Button } from "@kana-consultant/ui-kit";
import { Eye, Edit, Trash2, Play } from "lucide-react";
import { type AIModel } from "@/types/provider.types";

interface ModelTableProps {
    models: AIModel[];
    isLoading?: boolean;
    onEdit: (model: string) => void;
    onView: (model: AIModel) => void;
    onDelete?: (id: AIModel) => void;
    onTest?: (id: string) => void;
}

export function ModelTable({ models, isLoading, onEdit, onView, onDelete, onTest }: ModelTableProps) {
    console.log(models)
    if (isLoading) {
        return <div className="text-center py-8">Loading...</div>;
    }

    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Provider</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {models.map((model) => (
                    <TableRow key={model.id}>
                        <TableCell>{model.display_name}</TableCell>
                        <TableCell>{model.provider_id}</TableCell>
                        <TableCell>
                            {model.is_active ? (
                                <Badge tone="success">Active</Badge>
                            ) : (
                                <Badge tone="outline">Inactive</Badge>
                            )}
                        </TableCell>
                        {/* <TableCell>
                            {model.is_default && <Badge tone="info">Default</Badge>}
                        </TableCell> */}
                        <TableCell className="text-right space-x-1">
                            <Button variant="ghost" size="sm" onClick={() => onView(model)}>
                                <Eye className="h-4 w-4" />
                            </Button>
                            <Button variant="ghost" size="sm" onClick={() => onEdit(model.id as string)}>
                                <Edit className="h-4 w-4" />
                            </Button>
                            {onTest && (
                                <Button variant="ghost" size="sm" onClick={() => onTest(model.id)}>
                                    <Play className="h-4 w-4" />
                                </Button>
                            )}
                            {onDelete && (
                                <Button variant="ghost" size="sm" onClick={() => onDelete(model)}>
                                    <Trash2 className="h-4 w-4" />
                                </Button>
                            )}
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}