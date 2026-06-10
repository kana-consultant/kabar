// components/JsonBuilder.tsx

import { useState, useEffect, useRef, useCallback, memo } from "react";
import { Button } from "@kana-consultant/ui-kit";
import { Input } from "@kana-consultant/ui-kit";
import { Label } from "@kana-consultant/ui-kit";
import { Badge } from "@kana-consultant/ui-kit";
import {
    GripVertical, Plus, Trash2, ChevronDown, ChevronRight,
    Braces, Type, List, Hash, ToggleLeft, HelpCircle, X
} from "lucide-react";
import { cn } from "@/lib/utils";
import { type Field, jsonToFields, fieldsToJson, genId, getNextFieldNumber, getNextObjectNumber } from "@/types/JsonBuilder";

interface JsonBuilderProps {
    value: Record<string, any>;
    onChange: (jsonObject: Record<string, any>) => void;
    availableVariables?: string[];
}

const defaultVariables = ["{model}", "{prompt}", "{temperature}", "{max_token}", "{system_prompt}"];

// Komponen untuk menampilkan value dalam bentuk badge
const ValueBadge = ({ value, onRemove, onEdit }: { value: any; onRemove?: () => void; onEdit?: () => void }) => {
    const displayValue = typeof value === 'object' ? '{...}' : String(value);

    const getBadgeColor = () => {
        if (typeof value === 'number') return "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400";
        if (typeof value === 'boolean') return "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400";
        if (typeof value === 'object') return "bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400";
        return "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300";
    };

    const handleRemoveClick = (e: React.MouseEvent) => {
        e.stopPropagation();
        e.preventDefault();
        if (onRemove) {
            onRemove();
        }
    };

    const handleEditClick = (e: React.MouseEvent) => {
        e.stopPropagation();
        if (onEdit) {
            onEdit();
        }
    };

    return (
        <div
            className="group/badge inline-flex items-center gap-1 bg-white dark:bg-slate-900 rounded-md border border-slate-200 dark:border-slate-700 px-2 py-0.5 hover:border-blue-300 transition-colors cursor-pointer"
            onClick={handleEditClick}
        >
            {typeof value === 'number' && <Hash className="h-3 w-3 text-amber-500" />}
            {typeof value === 'boolean' && <ToggleLeft className="h-3 w-3 text-emerald-500" />}
            {typeof value === 'object' && !Array.isArray(value) && value !== null && <Braces className="h-3 w-3 text-purple-500" />}
            {typeof value === 'string' && !value.startsWith('{') && <Type className="h-3 w-3 text-slate-400" />}
            {typeof value === 'string' && value.startsWith('{') && value.endsWith('}') && (
                <span className="text-blue-500 font-mono text-xs font-bold">{'{ }'}</span>
            )}

            <span
                className={cn(
                    "text-xs font-mono",
                    getBadgeColor()
                )}
            >
                {displayValue.length > 30 ? displayValue.substring(0, 27) + "..." : displayValue}
            </span>

            {onRemove && (
                <button
                    type="button"
                    onClick={handleRemoveClick}
                    className="opacity-0 group-hover/badge:opacity-100 hover:bg-red-100 rounded p-0.5 transition-all"
                >
                    <X className="h-3 w-3 text-red-500" />
                </button>
            )}
        </div>
    );
};

