import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus, Trash2, GripVertical, ChevronRight, ChevronDown } from "lucide-react";
import { type Field } from "./types";

interface JsonBuilderFieldProps {
    field: Field;
    level: number;
    onUpdate: (id: string, key: string, value: string) => void;
    onDelete: (id: string) => void;
    onToggleExpand: (id: string) => void;
    onAddField: (parentId: string) => void;
    onAddObject: (parentId: string) => void;
    renderChildren: (children: Field[], level: number) => React.ReactNode;
}

export function JsonBuilderField({
    field,
    level,
    onUpdate,
    onDelete,
    onToggleExpand,
    onAddField,
    onAddObject,
    renderChildren
}: JsonBuilderFieldProps) {
    const indent = level * 20;
    const isObject = field.type === "object";
    const hasChildren = field.children.length > 0;

    return (
        <div style={{ marginLeft: indent }} className="mb-2">
            <div className="flex items-center gap-2 p-2 border rounded-lg bg-white dark:bg-slate-900">
                {/* Drag Handle */}
                <div className="cursor-move text-slate-400">
                    <GripVertical className="h-4 w-4" />
                </div>

                {/* Expand/Collapse */}
                {isObject && (
                    <button onClick={() => onToggleExpand(field.id)} className="w-5">
                        {field.expanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                    </button>
                )}
                {!isObject && <div className="w-5" />}

                {/* Key Input */}
                <Input
                    placeholder="nama field"
                    value={field.key}
                    onChange={(e) => onUpdate(field.id, e.target.value, field.value)}
                    className="w-48 h-8 text-sm"
                />

                {/* Value Input (only for field type) */}
                {!isObject && (
                    <>
                        <span className="text-slate-400">:</span>
                        <Input
                            placeholder="nilai"
                            value={field.value}
                            onChange={(e) => onUpdate(field.id, field.key, e.target.value)}
                            className="flex-1 h-8 text-sm"
                        />
                    </>
                )}

                {/* Actions for Object */}
                {isObject && (
                    <div className="flex gap-1">
                        <Button size="sm" variant="outline" onClick={() => onAddField(field.id)}>
                            <Plus className="h-3 w-3 mr-1" />
                            Field
                        </Button>
                        <Button size="sm" variant="outline" onClick={() => onAddObject(field.id)}>
                            <Plus className="h-3 w-3 mr-1" />
                            Object
                        </Button>
                    </div>
                )}

                {/* Delete Button */}
                <Button size="sm" variant="ghost" onClick={() => onDelete(field.id)}>
                    <Trash2 className="h-4 w-4 text-red-500" />
                </Button>
            </div>

            {/* Children */}
            {isObject && field.expanded && hasChildren && (
                <div className="ml-4 pl-4 border-l-2 border-blue-200 mt-1">
                    {renderChildren(field.children, level + 1)}
                </div>
            )}
        </div>
    );
}