// Memoized Field Row Component
const FieldRow = memo(({
    field,
    depth,
    dragOverId,
    onDragStart,
    onDragOver,
    onDrop,
    onDragLeave,
    onToggleExpand,
    onUpdateField,
    onRemoveField,
    onAddField,
    onAddObject,
    onEditValue,
    editingFieldId,
    onCancelEdit,
    onInsertVariable,
    availableVariables,
}: {
    field: Field;
    depth: number;
    dragOverId: string | null;
    onDragStart: (e: React.DragEvent, id: string) => void;
    onDragOver: (e: React.DragEvent, id: string) => void;
    onDrop: (e: React.DragEvent, id: string) => void;
    onDragLeave: () => void;
    onToggleExpand: (id: string) => void;
    onUpdateField: (id: string, key: string, value: any) => void;
    onRemoveField: (id: string) => void;
    onAddField: (id: string | null) => void;
    onAddObject: (id: string | null) => void;
    onEditValue: (id: string) => void;
    editingFieldId: string | null;
    onCancelEdit: () => void;
    onInsertVariable: (variable: string) => void;
    availableVariables: string[];
}) => {
    const isArrayObject = field.type === "object" &&
        field.children.length > 0 &&
        /^\[\d+\]$/.test(field.children[0]?.key || "");

    const isEditing = editingFieldId === field.id;
    const [editValue, setEditValue] = useState(String(field.value ?? ""));
    const inputRef = useRef<HTMLInputElement>(null);

    // Sync editValue saat field value berubah dari luar (misal: insert variable)
    useEffect(() => {
        setEditValue(String(field.value ?? ""));
    }, [field.value]);

    useEffect(() => {
        if (isEditing && inputRef.current) {
            inputRef.current.focus();
            // Pindah cursor ke akhir
            const len = inputRef.current.value.length;
            inputRef.current.setSelectionRange(len, len);
        }
    }, [isEditing, field.value]); // re-run saat value berubah (insert variable)

    const handleSaveEdit = () => {
        onUpdateField(field.id, "value", editValue);
        onCancelEdit();
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') {
            handleSaveEdit();
        } else if (e.key === 'Escape') {
            onCancelEdit();
        }
    };

    // Local state untuk input key agar tidak lag saat typing
    const [localKey, setLocalKey] = useState(field.key);

    useEffect(() => {
        setLocalKey(field.key);
    }, [field.key]);

    const handleKeyChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newValue = e.target.value;
        setLocalKey(newValue);
        onUpdateField(field.id, "key", newValue);
    };

    return (
        <div className="space-y-1">
            <div
                className={cn(
                    "flex items-center gap-2 p-2 rounded-lg group transition-all border",
                    "hover:bg-slate-50 dark:hover:bg-white/5 hover:border-slate-200 dark:hover:border-slate-700",
                    dragOverId === field.id
                        ? "border-dashed border-blue-500 bg-blue-50/50 dark:bg-blue-500/10 shadow-sm"
                        : "border-transparent"
                )}
                style={{ marginLeft: `${Math.max(depth * 20, 0)}px` }}
                draggable
                onDragStart={(e) => onDragStart(e, field.id)}
                onDragOver={(e) => onDragOver(e, field.id)}
                onDrop={(e) => onDrop(e, field.id)}
                onDragLeave={onDragLeave}
            >
                <GripVertical className="h-4 w-4 text-slate-400 opacity-0 group-hover:opacity-100 cursor-grab shrink-0 transition-opacity" />

                {field.type === "object" && (
                    <button
                        type="button"
                        onClick={() => onToggleExpand(field.id)}
                        className="shrink-0 p-0.5 hover:bg-muted rounded transition-colors"
                    >
                        {field.expanded ? <ChevronDown className="h-3.5 w-3.5" /> : <ChevronRight className="h-3.5 w-3.5" />}
                    </button>
                )}

                <div className="flex items-center gap-3 flex-1 min-w-0 flex-wrap sm:flex-nowrap">
                    {/* Key Field */}
                    <div className="flex items-center gap-1 bg-muted/50 rounded-md px-2 py-0.5 min-w-[100px]">
                        {isArrayObject ? (
                            <List className="h-3.5 w-3.5 text-green-500 shrink-0" />
                        ) : field.type === "object" ? (
                            <Braces className="h-3.5 w-3.5 text-blue-500 shrink-0" />
                        ) : (
                            <Type className="h-3.5 w-3.5 text-slate-400 shrink-0" />
                        )}
                        <Input
                            value={localKey}
                            onChange={handleKeyChange}
                            className="h-7 text-xs font-mono border-0 bg-transparent focus-visible:ring-0 focus-visible:ring-offset-0 px-0 w-28 font-medium"
                            placeholder="field_name"
                        />
                    </div>

                    <span className="text-slate-400 shrink-0 font-mono text-sm">:</span>

                    {/* Value Display */}
                    {field.type === "object" ? (
                        <Badge tone="outline" className="text-[10px] px-2 h-5 shrink-0 font-mono bg-blue-50 dark:bg-blue-950/30">
                            {isArrayObject ? `📋 Array [${field.children.length}]` : `📦 Object {${field.children.length}}`}
                        </Badge>
                    ) : isEditing ? (
                        <div className="flex-1 flex items-center gap-2 flex-wrap">
                            <Input
                                ref={inputRef}
                                value={editValue}
                                onChange={(e) => setEditValue(e.target.value)}
                                onKeyDown={handleKeyDown}
                                onBlur={handleSaveEdit}
                                className="h-7 text-xs font-mono flex-1"
                                placeholder='Value (text, 0.7, true, "{variable}")'
                            />
                            {/* Tombol variable mini di dalam edit mode */}
                            <div className="flex gap-1 flex-wrap shrink-0">
                                {availableVariables.map((v) => (
                                    <button
                                        key={v}
                                        type="button"
                                        onMouseDown={(e) => {
                                            // Pakai onMouseDown + preventDefault agar blur pada input
                                            // tidak terpicu sebelum insert terjadi
                                            e.preventDefault();
                                            setEditValue(prev => prev + v);
                                            // Update ke parent juga supaya ref lastEdited akurat
                                            onInsertVariable(v);
                                            // Kembalikan fokus ke input
                                            setTimeout(() => inputRef.current?.focus(), 0);
                                        }}
                                        className="text-[10px] px-1.5 py-0.5 rounded bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 hover:bg-blue-200 transition-colors font-mono"
                                    >
                                        {v}
                                    </button>
                                ))}
                            </div>
                        </div>
                    ) : (
                        <div className="flex-1">
                            <ValueBadge
                                value={field.value}
                                onEdit={() => onEditValue(field.id)}
                                onRemove={field.value !== "" && field.value !== undefined ? () => onUpdateField(field.id, "value", "") : undefined}
                            />
                        </div>
                    )}
                </div>

                <div className="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 shrink-0 transition-opacity">
                    {field.type === "object" && (
                        <>
                            <Button
                                variant="ghost"
                                size="icon"
                                type="button"
                                className="h-7 w-7"
                                onClick={() => onAddField(field.id)}
                                title="Add Field"
                            >
                                <Plus className="h-3.5 w-3.5" />
                            </Button>
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-7 w-7"
                                onClick={() => onAddObject(field.id)}
                                title="Add Nested Object"
                            >
                                <Braces className="h-3.5 w-3.5 text-blue-500" />
                            </Button>
                        </>
                    )}
                    <Button
                        variant="ghost"
                        type="button"
                        size="icon"
                        className="h-7 w-7 text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-500/10"
                        onClick={() => onRemoveField(field.id)}
                        title="Delete"
                    >
                        <Trash2 className="h-3.5 w-3.5" />
                    </Button>
                </div>
            </div>

            {field.expanded && field.type === "object" && field.children.length > 0 && (
                <div className="border-l-2 border-slate-200 dark:border-slate-700 ml-3 pl-2 my-1">
                    {field.children.map((child) => (
                        <FieldRow
                            key={child.id}
                            field={child}
                            depth={depth + 1}
                            dragOverId={dragOverId}
                            onDragStart={onDragStart}
                            onDragOver={onDragOver}
                            onDrop={onDrop}
                            onDragLeave={onDragLeave}
                            onToggleExpand={onToggleExpand}
                            onUpdateField={onUpdateField}
                            onRemoveField={onRemoveField}
                            onAddField={onAddField}
                            onAddObject={onAddObject}
                            onEditValue={onEditValue}
                            editingFieldId={editingFieldId}
                            onCancelEdit={onCancelEdit}
                            onInsertVariable={onInsertVariable}
                            availableVariables={availableVariables}
                        />
                    ))}
                </div>
            )}
        </div>
    );
});

FieldRow.displayName = "FieldRow";

export function JsonBuilder({ value, onChange, availableVariables = defaultVariables }: JsonBuilderProps) {
    const [fields, setFields] = useState<Field[]>(() => {
        return jsonToFields(value || {});
    });
    const [dragOverId, setDragOverId] = useState<string | null>(null);
    const [dragSourceId, setDragSourceId] = useState<string | null>(null);
    const [showHelper, setShowHelper] = useState(true);
    const [editingFieldId, setEditingFieldId] = useState<string | null>(null);
    // ── FIX: simpan field terakhir yang diedit agar insert variabel dari
    //         Quick Guide tetap bisa bekerja meski input sudah blur ──────────
    const lastEditedFieldIdRef = useRef<string | null>(null);
    // ── FIX: pesan error validasi variabel ───────────────────────────────────
    const [varError, setVarError] = useState<string | null>(null);

    const isFirstRender = useRef(true);
    const updateTimeoutRef = useRef<number | null>(null);
    const pendingUpdateRef = useRef<Field[] | null>(null);

    // Only sync from props on first render or when value changes externally
    useEffect(() => {
        if (isFirstRender.current) {
            isFirstRender.current = false;
            return;
        }

        const currentJson = fieldsToJson(fields);
        if (JSON.stringify(currentJson) !== JSON.stringify(value)) {
            const converted = jsonToFields(value || {});
            setFields(converted);
        }
    }, [value]);

    // Debounced update to parent
    const debouncedOnChange = useCallback((newFields: Field[]) => {
        if (updateTimeoutRef.current) {
            clearTimeout(updateTimeoutRef.current);
        }

        pendingUpdateRef.current = newFields;

        updateTimeoutRef.current = setTimeout(() => {
            if (pendingUpdateRef.current) {
                const obj = fieldsToJson(pendingUpdateRef.current);
                onChange(obj);
                pendingUpdateRef.current = null;
            }
        }, 150);
    }, [onChange]);

    function updateAndSync(newFields: Field[]) {
        setFields(newFields);
        debouncedOnChange(newFields);
    }

    // ── FIX: parseValueType — validasi semua token {xxx} harus ada di
    //         availableVariables. Jika ada yang tidak dikenal, tolak & beri error.
    function parseValueType(val: string): any {
        if (val.toLowerCase() === "true") return true;
        if (val.toLowerCase() === "false") return false;

        // Cari semua token { ... } dalam string
        const tokenRegex = /\{[^{}]+\}/g;
        const tokens = val.match(tokenRegex);

        if (tokens && tokens.length > 0) {
            const invalidTokens = tokens.filter(t => !availableVariables.includes(t));
            if (invalidTokens.length > 0) {
                setVarError(
                    `Variabel tidak dikenal: ${invalidTokens.join(", ")}. ` +
                    `Hanya ${availableVariables.join(", ")} yang diizinkan.`
                );
                // Auto-clear pesan error setelah 4 detik
                setTimeout(() => setVarError(null), 4000);
                return ""; // tolak value
            }
        }

        // Bersihkan error jika valid
        setVarError(null);

        if (val === "" || isNaN(Number(val)) || val.startsWith("{")) return val;
        return Number(val);
    }

    const addField = useCallback((parentId: string | null) => {
        if (!parentId) {
            const nextNum = getNextFieldNumber(fields);
            const newField: Field = {
                id: genId(),
                key: `field${nextNum}`,
                value: "",
                type: "field",
                children: [],
                expanded: false,
            };
            updateAndSync([...fields, newField]);
        } else {
            const addFieldToParent = (fields: Field[], parentId: string): Field[] | null => {
                for (let i = 0; i < fields.length; i++) {
                    if (fields[i].id === parentId && fields[i].type === "object") {
                        const newFields = [...fields];
                        const nextNum = getNextFieldNumber(newFields[i].children);
                        newFields[i] = {
                            ...newFields[i],
                            children: [...newFields[i].children, {
                                id: genId(),
                                key: `field${nextNum}`,
                                value: "",
                                type: "field",
                                children: [],
                                expanded: false,
                            }],
                        };
                        return newFields;
                    }
                    if (fields[i].children.length > 0) {
                        const updated = addFieldToParent(fields[i].children, parentId);
                        if (updated) {
                            const newFields = [...fields];
                            newFields[i] = { ...newFields[i], children: updated };
                            return newFields;
                        }
                    }
                }
                return null;
            };
            const updated = addFieldToParent(fields, parentId);
            if (updated) updateAndSync(updated);
        }
    }, [fields]);

    const addObject = useCallback((parentId: string | null) => {
        if (!parentId) {
            const nextNum = getNextObjectNumber(fields);
            const newField: Field = {
                id: genId(),
                key: `object${nextNum}`,
                value: "",
                type: "object",
                children: [],
                expanded: true,
            };
            updateAndSync([...fields, newField]);
        } else {
            const addObjectToParent = (fields: Field[], parentId: string): Field[] | null => {
                for (let i = 0; i < fields.length; i++) {
                    if (fields[i].id === parentId && fields[i].type === "object") {
                        const newFields = [...fields];
                        const nextNum = getNextObjectNumber(newFields[i].children);
                        newFields[i] = {
                            ...newFields[i],
                            children: [...newFields[i].children, {
                                id: genId(),
                                key: `object${nextNum}`,
                                value: "",
                                type: "object",
                                children: [],
                                expanded: true,
                            }],
                        };
                        return newFields;
                    }
                    if (fields[i].children.length > 0) {
                        const updated = addObjectToParent(fields[i].children, parentId);
                        if (updated) {
                            const newFields = [...fields];
                            newFields[i] = { ...newFields[i], children: updated };
                            return newFields;
                        }
                    }
                }
                return null;
            };
            const updated = addObjectToParent(fields, parentId);
            if (updated) updateAndSync(updated);
        }
    }, [fields]);

    const removeField = useCallback((fieldId: string) => {
        const removeFieldById = (fields: Field[], targetId: string): Field[] | null => {
            const filtered = fields.filter(f => f.id !== targetId);
            if (filtered.length !== fields.length) return filtered;

            for (let i = 0; i < fields.length; i++) {
                if (fields[i].children.length > 0) {
                    const updated = removeFieldById(fields[i].children, targetId);
                    if (updated) {
                        const newFields = [...fields];
                        newFields[i] = { ...newFields[i], children: updated };
                        return newFields;
                    }
                }
            }
            return null;
        };
        const updated = removeFieldById(fields, fieldId);
        if (updated) updateAndSync(updated);
    }, [fields]);

    const updateFieldValue = useCallback((fieldId: string, fieldKey: string, fieldValue: any) => {
        const finalValue = fieldKey === "value" ? parseValueType(fieldValue) : fieldValue;

        const updateFieldById = (fields: Field[], targetId: string, updater: (field: Field) => Field): Field[] | null => {
            for (let i = 0; i < fields.length; i++) {
                if (fields[i].id === targetId) {
                    const newFields = [...fields];
                    newFields[i] = updater(newFields[i]);
                    return newFields;
                }
                if (fields[i].children.length > 0) {
                    const updated = updateFieldById(fields[i].children, targetId, updater);
                    if (updated) {
                        const newFields = [...fields];
                        newFields[i] = { ...newFields[i], children: updated };
                        return newFields;
                    }
                }
            }
            return null;
        };

        const updated = updateFieldById(fields, fieldId, (field) => ({
            ...field,
            [fieldKey]: finalValue,
        }));
        if (updated) updateAndSync(updated);
    }, [fields]);

    const toggleExpand = useCallback((fieldId: string) => {
        const updateFieldById = (fields: Field[], targetId: string, updater: (field: Field) => Field): Field[] | null => {
            for (let i = 0; i < fields.length; i++) {
                if (fields[i].id === targetId) {
                    const newFields = [...fields];
                    newFields[i] = updater(newFields[i]);
                    return newFields;
                }
                if (fields[i].children.length > 0) {
                    const updated = updateFieldById(fields[i].children, targetId, updater);
                    if (updated) {
                        const newFields = [...fields];
                        newFields[i] = { ...newFields[i], children: updated };
                        return newFields;
                    }
                }
            }
            return null;
        };

        const updated = updateFieldById(fields, fieldId, (field) => ({
            ...field,
            expanded: !field.expanded,
        }));
        if (updated) setFields(updated);
    }, [fields]);

    const findFieldById = useCallback((fields: Field[], targetId: string): Field | null => {
        for (const field of fields) {
            if (field.id === targetId) return field;
            if (field.children.length > 0) {
                const found = findFieldById(field.children, targetId);
                if (found) return found;
            }
        }
        return null;
    }, []);

    // ── FIX: handleInsertVariable — gunakan lastEditedFieldIdRef sebagai
    //         fallback agar insert dari Quick Guide tetap bekerja ─────────────
    const handleInsertVariable = useCallback((variable: string) => {
        const targetId = editingFieldId ?? lastEditedFieldIdRef.current;
        if (!targetId) return;

        // Jika field belum dalam mode edit, buka dulu
        if (!editingFieldId) {
            setEditingFieldId(targetId);
        }

        const currentField = findFieldById(fields, targetId);
        const currentValue = String(currentField?.value ?? "");
        updateFieldValue(targetId, "value", currentValue + variable);
    }, [editingFieldId, fields, findFieldById, updateFieldValue]);

    // ── FIX: handleEditValue — simpan ke ref juga ────────────────────────────
    const handleEditValue = useCallback((id: string) => {
        setEditingFieldId(id);
        lastEditedFieldIdRef.current = id;
    }, []);

    const handleCancelEdit = useCallback(() => {
        setEditingFieldId(null);
        // Tidak reset lastEditedFieldIdRef agar Quick Guide masih bisa insert
    }, []);

    // Drag & Drop functions
    const findFieldPath = useCallback((fields: Field[], targetId: string): { parent: Field[] | null, index: number, field: Field | null } => {
        for (let i = 0; i < fields.length; i++) {
            if (fields[i].id === targetId) {
                return { parent: null, index: i, field: fields[i] };
            }
            if (fields[i].children.length > 0) {
                for (let j = 0; j < fields[i].children.length; j++) {
                    if (fields[i].children[j].id === targetId) {
                        return { parent: fields[i].children, index: j, field: fields[i].children[j] };
                    }
                }
                const deeper = findFieldPath(fields[i].children, targetId);
                if (deeper.field) return deeper;
            }
        }
        return { parent: null, index: -1, field: null };
    }, []);

    const removeFieldById = useCallback((fields: Field[], targetId: string): Field[] | null => {
        const filtered = fields.filter(f => f.id !== targetId);
        if (filtered.length !== fields.length) return filtered;

        for (let i = 0; i < fields.length; i++) {
            if (fields[i].children.length > 0) {
                const updated = removeFieldById(fields[i].children, targetId);
                if (updated) {
                    const newFields = [...fields];
                    newFields[i] = { ...newFields[i], children: updated };
                    return newFields;
                }
            }
        }
        return null;
    }, []);

    const insertFieldAt = useCallback((fields: Field[], targetId: string, sourceField: Field, position: 'before' | 'after' | 'inside'): Field[] | null => {
        for (let i = 0; i < fields.length; i++) {
            if (fields[i].id === targetId) {
                const newFields = [...fields];
                if (position === 'before') {
                    newFields.splice(i, 0, { ...sourceField, id: genId() });
                } else if (position === 'after') {
                    newFields.splice(i + 1, 0, { ...sourceField, id: genId() });
                } else if (position === 'inside' && fields[i].type === 'object') {
                    newFields[i] = {
                        ...newFields[i],
                        children: [...newFields[i].children, { ...sourceField, id: genId() }]
                    };
                }
                return newFields;
            }

            if (fields[i].children.length > 0) {
                for (let j = 0; j < fields[i].children.length; j++) {
                    if (fields[i].children[j].id === targetId) {
                        const newChildren = [...fields[i].children];
                        if (position === 'before') {
                            newChildren.splice(j, 0, { ...sourceField, id: genId() });
                        } else if (position === 'after') {
                            newChildren.splice(j + 1, 0, { ...sourceField, id: genId() });
                        } else if (position === 'inside' && fields[i].children[j].type === 'object') {
                            newChildren[j] = {
                                ...newChildren[j],
                                children: [...newChildren[j].children, { ...sourceField, id: genId() }]
                            };
                        }
                        const newFields = [...fields];
                        newFields[i] = { ...newFields[i], children: newChildren };
                        return newFields;
                    }
                }

                const updated = insertFieldAt(fields[i].children, targetId, sourceField, position);
                if (updated) {
                    const newFields = [...fields];
                    newFields[i] = { ...newFields[i], children: updated };
                    return newFields;
                }
            }
        }
        return null;
    }, []);

    const handleDragStart = useCallback((e: React.DragEvent, nodeId: string) => {
        e.dataTransfer.setData("text/plain", nodeId);
        e.dataTransfer.effectAllowed = "move";
        setDragSourceId(nodeId);
    }, []);

    const handleDragOver = useCallback((e: React.DragEvent, nodeId: string) => {
        e.preventDefault();
        if (dragSourceId === nodeId) return;
        setDragOverId(nodeId);
    }, [dragSourceId]);

    const handleDrop = useCallback((e: React.DragEvent, targetId: string) => {
        e.preventDefault();
        const sourceId = e.dataTransfer.getData("text/plain");

        if (!sourceId || sourceId === targetId) {
            setDragOverId(null);
            setDragSourceId(null);
            return;
        }

        const rect = e.currentTarget.getBoundingClientRect();
        const relativeY = (e.clientY - rect.top) / rect.height;

        let position: 'before' | 'after' | 'inside' = 'after';
        if (relativeY < 0.25) position = 'before';
        else if (relativeY > 0.75) position = 'after';
        else position = 'inside';

        const sourcePath = findFieldPath(fields, sourceId);
        if (!sourcePath.field) {
            setDragOverId(null);
            setDragSourceId(null);
            return;
        }

        let fieldsWithoutSource = removeFieldById(fields, sourceId) || [...fields];
        const sourceField = { ...sourcePath.field };
        const updatedFields = insertFieldAt(fieldsWithoutSource, targetId, sourceField, position);

        if (updatedFields) updateAndSync(updatedFields);

        setDragOverId(null);
        setDragSourceId(null);
    }, [fields, findFieldPath, removeFieldById, insertFieldAt]);

    const handleDragLeave = useCallback(() => {
        setDragOverId(null);
    }, []);

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between flex-wrap gap-3">
                <div>
                    <Label className="text-sm font-semibold flex items-center gap-2">
                        <Braces className="h-4 w-4 text-blue-500" />
                        Request Payload Builder
                    </Label>
                    <p className="text-xs text-muted-foreground mt-0.5">
                        Build your API request structure visually - click on value badges to edit
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button type="button" variant="outline" size="sm" className="h-8 text-xs" onClick={() => addObject(null)}>
                        <Braces className="h-3.5 w-3.5 mr-1 text-blue-500" /> + Object
                    </Button>
                    <Button type="button" variant="outline" size="sm" className="h-8 text-xs" onClick={() => addField(null)}>
                        <Type className="h-3.5 w-3.5 mr-1" /> + Field
                    </Button>
                </div>
            </div>

            {/* ── FIX: error banner untuk variabel tidak dikenal ── */}
            {varError && (
                <div className="flex items-start gap-2 bg-red-50 dark:bg-red-950/30 border border-red-200 dark:border-red-800/50 rounded-lg px-3 py-2 text-xs text-red-700 dark:text-red-400">
                    <X className="h-3.5 w-3.5 mt-0.5 shrink-0 text-red-500" />
                    <span>{varError}</span>
                </div>
            )}

            <div className="border rounded-xl p-4 min-h-[200px] bg-white dark:bg-slate-950 shadow-sm">
                {fields.length === 0 ? (
                    <div className="text-center py-12">
                        <div className="w-16 h-16 mx-auto mb-3 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center">
                            <Braces className="h-8 w-8 text-slate-400" />
                        </div>
                        <p className="font-medium text-slate-600 dark:text-slate-400">Empty Payload</p>
                        <p className="text-xs text-muted-foreground mt-1">Click "+ Object" or "+ Field" to start building</p>
                    </div>
                ) : (
                    <div className="space-y-1">
                        {fields.map((field) => (
                            <FieldRow
                                key={field.id}
                                field={field}
                                depth={0}
                                dragOverId={dragOverId}
                                onDragStart={handleDragStart}
                                onDragOver={handleDragOver}
                                onDrop={handleDrop}
                                onDragLeave={handleDragLeave}
                                onToggleExpand={toggleExpand}
                                onUpdateField={updateFieldValue}
                                onRemoveField={removeField}
                                onAddField={addField}
                                onAddObject={addObject}
                                onEditValue={handleEditValue}
                                editingFieldId={editingFieldId}
                                onCancelEdit={handleCancelEdit}
                                onInsertVariable={handleInsertVariable}
                                availableVariables={availableVariables}
                            />
                        ))}
                    </div>
                )}
            </div>

            {showHelper && (
                <div className="bg-blue-50 dark:bg-blue-950/30 rounded-lg p-3 border border-blue-200 dark:border-blue-800/50 space-y-2">
                    <div className="flex items-center gap-2">
                        <HelpCircle className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                        <span className="text-xs font-semibold text-blue-700 dark:text-blue-300 uppercase tracking-wide">Quick Guide</span>
                    </div>

                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 text-xs">
                        <div className="flex items-start gap-2">
                            <span className="text-blue-500 font-bold">🎨</span>
                            <span className="text-slate-600 dark:text-slate-400"><strong className="text-slate-800 dark:text-slate-200">Click on badges</strong> to edit values</span>
                        </div>
                        <div className="flex items-start gap-2">
                            <span className="text-blue-500 font-bold">📦</span>
                            <span className="text-slate-600 dark:text-slate-400"><strong className="text-slate-800 dark:text-slate-200">Objects</strong> contain nested fields</span>
                        </div>
                        <div className="flex items-start gap-2">
                            <span className="text-blue-500 font-bold">🔧</span>
                            <span className="text-slate-600 dark:text-slate-400"><strong className="text-slate-800 dark:text-slate-200">Drag & Drop</strong> fields to reorder</span>
                        </div>
                    </div>

                    <div className="border-t border-blue-200 dark:border-blue-800/50 pt-2 mt-1">
                        <span className="text-[10px] font-mono text-blue-600 dark:text-blue-400 block mb-1.5">
                            Available Variables — klik untuk insert ke field terakhir yang diedit:
                        </span>
                        <div className="flex flex-wrap gap-1.5">
                            {availableVariables.map((v) => (
                                <Badge
                                    key={v}
                                    tone="info"
                                    className="cursor-pointer hover:bg-blue-500 hover:text-white transition-colors py-0.5 font-mono text-xs select-none"
                                    onClick={() => handleInsertVariable(v)}
                                >
                                    {v}
                                </Badge>
                            ))}
                        </div>
                        <p className="text-[10px] text-blue-600/70 dark:text-blue-400/70 mt-1.5 italic">
                            💡 Tip: Klik badge value untuk edit, lalu klik variable untuk insert. Variable di luar daftar ini tidak akan diterima.
                        </p>
                    </div>
                </div>
            )}
        </div>
    );
}